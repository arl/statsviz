export default class Plot {
    constructor(div, name, dataFunc, layout, data, updateFreq, hasHorzEvents) {
        this._name = name;
        this._dataFunc = dataFunc;
        this._layout = layout;
        this._updateFreq = updateFreq;
        this._updateCount = 0;
        this._hasHorzEvents = hasHorzEvents;
        this._cfg = {
            displaylogo: false,
            modeBarButtonsToRemove: ['2D', 'zoom2d', 'pan2d', 'select2d', 'lasso2d', 'zoomIn2d', 'zoomOut2d', 'autoScale2d', 'resetScale2d', 'toggleSpikelines'],
            toImageButtonOptions: {
                format: 'png',
                filename: name,
            }
        }
        this._htmlElt = div;

        Plotly.newPlot(div, dataFunc(data), layout, this._cfg);
    }

    name() {
        return this._name;
    }

    extractData(data) {
        return this._dataFunc(data);
    }

    config() {
        return this._cfg
    }

    layout() {
        return this._layout
    }

    update(data, horzEvents) {
        this._updateCount++;
        if (this._updateFreq == 0 || (this._updateCount % this._updateFreq == 0)) {
            if (this._hasHorzEvents === true) {
                this._layout.shapes = horzEvents;
            }
            Plotly.react(this._htmlElt, this._dataFunc(data), this._layout, this._cfg);
        }
    }
};