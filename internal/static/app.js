import * as stats from './stats.js';
import * as plot from "./plot.js";
import PlotsDef from './plotsdef.js';

const buildWebsocketURI = () => {
    var loc = window.location,
        ws_prot = "ws:";
    if (loc.protocol === "https:") {
        ws_prot = "wss:";
    }
    return ws_prot + "//" + loc.host + loc.pathname + "ws"
}

const dataRetentionSeconds = 60;
var timeout = 250;

const clamp = (val, min, max) => {
    if (val < min) return min;
    if (val > max) return max;
    return val;
}

/* WebSocket connection handling */

const connect = () => {
    const uri = buildWebsocketURI();
    let ws = new WebSocket(uri);
    console.info(`Attempting websocket connection to server at ${uri}`);

    ws.onopen = () => {
        console.info("Successfully connected");
        timeout = 250; // reset connection timeout for next time
    };

    ws.onclose = event => {
        console.error(`Closed websocket connection: code ${event.code}`);
        setTimeout(connect, clamp(timeout += timeout, 250, 5000));
    };

    ws.onerror = err => {
        console.error(`Websocket error, closing connection.`);
        ws.close();
    };

    let initDone = false;
    ws.onmessage = event => {
        let allStats = JSON.parse(event.data)

        if (!initDone) {
            configurePlots(PlotsDef);
            stats.init(PlotsDef, dataRetentionSeconds);

            attachPlots();

            initDone = true;
            return;
        }

        stats.pushData(new Date(), allStats);
        if (isPaused()) {
            return
        }
        updatePlots(stats.slice(dataRetentionSeconds), PlotsDef.events);
    }
}

connect();

/* plots management */

// TODO(arl) not used for now
let paused = false;
const isPaused = () => { return paused; }
const togglePause = () => { paused = !paused; }
let plots = [];

const configurePlots = (plotdefs) => {
    plots = [];
    plotdefs.series.forEach(plotdef => {
        plots.push(new plot.Plot(plotdef));
    });
}

const attachPlots = () => {
    let row = null;
    let plotsDiv = $('#plots');
    plotsDiv.empty()

    for (let i = 0; i < plots.length; i++) {
        const plot = plots[i];
        let div = $('<div>');
        plot.createElement(div[0], i)
        plotsDiv.append(div);
    }
}

const updatePlots = (data) => {
    // Create shapes.
    let shapes = new Map();

    for (const [name, serie] of data.events) {
        shapes.set(name, plot.createHorizontalLines(serie));
    }

    plots.forEach(plot => {
        if (!plot.hidden) {
            plot.update(data, shapes);
        }
    });
}