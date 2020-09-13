var stats = (function () {
    var m = {};

    const numSeriesHeap = 4;
    const numSeries = 1 + numSeriesHeap; // timestamp serie + other series

    const idxTimes = 0;
    const idxHeapAlloc = 1;
    const idxHeapSys = 2;
    const idxHeapIdle = 3;
    const idxHeapInuse = 4;

    var data = {
        // timeseries [0]         -> timestamps
        // [1] to [numSeries - 1] -> timeseries datapoint
        series: new Array(numSeries), // TODO: rename to timeseries 
        lastGCs: new Array(),

        bySize: new Array(4), // [0] min [1] max [2] class sizes [3] live objects
    };

    m.init = function () {
        for (let i = 0; i < numSeries; i++) {
            data.series[i] = new Buffer(maxBufferLen, maxBufferCap);
        }
        for (let i = 0; i < data.bySize.length; i++) {
            data.bySize[i] = new Buffer(maxBufferLen, maxBufferCap);
        }
    };

    m.lastGCs = data.lastGCs;

    m.pushData = function (ts, memStats) {
        data.series[idxTimes].push(ts); // timestamp
        data.series[idxHeapAlloc].push(memStats.Mem.HeapAlloc);
        data.series[idxHeapSys].push(memStats.Mem.HeapSys);
        data.series[idxHeapIdle].push(memStats.Mem.HeapIdle);
        data.series[idxHeapInuse].push(memStats.Mem.HeapInuse);

        pushBySize(memStats.Mem.BySize);

        const nanoToSeconds = 1000 * 1000 * 1000;
        let lastGC = Math.floor(memStats.Mem.LastGC / nanoToSeconds);

        if (lastGC != data.lastGCs[data.lastGCs.length - 1]) {
            data.lastGCs.push(lastGC);
        }

        // Remove from the lastGCs array the timestamps which are prior to
        // the minimum timestamp in 'series'.
        // TODO: do this in a trimLastGC function
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

    // Slice data in order to keep the last nitems contained in the raw data.
    // TODO: rename sliceSeries or sliceHeap
    m.sliceData = function (nitems) {
        let d = new Array(numSeries);
        for (let i = 0; i < numSeries; i++) {
            d[i] = data.series[i].slice(nitems);
        }

        return d;
    }

    m.length = function () {
        return data.series[idxTimes].length();
    }

    m.sliceHeatmapData = function (nitems) {
        let d = new Array(5);
        // TODO : could reuse the timestamp slice since it has already been
        // sliced previously (in sliceData).
        d[0] = data.series[idxTimes].slice(nitems);
        d[1] = data.bySize[0].slice(nitems);
        d[2] = data.bySize[1].slice(nitems);
        d[3] = data.bySize[2].slice(nitems);
        d[4] = data.bySize[3].slice(nitems);

        return d;
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