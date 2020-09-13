var ui = (function () {
    var m = {};

    let paused = false;

    m.isPaused = function () { return paused; }
    m.togglePause = function () { paused = !paused; }
    m.plots = null;

    m.createPlots = function (opts, data, elts) {
        let heap = new uPlot(opts.heap, data.heap, elts.heap);
        let bySizes = new uPlot(opts.bySizes, data.bySizes, elts.bySizes);

        m.plots = {
            heap: heap,
            bySizes: bySizes,
        }
    }

    m.updatePlots = function (xScale, data) {
        m.plots.heap.batch(() => {
            m.plots.heap.setData(data.heap);
            m.plots.heap.setScale("x", xScale);
        });
        m.plots.bySizes.batch(() => {
            m.plots.bySizes.setData(data.bySizes);
            m.plots.bySizes.setScale("x", xScale);
        });
    }

    return m;
}());