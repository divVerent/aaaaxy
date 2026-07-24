const TARGET_FALSE_ALERT_PROBABILITY = 0.005;
const PERCENTILE_THRESHOLD = 1.0 - TARGET_FALSE_ALERT_PROBABILITY;

function parseExtra(extra) {
  out = {};
  for (const line of extra.split(/\n/)) {
    const m = /^(?<key>\w+): (?<value>.*)$/.exec(line);
    out[m.groups.key] = m.groups.value;
  }
  return out;
}

function cpuTypeOf(bench) {
  const extra = parseExtra(bench.extra);
  return `${extra.azure_vmsize} ${extra.cpu_model}`;
}

function normalCDF(x) {
  const sign = x < 0 ? -1 : 1;
  x = Math.abs(x) / Math.sqrt(2.0);
  const t = 1.0 / (1.0 + 0.3275911 * x);
  const y = 1.0 - (((((1.061405429 * t - 1.453152027) * t) + 1.421413741) * t - 0.284496736) * t + 0.254829592) * t * Math.exp(-x * x);
  return 0.5 * (1.0 + sign * y);
}

// --- Analytical Student's t Cumulative Distribution Function (CDF) ---
function tCDF(t, df) {
  if (df >= 100) return normalCDF(t); // Fall back to your shared normalCDF for large samples
  
  const theta = Math.atan2(t, Math.sqrt(df));
  if (df === 1) return 0.5 + theta / Math.PI;
  
  const cos = Math.cos(theta);
  const sin = Math.sin(theta);
  
  if (df % 2 === 0) {
    let sum = 1.0, term = 1.0;
    for (let i = 2; i <= df - 2; i += 2) {
      term *= ((i - 1) / i) * (cos * cos);
      sum += term;
    }
    return 0.5 + 0.5 * sin * sum;
  } else {
    let sum = cos, term = cos;
    for (let i = 3; i <= df - 2; i += 2) {
      term *= ((i - 1) / i) * (cos * cos);
      sum += term;
    }
    return 0.5 + (theta + sin * sum) / Math.PI;
  }
}

if (typeof module != 'undefined') {
  module.exports = {PERCENTILE_THRESHOLD, cpuTypeOf, normalCDF, tCDF};
}
