import RingBuffer from "./RingBuffer.js";

export default class StatsManager {
  #retention;
  #times;
  #plotData;
  #eventsData;

  constructor(retentionSeconds, config) {
    this.#retention = retentionSeconds;
    this.#initBuffers(config);
  }

  #initBuffers(config) {
    const cap = this.#retention;
    this.#times = new RingBuffer(cap);
    this.#plotData = new Map();
    this.#eventsData = new Map();

    for (const pd of config.series) {
      let dims = pd.type === "heatmap" ? pd.buckets.length : pd.subplots.length;
      const arr = Array.from({ length: dims }, () => new RingBuffer(cap));
      this.#plotData.set(pd.name, arr);
    }
    for (const evt of config.events) {
      this.#eventsData.set(evt, []);
    }
  }

  pushData(payload) {
    this.#times.push(payload.timestamp);

    for (const [name, buffers] of this.#plotData) {
      const values = payload.series[name];
      values.forEach((v, i) => buffers[i].push(v));
    }

    for (const [evtName, arr] of this.#eventsData) {
      const raw = payload.series[evtName][0];
      const ts = new Date(Math.floor(raw));
      if (!arr.length || ts.getTime() !== arr[arr.length - 1].getTime()) {
        arr.push(ts);
        // drop old events
        const oldest = this.#times.first;
        if (arr[0] < oldest) arr.shift();
      }
    }
  }

  slice(lastN) {
    const times = this.#times.slice(lastN);
    const series = new Map();
    for (const [name, buffers] of this.#plotData) {
      series.set(
        name,
        buffers.map((buf) => buf.slice(lastN))
      );
    }
    return { times, series, events: this.#eventsData };
  }

  reset(config) {
    this.#initBuffers(config);
  }
}
