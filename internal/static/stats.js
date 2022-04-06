// stats holds the data and function to modify it.
import Buffer from "./buffer.js";

var data = {
    times: null,
    // Array of the last relevant GC times
    lastGCs: new Array(),

    // TODO(arl) put plot data in a subproperty, so we can just loop on elements
    // when pushing (and not have to pass plotDefs to stats.pushData, nor to
    // stats.slice)
};

const lastGCs = data.lastGCs;

const init = (plotdefs, buflen) => {
    const extraBufferCapacity = 20; // 20% of extra (preallocated) buffer datapoints
    const bufcap = buflen + (buflen * extraBufferCapacity) / 100; // number of actual datapoints

    data.times = new Buffer(buflen, bufcap);

    plotdefs.forEach(plotdef => {
        let ndim;
        switch (plotdef.type) {
            case 'scatter':
                ndim = plotdef.subplots.length;
                break;
            case 'heatmap':
                ndim = plotdef.heatmap.buckets.length;
                break;
            default:
                console.error(`[statsviz]: unknown plot type "${plotdef.type}"`);
                return;
        }
        const arr = new Array(ndim);
        for (let i = 0; i < ndim; i++) {
            arr[i] = new Buffer(buflen, bufcap)
        }
        data[plotdef.name] = arr;
    });
};

const pushData = (plotdefs, ts, allStats) => {
    data.times.push(ts); // timestamp

    const memStats = allStats.Mem;

    plotdefs.forEach(plotdef => {
        const name = plotdef.name;
        switch (plotdef.type) {
            case 'scatter':
                for (let i = 0; i < data[name].length; i++) {
                    data[name][i].push(plotdef.subplots[i].datapath(allStats));
                };
                break;
            case 'heatmap':
                for (let i = 0; i < data[name].length; i++) {
                    data[name][i].push(plotdef.heatmap.datapath(allStats, i));
                };
                break;
        };
    });

    updateLastGC(memStats);
}

const updateLastGC = memStats => {
    const nanoToSeconds = 1000 * 1000 * 1000;
    let t = Math.floor(memStats.LastGC / nanoToSeconds);
    let lastGC = new Date(t * 1000);
    if (data.lastGCs.length == 0) {
        data.lastGCs.push(lastGC);
        return;
    }
    if (lastGC.getTime() != data.lastGCs[data.lastGCs.length - 1].getTime()) {
        data.lastGCs.push(lastGC);
        // We've added a GC timestamp, check if we can cut the front. We
        // don't need to keep track data.lastGCs[0] if it happened before
        // the oldest timestamp we're showing. 
        let mints = data.times._buf[0];
        if (data.lastGCs[0] < mints) {
            data.lastGCs.splice(0, 1);
        }
    }
}

const slice = (plotdefs, nitems) => {
    let sliced = {
        times: data.times.slice(nitems),
    };

    plotdefs.forEach(plotdef => {
        const name = plotdef.name;
        sliced[name] = new Array(data[name].length);
        for (let i = 0; i < data[name].length; i++) {
            sliced[name][i] = data[name][i].slice(nitems);
        }
    });
    return sliced;
}

export { init, lastGCs, pushData, slice };