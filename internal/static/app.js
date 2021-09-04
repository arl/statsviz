(function() {
    function $(id) {
        return document.getElementById(id);
    }

    function buildWebsocketURI() {
        var loc = window.location,
            ws_prot = "ws:";
        if (loc.protocol === "https:") {
            ws_prot = "wss:";
        }
        return ws_prot + "//" + loc.host + loc.pathname + "ws"
    }

    const dataRetentionSeconds = 60;
    var timeout = 250;

    function clamp(val, min, max) {
        if (val < min) return min;
        if (val > max) return max;
        return val;
    }

    /* WebSocket connection handling */

    function connect() {
        let ws = new WebSocket(buildWebsocketURI());
        console.log("Attempting websocket connection to statsviz server...");

        ws.onopen = () => {
            console.log("Successfully connected");
            timeout = 250; // reset connection timeout for next time
        };

        ws.onclose = event => {
            console.log("Closed websocket connection: ", event);
            setTimeout(connect, clamp(timeout += timeout, 250, 5000));
        };

        ws.onerror = error => {
            console.log("Websocket error: ", error);
            ws.close();
        };

        var initDone = false;
        ws.onmessage = event => {
            let allStats = JSON.parse(event.data)
            if (!initDone) {
                stats.init(dataRetentionSeconds, allStats);
                stats.pushData(new Date(), allStats);
                initDone = true;
                let data = stats.slice(dataRetentionSeconds);
                ui.createPlots(data);
                return;
            }

            stats.pushData(new Date(), allStats);
            if (ui.isPaused()) {
                return
            }
            let data = stats.slice(dataRetentionSeconds);
            ui.updatePlots(data);
        }
    }
    connect();
}());