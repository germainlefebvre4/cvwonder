// Section rail, skip pill, hero link and deep links. Everything here is a real anchor that
// works without script; this code upgrades the links to instant jumps that do not play the
// pinned scene on the way, and computes the section being viewed.
import { landing, LAYOUT_EVENT } from './state';

const IDS = ['intro', 'scene', 'why', 'start', 'get-started'] as const;
type Id = (typeof IDS)[number];

const headerHeight = () => document.querySelector('header.header')?.getBoundingClientRect().height ?? 0;
const docTop = (el: Element) => el.getBoundingClientRect().top + window.scrollY;

function targetY(id: Id): number {
  if (id === 'intro') return 0;
  if (id === 'scene' && landing.pin) return landing.pin.start + 1; // one pixel in, so the pin engages
  const el = document.getElementById(id);
  return el ? Math.max(0, docTop(el) - headerHeight()) : window.scrollY;
}

export function jumpTo(id: Id, updateHash = true) {
  window.scrollTo({ top: targetY(id), behavior: 'instant' });
  if (updateHash) history.replaceState(null, '', `#${id}`);
  // Reveals are scroll-triggered: make sure the target section is not left hidden.
  requestAnimationFrame(() => landing.flushReveals?.());
}

const isId = (v: string | null | undefined): v is Id => !!v && (IDS as readonly string[]).includes(v);

/** Last section whose top has crossed the upper third of the viewport. */
function currentSection(): Id {
  const line = window.scrollY + window.innerHeight / 3;
  // A last section too short to reach the upper third is current once the page ends.
  if (window.scrollY + window.innerHeight >= document.documentElement.scrollHeight - 2) return IDS[IDS.length - 1];
  let current: Id = 'intro';
  for (const id of IDS) {
    const el = document.getElementById(id);
    if (!el) continue;
    const top = id === 'scene' && landing.pin ? landing.pin.start : id === 'intro' ? 0 : docTop(el);
    if (top <= line) current = id;
  }
  return current;
}

export function setupNav() {
  const links = [...document.querySelectorAll<HTMLAnchorElement>('[data-rail] a[data-section]')];

  document.addEventListener('click', (event) => {
    const a = (event.target as Element).closest<HTMLAnchorElement>('a[data-section], a[data-jump], a[data-skip]');
    if (!a || event.defaultPrevented || event.button !== 0 || event.metaKey || event.ctrlKey || event.shiftKey) return;
    const id = a.dataset.section ?? a.dataset.jump ?? (a.hasAttribute('data-skip') ? 'why' : undefined);
    if (!isId(id)) return;
    event.preventDefault();
    jumpTo(id);
  });

  let queued = false;
  const update = () => {
    queued = false;
    const id = currentSection();
    for (const link of links) {
      if (link.dataset.section === id) link.setAttribute('aria-current', 'true');
      else link.removeAttribute('aria-current');
    }
  };
  const schedule = () => {
    if (!queued) {
      queued = true;
      requestAnimationFrame(update);
    }
  };
  window.addEventListener('scroll', schedule, { passive: true });
  window.addEventListener('resize', schedule);
  window.addEventListener('load', schedule);
  window.addEventListener(LAYOUT_EVENT, schedule);
  update();
}

/** Deep links: the browser jumped before the pin spacing existed, so jump again. */
export function applyHash() {
  const id = decodeURIComponent(location.hash.slice(1));
  if (isId(id)) jumpTo(id, false);
}
