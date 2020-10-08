// ui holds the user interface state
var ui = (function () {
    var m = {};

    let paused = false;

    m.isPaused = function () { return paused; }
    m.togglePause = function () { paused = !paused; }
    m.plots = null;

    function GCLines(data) {
        const gcs = stats.lastGCs;
        const mints = data.times[0];
        const maxts = data.times[data.times.length - 1];

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

    function heapData(data) {
        return [
            {
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

    function mspanMCacheData(data) {
        return [
            {
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

    let mspanMCacheLayout = {
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

    function sizeClassesData(data) {
        var ret = [
            {
                x: data.times,
                y: stats.classSizes,
                z: data.bySizes,
                type: 'heatmap',
                hovertemplate: '<br><b>size class</b>: %{y:} B' +
                    '<br><b>objects</b>: %{z}<br>',
                showlegend: false,
                colorscale: colorscale,
            }
        ];
        return ret;
    }

    let sizeClassesLayout = {
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

    function objectsData(data) {
        return [
            {
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

    let objectsLayout = {
        title: 'Objects',
        xaxis: {
            title: 'time',
            tickformat: '%H:%M:%S',
        },
        yaxis: {
            title: 'objects'
        }
    };

    function goroutinesData(data) {
        return [
            {
                x: data.times,
                y: data.goroutines,
                type: 'scatter',
                name: 'goroutines',
                hovertemplate: '<b>goroutines</b>: %{y}',
            },
        ]
    }

    let goroutinesLayout = {
        title: 'Goroutines',
        xaxis: {
            title: 'time',
            tickformat: '%H:%M:%S',
        },
        yaxis: {
            title: 'goroutines',
        }
    };

    function gcFractionData(data) {
        return [
            {
                x: data.times,
                y: data.gcfraction,
                type: 'scatter',
                name: 'gc/cpu',
                hovertemplate: '<b>gcc/CPU fraction</b>: %{y:,.4%}',
            },
        ]
    }

    let gcFractionLayout = {
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

    const config = {
        displayModeBar: false,
    }

    let heapElt = null;
    let mspanMCacheElt = null;
    let sizeClassElt = null;
    let objectsElt = null;
    let gcfractionElt = null;
    let goroutinesElt = null;

    m.createPlots = function (data) {
        // $(".ui.accordion").accordion();
        $('.ui.accordion').accordion({
            exclusive: false,
            onOpen: function () {
                this.firstElementChild.hidden = false;
            },
            onClose: function () {
                this.firstElementChild.hidden = true;
            }
        });

        heapElt = document.getElementById('heap');
        heapElt.hidden = false;

        mspanMCacheElt = document.getElementById('mspan-mcache');
        sizeClassesElt = document.getElementById('size-classes');
        objectsElt = document.getElementById('objects');
        gcfractionElt = document.getElementById('gcfraction');
        goroutinesElt = document.getElementById('goroutines');

        Plotly.plot(heapElt, heapData(data), heapLayout, config);
        Plotly.plot(mspanMCacheElt, mspanMCacheData(data), mspanMCacheLayout, config);
        Plotly.plot(sizeClassesElt, sizeClassesData(data), sizeClassesLayout, config);
        Plotly.plot(objectsElt, objectsData(data), objectsLayout, config);
        Plotly.plot(gcfractionElt, gcFractionData(data), gcFractionLayout, config);
        Plotly.plot(goroutinesElt, goroutinesData(data), goroutinesLayout, config);

        mspanMCacheElt.hidden = true;
        sizeClassesElt.hidden = true;
        objectsElt.hidden = true;
        gcfractionElt.hidden = true;
        goroutinesElt.hidden = true;
    }

    var updateIdx = 0;
    m.updatePlots = function (data) {
        let gcLines = GCLines(data);

        heapLayout.shapes = gcLines;
        if (!heapElt.hidden) {
            Plotly.react(heapElt, heapData(data), heapLayout, config);
            console.log("updating: heap");
        }

        mspanMCacheLayout.shapes = gcLines;
        if (!mspanMCacheElt.hidden) {
            Plotly.react(mspanMCacheElt, mspanMCacheData(data), mspanMCacheLayout, config);
            console.log("updating: mspan");
        }

        objectsLayout.shapes = gcLines;
        if (!objectsElt.hidden) {
            Plotly.react(objectsElt, objectsData(data), objectsLayout, config);
            console.log("updating: objects");
        }

        if (!gcfractionElt.hidden) {
            Plotly.react(gcfractionElt, gcFractionData(data), gcFractionLayout, config);
            console.log("updating: gcfracion");
        }
        if (!goroutinesElt.hidden) {
            Plotly.react(goroutinesElt, goroutinesData(data), goroutinesLayout, config);
            console.log("updating: goroutines");
        }
        if (!sizeClassesElt.hidden && updateIdx % 5 == 0) {
            // Update the size class heatmap 5 times less often since it's expensive. 
            Plotly.react(sizeClassesElt, sizeClassesData(data), sizeClassesLayout, config);
            console.log("updating: heatmap");
        }

        updateIdx++;
    }

    function traceInfo(traceName) {
        let traces = {
            'heap alloc': 'HeapAlloc',
            'heap sys': 'HeapSys',
            'heap idle': 'HeapIdle',
            'heap in-use': 'HeapInuse',
            'next gc': 'NextGC',

            'mspan in-use': 'MSpanInuse',
            'mspan sys': 'MSpanSys',
            'mcache in-use': 'MCacheInuse',
            'mcache sys': 'MCacheSys',

            'gcfraction': 'GCCPUFraction',

            'lookups': 'Lookups',
            'heap objects': 'HeapObjects',
        };

        let fieldName = traces[traceName];
        if (fieldName !== undefined) {
            return memStatsDoc(fieldName);
        }
        if (traceName == 'goroutines') {
            return "The number of goroutines"
        }
        if (traceName == 'live') {
            return "The number of live objects"
        }
        if (traceName == 'goroutines') {
            return "Number of the goroutines"
        }
        if (traceName == 'size classes') {
            return "Reports per-size class allocation statistics"
        }
    };

    return m;
}());