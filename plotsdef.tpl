export default {
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
            colorscale: [
                [0, "rgb(166,206,227, 0.5)"],
                [0.05, "rgb(31,120,180,0.5)"],
                [0.2, "rgb(178,223,138,0.5)"],
                [0.5, "rgb(51,160,44,0.5)"],
                [1, "rgb(227,26,28,0.5)"]
            ],
            buckets: [
                0,
                8,
                16,
                24,
                32,
                48,
                64,
                80,
                96,
                112,
                128,
                144,
                160,
                176,
                192,
                208,
                224,
                240,
                256,
                288,
                320,
                352,
                384,
                416,
                448,
                480,
                512,
                576,
                640,
                704,
                768,
                896,
                1024,
                1152,
                1280,
                1408,
                1536,
                1792,
                2048,
                2304,
                2688,
                3072,
                3200,
                3456,
                4096,
                4864,
                5376,
                6144,
                6528,
                6784,
                6912,
                8192,
                9472,
                9728,
                10240,
                10880,
                12288,
                13568,
                14336,
                16384,
                18432
            ],
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