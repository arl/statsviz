import { Plot, createVerticalLines } from "./plot.js";
import Plotly from "plotly.js-cartesian-dist";

function debounce(fn, delay) {
  let timer = null;
  return function () {
    const context = this;
    const args = arguments;
    clearTimeout(timer);
    timer = setTimeout(() => fn.apply(context, args), delay);
  };
}

export default class PlotManager {
  constructor(config) {
    this.container = document.getElementById("plots");
    this.plots = config.series.map((pd) => new Plot(pd));
    this.#attach();

    window.addEventListener(
      "resize",
      debounce(() => {
        this.#resizeAll();
      }, 100)
    );

    requestAnimationFrame(() => this.#resizeAll());
  }

  #attach() {
    this.container.replaceChildren();
    this.plots.forEach((p) => {
      const div = document.createElement("div");
      div.id = p.name();
      div.className = "plot-card";
      // It's important to append the <div> to its container before creating the
      // plot so Plotly knows the <div> dimensions.
      this.container.appendChild(div);
      p.createElement(div);
    });
  }

  update(data, gcEnabled, timeRange, force = false) {
    // Create GC vertical lines.
    const shapes = new Map();
    if (gcEnabled) {
      for (const [name, serie] of data.events) {
        shapes.set(name, createVerticalLines(serie));
      }
    }

    // X-axis range.
    const now = data.times[data.times.length - 1];
    const xrange = [now - timeRange * 1000, now];

    this.plots.forEach((p) => {
      if (p.isVisible()) p.update(xrange, data, shapes, force);
    });
  }

  #resizeAll() {
    this.plots.forEach((p) => {
      const gd = document.getElementById(p.name());
      // We're being super defensive here to ensure that the div is
      // actually there (or Plotly.resize would fail).
      if (!gd) return;
      if (!p.isVisible()) return;
      const { offsetWidth: w, offsetHeight: h } = gd;
      if (w === 0 || h === 0) return;

      p.updateCachedWidth();
      Plotly.Plots.resize(gd);
    });
  }
}
