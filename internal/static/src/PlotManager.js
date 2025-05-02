import { Plot, createVerticalLines } from "./plot.js";

export default class PlotManager {
  constructor(config) {
    this.container = document.getElementById("plots");
    // Build plot instances.
    this.plots = config.series.map((pd) => new Plot(pd));
    this.#attach();
  }

  #attach() {
    // Render plot placeholders into DOM.
    this.container.replaceChildren(
      ...this.plots.map((p) => {
        const div = document.createElement("div");
        div.id = p.name();
        p.createElement(div);
        return div;
      })
    );
  }

  update(data, gcEnabled, timeRange, force = false) {
    // Compute GC vertical lines.
    const shapes = new Map();
    if (gcEnabled) {
      for (const [name, serie] of data.events) {
        shapes.set(name, createVerticalLines(serie));
      }
    }
    // X-axis range.
    const now = data.times[data.times.length - 1];
    const xrange = [now - timeRange * 1000, now];

    // Delegate to each Plot.
    this.plots.forEach((p) => {
      if (!p.hidden) p.update(xrange, data, shapes, force);
    });
  }
}
