// ui holds the user interface state
import { classSizes, lastGCs } from './stats.js';
import Plot from "./plot.js";

const GCLines = data => {
    const gcs = lastGCs;
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

const colorscale = [
    [0, 'rgb(166,206,227, 0.5)'],
    [0.05, 'rgb(31,120,180,0.5)'],
    [0.2, 'rgb(178,223,138,0.5)'],
    [0.5, 'rgb(51,160,44,0.5)'],
    [1, 'rgb(227,26,28,0.5)']
];

let plots = [];

const createPlots = (data) => {
    const plotDefs = [{
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
                buckets: classSizes,
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

    let curRow = null;
    let container = $('#plots');

    for (let i = 0; i < plotDefs.length; i++) {
        const plotDef = plotDefs[i];
        if (i % 2 == 0) {
            curRow = $('<div>', { class: 'row' });
            container.append(curRow);
        }
        let col = $('<div>', { class: 'col' });
        let plotDiv = $('<div>', { id: plotDef.config.name });

        let plot = new Plot(plotDef.config, plotDiv[0], data);
        plots.push(plot);

        col.append(plotDiv);
        curRow.append(col);
    };
}

const updatePlots = data => {
    let gcLines = GCLines(data);

    plots.forEach(plot => {
        if (!plot.hidden) {
            plot.update(data, gcLines);
        }
    });
}

let paused = false;
const isPaused = () => { return paused; }
const togglePause = () => { paused = !paused; }

export { isPaused, togglePause, createPlots, updatePlots };