// Asserts on the built site (dist/) that the animated landing page stays healthy:
//  - every rail and skip target id exists in both landing pages,
//  - the animation script is referenced only by the two landing pages,
//  - the script the landing adds stays within the gzip budget.
// Run after `npm run build`: `npm run check:landing`.
import { readFileSync, readdirSync, statSync, existsSync } from 'node:fs';
import { gzipSync } from 'node:zlib';
import { join, relative, resolve, dirname } from 'node:path';
import { fileURLToPath } from 'node:url';

const MAX_GZIP = Number(process.env.LANDING_MAX_GZIP ?? 60 * 1024);
const dist = fileURLToPath(new URL('../dist/', import.meta.url));
const LANDINGS = ['index.html', 'fr/index.html'];
// Set by src/scripts/landing.ts; used to recognize the landing entry in the build output.
const MARKER = 'cvwLanding';

const failures = [];
const fail = (msg) => failures.push(msg);

if (!existsSync(dist)) {
  console.error('dist/ not found: run `npm run build` first.');
  process.exit(1);
}

function* htmlFiles(dir) {
  for (const name of readdirSync(dir)) {
    const path = join(dir, name);
    if (statSync(path).isDirectory()) yield* htmlFiles(path);
    else if (name.endsWith('.html')) yield path;
  }
}

const read = (path) => readFileSync(path, 'utf8');
const scriptSrcs = (html) => [...html.matchAll(/<script[^>]*\ssrc="([^"]+)"/g)].map((m) => m[1]);
const toFile = (src) => join(dist, src.replace(/^\//, ''));

// Static and dynamic relative imports of a built chunk.
function dependencies(file) {
  const code = read(file);
  const found = new Set();
  for (const m of code.matchAll(/(?:from|import)\s*\(?\s*["'`](\.{1,2}\/[^"'`]+\.js)["'`]/g)) {
    found.add(resolve(dirname(file), m[1]));
  }
  return [...found];
}

function closure(entry) {
  const seen = new Set();
  const walk = (file) => {
    if (seen.has(file) || !existsSync(file)) return;
    seen.add(file);
    dependencies(file).forEach(walk);
  };
  walk(entry);
  return seen;
}

const landingFiles = new Set(); // every file reachable from a landing entry
const landingEntries = new Set();
for (const page of LANDINGS) {
  const path = join(dist, page);
  if (!existsSync(path)) {
    fail(`${page}: not built`);
    continue;
  }
  const html = read(path);
  const entries = scriptSrcs(html).map(toFile).filter((f) => existsSync(f) && read(f).includes(MARKER));
  if (!entries.length) fail(`${page}: landing script not referenced`);
  entries.forEach((e) => {
    landingEntries.add(e);
    closure(e).forEach((f) => landingFiles.add(f));
  });

  // Rail and skip targets must exist as ids in the page.
  const targets = new Set([...html.matchAll(/data-section="([^"]+)"/g)].map((m) => m[1]));
  for (const m of html.matchAll(/<a[^>]*data-skip[^>]*href="#([^"]+)"|<a[^>]*href="#([^"]+)"[^>]*data-skip/g)) targets.add(m[1] ?? m[2]);
  for (const m of html.matchAll(/data-jump="([^"]+)"/g)) targets.add(m[1]);
  if (!targets.size) fail(`${page}: no rail or skip targets found`);
  for (const id of targets) {
    if (!new RegExp(`\\sid="${id}"`).test(html)) fail(`${page}: target id "${id}" does not exist`);
  }
}

// Only the landing pages may reference the landing script (or any chunk only it pulls in).
const landingNames = [...landingEntries].map((f) => f.split('/').pop());
for (const file of htmlFiles(dist)) {
  const rel = relative(dist, file);
  if (LANDINGS.includes(rel)) continue;
  const html = read(file);
  for (const name of landingNames) {
    if (html.includes(name)) fail(`${rel}: references the landing script ${name}`);
  }
}

// Size of the script the landing adds: the animation chunks, not Starlight's shared scripts.
let total = 0;
const shared = new Set();
for (const page of ['start/install/index.html']) {
  const path = join(dist, page);
  if (existsSync(path)) scriptSrcs(read(path)).map(toFile).forEach((f) => closure(f).forEach((c) => shared.add(c)));
}
for (const file of landingFiles) {
  if (shared.has(file)) continue;
  total += gzipSync(readFileSync(file)).length;
}
console.log(`landing script: ${(total / 1024).toFixed(1)} KB gzip (limit ${(MAX_GZIP / 1024).toFixed(0)} KB)`);
if (total > MAX_GZIP) fail(`landing script is ${(total / 1024).toFixed(1)} KB gzip, above ${(MAX_GZIP / 1024).toFixed(0)} KB`);

if (failures.length) {
  console.error(`\nLanding check failed:\n${failures.map((f) => `  - ${f}`).join('\n')}`);
  process.exit(1);
}
console.log('PASS landing page checks');
