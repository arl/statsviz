import StatsManager from "./StatsManager.js";
import PlotManager from "./PlotManager.js";
import { initNav, running, gcEnabled, timerange } from "./nav.js";
import { buildWebsocketURI } from "./utils.js";
import WebSocketClient from "./socket.js";
import "bootstrap/dist/js/bootstrap.min.js";

export let statsMgr;
export let plotMgr;

export const drawPlots = (force) => {
  if (running) {
    const data = statsMgr.slice(timerange);
    plotMgr.update(data, gcEnabled, timerange, force);
  }
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
