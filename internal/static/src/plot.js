import * as theme from "./theme.js";
import {
  defaultPlotHeight,
  newConfigObject,
  newLayoutObject,
  themeColors,
} from "./plotConfig.js";
import { formatFunction } from "./utils.js";
import Plotly from "plotly.js-cartesian-dist";
import tippy, { followCursor } from "tippy.js";
import "tippy.js/dist/tippy.css";
import "bootstrap-icons/font/bootstrap-icons.min.css";

const plotsDiv = document.getElementById("plots");

class Plot {
  #htmlElt;
  #plotlyLayout;
  #plotlyConfig;
  #lastData;
  #updateCount;
  #maximized;
  #cfg;
  #dataTemplate;
  #cachedWidth;
  #inViewport = false;
  #observer;

  constructor(cfg) {
    cfg.layout.paper_bgcolor = themeColors[theme.getThemeMode()].paper_bgcolor;
    cfg.layout.plot_bgcolor = themeColors[theme.getThemeMode()].plot_bgcolor;
    cfg.layout.font_color = themeColors[theme.getThemeMode()].font_color;

    this.#maximized = false;
    this.#cfg = cfg;
    this.#updateCount = 0;
    this.#dataTemplate = [];
    this.#lastData = [{ x: new Date() }];
    this.#cachedWidth = null;

    if (this.#cfg.type == "heatmap") {
      this.#dataTemplate.push({
        type: "heatmap",
        x: null,
        y: this.#cfg.buckets,
        z: null,
        showlegend: false,
        colorscale: this.#cfg.colorscale,
        custom_data: this.#cfg.custom_data,
      });
    } else {
      this.#dataTemplate = this.#cfg.subplots.map((subplot) => {
        return {
          type: this.#cfg.type,
          x: null,
          y: null,
          name: subplot.name,
          hovertemplate: `<b>${subplot.unitfmt}</b>`,
        };
      });
    }

    this.#plotlyLayout = newLayoutObject(cfg, this.#maximized);
    this.#plotlyConfig = newConfigObject(cfg, this.#maximized);
  }

  name() {
    return this.#cfg.name;
  }

  hasTag(tag) {
    return this.#cfg.tags.includes(tag);
  }

  matches(query) {
    if (!query) return true;
    if (!this.#cfg.metrics) return false;
    const q = query.toLowerCase();
    return this.#cfg.metrics.some((m) => m.toLowerCase().includes(q));
  }

  setVisible(visible) {
    this.#htmlElt.hidden = !visible;
  }

  isVisible() {
    return !this.#htmlElt.hidden;
  }

  createElement(div) {
    this.#htmlElt = div;

    this.#observer = new IntersectionObserver((entries) => {
      entries.forEach((entry) => {
        this.#inViewport = entry.isIntersecting;
        if (this.#inViewport) {
          this.#react();
        }
      });
    });
    this.#observer.observe(this.#htmlElt);

    // Measure the final CSS width.
    this.#cachedWidth = div.clientWidth;
    this.#plotlyLayout.width = this.#cachedWidth;
    this.#plotlyLayout.height = defaultPlotHeight;

    // Pass a single data with no data to create an empty plot (this removes
    // the 'bad time formatting' warning at startup).

    Plotly.newPlot(
      this.#htmlElt,
      this.#lastData,
      this.#plotlyLayout,
      this.#plotlyConfig
    );

    if (this.#cfg.type == "heatmap") {
      this._installHeatmapTooltip();
    }

    this.#htmlElt.infoText = this.#cfg.infoText
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
    const hover = this.#cfg.hover;
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
    const onUnhover = (_data) => {
      instance.hide();
    };

    this.#htmlElt.on("plotly_hover", onHover).on("plotly_unhover", onUnhover);
  }

  #extractData(data) {
    const serie = data.series.get(this.#cfg.name);

    if (this.#cfg.type == "heatmap") {
      this.#dataTemplate[0].x = data.times;
      this.#dataTemplate[0].z = serie;
      this.#dataTemplate[0].hoverinfo = "none";
      this.#dataTemplate[0].colorbar = { len: "350", lenmode: "pixels" };
    } else {
      for (let i = 0; i < this.#dataTemplate.length; i++) {
        this.#dataTemplate[i].x = data.times;
        this.#dataTemplate[i].y = serie[i];

        this.#dataTemplate[i].stackgroup = this.#cfg.subplots[i].stackgroup;
        this.#dataTemplate[i].hoveron = this.#cfg.subplots[i].hoveron;
        this.#dataTemplate[i].type =
          this.#cfg.subplots[i].type || this.#cfg.type;
        this.#dataTemplate[i].marker = {
          color: this.#cfg.subplots[i].color,
        };
      }
    }
    return this.#dataTemplate;
  }

  update(xrange, data, shapes, force) {
    this.#lastData = this.#extractData(data);
    this.#updateCount++;
    if (
      force ||
      this.#cfg.updateFreq == 0 ||
      this.#updateCount % this.#cfg.updateFreq == 0
    ) {
      // Update layout with vertical shapes if necessary.
      if (this.#cfg.events != "") {
        this.#plotlyLayout.shapes = shapes.get(this.#cfg.events);
      }

      // Move the xaxis time range.
      this.#plotlyLayout.xaxis.range = xrange;

      if (this.#maximized) {
        this.#plotlyConfig.responsive = true;
      } else {
        this.#plotlyLayout.height = defaultPlotHeight;
        this.#plotlyConfig.responsive = false;
      }

      // Use cached width - only recalculated on resize
      this.#plotlyLayout.width = this.#cachedWidth;

      if (this.#inViewport) {
        this.#react();
      }
    }
  }

  isMaximized() {
    return this.#maximized;
  }

  maximize() {
    this.#maximized = true;
    const plotsDiv = document.getElementById("plots", this.#maximized);

    this.#plotlyLayout = newLayoutObject(this.#cfg, this.#maximized);
    this.#plotlyConfig = newConfigObject(this.#cfg, this.#maximized);

    this.#cachedWidth = plotsDiv.clientWidth;
    this.#plotlyLayout.width = this.#cachedWidth;
    this.#plotlyLayout.height = plotsDiv.parentElement.clientHeight - 50;

    this.#react();
  }

  minimize() {
    this.#maximized = false;

    this.#plotlyLayout = newLayoutObject(this.#cfg, this.#maximized);
    this.#plotlyConfig = newConfigObject(this.#cfg, this.#maximized);

    this.#cachedWidth = this.#htmlElt.clientWidth;
    this.#plotlyLayout.width = this.#cachedWidth;

    this.#react();
  }

  resize() {
    this.updateCachedWidth();
    const layoutUpdate = { width: this.#cachedWidth };
    if (this.#maximized) {
      const plotsDiv = document.getElementById("plots");
      layoutUpdate.height = plotsDiv.parentElement.clientHeight - 50;
    } else {
      layoutUpdate.height = defaultPlotHeight;
    }
    Plotly.relayout(this.#htmlElt, layoutUpdate);
  }

  updateCachedWidth() {
    if (this.#maximized) {
      this.#cachedWidth = plotsDiv.clientWidth;
    } else {
      this.#cachedWidth = this.#htmlElt.clientWidth;
    }
  }

  #react() {
    Plotly.react(
      this.#htmlElt,
      this.#lastData,
      this.#plotlyLayout,
      this.#plotlyConfig
    );
  }

  updateTheme() {
    const mode = theme.getThemeMode();
    const { paper_bgcolor, plot_bgcolor, font_color } = themeColors[mode];

    this.#cfg.layout.paper_bgcolor = paper_bgcolor;
    this.#cfg.layout.plot_bgcolor = plot_bgcolor;
    this.#cfg.layout.font_color = font_color;

    Plotly.relayout(this.#htmlElt, {
      paper_bgcolor: paper_bgcolor,
      plot_bgcolor: plot_bgcolor,
      "font.color": font_color,
    });
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
