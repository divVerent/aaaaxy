((data, filterJSON) => {
  loadConfigFromString(filterJSON);
  const filter = CONFIG.filter;

  allFilters = {
    '': new Set()
  };

  function keepBench(subentry, bench) {
    cpuType = cpuTypeOf(bench);
    if (!(cpuType in allFilters)) {
      allFilters[cpuType] = new Set();
    }
    allFilters[cpuType].add(subentry.date);
    allFilters[''].add(subentry.date);
    return !filter || cpuType == filter;
  }

  function filterEntries(entries) {
    for (const subentries of Object.values(entries)) {
      for (const subentry of subentries) {
        subentry.benches = subentry.benches.filter((bench) => keepBench(subentry, bench));
      }
    }
  }
  filterEntries(data.entries);

  let html = `
    <h2>Available CPU Types</h2>
    <div>
      <table>
        <tr>
          <th>CPU Type</th>
          <th>Count</th>
        </tr>
  `;
  for (const [filterStr, set] of Object.entries(allFilters).toSorted(([aName, aSet], [bName, bSet]) => {
      if (aSet != bSet) {
        return bSet.size - aSet.size;
      }
      return aName.localeCompare(bName);
    })) {
    const filterJSON = configToString({'filter': filterStr.length ? filterStr : null});
    const url = `JavaScript:location.hash = '#${escape(filterJSON)}'; location.reload(); false;`;
    html += `
        <tr>
          <td><a href="${url}">${filterStr.length ? filterStr : '(all)'}</a></td>
          <td>${set.size}</td>
        </tr>
    `;
  }
  html += `
      </table>
    </div>
  `;
  document.write(html);
})(window.BENCHMARK_DATA, unescape(location.hash.substring(1)));
