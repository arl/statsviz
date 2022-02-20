// ui holds the user interface state
import { classSizes, lastGCs } from './stats.js';

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

const heapData = data => {
    return [{
            x: data.times,
            y: data.heap[0],
            type: 'scatter',
            name: 'heap alloc',
            hovertemplate: '<b>heap alloc</b>: %{y:.4s}B',
        },
        {
            x: data.times,
            y: data.heap[1],
            type: 'scatter',
            name: 'heap sys',
            hovertemplate: '<b>heap sys</b>: %{y:.4s}B',
        },
        {
            x: data.times,
            y: data.heap[2],
            type: 'scatter',
            name: 'heap idle',
            hovertemplate: '<b>heap idle</b>: %{y:.4s}B',
        },
        {
            x: data.times,
            y: data.heap[3],
            type: 'scatter',
            name: 'heap in-use',
            hovertemplate: '<b>heap in-use</b>: %{y:.4s}B',
        },
        {
            x: data.times,
            y: data.heap[4],
            type: 'scatter',
            name: 'next gc',
            hovertemplate: '<b>next gc</b>: %{y:.4s}B',
        },
    ]
}

// https://plotly.com/javascript/reference/layout
const heapLayout = {
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

const mspanMCacheData = data => {
    return [{
            x: data.times,
            y: data.mspanMCache[0],
            type: 'scatter',
            name: 'mspan in-use',
            hovertemplate: '<b>mspan in-use</b>: %{y:.4s}B',
        },
        {
            x: data.times,
            y: data.mspanMCache[1],
            type: 'scatter',
            name: 'mspan sys',
            hovertemplate: '<b>mspan sys</b>: %{y:.4s}B',
        },
        {
            x: data.times,
            y: data.mspanMCache[2],
            type: 'scatter',
            name: 'mcache in-use',
            hovertemplate: '<b>mcache in-use</b>: %{y:.4s}B',
        },
        {
            x: data.times,
            y: data.mspanMCache[3],
            type: 'scatter',
            name: 'mcache sys',
            hovertemplate: '<b>mcache sys</b>: %{y:.4s}B',
        },
    ]
}

const mspanMCacheLayout = {
    title: 'MSpan/MCache',
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

const colorscale = [
    [0, 'rgb(166,206,227, 0.5)'],
    [0.05, 'rgb(31,120,180,0.5)'],
    [0.2, 'rgb(178,223,138,0.5)'],
    [0.5, 'rgb(51,160,44,0.5)'],
    [1, 'rgb(227,26,28,0.5)']
];

const sizeClassesData = data => {
    var ret = [{
        x: data.times,
        y: classSizes,
        z: data.bySizes,
        type: 'heatmap',
        hovertemplate: '<br><b>size class</b>: %{y:} B' +
            '<br><b>objects</b>: %{z}<br>',
        showlegend: false,
        colorscale: colorscale,
    }];
    return ret;
}

const sizeClassesLayout = {
    title: 'Size Classes',
    xaxis: {
        title: 'time',
        tickformat: '%H:%M:%S',
    },
    yaxis: {
        title: 'size classes',
        exponentformat: 'SI',
    }
};

const objectsData = data => {
    return [{
            x: data.times,
            y: data.objects[0],
            type: 'scatter',
            name: 'live',
            hovertemplate: '<b>live objects</b>: %{y}',
        },
        {
            x: data.times,
            y: data.objects[1],
            type: 'scatter',
            name: 'lookups',
            hovertemplate: '<b>pointer lookups</b>: %{y}',
        },
        {
            x: data.times,
            y: data.objects[2],
            type: 'scatter',
            name: 'heap',
            hovertemplate: '<b>heap objects</b>: %{y}',
        },
    ]
}

const objectsLayout = {
    title: 'Objects',
    xaxis: {
        title: 'time',
        tickformat: '%H:%M:%S',
    },
    yaxis: {
        title: 'objects'
    }
};

const goroutinesData = data => {
    return [{
        x: data.times,
        y: data.goroutines,
        type: 'scatter',
        name: 'goroutines',
        hovertemplate: '<b>goroutines</b>: %{y}',
    }, ]
}

const goroutinesLayout = {
    title: 'Goroutines',
    xaxis: {
        title: 'time',
        tickformat: '%H:%M:%S',
    },
    yaxis: {
        title: 'goroutines',
    }
};

const gcFractionData = data => {
    return [{
        x: data.times,
        y: data.gcfraction,
        type: 'scatter',
        name: 'gc/cpu',
        hovertemplate: '<b>gcc/CPU fraction</b>: %{y:,.4%}',
    }, ]
}

const gcFractionLayout = {
    title: 'GC CPU fraction',
    xaxis: {
        title: 'time',
        tickformat: '%H:%M:%S',
    },
    yaxis: {
        title: 'gc/cpu (%)',
        tickformat: ',.5%',
    }
};


const configs = () => {
    const plots = ['heap', 'mspan-mcache', 'size-classes', 'objects', 'gcfraction', 'goroutines'];
    const cfgs = {};

    plots.forEach(plotName => {
        // Create plot config where only 'save image' and 'show on hover' toggles are enabled.
        const config = {
            displaylogo: false,
            modeBarButtonsToRemove: ['2D', 'zoom2d', 'pan2d', 'select2d', 'lasso2d', 'zoomIn2d', 'zoomOut2d', 'autoScale2d', 'resetScale2d', 'toggleSpikelines'],
            toImageButtonOptions: {
                format: 'png',
                filename: plotName
            }
        }

        cfgs[plotName] = config;
    });

    return cfgs;
};

const heapElt = $('#heap')[0];
const mspanMCacheElt = $('#mspan-mcache')[0];
const sizeClassesElt = $('#size-classes')[0];
const objectsElt = $('#objects')[0];
const gcfractionElt = $('#gcfraction')[0];
const goroutinesElt = $('#goroutines')[0];

const createPlots = (data) => {
    $('.ui.accordion').accordion({
        exclusive: false,
        onOpen: function() {
            this.firstElementChild.hidden = false;
        },
        onClose: function() {
            this.firstElementChild.hidden = true;
        }
    });

    Plotly.newPlot(heapElt, heapData(data), heapLayout, configs['heap']);
    Plotly.newPlot(mspanMCacheElt, mspanMCacheData(data), mspanMCacheLayout, configs['mspan-mcache']);
    Plotly.newPlot(sizeClassesElt, sizeClassesData(data), sizeClassesLayout, configs['size-classes']);
    Plotly.newPlot(objectsElt, objectsData(data), objectsLayout, configs['objects']);
    Plotly.newPlot(gcfractionElt, gcFractionData(data), gcFractionLayout, configs['gcfraction']);
    Plotly.newPlot(goroutinesElt, goroutinesData(data), goroutinesLayout, configs['goroutines']);
}

var updateIdx = 0;

const updatePlots = data => {
    let gcLines = GCLines(data);

    heapLayout.shapes = gcLines;
    if (!heapElt.hidden) {
        Plotly.react(heapElt, heapData(data), heapLayout, configs['heap']);
    }

    mspanMCacheLayout.shapes = gcLines;
    if (!mspanMCacheElt.hidden) {
        Plotly.react(mspanMCacheElt, mspanMCacheData(data), mspanMCacheLayout, configs['mspan-mcache']);
    }

    objectsLayout.shapes = gcLines;
    if (!objectsElt.hidden) {
        Plotly.react(objectsElt, objectsData(data), objectsLayout, configs['objects']);
    }

    if (!gcfractionElt.hidden) {
        Plotly.react(gcfractionElt, gcFractionData(data), gcFractionLayout, configs['gcfraction']);
    }

    if (!goroutinesElt.hidden) {
        Plotly.react(goroutinesElt, goroutinesData(data), goroutinesLayout, configs['goroutines']);
    }

    if (!sizeClassesElt.hidden && updateIdx % 5 == 0) {
        // Update the size class heatmap 5 times less often since it's expensive. 
        Plotly.react(sizeClassesElt, sizeClassesData(data), sizeClassesLayout, configs['size-classes']);
    }

    updateIdx++;
}

let paused = false;
const isPaused = () => { return paused; }
const togglePause = () => { paused = !paused; }

export { isPaused, togglePause, createPlots, updatePlots };