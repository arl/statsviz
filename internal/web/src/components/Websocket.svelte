<script>
  import { onMount, onDestroy } from "svelte";
  import { initMetrics, metricsStore } from "../stores/metrics";

  let timeout = 250; // initial timeout
  const MAX_TIMEOUT = 5000; // max timeout

  let ws = null;

  function buildWebsocketURI() {
    // TODO check how to read env var with rollup
    // const wsUrl = import.meta.env.VITE_WEBSOCKET_URL;
    const wsUrl = "ws://localhost:9090/debug/statsviz/ws";
    if (wsUrl) {
      return wsUrl;
    }

    const proto = window.location.protocol === "https:" ? "wss" : "ws";
    const host = window.location.host;
    const pathname = window.location.pathname;

    return `${proto}://${host}${pathname}ws`;
  }

  function clamp(value, min, max) {
    return Math.min(Math.max(value, min), max);
  }

  const connect = () => {
    const uri = buildWebsocketURI();
    ws = new WebSocket(uri);
    console.info(`Attempting WebSocket connection to server at ${uri}`);

    ws.onopen = () => {
      console.info("Successfully connected to WebSocket");
      timeout = 250; // Reset for next close/connect cycle.
    };

    ws.onclose = (event) => {
      console.error(`WebSocket connection closed: code ${event.code}`);
      ws = null;
      setTimeout(connect, clamp((timeout += timeout), 250, MAX_TIMEOUT));
    };

    ws.onerror = (err) => {
      console.error("WebSocket encountered an error:", err);
      if (ws) ws.close();
    };

    ws.onmessage = (event) => {
      try {
        const payload = JSON.parse(event.data);

        switch (payload.event) {
          case "config":
            initMetrics(payload.data);
            break;
          case "metrics":
            metricsStore.update(payload.data);
            break;
        }
      } catch (error) {
        console.error("Error while websocket message:", error);
      }
    };
  };

  onMount(connect);
  onDestroy(() => {
    if (ws) ws.close();
  });
</script>
