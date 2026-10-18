import StatsManager from "./StatsManager.js";
import PlotManager from "./PlotManager.js";
import { initNav, navState } from "./nav.js";
import { buildWebsocketURI } from "./utils.js";
import WebSocketClient from "./socket.js";

export const state = {
  statsMgr: null,
  plotMgr: null,
  ws: null,
};

// RAF-based throttling for plot updates
let rafId = null;
let pendingUpdate = false;
let forceNextUpdate = false;

const scheduleUpdate = () => {
  if (rafId !== null) return; // Already scheduled

  rafId = requestAnimationFrame(() => {
    rafId = null;
    if (pendingUpdate && navState.running) {
      const data = state.statsMgr.slice(navState.timerange);
      state.plotMgr.update(
        data,
        navState.gcEnabled,
        navState.timerange,
        forceNextUpdate,
      );
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

  state.ws = new WebSocketClient(
    uri,
    // onConfig
    (cfg) => {
      state.plotMgr = new PlotManager(cfg);
      state.statsMgr = new StatsManager(600, cfg);

      initNav(() => {
        drawPlots(false);
      });
    },
    // onData
    (msg) => {
      state.statsMgr.pushData(msg);
      drawPlots(true);
    }
  );
};
