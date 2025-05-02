export const clamp = (val, min, max) => Math.min(Math.max(val, min), max);

export const buildWebsocketURI = () => {
  const wsUrl = import.meta.env.VITE_WEBSOCKET_URL;
  if (wsUrl) return wsUrl;
  const { protocol, host, pathname } = window.location;
  const wsProt = protocol === "https:" ? "wss:" : "ws:";
  return `${wsProt}//${host}${pathname}ws`;
};

const durUnits = ["w", "d", "h", "m", "s", "ms", "µs", "ns"];
const durVals = [6048e11, 864e11, 36e11, 6e10, 1e9, 1e6, 1e3, 1];

// Formats a time duration provided in second.
export const formatDuration = (sec) => {
  let ns = sec * 1e9;
  for (let i = 0; i < durUnits.length; i++) {
    let inc = ns / durVals[i];

    if (inc < 1) continue;
    return Math.round(inc) + durUnits[i];
  }
  return res.trim();
};

const bytesUnits = ["B", "KB", "MB", "GB", "TB", "PB", "EB"];

// Formats a size in bytes.
export const formatBytes = (bytes) => {
  let i = 0;
  while (bytes > 1000) {
    bytes /= 1000;
    i++;
  }
  const res = Math.trunc(bytes);
  return `${res}${bytesUnits[i]}`;
};

// Returns a format function based on the provided unit.
export const formatFunction = (unit) => {
  switch (unit) {
    case "duration":
      return formatDuration;
    case "bytes":
      return formatBytes;
  }
  // Default formatting
  return (y) => {
    `${y} ${hover.yunit}`;
  };
};
