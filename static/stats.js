var stats = (function () {
    var m = {};

    const idxTimes = 0;
    const idxHeapAlloc = 1;
    const idxHeapSys = 2;
    const idxHeapIdle = 3;
    const idxHeapInuse = 4;

    const numSeriesHeap = 4;
    const totalSeries = 1 + numSeriesHeap; // times + other series

    var data = {
        series: new Array(totalSeries), // TODO: rename to timeseries 
        lastGCs: new Array(),

        bySize: new Array(4), // [0] min [1] max [2] class sizes [3] live objects
    };

    m.init = function (buflen) {
        const extraBufferCapacity = 20; // 20% of extra (preallocated) buffer datapoints

        const bufcap = buflen + (buflen * extraBufferCapacity) / 100; // number of actual datapoints


        for (let i = 0; i < totalSeries; i++) {
            data.series[i] = new Buffer(buflen, bufcap);
        }
        for (let i = 0; i < data.bySize.length; i++) {
            data.bySize[i] = new Buffer(buflen, bufcap);
        }
    };

    // Array of the last relevant GC times
    m.lastGCs = data.lastGCs;

    // Contain indexed class sizes, this is initialized after reception of the first message.
    m.classSizes = new Array();

    m.initClassSizes = function (bySize) {
        for (let i = 0; i < bySize.length; i++) {
            m.classSizes.push(bySize[i].Size);
        }
    }

    function updateLastGC(memStats) {
        const nanoToSeconds = 1000 * 1000 * 1000;
        let lastGC = Math.floor(memStats.Mem.LastGC / nanoToSeconds);

        if (lastGC != data.lastGCs[data.lastGCs.length - 1]) {
            data.lastGCs.push(lastGC);
        }

        // Remove from the lastGCs array the timestamps which are prior to
        // the minimum timestamp in 'series'.
        let mints = data.series[idxTimes]._buf[0];
        let mingc = 0;
        for (let i = 0, n = data.lastGCs.length; i < n; i++) {
            if (data.lastGCs[i] > mints) {
                break;
            }
            mingc = i;
        }
        data.lastGCs.splice(0, mingc);
    }

    m.pushData = function (ts, memStats) {
        data.series[idxTimes].push(ts); // timestamp
        data.series[idxHeapAlloc].push(memStats.Mem.HeapAlloc);
        data.series[idxHeapSys].push(memStats.Mem.HeapSys);
        data.series[idxHeapIdle].push(memStats.Mem.HeapIdle);
        data.series[idxHeapInuse].push(memStats.Mem.HeapInuse);

        pushBySize(memStats.Mem.BySize);

        updateLastGC(memStats);
    }

    m.length = function () {
        return data.series[idxTimes].length();
    }

    m.slice = function (nitems) {
        // Time data
        let times = data.series[idxTimes].slice(nitems);

        // Heap plot data
        let heap = new Array(numSeriesHeap);
        heap[0] = times;
        for (let i = 1; i <= numSeriesHeap; i++) {
            heap[i] = data.series[i].slice(nitems);
        }

        // BySizes heatmap data
        let bySizes = new Array(5);
        bySizes[0] = times
        bySizes[1] = data.bySize[0].slice(nitems);
        bySizes[2] = data.bySize[1].slice(nitems);
        bySizes[3] = data.bySize[2].slice(nitems);
        bySizes[4] = data.bySize[3].slice(nitems);

        return {
            heap: heap,
            bySizes: bySizes,
        }
    }

    function pushBySize(bySize) {
        let sizesIndices = new Array(bySize.length);
        let counts = new Array(bySize.length);

        for (let i = 0; i < bySize.length; i++) {
            const size = bySize[i];
            sizesIndices[i] = i;
            counts[i] = size.Mallocs - size.Frees;
        }

        // TODO data.bySize[0] [1] and [2] in theory never change so there 
        // should be no need to recreate it each time.
        data.bySize[0].push(0);
        data.bySize[1].push(sizesIndices.length - 1);
        data.bySize[2].push(sizesIndices);
        data.bySize[3].push(counts);
    }

    return m;
}());