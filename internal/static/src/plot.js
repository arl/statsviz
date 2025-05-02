import * as theme from "./theme.js";
import { newConfigObject, newLayoutObject, themeColors } from "./plotConfig.js";
import { formatFunction } from "./utils.js";
import { plotWidth, plotHeight } from "./plotConfig.js";
import Plotly from "plotly.js-cartesian-dist-min";
import tippy, { followCursor } from "tippy.js";
import "tippy.js/dist/tippy.css";

const plotsDiv = document.getElementById("plots");

class Plot {
  constructor(cfg) {
    cfg.layout.paper_bgcolor = themeColors[theme.getThemeMode()].paper_bgcolor;
    cfg.layout.plot_bgcolor = themeColors[theme.getThemeMode()].plot_bgcolor;
    cfg.layout.font_color = themeColors[theme.getThemeMode()].font_color;

    this._hidden = false;
    this._maximized = false;
    this._cfg = cfg;
    this._updateCount = 0;
    this._dataTemplate = [];
    this._lastData = [{ x: new Date() }];

    if (this._cfg.type == "heatmap") {
      this._dataTemplate.push({
        type: "heatmap",
        x: null,
        y: this._cfg.buckets,
        z: null,
        showlegend: false,
        colorscale: this._cfg.colorscale,
        custom_data: this._cfg.custom_data,
      });
    } else {
      this._dataTemplate = this._cfg.subplots.map((subplot) => {
        return {
          type: this._cfg.type,
          x: null,
          y: null,
          name: subplot.name,
          hovertemplate: `<b>${subplot.unitfmt}</b>`,
        };
      });
    }

    this._plotlyLayout = newLayoutObject(cfg, this._maximized);
    this._plotlyConfig = newConfigObject(cfg, this._maximized);
  }

  name() {
    return this._cfg.name;
  }

  hide() {
    this._htmlElt.hidden = true;
  }

  show() {
    this._htmlElt.hidden = false;
  }

  isVisible() {
    return !this._htmlElt.hidden;
  }

  createElement(div) {
    this._htmlElt = div;
    // Pass a single data with no data to create an empty plot, this removes
    // the 'bad time formatting' warning at startup.
    Plotly.newPlot(
      this._htmlElt,
      this._lastData,
      this._plotlyLayout,
      this._plotlyConfig
    );
    if (this._cfg.type == "heatmap") {
      this._installHeatmapTooltip();
    }

    this._htmlElt.infoText = this._cfg.infoText
      .split("\n")
      .map((line) => `<p>${line}</p>`)
      .join("");
  }

  _installHeatmapTooltip() {
    const options = {
      followCursor: true,
      trigger: "manual",
      allowHTML: true,
      plugins: [followCursor],
    };
    const instance = tippy(document.body, options);
    const hover = this._cfg.hover;
    const formatYUnit = formatFunction(hover.yunit);

    const onHover = (data) => {
      const pt2txt = (d) => {
        let bucket;
        if (d.y == 0) {
          const yhigh = formatYUnit(d.data.custom_data[d.y]);
          bucket = `(-Inf, ${yhigh})`;
        } else if (d.y == d.data.custom_data.length - 1) {
          const ylow = formatYUnit(d.data.custom_data[d.y]);
          bucket = `[${ylow}, +Inf)`;
        } else {
          const ylow = formatYUnit(d.data.custom_data[d.y - 1]);
          const yhigh = formatYUnit(d.data.custom_data[d.y]);
          bucket = `[${ylow}, ${yhigh})`;
        }

        return `
<div class="tooltip-table tooltip-style">
    <div class="tooltip-row">
        <div class="tooltip-label">${hover.yname}</div>
        <div class="tooltip-value">${bucket}</div>
    </div>
    <div class="tooltip-row">
        <div class="tooltip-label">${hover.zname}</div>
        <div class="tooltip-value">${d.z}</div>
    </div>
</div> `;
      };
      instance.setContent(data.points.map(pt2txt)[0]);
      instance.show();
    };
    const onUnhover = (data) => {
      instance.hide();
    };

    this._htmlElt.on("plotly_hover", onHover).on("plotly_unhover", onUnhover);
  }

  _extractData(data) {
    const serie = data.series.get(this._cfg.name);
    if (this._cfg.type == "heatmap") {
      this._dataTemplate[0].x = data.times;
      this._dataTemplate[0].z = serie;
      this._dataTemplate[0].hoverinfo = "none";
      this._dataTemplate[0].colorbar = { len: "350", lenmode: "pixels" };
    } else {
      for (let i = 0; i < this._dataTemplate.length; i++) {
        this._dataTemplate[i].x = data.times;
        this._dataTemplate[i].y = serie[i];
        this._dataTemplate[i].stackgroup = this._cfg.subplots[i].stackgroup;
        this._dataTemplate[i].hoveron = this._cfg.subplots[i].hoveron;
        this._dataTemplate[i].type =
          this._cfg.subplots[i].type || this._cfg.type;
        this._dataTemplate[i].marker = {
          color: this._cfg.subplots[i].color,
        };
      }
    }
    return this._dataTemplate;
  }

  update(xrange, data, shapes, force) {
    this._lastData = this._extractData(data);
    this._updateCount++;
    if (
      force ||
      this._cfg.updateFreq == 0 ||
      this._updateCount % this._cfg.updateFreq == 0
    ) {
      // Update layout with vertical shapes if necessary.
      if (this._cfg.events != "") {
        this._plotlyLayout.shapes = shapes.get(this._cfg.events);
      }

      // Move the xaxis time range.
      this._plotlyLayout.xaxis.range = xrange;

      if (this._maximized) {
        this._plotlyLayout.width = plotsDiv.clientWidth;
        this._plotlyLayout.height = null;
        this._plotlyConfig.responsive = true;
      } else {
        this._plotlyLayout.width = plotWidth;
        this._plotlyLayout.height = plotHeight;
        this._plotlyConfig.responsive = false;
      }

      Plotly.react(
        this._htmlElt,
        this._lastData,
        this._plotlyLayout,
        this._plotlyConfig
      );
    }
  }

  maximize() {
    this._maximized = true;
    const plotsDiv = document.getElementById("plots");

    this._plotlyLayout = newLayoutObject(this._cfg, this._maximized);
    this._plotlyConfig = newConfigObject(this._cfg, this._maximized);

    this._plotlyLayout.width = plotsDiv.clientWidth;
    // this._plotlyLayout.height = plotsDiv.clientHeight;
    this._plotlyLayout.height = 2 * plotHeight;
    this._plotlyConfig.responsive = true;
    Plotly.react(
      this._htmlElt,
      this._lastData,
      this._plotlyLayout,
      this._plotlyConfig
    );
  }

  minimize() {
    this._maximized = false;

    this._plotlyLayout = newLayoutObject(this._cfg, this._maximized);
    this._plotlyConfig = newConfigObject(this._cfg, this._maximized);

    this._plotlyLayout.width = plotWidth;
    this._plotlyLayout.height = plotHeight;
    this._plotlyConfig.responsive = false;
    Plotly.react(
      this._htmlElt,
      this._lastData,
      this._plotlyLayout,
      this._plotlyConfig
    );
  }

  /**
   * update theme color and immediately force plot redraw to apply the new theme
   */
  updateTheme() {
    const themeMode = theme.getThemeMode();
    this._cfg.layout.paper_bgcolor = themeColors[themeMode].paper_bgcolor;
    this._cfg.layout.plot_bgcolor = themeColors[themeMode].plot_bgcolor;
    this._cfg.layout.font_color = themeColors[themeMode].font_color;

    this._plotlyLayout = newLayoutObject(this._cfg, this._maximized);
    this._plotlyConfig = newConfigObject(this._cfg, this._maximized);

    Plotly.react(
      this._htmlElt,
      this._lastData,
      this._plotlyLayout,
      this._plotlyConfig
    );
  }
}

// Create 'vertical lines' shapes for each of the given timestamps.
const createVerticalLines = (tss) => {
  const shapes = [];
  for (let i = 0, n = tss.length; i < n; i++) {
    const d = tss[i];
    shapes.push({
      type: "line",
      x0: d,
      x1: d,
      yref: "paper",
      y0: 0,
      y1: 1,
      line: {
        color: "rgb(55, 128, 191)",
        width: 1,
        dash: "longdashdot",
      },
    });
  }
  return shapes;
};

export { createVerticalLines, Plot };
