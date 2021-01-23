(function () {
    function $(id) {
        return document.getElementById(id);
    }

    function buildWebsocketURI() {
        var loc = window.location, ws_prot = "ws:";
        if (loc.protocol === "https:") {
            ws_prot = "wss:";
        }
        return ws_prot + "//" + loc.host + loc.pathname + "ws"
    }

    const dataRetentionSeconds = 60;

    /* WebSocket callbacks */

    let socket = new WebSocket(buildWebsocketURI());
    console.log("Attempting Connection...");

    socket.onopen = () => {
        console.log("Successfully Connected");
    };

    socket.onclose = event => {
        console.log("Socket Closed Connection: ", event);
        socket.send("Client Closed!")
    };

    socket.onerror = error => {
        console.log("Socket Error: ", error);
    };

    var initDone = false;
    socket.onmessage = event => {
        let allStats = JSON.parse(event.data)
        if (!initDone) {
            stats.init(dataRetentionSeconds, allStats);
            initDone = true;
            return;
        }

        console.log(allStats);

        updateStats(allStats);
    }

    function updateStats(allStats) {
        stats.pushData(new Date(), allStats);

        if (ui.isPaused()) {
            return
        }

        let data = stats.slice(dataRetentionSeconds);
        if (data.heap[0].length == 1) {
            ui.createPlots(data);
        }
        ui.updatePlots(data);
    }

}());
