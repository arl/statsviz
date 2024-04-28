import * as stats from './stats.js';
import * as plot from "./plot.js";
import * as theme from "./theme.js";
import PlotsDef from './plotsdef.js';

const dataRetentionSeconds = 600;
var timeout = 250;

const clamp = (val, min, max) => {
    if (val < min) return min;
    if (val > max) return max;
    return val;
}

/* nav bar ui management */
let paused = false;
let show_gc = true;
let timerange = 60;

const dataProcessor = {
    initDone: false,
    close: (e) => {dataProcessor.onclose(e||"connection fail")},
    connected: false,
    retrying: false,
    onopen: () => {
        dataProcessor.initDone = false;
        dataProcessor.connected = true;
        console.info("Successfully connected");
        timeout = 250; // reset connection timeout for next time
    },
    onclose: event => {
        dataProcessor.connected = false;
        console.error(`Closed connection: code ${event.code || event}`);
        if (dataProcessor.retrying) {
            return
        }
        dataProcessor.retrying = true
        setTimeout(() => {
            connect()
            dataProcessor.retrying = false
        }, clamp(timeout += timeout, 250, 5000));
    },
    onerror: err => {
        console.error(`error, closing connection.`, err);
        dataProcessor.close();
    },
    onmessage: event => {
        let data = JSON.parse(event.data)
        if (!dataProcessor.initDone) {
            configurePlots(PlotsDef);
            stats.init(PlotsDef, dataRetentionSeconds);

            attachPlots();

            $('#play_pause').change(() => {
                paused = !paused;
            });
            $('#show_gc').change(() => {
                show_gc = !show_gc;
                updatePlots();
            });
            $('#select_timerange').click(() => {
                const val = parseInt($("#select_timerange option:selected").val(), 10);
                timerange = val;
                updatePlots();
            });
            dataProcessor.initDone = true;
        }
        dataProcessor.onData(data);
    },
    onData: data => {
        stats.pushData(data);
        if (paused || !dataProcessor.connected) {
            return
        }
        updatePlots()
    }
}
/* WebSocket connection handling */
const connect = () => {

    // compatible with the following writing methods
    // mux.HandleFunc("/debug/statsviz/ws", srv.Metrics())
    let path = window.location.pathname+(PlotsDef.metricsPath || "ws");
    const eventSource = new EventSource(path);
    console.info(`Attempting metrics connection to server at ${path}`);
    for (let event in dataProcessor) {
        eventSource[event] = dataProcessor[event];
    }
    dataProcessor.close = eventSource.close
}

connect();

let plots = [];

const configurePlots = (plotdefs) => {
    plots = [];
    plotdefs.series.forEach(plotdef => {
        plots.push(new plot.Plot(plotdef));
    });
}

const attachPlots = () => {
    let plotsDiv = $('#plots');
    plotsDiv.empty();

    for (let i = 0; i < plots.length; i++) {
        const plot = plots[i];
        let div = $(`<div id="${plot.name()}">`);
        plot.createElement(div[0], i)
        plotsDiv.append(div);
    }
}

function throttle(func, delay) {
    const context = this;
    let timerFlag = null;
    return function () {
        if (timerFlag === null) {
            func.apply(context,arguments);
            timerFlag = setTimeout(() => {
                timerFlag = null;
            }, delay);
        }
    };
}

const updatePlots = throttle(() => {
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

    plots.forEach(plot => {
        if (!plot.hidden) {
            plot.update(xrange, data, shapes);
        }
    });
}, (PlotsDef.sendFrequency || 1000) / 10)

const updatePlotsLayout = () => {
    plots.forEach(plot => {
        plot.updateTheme();
    });
}

theme.updateThemeMode();

/**
 * Change color theme when the user presses the theme switch button
 */
$('#color_theme_sw').change(() => {
    const themeMode = theme.getThemeMode();
    const newTheme = themeMode === "dark" && "light" || "dark";
    localStorage.setItem("theme-mode", newTheme);
    theme.updateThemeMode();
    updatePlotsLayout();
});
