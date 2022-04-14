// stats holds the data and function to modify it.
import Buffer from "./buffer.js";

var series = {
    times: null,
    eventsData: new Map(),
    plotData: new Map(),
};

const init = (plotdefs, buflen) => {
    const extraBufferCapacity = 20; // 20% of extra (preallocated) buffer datapoints
    const bufcap = buflen + (buflen * extraBufferCapacity) / 100; // number of actual datapoints

    series.times = new Buffer(buflen, bufcap);
    series.plotData.clear();
    plotdefs.series.forEach(plotdef => {
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

        let data = new Array(ndim);
        for (let i = 0; i < ndim; i++) {
            data[i] = new Buffer(buflen, bufcap);
        }
        series.plotData.set(plotdef.name, data);
    });

    plotdefs.events.forEach(event => {
        series.eventsData.set(event, new Array());
    });
}

const pushData = (ts, data) => {
    series.times.push(ts); // timestamp

    // Update time series.
    for (const [name, plotData] of series.plotData) {
        const curdata = data[name];
        for (let i = 0; i < curdata.length; i++) {
            plotData[i].push(curdata[i]);
        }
    }

    for (const [name, event] of series.eventsData) {
        const eventTs = data[name];
        if (event.length == 0) {
            event.push(eventTs);
            return;
        }
        if (eventTs.getTime() != event[event.length - 1].getTime()) {
            event.push(eventTs);
            // We've added a new timestamp, check if we can cut the front. We
            // don't need to keep track of event[0] if it happened before
            // the oldest timestamp we're showing. 
            let mints = series.times._buf[0];
            if (event[0] < mints) {
                event.splice(0, 1);
            }
        }
    }
}

const slice = (nitems) => {
    let sliced = {
        times: series.times.slice(nitems),
        series: new Map(),
        events: series.eventsData,
    };

    for (const [name, plotData] of series.plotData) {
        const arr = new Array(plotData.length);
        for (let i = 0; i < plotData.length; i++) {
            arr[i] = plotData[i].slice(nitems);
        }
        sliced.series.set(name, arr);
    }
    return sliced;
}

export { init, pushData, slice };