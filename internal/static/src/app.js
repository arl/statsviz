import * as stats from "./stats.js";
import * as plot from "./plot.js";
import * as theme from "./theme.js";

const buildWebsocketURI = () => {
  const wsUrl = import.meta.env.VITE_WEBSOCKET_URL;
  if (wsUrl) {
    return wsUrl;
  }
  console.log(`wsUrl: ${wsUrl}`);

  var loc = window.location,
    ws_prot = "ws:";
  if (loc.protocol === "https:") {
    ws_prot = "wss:";
  }
  return ws_prot + "//" + loc.host + loc.pathname + "ws";
};

const dataRetentionSeconds = 600;
var timeout = 250;

const clamp = (val, min, max) => {
  if (val < min) return min;
  if (val > max) return max;
  return val;
};

/* nav bar ui management */
let paused = false;
let show_gc = true;
let timerange = 60;
let config;
export let plots;

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
    let data = JSON.parse(event.data);
    if (data.event == "config") {
      config = data.data;
      plots = configurePlots(config);
      stats.init(config, dataRetentionSeconds);

      attachPlots(plots);

      installEventHandlers(plots);
    } else {
      stats.pushData(data.data);
      if (paused) {
        return;
      }
      if (!paused) updatePlots(plots);
    }
  };
};

const configurePlots = (config) => {
  const plots = [];
  config.series.forEach((plotdef) => {
    plots.push(new plot.Plot(plotdef));
  });
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

const installEventHandlers = (plots) => {
  document
    .getElementById("play_pause_switch")
    .addEventListener("change", (e) => {
      paused = !paused;
    });

  document.getElementById("show_gc_switch").addEventListener("change", (e) => {
    show_gc = !show_gc;
    updatePlots(plots);
  });

  document
    .getElementById("select-timerange")
    .addEventListener("change", (e) => {
      const val = parseInt(e.target.value, 10);
      timerange = val;
      updatePlots(plots);
    });

  document
    .getElementById("dark_mode_switch")
    .addEventListener("change", (e) => {
      const themeMode = theme.getThemeMode();
      const newTheme = (themeMode === "dark" && "light") || "dark";
      localStorage.setItem("theme-mode", newTheme);
      theme.updateThemeMode();
      plots.forEach((plot) => {
        plot.updateTheme();
      });
    });
};

const updatePlots = (plots) => {
  // Create shapes.
  let shapes = new Map();

  let data = stats.slice(timerange);

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
      plot.update(xrange, data, shapes);
    }
  });
};
