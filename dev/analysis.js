((data) => {
  const TARGET_FALSE_ALERT_PROBABILITY = 0.005;
  const PERCENTILE_THRESHOLD = 1.0 - TARGET_FALSE_ALERT_PROBABILITY;
  
  if (!data || !data.entries) {
    alert('Benchmark data not found!');
    return;
  }
  
  const stats = {};  // Benchmark -> raw values seen.
  const ratioStats = {};  // Benchmark -> ratios of curr/prev seen.
  
  // Load benchmark results.
  for (const env in data.entries) {
    const runs = data.entries[env].sort((a, b) => a.date - b.date);
    let previousRun = {};
    runs.forEach(run => {
      if (!run.benches || run.benches.length == 0) return;
      const currentRun = {};
      run.benches.forEach(bench => {
        const name = env + " | " + bench.name;
        currentRun[bench.name] = bench.value;
        if (!stats[name]) stats[name] = [];
        stats[name].push(bench.value);
        if (previousRun[bench.name] !== undefined) {
          const prevVal = previousRun[bench.name];
          if (prevVal > 0 && bench.value > 0) {
            if (!ratioStats[name]) ratioStats[name] = [];
            ratioStats[name].push(bench.value / prevVal);
          }
        }
      });
      previousRun = currentRun;
    });
  }
  
  const benchmarkModels = [];
  for (const name in ratioStats) {
    const ratios = ratioStats[name];
    const n = ratios.length;
    if (n < 2) continue;
    const logRatios = ratios.map(r => Math.log(r));
    const logMean = logRatios.reduce((sum, v) => sum + v, 0) / n;
    const logVar = logRatios.reduce((sum, v) => sum + Math.pow(v - logMean, 2), 0) / (n - 1);
    const logStd = Math.sqrt(logVar);
    if (logStd > 0) {
      benchmarkModels.push({ name, mu: logMean, sigma: logStd });
    }
  }
  
  function pNoAlertOf(model, f) {
    const lnF = Math.log(f);
    const z = (lnF - model.mu) / model.sigma;
    return normalCDF(z);
  }
  
  function pNoAlert(models, f) {
    let p = 1.0;
    for (const model of models) {
      p *= pNoAlertOf(model, f);
    }
    return p;
  }
  
  function searchFactor(models, threshold) {
    let low = 1.0001;
    let high = 2.0;
    while (pNoAlert(models, high) < threshold) {
      high *= 2.0;
      if (high > 1e6) {
        return null;
      }
    }
  
    for (let i = 0; i < 60; i++) {
      const mid = (low + high) / 2;
      if (pNoAlert(models, mid) > threshold) {
        high = mid;
      } else {
        low = mid;
      }
    }
    return high;
  }
  
  // 3. Compute Global Factor 'f' accounting for INDIVIDUAL variances
  const globalF = searchFactor(benchmarkModels, PERCENTILE_THRESHOLD);
  if (globalF == null) {
    alert('Data is too noisy to find a bound.');
    return;
  }
  
  const singleThreshold = Math.pow(PERCENTILE_THRESHOLD, 1.0/benchmarkModels.length);
  
  // Prepare HTML Output
  let html = `
        <h2>Global Alert Factor Analysis</h2>
        <p>Number of valid benchmarks modeled: ${benchmarkModels.length}</p>
        <p>Recommended global alert Threshold: ${globalF.toFixed(4)}</p>
  
        <h2>Benchmark Statistics</h2>
        <div>
          <table>
            <tr>
              <th>Benchmark</th>
              <th>Count</th>
              <th>Mean (µs)</th>
              <th>StdDev (µs)</th>
              <th>StdDev %</th>
              <th>Ideal Factor for ${((1.0 - singleThreshold) * 100).toFixed(2)}%</th>
              <th>False Positive Risk at ${globalF.toFixed(4)}</th>
            </tr>
            `;
  
  // 2. Compute Mean and StdDev per benchmark (Absolute Values)
  for (const model of benchmarkModels) {
    const name = model.name;
    const values = stats[name];
    const n = values.length;
    if (n < 2) continue;
    const mean = values.reduce((sum, v) => sum + v, 0) / n;
    const variance = values.reduce((sum, v) => sum + Math.pow(v - mean, 2), 0) / (n - 1);
    const stddev = Math.sqrt(variance);
    const f = searchFactor([model], singleThreshold);
    const pAlert = 1.0 - pNoAlertOf(model, globalF);
    html += `
            <tr>
              <td>${name}</td>
              <td>${n}</td>
              <td>${mean.toFixed(4)}</td>
              <td>${stddev.toFixed(4)}</td>
              <td>${((stddev / mean) * 100).toFixed(2)}%</td>
              <td>${f == null ? '(undefined)' : f.toFixed(4)}</td>
              <td>${pAlert < 0.00005 ? '' : (pAlert * 100).toFixed(2) + '%'}</td>
            </tr>
            `;
  }
  html += `</table>`;
  
  html += `</div>`; // end container
  document.write(html);
})(window.BENCHMARK_DATA);
