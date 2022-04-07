// stats holds the data and function to modify it.
import Buffer from "./buffer.js";

var data = {
    times: null,
    // Array of the last relevant GC times
    lastGCs: new Array(),
    // Where we'll store objects {data: Array(), type: string, datafunc: => )
    series: new Map(),
};

const lastGCs = data.lastGCs;

const init = (plotdefs, buflen) => {
    const extraBufferCapacity = 20; // 20% of extra (preallocated) buffer datapoints
    const bufcap = buflen + (buflen * extraBufferCapacity) / 100; // number of actual datapoints

    data.times = new Buffer(buflen, bufcap);
    data.series.clear();
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
        };

        const serie = {
            data: new Array(ndim),
            type: plotdef.type,
            datafunc: new Array(ndim),
        }
        if (serie.type == 'heatmap') {
            serie.datafunc = plotdef.heatmap.datapath;
        }

        for (let i = 0; i < ndim; i++) {
            serie.data[i] = new Buffer(buflen, bufcap);
            if (serie.type == 'scatter') {
                serie.datafunc[i] = plotdef.subplots[i].datapath;
            }
        }

        data.series.set(plotdef.name, serie);
    });
};

const pushData = (ts, allStats) => {
    const memStats = allStats.Mem;
    data.times.push(ts); // timestamp

    for (const [name, serie] of data.series) {
        switch (serie.type) {
            case 'scatter':
                for (let i = 0; i < serie.data.length; i++) {
                    serie.data[i].push(serie.datafunc[i](allStats));
                };
                break;
            case 'heatmap':
                for (let i = 0; i < serie.data.length; i++) {
                    serie.data[i].push(serie.datafunc(allStats, i));
                };
                break;
        };
    }

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

const slice = (nitems) => {
    let sliced = {
        times: data.times.slice(nitems),
        series: new Map(),
    };

    for (const [name, serie] of data.series) {
        const arr = new Array(serie.data.length);
        for (let i = 0; i < serie.data.length; i++) {
            arr[i] = serie.data[i].slice(nitems);
        }
        sliced.series.set(name, arr);
    }
    return sliced;
}

export { init, lastGCs, pushData, slice };