// Asserts that every previously published URL has a redirect, and that every
// redirect target exists in the built site. Run after `npm run build`.
import { existsSync, readFileSync } from 'node:fs';
import { fileURLToPath } from 'node:url';
import { redirects } from '../redirects.ts';

const root = fileURLToPath(new URL('..', import.meta.url));
const fixture = JSON.parse(readFileSync(`${root}redirects.fixture.json`, 'utf8'));
const norm = (u) => (u.endsWith('/') ? u : `${u}/`);
const mapping = Object.fromEntries(Object.entries(redirects).map(([k, v]) => [norm(k), v]));

const errors = [];
for (const old of [...fixture, '/docs/']) {
  const target = mapping[norm(old)];
  if (!target) {
    errors.push(`no redirect for ${old}`);
    continue;
  }
  if (!existsSync(`${root}dist${target}index.html`)) errors.push(`${old} -> ${target}: target missing in dist/`);
  if (!existsSync(`${root}dist${norm(old)}index.html`)) errors.push(`${old}: redirect page missing in dist/`);
}

if (errors.length) {
  console.error(`Redirect check failed:\n- ${errors.join('\n- ')}`);
  process.exit(1);
}
console.log(`Redirect check passed: ${fixture.length + 1} legacy URLs redirect to existing pages.`);
