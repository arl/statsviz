export default class RingBuffer {
  constructor(capacity) {
    if (capacity < 1) throw new Error("Capacity must be > 0");
    this._buf = new Array(capacity);
    this._capacity = capacity;
    this._size = 0;
    this._start = 0;
  }

  push(item) {
    const end = (this._start + this._size) % this._capacity;
    this._buf[end] = item;
    if (this._size < this._capacity) {
      this._size++;
    } else {
      this._start = (this._start + 1) % this._capacity;
    }
  }

  slice(lastN) {
    const n = Math.min(lastN, this._size);
    const result = [];
    for (let i = this._size - n; i < this._size; i++) {
      result.push(this._buf[(this._start + i) % this._capacity]);
    }
    return result;
  }

  get last() {
    if (this._size === 0) return undefined;
    return this._buf[(this._start + this._size - 1) % this._capacity];
  }
}
