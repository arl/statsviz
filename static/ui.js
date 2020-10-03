// ui holds the user interface state
var ui = (function () {
    var m = {};

    let paused = false;

    m.isPaused = function () { return paused; }
    m.togglePause = function () { paused = !paused; }
    m.plots = null;

    function heapData(data) {
        return [
            {
                x: data.heap[0],
                y: data.heap[1],
                type: 'scatter',
                name: 'heap alloc'
            },
            {
                x: data.heap[0],
                y: data.heap[2],
                type: 'scatter',
                name: 'heap sys'
            },
            {
                x: data.heap[0],
                y: data.heap[3],
                type: 'scatter',
                name: 'heap idle'
            },
            {
                x: data.heap[0],
                y: data.heap[4],
                type: 'scatter',
                name: 'heap in-use'
            },
        ]
    }

    // https://plotly.com/javascript/reference/layout
    let heapLayout = {
        title: 'Heap',
        xaxis: {
            title: 'time',
            tickformat: '%H:%M:%S',
        },
        yaxis: {
            title: 'bytes',
            ticksuffix: 'B',
            // tickformat: ' ',
            exponentformat: 'SI',
        }
    };

    m.createPlots = function (opts, data, elts) {
        Plotly.plot('heap', heapData(data), heapLayout);
    }


    function dateFromTimestamp(ts) {
        return
    }

    function GCLines(data) {
        const gcs = stats.lastGCs;
        const mints = data.heap[0][0];
        const maxts = data.heap[0][data.heap[0].length - 1];

        let shapes = [];

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

    m.updatePlots = function (xScale, data) {
        heapLayout.shapes = GCLines(data);
        Plotly.react('heap', heapData(data), heapLayout)
    }

    return m;
}());