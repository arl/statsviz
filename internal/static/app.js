import * as stats from './stats.js';
import Plot from "./plot.js";


const buildWebsocketURI = () => {
    var loc = window.location,
        ws_prot = "ws:";
    if (loc.protocol === "https:") {
        ws_prot = "wss:";
    }
    return ws_prot + "//" + loc.host + loc.pathname + "ws"
}

const dataRetentionSeconds = 60;
var timeout = 250;

const clamp = (val, min, max) => {
    if (val < min) return min;
    if (val > max) return max;
    return val;
}

/* WebSocket connection handling */

const connect = () => {
    const uri = buildWebsocketURI();
    let ws = new WebSocket(uri);
    console.info(`Attempting websocket connection to server at ${uri}`);

    ws.onopen = () => {
        console.info("Successfully connected");
        timeout = 250; // reset connection timeout for next time
    };

    ws.onclose = event => {
        console.error(`Closed websocket connection: code ${event.code}`);
        setTimeout(connect, clamp(timeout += timeout, 250, 5000));
    };

    ws.onerror = err => {
        console.error(`Websocket error, closing connection.`);
        ws.close();
    };

    let initDone = false;
    let plotdefs = null;
    ws.onmessage = event => {
        let allStats = JSON.parse(event.data)

        if (!initDone) {
            // TODO: size classes should be defined in the 'init' message.
            const sizeClasses = extractSizeClasses(allStats);
            plotdefs = createPlotDefs(sizeClasses);
            configurePlots(plotdefs);

            stats.init(plotdefs, dataRetentionSeconds);

            attachPlots();

            initDone = true;
            return;
        }

        const converted = convertData(allStats)
        stats.pushData(new Date(), converted);
        if (isPaused()) {
            return
        }
        updatePlots(stats.slice(dataRetentionSeconds), plotdefs.events);
    }
}

connect();

/* plots management */

// TODO(arl) not used for now
let paused = false;
const isPaused = () => { return paused; }
const togglePause = () => { paused = !paused; }
let plots = [];

const configurePlots = (plotdefs) => {
    plots = [];
    plotdefs.series.forEach(plotdef => {
        plots.push(new Plot(plotdef));
    });
}

const attachPlots = () => {
    let row = null;
    let plotsDiv = $('#plots');
    plotsDiv.empty()

    for (let i = 0; i < plots.length; i++) {
        const plot = plots[i];
        if (i % 2 == 0) {
            row = $('<div>', { class: 'row' });
            plotsDiv.append(row);
        }

        let col = $('<div>', { class: 'col' });
        let div = $('<div>', { id: plot.name() });

        plot.createElement(div[0])
        col.append(div);
        row.append(col);
    }
}

const updatePlots = (data) => {
    // Create shapes.
    let shapes = new Map();

    for (const [eventName, eventSerie] of data.events) {
        shapes.set(eventName, createEventShape(data, eventSerie));
    }

    plots.forEach(plot => {
        if (!plot.hidden) {
            plot.update(data, shapes);
        }
    });
}

const createEventShape = (data, eventSerie) => {
    // TODO(arl): do we really need to pass 'data' to extract mints and maxtx?
    // aren't event serie already clamped to the visible time range?
    const mints = data.times[0];
    const maxts = data.times[data.times.length - 1];

    const shapes = [];
    for (let i = 0, n = eventSerie.length; i < n; i++) {
        let d = eventSerie[i];
        // Clamp GC times which are out of bounds
        if (d < mints || d > maxts) {
            continue;
        }

        shapes.push({
            type: 'line',
            x0: d,
            x1: d,
            yref: 'paper',
            y0: 0,
            y1: 1,
            line: {
                color: 'rgb(55, 128, 191)',
                width: 1,
                dash: 'longdashdot',
            }
        })
    }
    return shapes;
}

/* plots definition
 * 
 * (TODO(arl) -> will be defined in Go and read in the 'init' ws message
 */

const extractSizeClasses = (allStats) => {
    const sizeClasses = new Array(allStats.Mem.BySize.length);
    for (let i = 0; i < sizeClasses.length; i++) {
        sizeClasses[i] = allStats.Mem.BySize[i].Size;
    }
    return sizeClasses;
}

const colorscale = [
    [0, 'rgb(166,206,227, 0.5)'],
    [0.05, 'rgb(31,120,180,0.5)'],
    [0.2, 'rgb(178,223,138,0.5)'],
    [0.5, 'rgb(51,160,44,0.5)'],
    [1, 'rgb(227,26,28,0.5)']
];

const convertData = (raw) => {

    let bySizes = new Array(raw.Mem.BySize.length);
    for (let i = 0; i < raw.Mem.BySize.length; i++) {
        const size = raw.Mem.BySize[i];
        bySizes[i] = size.Mallocs - size.Frees;
    }

    return {
        'heap': [
            raw.Mem.HeapAlloc,
            raw.Mem.HeapSys,
            raw.Mem.HeapIdle,
            raw.Mem.HeapInuse,
            raw.Mem.NextGC,
        ],
        'objects': [
            raw.Mem.Mallocs - raw.Mem.Freed,
            raw.Mem.Lookups,
            raw.Mem.HeapObjects,
        ],
        'mspanMCache': [
            raw.Mem.MSpanInuse,
            raw.Mem.MSpanSys,
            raw.Mem.MCacheInuse,
            raw.Mem.MCacheSys,
        ],
        'goroutines': [
            raw.NumGoroutine,
        ],
        'bySizes': bySizes,
        'gcfraction': [
            raw.Mem.GCCPUFraction,
        ],
        // Event serie, used for vertical lines on plots (via plotly 'shapes').
        // This get automatically deduplicated in javascript.
        'lastgc': [
            raw.Mem.LastGC,
        ],
    };
}

const createPlotDefs = (sizeClasses) => {
    return {
        "events": ["lastgc"],
        "series": [{
            name: "heap",
            title: 'Heap',
            type: 'scatter',
            updateFreq: 0,
            horzEvents: 'lastgc',
            layout: {
                yaxis: {
                    title: 'bytes',
                    ticksuffix: 'B',
                },
            },
            subplots: [{
                name: 'heap alloc',
                unitfmt: '%{y:.4s}B',
            }, {
                name: 'heap sys',
                unitfmt: '%{y:.4s}B',
            }, {
                name: 'heap idle',
                unitfmt: '%{y:.4s}B',
            }, {
                name: 'heap in-use',
                unitfmt: '%{y:.4s}B',
            }, {
                name: 'heap next gc',
                unitfmt: '%{y:.4s}B',
            }, ],
        }, {
            name: "objects",
            title: 'Objects',
            type: 'scatter',
            updateFreq: 0,
            horzEvents: 'lastgc',
            layout: {
                yaxis: {
                    title: 'objects',
                },
            },
            subplots: [{
                name: 'live',
                hover: 'live objects',
                unitfmt: '%{y}',
            }, {
                name: 'lookups',
                hover: 'pointer lookups',
                unitfmt: '%{y}',
            }, {
                name: 'heap',
                hover: 'heap objects',
                unitfmt: '%{y}',
            }, ],
        }, {
            name: 'mspanMCache',
            title: 'MSpan/MCache',
            type: 'scatter',
            updateFreq: 0,
            horzEvents: 'lastgc',
            layout: {
                yaxis: {
                    title: 'bytes',
                    ticksuffix: 'B',
                },
            },
            subplots: [{
                name: 'mspan in-use',
                unitfmt: '%{y:.4s}B',
            }, {
                name: 'mspan sys',
                unitfmt: '%{y:.4s}B',
            }, {
                name: 'mcache in-use',
                unitfmt: '%{y:.4s}B',
            }, {
                name: 'mcache sys',
                unitfmt: '%{y:.4s}B',
            }, ],
        }, {
            name: 'goroutines',
            title: 'Goroutines',
            type: 'scatter',
            updateFreq: 0,
            horzEvents: '',
            layout: {
                yaxis: {
                    title: 'goroutines',
                },
            },
            subplots: [{
                name: 'goroutines',
                unitfmt: '%{y}',
            }],
        }, {
            name: 'bySizes',
            title: 'Size Classes',
            type: 'heatmap',
            updateFreq: 5,
            horzEvents: '',
            layout: {
                yaxis: {
                    title: 'size classes',
                    // TODO(arl) try also with log2 (not supported but we could recreate the ticks ourselves).
                    // see https://github.com/plotly/plotly.js/issues/4147#issuecomment-524378823
                    // type: 'log', 
                },
            },
            heatmap: {
                // TODO(arl) refine this, we should not pass all of that but probably have 
                // one hover and one unit for each of the 2 dimensions.
                hover: '<br><b>size class</b>: %{y:} B' +
                    '<br><b>objects</b>: %{z}<br>',
                colorscale: colorscale,
                buckets: sizeClasses,
            },
        }, {
            name: "gcfraction",
            title: 'GC CPU fraction',
            type: 'scatter',
            updateFreq: 0,
            horzEvents: '',
            layout: {
                yaxis: {
                    title: 'gc/cpu (%)',
                    tickformat: ',.5%',
                },
            },
            subplots: [{
                name: 'gc/cpu',
                hover: 'gc/cpu fraction',
                unitfmt: '%{y:,.4%}',
            }],
        }, ]
    };
}