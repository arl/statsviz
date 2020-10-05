// stats holds the data and function to modify it.
var stats = (function () {
    var m = {};

    const idxHeapAlloc = 0;
    const idxHeapSys = 1;
    const idxHeapIdle = 2;
    const idxHeapInuse = 3;
    const idxHeapNextGC = 4;
    const numSeriesHeap = 5;

    const idxMSpanMCacheMSpanInUse = 0;
    const idxMSpanMCacheMSpanSys = 1;
    const idxMSpanMSpanMSCacheInUse = 2;
    const idxMSpanMSpanMSCacheSys = 3;
    const numSeriesMSpanMCache = 4;

    const idxObjectsLive = 0;
    const idxObjectsLookups = 1;
    const idxObjectsHeap = 2;
    const numSeriesObjects = 3;

    var data = {
        times: null,
        heap: new Array(numSeriesHeap),
        mspanMCache: new Array(numSeriesMSpanMCache),
        objects: new Array(numSeriesObjects),
        goroutines: null,
        gcfraction: null,
        lastGCs: new Array(),
        bySize: null,
    };

    m.init = function (buflen, allStats) {
        const extraBufferCapacity = 20; // 20% of extra (preallocated) buffer datapoints
        const bufcap = buflen + (buflen * extraBufferCapacity) / 100; // number of actual datapoints

        const memStats = allStats.Mem;

        data.times = new Buffer(buflen, bufcap);
        data.goroutines = new Buffer(buflen, bufcap);
        data.gcfraction = new Buffer(buflen, bufcap);

        for (let i = 0; i < numSeriesHeap; i++) {
            data.heap[i] = new Buffer(buflen, bufcap);
        }

        for (let i = 0; i < numSeriesMSpanMCache; i++) {
            data.mspanMCache[i] = new Buffer(buflen, bufcap);
        }

        for (let i = 0; i < numSeriesObjects; i++) {
            data.objects[i] = new Buffer(buflen, bufcap);
        }

        // size classes heatmap
        for (let i = 0; i < memStats.BySize.length; i++) {
            m.classSizes.push(memStats.BySize[i].Size);
        }

        data.bySize = new Array(m.classSizes.length);
        for (let i = 0; i < data.bySize.length; i++) {
            data.bySize[i] = new Buffer(buflen, bufcap);
        }
    };

    // Array of the last relevant GC times
    m.lastGCs = data.lastGCs;

    function updateLastGC(memStats) {
        const nanoToSeconds = 1000 * 1000 * 1000;
        let t = Math.floor(memStats.LastGC / nanoToSeconds);

        let lastGC = new Date(t * 1000);

        if (lastGC != data.lastGCs[data.lastGCs.length - 1]) {
            data.lastGCs.push(lastGC);
        }

        // Remove from the lastGCs array the timestamps which are prior to
        // the minimum timestamp in 'series'.
        let mints = data.times._buf[0];
        let mingc = 0;
        for (let i = 0, n = data.lastGCs.length; i < n; i++) {
            if (data.lastGCs[i] > mints) {
                break;
            }
            mingc = i;
        }
        data.lastGCs.splice(0, mingc);
    }

    // Contain indexed class sizes, this is initialized after reception of the first message.
    m.classSizes = new Array();

    m.pushData = function (ts, allStats) {
        data.times.push(ts); // timestamp

        const memStats = allStats.Mem;

        data.gcfraction.push(memStats.GCCPUFraction);
        data.goroutines.push(allStats.NumGoroutine);

        data.heap[idxHeapAlloc].push(memStats.HeapAlloc);
        data.heap[idxHeapSys].push(memStats.HeapSys);
        data.heap[idxHeapIdle].push(memStats.HeapIdle);
        data.heap[idxHeapInuse].push(memStats.HeapInuse);
        data.heap[idxHeapNextGC].push(memStats.NextGC);

        data.mspanMCache[idxMSpanMCacheMSpanInUse].push(memStats.MSpanInuse);
        data.mspanMCache[idxMSpanMCacheMSpanSys].push(memStats.MSpanSys);
        data.mspanMCache[idxMSpanMSpanMSCacheInUse].push(memStats.MCacheInuse);
        data.mspanMCache[idxMSpanMSpanMSCacheSys].push(memStats.MCacheSys);

        data.objects[idxObjectsLive].push(memStats.Mallocs - memStats.Frees);
        data.objects[idxObjectsLookups].push(memStats.Lookups);
        data.objects[idxObjectsHeap].push(memStats.HeapObjects);

        for (let i = 0; i < memStats.BySize.length; i++) {
            const size = memStats.BySize[i];
            data.bySize[i].push(size.Mallocs - size.Frees);
        }

        updateLastGC(memStats);
    }

    m.length = function () {
        return data.times.length();
    }

    m.slice = function (nitems) {
        const times = data.times.slice(nitems);
        const gcfraction = data.gcfraction.slice(nitems);
        const goroutines = data.goroutines.slice(nitems);

        // Heap plot data
        let heap = new Array(numSeriesHeap);
        for (let i = 0; i < numSeriesHeap; i++) {
            heap[i] = data.heap[i].slice(nitems);
        }

        // MSpan/MCache plot data
        let mspanMCache = new Array(numSeriesMSpanMCache);
        for (let i = 0; i < numSeriesMSpanMCache; i++) {
            mspanMCache[i] = data.mspanMCache[i].slice(nitems);
        }

        // Objects plot data
        let objects = new Array(numSeriesObjects);
        for (let i = 0; i < numSeriesObjects; i++) {
            objects[i] = data.objects[i].slice(nitems);
        }

        // BySizes heatmap data
        let bySizes = new Array(data.bySize.length);
        for (let i = 0; i < data.bySize.length; i++) {
            const size = data.bySize[i];
            bySizes[i] = data.bySize[i].slice(nitems);
        }

        return {
            times: times,
            gcfraction: gcfraction,
            goroutines: goroutines,
            heap: heap,
            mspanMCache: mspanMCache,
            objects: objects,
            bySizes: bySizes,
        }
    }

    return m;
}());