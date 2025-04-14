import { writable, derived } from "svelte/store";
import { configStore } from "./config";
import Buffer from "../lib/buffer";

// 20% of extra (preallocated) buffer datapoints
const extraBufferCapacity = 20;
const MaxDatapoints = 600;
const bufcap = MaxDatapoints + (MaxDatapoints * extraBufferCapacity) / 100; // number of actual datapoints

const _metricsStore = writable(null);

export function initMetrics(cfg) {
  configStore.set(cfg);
  console.log("received new config, reset metrics stored in-memory");

  const metrics = {
    timestamps: new Buffer(MaxDatapoints, bufcap),
    events: new Map(),
    series: new Map(),
  };

  cfg.events.forEach((event) => {
    metrics.events.set(event, []);
  });

  cfg.series.forEach((serie) => {
    let ndim =
      serie.type === "heatmap" ? serie.buckets.length : serie.subplots.length;
    const bufs = Array.from(
      { length: ndim },
      () => new Buffer(MaxDatapoints, bufcap)
    );
    metrics.series.set(serie.name, bufs);
  });

  _metricsStore.set(metrics);
}

function update(snapshot) {
  const { series, timestamp } = snapshot;

  _metricsStore.update((metrics) => {
    metrics.timestamps.push(timestamp);

    // Update events series, deduplicating event timestamps and cutting the ones
    // that are older than the oldest timestamp we're tracking.
    for (const [name, event] of metrics.events) {
      if (event.length == 0) {
        if (series[name].length != 0) {
          const eventTs = new Date(Math.floor(series[name][0]));
          event.push(eventTs);
        }
        break;
      }

      const eventTs = new Date(Math.floor(series[name][0]));
      if (eventTs.getTime() != event[event.length - 1]) {
        event.push(eventTs);
        let oldest = metrics.timestamps.first();
        if (event[0] < oldest) {
          event.splice(0, 1);
          console.debug("Trimming oldest event");
        }
      }
    }

    // TODO: could use a $derived store for this?
    const now = metrics.timestamps.last();
    const xrange = [now - timerange * 1000, now];

    for (const [key, values] of Object.entries(series)) {
      const serie = metrics.series.get(key);
      if (!serie) continue; // Ignore 'events' series.
      for (let i = 0; i < values.length; i++) {
        serie[i].push(values[i]);
      }
    }

    return metrics;
  });
}

export const metricsStore = {
  subscribe: _metricsStore.subscribe,
  update: update,
};

const timerange = 60; // TODO: should come from a store

export const metricStore = (key) => {
  return derived(metricsStore, ($metricsStore) => {
    if (!$metricsStore || $metricsStore.timestamps.length() === 0)
      return {
        times: [],
        data: [],
        events: new Map(),
        xrange: [0, 0],
      };
    const now = $metricsStore.timestamps.last();
    return {
      times: $metricsStore.timestamps.all(),
      data: $metricsStore.series.get(key).map((buf) => buf.all()),
      events: $metricsStore.events,
      xrange: [now - timerange * 1000, now],
    };
  });
};
