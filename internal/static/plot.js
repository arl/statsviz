const plotWidth = 620;
const plotHeight = 420;

// https://plotly.com/javascript/configuration-options/
const plotlyConfigBase = {
    displaylogo: false,
    modeBarButtonsToRemove: ['2D', 'zoom2d', 'pan2d', 'select2d', 'lasso2d', 'zoomIn2d', 'zoomOut2d', 'autoScale2d', 'resetScale2d', 'toggleSpikelines'],
    toImageButtonOptions: {
        format: 'png',
    }
}

// https://plotly.com/javascript/reference/layout
const plotlyLayoutBase = {
    width: plotWidth,
    height: plotHeight,
    xaxis: {
        title: 'time',
        tickformat: '%H:%M:%S',
    },
    yaxis: {
        exponentformat: 'SI',
    }
};

export default class Plot {

    /* Constructs and configures a new Plot object, which wraps a Plotly chart
        and becomes the only reference to plotly the external world will need.
        - cfg is an object expecting the following fields to be present:
        {
            name: String         // plot 'internal' identifier
            title: String        // displayed plot title
            type: String         // 'scatter' or 'heatmap'
            updateFreq: Integer  // freq of update:
                                 //  - 1: update each time we receive new metrics
                                 //  - 2: update half the time, etc.
            hasHorzEvents: bool  // show events as horizontal bars on this chart
            layout: Object       // Plotly-specific: gets merged over 'plotlyLayoutBase'
        }
        - div: HTMLElement is the div html element passed to Plotly.newPlot
        - dataFunc: function is the function used to extract and fill the plot data from the incoming stats data
        - data: Object is the actual data, used to initialize chart
    */


    constructor(cfg, div, dataFunc, data) {
        this._cfg = cfg;
        this._dataFunc = dataFunc;
        this._updateCount = 0;
        this._htmlElt = div;

        this._plotlyLayout = {...plotlyLayoutBase, ...this._cfg.layout };
        this._plotlyLayout.title = this._cfg.title;

        this._plotlyConfig = {...plotlyConfigBase }
        this._plotlyConfig.toImageButtonOptions.filename = this._cfg.name

        Plotly.newPlot(this._htmlElt, this._dataFunc(data), this._plotlyLayout, this._plotlyConfig);
    }

    name() {
        return this._cfg.name;
    }

    extractData(data) {
        return this._dataFunc(data);
    }

    layout() {
        return this._layout
    }

    update(data, horzEvents) {
        this._updateCount++;
        if (this._cfg.updateFreq == 0 || (this._updateCount % this._cfg.updateFreq == 0)) {
            if (this._cfg.hasHorzEvents === true) {
                this._plotlyLayout.shapes = horzEvents;
            }
            Plotly.react(this._htmlElt, this._dataFunc(data), this._plotlyLayout, this._plotlyConfig);
        }
    }
};