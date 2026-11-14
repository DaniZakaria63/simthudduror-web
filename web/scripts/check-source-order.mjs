import { readdir, readFile } from "node:fs/promises";
import { dirname, extname, join, relative } from "node:path";
import { fileURLToPath } from "node:url";

const root = join(dirname(fileURLToPath(import.meta.url)), "..");
const sourceRoot = join(root, "src");
const violations = [];

async function walk(dir) {
  const entries = await readdir(dir, { withFileTypes: true });
  const files = [];
  for (const entry of entries) {
    const full = join(dir, entry.name);
    if (entry.isDirectory()) {
      files.push(...(await walk(full)));
    } else if (extname(entry.name) === ".tsx") {
      files.push(full);
    }
  }
  return files;
}

const importEnd = /from\s+["'][^"']+["'];?\s*$/;
const bareImport = /^import\s+["'][^"']+["'];?\s*$/;
const cssImport = /^import\s+["'][^"']+\.css["'];?\s*$/;

for (const file of await walk(sourceRoot)) {
  const label = relative(root, file).replaceAll("\\", "/");
  if (label.includes(".test.") || label.endsWith("src/main.tsx")) {
    continue;
  }
  const lines = (await readFile(file, "utf8")).split("\n");
  const nonEmpty = lines
    .map((line, index) => ({ line: line.trim(), index }))
    .filter((item) => item.line);
  const hasClassName = lines.some((line) => line.includes("className="));
  const cssImports = nonEmpty.filter((item) => cssImport.test(item.line));

  if (hasClassName) {
    if (cssImports.length !== 1) {
      violations.push(`${label}: expected exactly one colocated stylesheet import`);
    } else if (nonEmpty[nonEmpty.length - 1].index !== cssImports[0].index) {
      violations.push(`${label}: stylesheet import must be the last non-empty line`);
    }
  } else if (cssImports.length !== 0) {
    violations.push(`${label}: unexpected stylesheet import`);
  }

  let inImport = false;
  let codeStarted = false;
  let inBlockComment = false;
  for (const { line, index } of nonEmpty) {
    if (inBlockComment) {
      if (line.includes("*/")) {
        inBlockComment = false;
      }
      continue;
    }
    if (inImport) {
      if (importEnd.test(line) || (line.endsWith('";') && line.startsWith('"'))) {
        inImport = false;
      }
      continue;
    }
    if (line.startsWith("/*")) {
      if (!line.includes("*/")) {
        inBlockComment = true;
      }
      continue;
    }
    if (line.startsWith("//")) {
      continue;
    }
    if (line.startsWith("import ") || line === "import") {
      if (codeStarted && !cssImport.test(line)) {
        violations.push(
          `${label}:${index + 1}: import after component code (only the trailing stylesheet may follow)`,
        );
      }
      if (!importEnd.test(line) && !bareImport.test(line)) {
        inImport = true;
      }
      continue;
    }
    codeStarted = true;
  }
}

if (violations.length > 0) {
  console.error("Source order violations:\n" + violations.map((item) => `- ${item}`).join("\n"));
  process.exit(1);
}
console.log("Source order OK");
