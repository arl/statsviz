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
            name: String            // plot 'internal' identifier
                                    // TODO(arl): also temporarily used as key in data
            title: String           // displayed plot title
            type: String            // 'scatter' or 'heatmap'
            updateFreq: Integer     // freq of update:
                                    //  - 1: update each time we receive new metrics
                                    //  - 2: update half the time, etc.
            horzEvents: string|''   // show an 'event' serie as horizontal lines
            layout: Object          // Plotly-specific: gets merged over 'plotlyLayoutBase'

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
                hover: {
                    yunit: String // Y unit
                    yname: String // label for the y value in hover tooltip
                    zname: String // label for the z value in hover tooltip
                },
                colorscale: [ [Number: 'color'] ] // heatmap colorscale palette, numbers go from 0 to 1
                buckets: classSizes,              // heatmap list of buckets
            },

        }
    */
    constructor(cfg) {
        this._cfg = cfg;
        this._updateCount = 0;
        this._dataTemplate = [];

        if (['scatter', 'bar'].includes(this._cfg.type)) {
            this._cfg.subplots.forEach(subplot => {
                const hover = subplot.hover || subplot.name;
                const unitfmt = subplot.unitfmt;

                this._dataTemplate.push({
                    type: this._cfg.type,
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
                custom_data: this._cfg.heatmap.custom_data,
            });
        }

        this._plotlyLayout = {...plotlyLayoutBase, ...this._cfg.layout };
        this._plotlyLayout.title = this._cfg.title;

        this._plotlyConfig = {...plotlyConfigBase }
        this._plotlyConfig.toImageButtonOptions.filename = this._cfg.name
    }

    createElement(div, idx) {
        this._htmlElt = div;
        this._plotIdx = idx;
        Plotly.newPlot(this._htmlElt, null, this._plotlyLayout, this._plotlyConfig);
    }

    // Install callbacks for showing info about the rectangle area under the cursor.
    installHover(hoverinfo) {
        const options = {
            arrow: true,
            followCursor: true,
            popperOptions: {
                placement: "auto"
            },
            interactive: true,
            trigger: "manual",
            allowHTML: true
        };
        const instance = tippy(document.body, options);
        if (this._cfg.type == 'heatmap') {
            const hover = this._cfg.heatmap.hover;
            const formatYUnit = formatFunction(hover.yunit);
            this._htmlElt.on('plotly_hover', function(data) {
                    var infotext = data.points.map(function(d) {
                        const yval = formatYUnit(d.data.custom_data[d.y]);
                        return `${hover.yname}: <b>${yval}<b/><br/>${hover.zname}: <b>${d.z}</b>`;
                    });

                    let info = document.createElement('div');
                    info.innerHTML = infotext;
                    instance.setContent(info);
                    instance.show();
                })
                .on('plotly_unhover', function(data) {
                    instance.hide();
                });
        }
    }

    _extractData(data) {
        const serie = data.series.get(this._cfg.name);
        if (['scatter', 'bar'].includes(this._cfg.type)) {
            for (let i = 0; i < this._dataTemplate.length; i++) {
                this._dataTemplate[i].x = data.times;
                this._dataTemplate[i].y = serie[i];
                this._dataTemplate[i].stackgroup = this._cfg.subplots[i].stackgroup;
                this._dataTemplate[i].hoveron = this._cfg.subplots[i].hoveron;
                this._dataTemplate[i].marker = {
                    color: this._cfg.subplots[i].color,
                };
            }
        } else if (this._cfg.type == 'heatmap') {
            this._dataTemplate[0].x = data.times;
            this._dataTemplate[0].z = serie;
            this._dataTemplate[0].hoverinfo = 'none';
        }
        return this._dataTemplate;
    }

    update(data, shapes) {
        this._updateCount++;
        if (this._cfg.updateFreq == 0 || (this._updateCount % this._cfg.updateFreq == 0)) {
            if (this._cfg.horzEvents != '') {
                this._plotlyLayout.shapes = shapes.get(this._cfg.horzEvents);
            }
            Plotly.react(this._htmlElt, this._extractData(data), this._plotlyLayout, this._plotlyConfig);
        }
    }
};

const durUnits = ['w', 'd', 'h', 'm', 's', 'ms', 'µs', 'ns'];
const durVals = [6048e11, 864e11, 36e11, 6e10, 1e9, 1e6, 1e3, 1];

// Formats a time duration provided in second.
const formatDuration = sec => {
    let ns = sec * 1e9;
    for (let i = 0; i < durUnits.length; i++) {
        let inc = ns / durVals[i];

        if (inc < 1) continue;
        return Math.round(inc) + durUnits[i];
    }
    return res.trim();
};

const bytesUnits = ['B', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB'];

// Formats a size in bytes.
const formatBytes = bytes => {
    let i = 0;
    while (bytes > 1000) {
        bytes /= 1000;
        i++;
    }
    const res = Math.trunc(bytes);
    return `${res}${bytesUnits[i]}`;
};

// Returns a format function based on the provided unit.
const formatFunction = unit => {
    switch (unit) {
        case 'duration':
            return formatDuration;
        case 'bytes':
            return formatBytes;
    }
    // Default formatting
    return (y) => { `${y} ${hover.yunit}` };
};