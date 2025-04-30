import * as stats from "./stats.js";
import * as plot from "./plot.js";
import * as theme from "./theme.js";
import "bootstrap/dist/js/bootstrap.min.js";

const buildWebsocketURI = () => {
  const wsUrl = import.meta.env.VITE_WEBSOCKET_URL;
  if (wsUrl) {
    return wsUrl;
  }

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
export let allPlots;

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
      allPlots = configurePlots(config);
      stats.init(config, dataRetentionSeconds);

      attachPlots(allPlots);

      installEventHandlers(allPlots);
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
  // Show GC toggle.
  const gcToggle = document.getElementById("gcToggle");

  gcToggle.checked = show_gc;
  gcToggle.addEventListener("change", (e) => {
    show_gc = !show_gc;
    gcToggle.checked = show_gc;
    updatePlots(plots, true);
  });

  // Pause/Resume button.
  const pauseBtn = document.getElementById("pauseBtn");
  pauseBtn.addEventListener("click", (e) => {
    paused = !paused;
    pauseBtn.textContent = paused ? "Resume" : "Pause";
    pauseBtn.classList.toggle("active", paused);
    updatePlots(plots, true);
  });

  // Dark mode toggle.
  const themeToggle = document.getElementById("themeToggle");
  themeToggle.addEventListener("change", (e) => {
    const themeMode = theme.getThemeMode();
    const newTheme = (themeMode === "dark" && "light") || "dark";
    localStorage.setItem("theme-mode", newTheme);

    theme.updateThemeMode();

    plots.forEach((plot) => {
      plot.updateTheme();
    });
  });

  // Time range selection
  const rangeInputs = document.querySelectorAll('input[name="range"]');

  rangeInputs.forEach((r, i) =>
    r.addEventListener("change", () => {
      if (r.checked) {
        rangeInputs[i].checked = true;
        const val = 60 * parseInt(rangeInputs[i].value, 10);
        timerange = val;
        updatePlots(plots, true);
      }
    })
  );
  document.getElementById("range1").checked = true;
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
