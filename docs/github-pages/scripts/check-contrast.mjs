// Computes the WCAG contrast ratio of every accent/background pair declared in
// src/styles/theme.css and fails below 4.5:1. Run: `npm run check:contrast`.
import { readFileSync } from 'node:fs';
import { fileURLToPath } from 'node:url';

const MIN = 4.5;
const css = readFileSync(fileURLToPath(new URL('../src/styles/theme.css', import.meta.url)), 'utf8');

// Collect `--name: value;` declarations per theme scope.
function tokens(selector) {
  const start = css.indexOf(`${selector} {`);
  if (start === -1) throw new Error(`selector ${selector} not found in theme.css`);
  const body = css.slice(start, css.indexOf('}', start));
  return Object.fromEntries([...body.matchAll(/(--[\w-]+):\s*([^;]+);/g)].map((m) => [m[1], m[2].trim()]));
}

const scopes = {
  dark: tokens(':root'),
  light: { ...tokens(':root'), ...tokens(":root[data-theme='light']") },
};

function resolve(scope, value) {
  const m = /^var\((--[\w-]+)\)$/.exec(value);
  return m ? resolve(scope, scope[m[1]]) : value;
}

function luminance(hex) {
  const h = hex.replace('#', '');
  const [r, g, b] = [0, 2, 4].map((i) => parseInt(h.slice(i, i + 2), 16) / 255).map((c) => (c <= 0.03928 ? c / 12.92 : ((c + 0.055) / 1.055) ** 2.4));
  return 0.2126 * r + 0.7152 * g + 0.0722 * b;
}

function ratio(a, b) {
  const [hi, lo] = [luminance(a), luminance(b)].sort((x, y) => y - x);
  return (hi + 0.05) / (lo + 0.05);
}

// [foreground token, background token, usage]
const pairs = [
  ['--cvw-accent', '--cvw-bg', 'accent text/links on page background'],
  ['--cvw-accent', '--cvw-bg-raised', 'accent text/links on raised surface'],
  ['--cvw-on-accent', '--cvw-accent', 'text on accent button'],
  ['--cvw-text', '--cvw-bg', 'body text on page background'],
  // Landing title gradient: every stop must be readable on its own.
  ...Object.keys(scopes.dark)
    .filter((name) => /^--cvw-grad-\d+$/.test(name))
    .map((name) => [name, '--cvw-bg', `gradient text stop ${name} on page background`]),
];

let failed = false;

// Light mode inherits dark values by default; the new landing tokens must be set in both.
const lightOwn = tokens(":root[data-theme='light']");
for (const name of Object.keys(scopes.dark).filter((n) => /^--cvw-(grad|aurora)-\d+$/.test(n))) {
  if (!(name in lightOwn)) {
    failed = true;
    console.log(`FAIL light token ${name} is not defined for light mode`);
  }
}
for (const [mode, scope] of Object.entries(scopes)) {
  for (const [fg, bg, usage] of pairs) {
    const r = ratio(resolve(scope, scope[fg]), resolve(scope, scope[bg]));
    const ok = r >= MIN;
    failed ||= !ok;
    console.log(`${ok ? 'PASS' : 'FAIL'} ${mode.padEnd(5)} ${r.toFixed(2).padStart(5)}:1  ${usage}`);
  }
}
if (failed) {
  console.error(`\nContrast check failed: at least one pair is below ${MIN}:1.`);
  process.exit(1);
}
