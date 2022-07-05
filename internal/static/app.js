import * as stats from './stats.js';
import Plot from "./plot.js";
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
        plots.push(new Plot(plotdef));
    });
}

const attachPlots = () => {
    let row = null;
    let plotsDiv = $('#plots');
    plotsDiv.empty()

    for (let i = 0; i < plots.length; i++) {
        const plot = plots[i];
        if (i % 2 == 0) {
            row = $('<div>', { class: 'row' });
            plotsDiv.append(row);
        }

        let col = $('<div>', { class: 'col' });
        let div = $('<div>');
        let hoverinfo = $('<div>', {
            id: `hoverinfo-${i}`,
            style: "margin-left:80px;",
        });

        plot.createElement(div[0], i)
        col.append(div);
        col.append(hoverinfo);
        row.append(col);
        plot.installHover(hoverinfo[0]);
    }
}

const updatePlots = (data) => {
    // Create shapes.
    let shapes = new Map();

    for (const [eventName, eventSerie] of data.events) {
        shapes.set(eventName, createEventShape(data, eventSerie));
    }

    plots.forEach(plot => {
        if (!plot.hidden) {
            plot.update(data, shapes);
        }
    });
}

const createEventShape = (data, eventSerie) => {
    // TODO(arl): do we really need to pass 'data' to extract mints and maxtx?
    // aren't event serie already clamped to the visible time range?
    const mints = data.times[0];
    const maxts = data.times[data.times.length - 1];

    const shapes = [];
    for (let i = 0, n = eventSerie.length; i < n; i++) {
        let d = eventSerie[i];
        // Clamp GC times which are out of bounds
        if (d < mints || d > maxts) {
            continue;
        }

        shapes.push({
            type: 'line',
            x0: d,
            x1: d,
            yref: 'paper',
            y0: 0,
            y1: 1,
            line: {
                color: 'rgb(55, 128, 191)',
                width: 1,
                dash: 'longdashdot',
            }
        })
    }
    return shapes;
}