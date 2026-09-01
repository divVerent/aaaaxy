const DEFAULT_CONFIG = {
  // Anomaly detection.
  'minFraction': 0.75, // Time fraction in which to not show anomalies.
  'noiseFloorAbsolute': 0.01, // Even if not seen, observe a min stddev of 10ns.
  'noiseFloorRelative': 0.01, // Even if not seen, observe a min stddev of 1%.
  'targetFalseAlertProbability': 0.005,

  // Graph filtering.
  'graphFalloff': 0.98,
  'graphMinScore': -1,

  // Custom filtering.
  'cpuFilter': '',
  'nameFilter': '',

  // Graph appearance.
  'graphColorMax': 223,
  'graphColorMin': 32,
  'graphColorStep': 0.25,
};

let CONFIG = Object.assign({}, DEFAULT_CONFIG);

function loadConfigFromString(str) {
  if (str == '') {
    return;
  }
  Object.assign(CONFIG, JSON.parse(str));
}

function configToString(extra) {
  const obj = Object.assign(Object.assign({}, CONFIG), extra);
  let diff = {};
  for (const [key, val] of Object.entries(obj)) {
    if (DEFAULT_CONFIG[key] === obj[key]) {
      continue;
    }
    if (val == null) {
      continue;
    }
    diff[key] = val;
  }
  return JSON.stringify(diff);
}

function parseExtra(extra) {
  let out = {};
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
    let sum = 1.0,
      term = 1.0;
    for (let i = 2; i <= df - 2; i += 2) {
      term *= ((i - 1) / i) * (cos * cos);
      sum += term;
    }
    return 0.5 + 0.5 * sin * sum;
  } else {
    let sum = cos,
      term = cos;
    for (let i = 3; i <= df - 2; i += 2) {
      term *= ((i - 1) / i) * (cos * cos);
      sum += term;
    }
    return 0.5 + (theta + sin * sum) / Math.PI;
  }
}

function colorForGraph(goodScore, badScore) {
  const score = goodScore + badScore;
  const goodColorVal = Math.round(CONFIG.graphColorMin + (CONFIG.graphColorMax - CONFIG.graphColorMin) / (1 + goodScore * CONFIG.graphColorStep));
  const badColorVal = Math.round(CONFIG.graphColorMin + (CONFIG.graphColorMax - CONFIG.graphColorMin) / (1 + badScore * CONFIG.graphColorStep));
  const colorVal = Math.round(CONFIG.graphColorMin + (CONFIG.graphColorMax - CONFIG.graphColorMin) / (1 + score * CONFIG.graphColorStep));
  const colorR = goodColorVal;
  const colorG = badColorVal;
  const colorB = colorVal;
  return color = '#' + colorR.toString(16).padStart(2, '0') + colorG.toString(16).padStart(2, '0') + colorB.toString(16).padStart(2, '0');
}

function scoreForGraph(dataset) {
  let score = 0;
  let goodScore = 0;
  let badScore = 0;
  for (const d of dataset) {
    score *= CONFIG.graphFalloff;
    goodScore *= CONFIG.graphFalloff;
    badScore *= CONFIG.graphFalloff;
    if (d.bench.anomaly > 0) {
      score += 1.0;
      badScore += 1.0;
    } else if (d.bench.anomaly < 0) {
      score += 1.0;
      goodScore += 1.0;
    }
  }
  return [score, goodScore, badScore];
}

const DEFAULT_GRAPH_MIN_SCORE = 1e-10;

function graphMinScore() {
  return CONFIG.graphMinScore >= 0 ? CONFIG.graphMinScore :
         CONFIG.nameFilter         ? 0                    :
                                     DEFAULT_GRAPH_MIN_SCORE;
}

if (typeof module != 'undefined') {
  module.exports = {
    CONFIG,
    loadConfigFromString,
    configToString,
    parseExtra,
    cpuTypeOf,
    normalCDF,
    tCDF,
    colorForGraph,
    scoreForGraph,
    graphMinScore,
  };
}
