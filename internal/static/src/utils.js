export const clamp = (val, min, max) => Math.min(Math.max(val, min), max);

export const buildWebsocketURI = () => {
  const wsUrl = import.meta.env.VITE_WEBSOCKET_URL;
  if (wsUrl) return wsUrl;
  const { protocol, host, pathname } = window.location;
  const wsProt = protocol === "https:" ? "wss:" : "ws:";
  return `${wsProt}//${host}${pathname}ws`;
};
