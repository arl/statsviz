import StatsManager from "./StatsManager.js";
import * as plot from "./plot.js";
import { initNav, paused, show_gc, timerange } from "./nav.js";
import { clamp, buildWebsocketURI } from "./utils.js";
import "bootstrap/dist/js/bootstrap.min.js";

const dataRetentionSeconds = 600;
var stats;
var config;
export var allPlots;
var timeout = 250;

/* WebSocket connection handling */
export const connect = () => {
  const uri = buildWebsocketURI();
  let ws = new WebSocket(uri);
  console.info(`Attempting websocket connection to server at ${uri}`);

  ws.onopen = () => {
    console.info("Successfully connected");
    timeout = 250; // reset connection timeout for next time
  };

  ws.onclose = (event) => {
    console.error(`Closed websocket connection: code ${event.code}`);
    setTimeout(connect, clamp((timeout += timeout), 250, 5000));
  };

  ws.onerror = (err) => {
    console.error("WebSocket error:", err);
    ws.close();
  };

  ws.onmessage = (event) => {
    const data = JSON.parse(event.data);
    if (data.event == "config") {
      config = data.data;
      allPlots = configurePlots(config);
      stats = new StatsManager(dataRetentionSeconds, config);
      attachPlots(allPlots);
      initNav(allPlots);
    } else {
      stats.pushData(data.data);
      if (paused) {
        return;
      }
      if (!paused) updatePlots(allPlots, false);
    }
  };
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
  const data = stats.slice(timerange);
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
