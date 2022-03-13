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

const plotWidth = 620;
const plotHeight = 420;


// https://plotly.com/javascript/reference/layout
const heapLayout = {
    title: 'Heap',
    width: plotWidth,
    height: plotHeight,
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
    width: plotWidth,
    height: plotHeight,
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
    width: plotWidth,
    height: plotHeight,
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
    width: plotWidth,
    height: plotHeight,
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
    width: plotWidth,
    height: plotHeight,
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
    width: plotWidth,
    height: plotHeight,
    xaxis: {
        title: 'time',
        tickformat: '%H:%M:%S',
    },
    yaxis: {
        title: 'gc/cpu (%)',
        tickformat: ',.5%',
    }
};

let plots = [];

const createPlots = (data) => {
    const plotDefs = [
        { name: "heap", dataFunc: heapData, layout: heapLayout, updateFreq: 0, hasHorsEvents: true },
        { name: "objects", dataFunc: objectsData, layout: objectsLayout, updateFreq: 0, hasHorsEvents: true },
        { name: "mspan-mcache", dataFunc: mspanMCacheData, layout: mspanMCacheLayout, updateFreq: 0, hasHorsEvents: true },
        { name: "goroutines", dataFunc: goroutinesData, layout: goroutinesLayout, updateFreq: 0, hasHorsEvents: false },
        { name: "size-classes", dataFunc: sizeClassesData, layout: sizeClassesLayout, updateFreq: 5, hasHorsEvents: false },
        { name: "gcfraction", dataFunc: gcFractionData, layout: gcFractionLayout, updateFreq: 0, hasHorsEvents: false },
    ];

    let curRow = null;
    let container = $('#plots');

    for (let i = 0; i < plotDefs.length; i++) {
        const plotDef = plotDefs[i];
        if (i % 2 == 0) {
            curRow = $('<div>', { class: 'row' });
            container.append(curRow);
        }
        let col = $('<div>', { class: 'col' });
        let plotDiv = $('<div>', { id: plotDef.name });

        let plot = new Plot(plotDiv[0], plotDef.name, plotDef.dataFunc, plotDef.layout, data, plotDef.updateFreq, plotDef.hasHorsEvents);
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