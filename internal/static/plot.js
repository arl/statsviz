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
                                 // TODO(arl): also temporarily used as key in data
            title: String        // displayed plot title
            type: String         // 'scatter' or 'heatmap'
            updateFreq: Integer  // freq of update:
                                 //  - 1: update each time we receive new metrics
                                 //  - 2: update half the time, etc.
            hasHorzEvents: bool  // show events as horizontal bars on this chart
            layout: Object       // Plotly-specific: gets merged over 'plotlyLayoutBase'

            // If type = 'scatter':
            subplots: [
                {
                    name: String    // subplot name
                    unitfmt: String // unit format string
                    hover: String   // hover title (optional, defaults to name)
                },
            ]

            // If type = 'heatmap':
            heatmap:
            {
                hover: 'TO REFINE (see TODO)',
                colorscale: [ [Number: 'color'] ] // heatmap colorscale palette, numbers go from 0 to 1
                buckets: classSizes,              // heatmap list of buckets
            },

        }
    */
    constructor(cfg) {
        this._cfg = cfg;
        this._updateCount = 0;
        this._dataTemplate = [];

        if (this._cfg.type == 'scatter') {
            this._cfg.subplots.forEach(subplot => {
                const hover = subplot.hover || subplot.name;
                const unitfmt = subplot.unitfmt;

                this._dataTemplate.push({
                    type: 'scatter',
                    x: null,
                    y: null,
                    name: subplot.name,
                    hovertemplate: `<b>${unitfmt}</b>`,
                })
            });
        } else if (this._cfg.type == 'heatmap') {
            this._dataTemplate.push({
                type: 'heatmap',
                x: null,
                y: this._cfg.heatmap.buckets,
                z: null,
                showlegend: false,
                colorscale: this._cfg.heatmap.colorscale,
                hovertemplate: this._cfg.heatmap.hover,
            });
        }

        this._plotlyLayout = {...plotlyLayoutBase, ...this._cfg.layout };
        this._plotlyLayout.title = this._cfg.title;

        this._plotlyConfig = {...plotlyConfigBase }
        this._plotlyConfig.toImageButtonOptions.filename = this._cfg.name
    }

    name() {
        return this._cfg.name;
    }

    createElement(div) {
        this._htmlElt = div;
        Plotly.newPlot(this._htmlElt, null, this._plotlyLayout, this._plotlyConfig);
    }

    extractData(data) {
        if (this._cfg.type == 'scatter') {
            for (let i = 0; i < this._dataTemplate.length; i++) {
                this._dataTemplate[i].x = data.times;
                this._dataTemplate[i].y = data.series.get(this._cfg.name)[i];
            }
        } else if (this._cfg.type == 'heatmap') {
            this._dataTemplate[0].x = data.times;
            this._dataTemplate[0].z = data.series.get(this._cfg.name);
        }
        return this._dataTemplate;
    }

    update(data, horzEvents) {
        this._updateCount++;
        if (this._cfg.updateFreq == 0 || (this._updateCount % this._cfg.updateFreq == 0)) {
            if (this._cfg.hasHorzEvents === true) {
                this._plotlyLayout.shapes = horzEvents;
            }
            Plotly.react(this._htmlElt, this.extractData(data), this._plotlyLayout, this._plotlyConfig);
        }
    }
};