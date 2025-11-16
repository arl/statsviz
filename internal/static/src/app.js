import StatsManager from "./StatsManager.js";
import PlotManager from "./PlotManager.js";
import { initNav, running, gcEnabled, timerange } from "./nav.js";
import { buildWebsocketURI } from "./utils.js";
import WebSocketClient from "./socket.js";
import "bootstrap/dist/js/bootstrap.min.js";

export let statsMgr;
export let plotMgr;

// RAF-based throttling for plot updates
let rafId = null;
let pendingUpdate = false;
let forceNextUpdate = false;

const scheduleUpdate = () => {
  if (rafId !== null) return; // Already scheduled

  rafId = requestAnimationFrame(() => {
    rafId = null;
    if (pendingUpdate && running) {
      const data = statsMgr.slice(timerange);
      plotMgr.update(data, gcEnabled, timerange, forceNextUpdate);
      pendingUpdate = false;
      forceNextUpdate = false;
    }
  });
};

export const drawPlots = (force) => {
  pendingUpdate = true;
  if (force) {
    forceNextUpdate = true;
  }
  scheduleUpdate();
};

export const connect = () => {
  const uri = buildWebsocketURI();

  new WebSocketClient(
    uri,
    // onConfig
    (cfg) => {
      plotMgr = new PlotManager(cfg);
      statsMgr = new StatsManager(600, cfg);

      initNav(() => {
        drawPlots(false);
      });
    },
    // onData
    (msg) => {
      statsMgr.pushData(msg);
      drawPlots(true);
    }
  );
};
