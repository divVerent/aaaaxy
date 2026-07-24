if (typeof module != 'undefined') {
  const common = require("./common.js");
  for (const [key, value] of Object.entries(common)) {
    global[key] = value;
  }
}

const MIN_FRACTION = 0.5; // Time fraction in which to not show anomalies.
const NOISE_FLOOR = 0.01; // Even if not seen, observe a min stddev of 1%.

((data, filter) => {
  for (const [platform, runs] of Object.entries(data.entries)) {
    if (!Array.isArray(runs)) continue;

    // Iterate chronologically through every single run
    for (let i = 0; i < runs.length; i++) {
      const currentRun = runs[i];
      for (const bench of currentRun.benches) {
        const benchName = bench.name;
        const newVal = bench.value;
        const cpu = cpuTypeOf(bench);
        const baselineVals = [];
        let otherVals = 0;
        for (let h = 0; h < runs.length; h++) {
          if (h == i) {
            continue;
          }
          const prevRun = runs[h];
          for (const prevBench of prevRun.benches) {
            if (prevBench.name === benchName) {
              const prevCPU = cpuTypeOf(prevBench);
              if (cpu == prevCPU) {
                if (h < i) {
                  baselineVals.push(prevBench.value);
                }
                ++otherVals;
              }
            }
          }
        }

        if (baselineVals.length < 2 || baselineVals.length < otherVals * MIN_FRACTION) {
          bench.anomaly = null;
          continue;
        }
        const n = baselineVals.length;
        const mean = baselineVals.reduce((a, b) => a + b, 0) / n;
        const sqSum = baselineVals.reduce((acc, val) => acc + Math.pow(val - mean, 2), 0);
        let stdDev = Math.sqrt(sqSum / (n - 1));
        const noiseFloor = mean * NOISE_FLOOR;
        if (stdDev < noiseFloor) {
          stdDev = noiseFloor;
        }

        // Compute Prediction Standard Error & t-statistic
        const se = stdDev * Math.sqrt(1 + 1 / n);
        const tStat = se > 0 ? (newVal - mean) / se : (newVal === mean ? 0 : Infinity);
        const df = n - 1;

        // Determine two-tailed p-value
        const pVal = 2.0 * (1.0 - tCDF(Math.abs(tStat), df));

        // Evaluate anomaly threshold & set field in-place
        const alpha = 1.0 - Math.pow(PERCENTILE_THRESHOLD, 1.0 / currentRun.benches.length);
        if (pVal < alpha) {
          bench.anomaly = tStat;
          if (i == runs.length - 1) {
            console.info(
              `🚨 ANOMALY DETECTED: [${platform}] "${benchName}" changed significantly at run index ${i}.\n` +
              `   New Value:  ${newVal.toFixed(4)} ${bench.unit || ''}\n` +
              `   Baseline:   ${mean.toFixed(4)} ± ${stdDev.toFixed(4)} (n = ${n} samples)\n` +
              `   Stats:    t-stat = ${tStat.toFixed(3)}, p-val = ${pVal.toExponential(3)}`
            );
          }
        } else {
          bench.anomaly = null;
        }
      }
    }
  }
})(window.BENCHMARK_DATA);
