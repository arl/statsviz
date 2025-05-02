import StatsManager from "./StatsManager.js";
import * as plot from "./plot.js";
import { initNav, paused, show_gc, timerange } from "./nav.js";
import { buildWebsocketURI } from "./utils.js";
import WebSocketClient from "./socket.js";
import "bootstrap/dist/js/bootstrap.min.js";

const dataRetentionSeconds = 600;
let config;
export var allPlots;

let statsMgr;
export const connect = () => {
  const uri = buildWebsocketURI();
  const client = new WebSocketClient(
    uri,
    // onConfig
    (cfg) => {
      config = cfg;
      allPlots = configurePlots(cfg);
      statsMgr = new StatsManager(dataRetentionSeconds, cfg);
      attachPlots(allPlots);
      initNav(allPlots);
    },
    // onData
    (msg) => {
      statsMgr.pushData(msg);
      if (!paused) updatePlots(allPlots, false);
    }
  );
};

const configurePlots = (config) => {
  const plots = config.series.map((pd) => new plot.Plot(pd));
  return plots;
};

const attachPlots = (plots) => {
  const plotsDiv = document.getElementById("plots");
  plotsDiv.replaceChildren(
    ...plots.map((plot) => {
      const div = document.createElement("div");
      div.id = plot.name();
      plot.createElement(div);
      return div;
    })
  );
};

export const updatePlots = (plots, force = false) => {
  const data = statsMgr.slice(timerange);
  const shapes = new Map();

  if (show_gc) {
    for (const [name, serie] of data.events) {
      shapes.set(name, plot.createVerticalLines(serie));
    }
  }

  // Always show the full range on x axis.
  const now = data.times[data.times.length - 1];
  let xrange = [now - timerange * 1000, now];

  plots.forEach((plot) => {
    if (!plot.hidden) {
      plot.update(xrange, data, shapes, force);
    }
  });
};
