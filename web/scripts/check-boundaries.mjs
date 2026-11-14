import { readdir, readFile } from "node:fs/promises";
import { dirname, extname, join, normalize, relative, sep } from "node:path";
import { fileURLToPath } from "node:url";

const root = normalize(join(dirname(fileURLToPath(import.meta.url)), ".."));
const sourceRoot = join(root, "src");
const violations = [];

async function walk(dir) {
  const entries = await readdir(dir, { withFileTypes: true });
  const files = [];
  for (const entry of entries) {
    const full = join(dir, entry.name);
    if (entry.isDirectory()) {
      files.push(...(await walk(full)));
    } else if ([".ts", ".tsx"].includes(extname(entry.name))) {
      files.push(full);
    }
  }
  return files;
}

function featureOf(path) {
  const parts = relative(sourceRoot, path).split(sep);
  const index = parts.indexOf("features");
  return index >= 0 ? parts[index + 1] : null;
}

function resolvedTarget(file, specifier) {
  let target;
  if (specifier.startsWith("@/")) {
    target = join(sourceRoot, specifier.slice(2));
  } else if (specifier.startsWith(".")) {
    target = normalize(join(dirname(file), specifier));
  } else {
    return null;
  }
  return target;
}

function isDeepFeatureImport(target) {
  const parts = relative(sourceRoot, target).split(sep);
  const index = parts.indexOf("features");
  if (index < 0) {
    return false;
  }
  const feature = parts[index + 1];
  const rest = parts.slice(index + 2).filter((part) => part !== "" && part !== "index.ts" && part !== "index");
  return { feature, deep: rest.length > 0 };
}

for (const file of await walk(sourceRoot)) {
  const source = await readFile(file, "utf8");
  const ownFeature = featureOf(file);
  const label = relative(root, file);

  if (/dangerouslySetInnerHTML/.test(source)) {
    violations.push(`${label}: dangerouslySetInnerHTML is not allowed`);
  }
  if (/\b(?:window\.)?(?:localStorage|sessionStorage)\.(?:getItem|setItem|removeItem|clear)\b/.test(source)) {
    violations.push(`${label}: browser storage is not allowed for credentials`);
  }
  if (/document\.cookie/.test(source)) {
    violations.push(`${label}: document.cookie is not allowed`);
  }
  if (!label.endsWith(".test.ts") && !label.endsWith(".test.tsx") && label !== "src/shared/api/client.ts" && /\bfetch\s*\(/.test(source)) {
    violations.push(`${label}: use the shared API client instead of raw fetch`);
  }

  const imports = [...source.matchAll(/(?:from\s+|import\s*)["']([^"']+)["']/g)].map((match) => match[1]);
  for (const specifier of imports) {
    if (!specifier.startsWith(".") && !specifier.startsWith("@/")) {
      continue;
    }
    if (specifier === "@/features" || specifier.startsWith("@/features/")) {
      const parts = specifier.split("/");
      if (!ownFeature && label.startsWith("src/shared/")) {
        violations.push(`${label}: shared code must not import feature ${parts[2]}`);
      }
      if (ownFeature && parts[2] !== ownFeature && parts.length > 3) {
        violations.push(`${label}: import feature ${parts[2]} through its public entry point`);
      }
      continue;
    }
    if (!ownFeature && !label.startsWith("src/shared/") && !label.startsWith("src/App.tsx")) {
      continue;
    }
    const target = resolvedTarget(file, specifier);
    if (!target) {
      continue;
    }
    const info = isDeepFeatureImport(target);
    if (!info) {
      continue;
    }
    if (!ownFeature && label.startsWith("src/shared/")) {
      violations.push(`${label}: shared code must not import feature ${info.feature}`);
      continue;
    }
    if (ownFeature && info.feature !== ownFeature && info.deep) {
      violations.push(`${label}: import feature ${info.feature} through its public entry point`);
    }
  }
}

if (violations.length > 0) {
  console.error("Feature boundary violations:\n" + violations.map((item) => `- ${item}`).join("\n"));
  process.exit(1);
}
console.log("Feature boundaries OK");
