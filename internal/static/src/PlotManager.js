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
  #shapesCache;
  #lastGcEnabled;
  #staggerHandle = null;

  constructor(config) {
    this.container = document.getElementById("plots");
    this.plots = config.series.map((pd) => new Plot(pd));
    this.#shapesCache = new Map();
    this.#lastGcEnabled = null;
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
    // Create GC vertical lines - only if needed.
    const shapes = new Map();
    if (gcEnabled) {
      // Only recreate shapes if GC state changed or events changed
      const gcStateChanged = this.#lastGcEnabled !== gcEnabled;

      for (const [name, serie] of data.events) {
        // Check if we need to regenerate shapes for this event
        const cached = this.#shapesCache.get(name);
        const eventsChanged =
          !cached ||
          cached.length !== serie.length ||
          (serie.length > 0 &&
            cached[cached.length - 1]?.x0?.getTime() !==
              serie[serie.length - 1]?.getTime());

        if (gcStateChanged || eventsChanged || !this.#shapesCache.has(name)) {
          const newShapes = createVerticalLines(serie);
          this.#shapesCache.set(name, newShapes);
          shapes.set(name, newShapes);
        } else {
          shapes.set(name, this.#shapesCache.get(name));
        }
      }
    } else {
      // GC disabled, clear all shapes
      if (this.#lastGcEnabled !== false) {
        this.#shapesCache.clear();
      }
    }

    this.#lastGcEnabled = gcEnabled;

    // X-axis range.
    const now = data.times[data.times.length - 1];
    const xrange = [now - timeRange * 1000, now];

    // Cancel any pending update to avoid overlapping updates
    if (this.#staggerHandle !== null) {
      cancelAnimationFrame(this.#staggerHandle);
      this.#staggerHandle = null;
    }

    const visiblePlots = this.plots.filter((p) => p.isVisible());
    let index = 0;

    const processBatch = () => {
      const start = performance.now();
      // Process plots for up to 12ms per frame to leave time for UI
      while (index < visiblePlots.length && performance.now() - start < 12) {
        visiblePlots[index].update(xrange, data, shapes, force);
        index++;
      }

      if (index < visiblePlots.length) {
        this.#staggerHandle = requestAnimationFrame(processBatch);
      } else {
        this.#staggerHandle = null;
      }
    };

    processBatch();
  }

  #resizeAll() {
    this.plots.forEach((p) => {
      const gd = document.getElementById(p.name());
      // We're being super defensive here to ensure that the div is
      // actually there (or Plotly.resize would fail).
      if (!gd || !p.isVisible()) return;
      const { offsetWidth: w, offsetHeight: h } = gd;
      if (w === 0 || h === 0) return;

      p.resize();
    });
  }
}
