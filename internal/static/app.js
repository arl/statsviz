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
            const sizeClasses = extractSizeClasses(allStats);
            plotdefs = createPlotDefs(sizeClasses);
            configurePlots(plotdefs);

            stats.init(plotdefs, dataRetentionSeconds);

            attachPlots();

            initDone = true;
            return;
        }

        stats.pushData(new Date(), allStats);
        if (isPaused()) {
            return
        }
        updatePlots(stats.slice(dataRetentionSeconds));
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
    plotdefs.forEach(plotdef => {
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

const updatePlots = data => {
    let gcLines = GCLines(data);

    plots.forEach(plot => {
        if (!plot.hidden) {
            plot.update(data, gcLines);
        }
    });
}

const GCLines = data => {
    const gcs = stats.lastGCs;
    const mints = data.times[0];
    const maxts = data.times[data.times.length - 1];

    const shapes = [];
    for (let i = 0, n = gcs.length; i < n; i++) {
        let d = gcs[i];
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

const createPlotDefs = (sizeClasses) => {
    return [{
        name: "heap",
        title: 'Heap',
        type: 'scatter',
        updateFreq: 0,
        hasHorsEvents: true,
        layout: {
            yaxis: {
                title: 'bytes',
                ticksuffix: 'B',
            },
        },
        subplots: [{
            name: 'heap alloc',
            unitfmt: '%{y:.4s}B',
            datapath: (raw) => { return raw.Mem.HeapAlloc; },
        }, {
            name: 'heap sys',
            unitfmt: '%{y:.4s}B',
            datapath: (raw) => { return raw.Mem.HeapSys; },
        }, {
            name: 'heap idle',
            unitfmt: '%{y:.4s}B',
            datapath: (raw) => { return raw.Mem.HeapIdle; },
        }, {
            name: 'heap in-use',
            unitfmt: '%{y:.4s}B',
            datapath: (raw) => { return raw.Mem.HeapInuse; },
        }, {
            name: 'heap next gc',
            unitfmt: '%{y:.4s}B',
            datapath: (raw) => { return raw.Mem.NextGC; },
        }, ],
    }, {
        name: "objects",
        title: 'Objects',
        type: 'scatter',
        updateFreq: 0,
        hasHorsEvents: true,
        layout: {
            yaxis: {
                title: 'objects',
            },
        },
        subplots: [{
            name: 'live',
            hover: 'live objects',
            unitfmt: '%{y}',
            datapath: (raw) => { return raw.Mem.Mallocs - raw.Mem.Freed; },
        }, {
            name: 'lookups',
            hover: 'pointer lookups',
            unitfmt: '%{y}',
            datapath: (raw) => { return raw.Mem.Lookups; },
        }, {
            name: 'heap',
            hover: 'heap objects',
            unitfmt: '%{y}',
            datapath: (raw) => { return raw.Mem.HeapObjects; },
        }, ],
    }, {
        name: 'mspanMCache',
        title: 'MSpan/MCache',
        type: 'scatter',
        updateFreq: 0,
        hasHorsEvents: true,
        layout: {
            yaxis: {
                title: 'bytes',
                ticksuffix: 'B',
            },
        },
        subplots: [{
            name: 'mspan in-use',
            unitfmt: '%{y:.4s}B',
            datapath: (raw) => { return raw.Mem.MSpanInuse; },
        }, {
            name: 'mspan sys',
            unitfmt: '%{y:.4s}B',
            datapath: (raw) => { return raw.Mem.MSpanSys; },
        }, {
            name: 'mcache in-use',
            unitfmt: '%{y:.4s}B',
            datapath: (raw) => { return raw.Mem.MCacheInuse; },
        }, {
            name: 'mcache sys',
            unitfmt: '%{y:.4s}B',
            datapath: (raw) => { return raw.Mem.MCacheSys; },
        }, ],
    }, {
        name: 'goroutines',
        title: 'Goroutines',
        type: 'scatter',
        updateFreq: 0,
        hasHorsEvents: false,
        layout: {
            yaxis: {
                title: 'goroutines',
            },
        },
        subplots: [{
            name: 'goroutines',
            unitfmt: '%{y}',
            datapath: (raw) => { return raw.NumGoroutine; },
        }],
    }, {
        name: 'bySizes',
        title: 'Size Classes',
        type: 'heatmap',
        updateFreq: 5,
        hasHorsEvents: false,
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
            datapath: (raw, i) => {
                // TODO(arl) must receive an already computed array
                const size = raw.Mem.BySize[i];
                return size.Mallocs - size.Frees;
            },
        },
    }, {
        name: "gcfraction",
        title: 'GC CPU fraction',
        type: 'scatter',
        updateFreq: 0,
        hasHorsEvents: false,
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
            datapath: (raw) => { return raw.Mem.GCCPUFraction; },
        }],
    }, ];
}