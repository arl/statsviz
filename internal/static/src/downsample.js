/**
 * Downsample data using Largest-Triangle-Three-Buckets (LTTB) algorithm
 * for time-series data visualization.
 *
 * @param {Array} xData - Time series data (timestamps)
 * @param {Array} yData - Value series data
 * @param {number} threshold - Target number of data points
 * @returns {Object} - { x: downsampled x values, y: downsampled y values, indices: selected indices }
 */
export function downsampleLTTB(xData, yData, threshold) {
  const dataLength = xData.length;

  // Return original data if already below threshold or invalid
  if (dataLength <= threshold || threshold <= 2) {
    return { x: xData, y: yData, indices: null };
  }

  const sampledX = new Array(threshold);
  const sampledY = new Array(threshold);
  const sampledIndices = new Array(threshold);

  // Always include first point
  sampledX[0] = xData[0];
  sampledY[0] = yData[0];
  sampledIndices[0] = 0;

  // Always include last point
  sampledX[threshold - 1] = xData[dataLength - 1];
  sampledY[threshold - 1] = yData[dataLength - 1];
  sampledIndices[threshold - 1] = dataLength - 1;

  // Bucket size
  const bucketSize = (dataLength - 2) / (threshold - 2);

  let a = 0; // Initially a is the first point in the triangle

  for (let i = 1; i < threshold - 1; i++) {
    // Calculate point average for next bucket (for triangle)
    let avgX = 0;
    let avgY = 0;
    let avgRangeStart = Math.floor((i + 1) * bucketSize) + 1;
    let avgRangeEnd = Math.floor((i + 2) * bucketSize) + 1;
    avgRangeEnd = avgRangeEnd < dataLength ? avgRangeEnd : dataLength;

    const avgRangeLength = avgRangeEnd - avgRangeStart;

    for (; avgRangeStart < avgRangeEnd; avgRangeStart++) {
      avgX += xData[avgRangeStart].getTime
        ? xData[avgRangeStart].getTime()
        : xData[avgRangeStart];
      avgY += yData[avgRangeStart];
    }
    avgX /= avgRangeLength;
    avgY /= avgRangeLength;

    // Get the range for this bucket
    let rangeOffs = Math.floor(i * bucketSize) + 1;
    const rangeTo = Math.floor((i + 1) * bucketSize) + 1;

    // Point a
    const pointAX = xData[a].getTime ? xData[a].getTime() : xData[a];
    const pointAY = yData[a];

    let maxArea = -1;
    let maxAreaPoint = rangeOffs;

    for (; rangeOffs < rangeTo; rangeOffs++) {
      const pointX = xData[rangeOffs].getTime
        ? xData[rangeOffs].getTime()
        : xData[rangeOffs];
      const pointY = yData[rangeOffs];

      // Calculate triangle area over three buckets
      const area =
        Math.abs(
          (pointAX - avgX) * (pointY - pointAY) -
            (pointAX - pointX) * (avgY - pointAY)
        ) * 0.5;

      if (area > maxArea) {
        maxArea = area;
        maxAreaPoint = rangeOffs;
      }
    }

    sampledX[i] = xData[maxAreaPoint];
    sampledY[i] = yData[maxAreaPoint];
    sampledIndices[i] = maxAreaPoint;
    a = maxAreaPoint; // This point is the next a
  }

  return { x: sampledX, y: sampledY, indices: sampledIndices };
}

/**
 * Simple downsampling by picking every Nth point
 * Faster but less visually accurate than LTTB
 */
export function downsampleEveryNth(xData, yData, threshold) {
  const dataLength = xData.length;

  if (dataLength <= threshold) {
    return { x: xData, y: yData };
  }

  const step = Math.ceil(dataLength / threshold);
  const sampledX = [];
  const sampledY = [];

  for (let i = 0; i < dataLength; i += step) {
    sampledX.push(xData[i]);
    sampledY.push(yData[i]);
  }

  // Always include last point
  if (sampledX[sampledX.length - 1] !== xData[dataLength - 1]) {
    sampledX.push(xData[dataLength - 1]);
    sampledY.push(yData[dataLength - 1]);
  }

  return { x: sampledX, y: sampledY };
}

/**
 * Determine if downsampling should be applied based on data size
 * and current viewport/zoom level
 */
export function shouldDownsample(dataLength, maxPointsPerPlot = 500) {
  return dataLength > maxPointsPerPlot;
}
