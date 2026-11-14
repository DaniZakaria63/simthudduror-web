import { readFile } from "node:fs/promises";

const summary = JSON.parse(await readFile(new URL("../coverage/coverage-summary.json", import.meta.url), "utf8"));
const total = summary.total;
const metrics = ["statements", "branches", "functions", "lines"];
let failed = false;

for (const metric of metrics) {
  const pct = total[metric]?.pct ?? 0;
  const passed = pct > 80;
  console.log(`${metric}: ${pct}% ${passed ? "OK" : "FAIL (must be strictly > 80%)"}`);
  if (!passed) {
    failed = true;
  }
}

if (failed) {
  process.exit(1);
}
