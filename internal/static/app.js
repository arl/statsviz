import * as stats from './stats.js';
import * as ui from './ui.js';

const $ = id => {
    return document.getElementById(id);
}

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
    let ws = new WebSocket(buildWebsocketURI());
    console.info("Attempting websocket connection to statsviz server...");

    ws.onopen = () => {
        console.info("Successfully connected");
        timeout = 250; // reset connection timeout for next time
    };

    ws.onclose = event => {
        console.info("Closed websocket connection: ", event);
        setTimeout(connect, clamp(timeout += timeout, 250, 5000));
    };

    ws.onerror = error => {
        console.error("Websocket error: ", error);
        ws.close();
    };

    var initDone = false;
    ws.onmessage = event => {
        let allStats = JSON.parse(event.data)
        if (!initDone) {
            const sizeClasses = extractSizeClasses(allStats);
            const plotdefs = createPlotDefs(sizeClasses);
            ui.configurePlots(plotdefs);

            stats.init(plotdefs, dataRetentionSeconds);
            stats.pushData(new Date(), allStats);

            const data = stats.slice(dataRetentionSeconds);
            ui.attachPlots(data);

            initDone = true;
            return;
        }

        stats.pushData(new Date(), allStats);
        if (ui.isPaused()) {
            return
        }
        let data = stats.slice(dataRetentionSeconds);
        ui.updatePlots(data);
    }
}

connect();


// TODO(arl) -> this should be defined in Go in the init message.
const extractSizeClasses = (allStats) => {
    const sizeClasses = new Array(allStats.Mem.BySize.length);
    for (let i = 0; i < allStats.Mem.BySize.length; i++) {
        sizeClasses.push(allStats.Mem.BySize[i].Size);
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
        config: {
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
                'name': 'heap alloc',
                'unitfmt': '%{y:.4s}B',
            }, {
                'name': 'heap sys',
                'unitfmt': '%{y:.4s}B',
            }, {
                'name': 'heap idle',
                'unitfmt': '%{y:.4s}B',
            }, {
                'name': 'heap in-use',
                'unitfmt': '%{y:.4s}B',
            }, {
                'name': 'heap next gc',
                'unitfmt': '%{y:.4s}B',
            }, ],
        },
    }, {
        config: {
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
                'name': 'live',
                'hover': 'live objects',
                'unitfmt': '%{y}',
            }, {
                'name': 'lookups',
                'hover': 'pointer lookups',
                'unitfmt': '%{y}',
            }, {
                'name': 'heap',
                'hover': 'heap objects',
                'unitfmt': '%{y}',
            }, ],
        },
    }, {
        config: {
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
                'name': 'mspan in-use',
                'unitfmt': '%{y:.4s}B',
            }, {
                'name': 'mspan sys',
                'unitfmt': '%{y:.4s}B',
            }, {
                'name': 'mcache in-use',
                'unitfmt': '%{y:.4s}B',
            }, {
                'name': 'mcache sys',
                'unitfmt': '%{y:.4s}B',
            }, ],
        },
    }, {
        config: {
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
                'name': 'goroutines',
                'unitfmt': '%{y}',
            }],
        },
    }, {
        config: {
            name: 'bySizes',
            title: 'Size Classes',
            type: 'heatmap',
            updateFreq: 5,
            hasHorsEvents: false,
            layout: {
                yaxis: {
                    title: 'size classes',
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
        },
    }, {
        config: {
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
                'name': 'gc/cpu',
                'hover': 'gc/cpu fraction',
                'unitfmt': '%{y:,.4%}',
            }],
        },
    }, ];
}