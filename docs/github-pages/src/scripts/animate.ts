// Animation chunk (GSAP + ScrollTrigger): hero entrance, aurora drift, pointer effects,
// scroll reveals and the pinned YAML-to-CV scene. Loaded lazily by landing.ts.
import { gsap } from 'gsap';
import { ScrollTrigger } from 'gsap/ScrollTrigger';
import { landing, LAYOUT_EVENT } from './state';

gsap.registerPlugin(ScrollTrigger);

const $ = <T extends Element = HTMLElement>(sel: string, root: ParentNode = document) => root.querySelector<T>(sel);
const $$ = <T extends Element = HTMLElement>(sel: string, root: ParentNode = document) => [...root.querySelectorAll<T>(sel)];
const headerHeight = () => $('header.header')?.getBoundingClientRect().height ?? 0;

/** Named moments of the scene on its 0..100 timeline. Tabs, skip and step highlight derive from these. */
export const MOMENTS = {
  // Acts: [start, end]. A tab shows its act fully played, i.e. just before the next one begins.
  write: { start: 0, end: 25.5 },
  generate: { start: 26, end: 60.5 },
  theme: { start: 61, end: 99.5 },
  // Theme changes: the wipe is centered on these values.
  wipe1: 72,
  wipe2: 90,
  // Horizon timeline is fully shown (first wipe done at 77.5) and its capture has not started to scroll (78).
  horizonShown: 77.7,
} as const;
const actAt = (p: number) => (p < MOMENTS.generate.start ? 0 : p < MOMENTS.theme.start ? 1 : 2);
const themeAt = (p: number) => (p < MOMENTS.wipe1 ? 0 : p < MOMENTS.wipe2 ? 1 : 2);

export async function start(scene: HTMLElement | null, setupThemeTabs: (scene: HTMLElement) => void) {
  document.documentElement.classList.add('cvw-anim');
  hero();
  aurora();
  reveals();
  pointerEffects();

  if (scene) {
    const mm = gsap.matchMedia();
    mm.add('(min-width: 900px)', () => pinnedScene(scene));
    mm.add('(max-width: 899.98px)', () => {
      // No pin below 900px: stacked layout with the theme chips acting as tabs.
      if (scene.dataset.mode !== 'tabs') setupThemeTabs(scene);
    });
  }
  ScrollTrigger.refresh();
  window.dispatchEvent(new Event(LAYOUT_EVENT));
}

/* ------------------------------------------------------------------ hero */

function hero() {
  const words = $$('[data-word]');
  const tagline = $('[data-tagline]');
  const actions = $('[data-actions]');
  const term = $('[data-term]');
  const cv = $('[data-cv-float]');
  const rows = $$('[data-row]');

  // Terminal: reset to "nothing typed" without changing any box size.
  const typed = rows.map((row) => {
    const cmd = $('[data-cmd]', row)!;
    const full = cmd.dataset.cmd ?? '';
    const on = $('.on', cmd)!;
    const off = $('.off', cmd)!;
    on.textContent = '';
    off.textContent = full;
    gsap.set($('[data-out]', row), { opacity: 0 });
    return { on, off, full, out: $('[data-out]', row)! };
  });

  const tl = gsap.timeline({ defaults: { ease: 'power3.out' } });
  tl.from(words, { opacity: 0, y: 26, duration: 0.7, stagger: 0.09 })
    .from(tagline, { opacity: 0, y: 16, duration: 0.6 }, '-=0.35')
    .from(actions, { opacity: 0, y: 16, duration: 0.6 }, '-=0.4')
    .from(term, { opacity: 0, y: 36, duration: 0.8 }, '-=0.3')
    .from(cv, { opacity: 0, y: 48, duration: 0.9 }, '-=0.6');
  for (const row of typed) {
    const counter = { n: 0 };
    tl.to(counter, {
      n: row.full.length,
      duration: Math.max(0.4, row.full.length * 0.045),
      ease: 'none',
      onUpdate() {
        const n = Math.round(counter.n);
        row.on.textContent = row.full.slice(0, n);
        row.off.textContent = row.full.slice(n);
      },
    }).to(row.out, { opacity: 1, duration: 0.25 }, '+=0.1');
    tl.to({}, { duration: 0.35 });
  }

  // The floating CV drifts gently.
  const inner = $('[data-cv-inner]');
  if (inner) gsap.to(inner, { y: -10, duration: 3.2, ease: 'sine.inOut', yoyo: true, repeat: -1 });
}

function aurora() {
  const blobs = $$('[data-aurora] i');
  const drift: Array<[number, number, number]> = [
    [6, 4, 22],
    [-5, 6, 26],
    [4, -5, 30],
  ];
  blobs.forEach((el, i) => {
    const [x, y, d] = drift[i];
    gsap.to(el, { xPercent: x * 3, yPercent: y * 3, duration: d, ease: 'sine.inOut', yoyo: true, repeat: -1 });
  });
}

/* ------------------------------------------------------- pointer effects */

function pointerEffects() {
  const fine = matchMedia('(hover: hover) and (pointer: fine)').matches;
  if (!fine) return;

  // Terminal and CV tilt and shift with the pointer; they rest again when it leaves.
  const stage = $('[data-hero-stage]');
  const term = $('[data-term]');
  const cv = $('[data-cv-float]');
  if (stage && term && cv) {
    gsap.set([term, cv], { transformPerspective: 1100 });
    const t = { ry: gsap.quickTo(term, 'rotationY', { duration: 0.6 }), rx: gsap.quickTo(term, 'rotationX', { duration: 0.6 }) };
    const c = {
      ry: gsap.quickTo(cv, 'rotationY', { duration: 0.8 }),
      rx: gsap.quickTo(cv, 'rotationX', { duration: 0.8 }),
      x: gsap.quickTo(cv, 'x', { duration: 0.8 }),
      y: gsap.quickTo(cv, 'y', { duration: 0.8 }),
    };
    stage.addEventListener('pointermove', (e) => {
      if (e.pointerType !== 'mouse') return;
      const r = stage.getBoundingClientRect();
      const nx = (e.clientX - r.left) / r.width - 0.5;
      const ny = (e.clientY - r.top) / r.height - 0.5;
      t.ry(nx * 8);
      t.rx(-ny * 6);
      c.ry(nx * 12);
      c.rx(-ny * 8);
      c.x(nx * 24);
      c.y(ny * 16);
    });
    stage.addEventListener('pointerleave', () => {
      for (const f of [t.ry, t.rx, c.ry, c.rx, c.x, c.y]) f(0);
    });
  }

  // Card spotlight follows the pointer.
  for (const card of $$('.card')) {
    card.addEventListener('pointermove', (e) => {
      const r = card.getBoundingClientRect();
      card.style.setProperty('--mx', `${e.clientX - r.left}px`);
      card.style.setProperty('--my', `${e.clientY - r.top}px`);
    });
  }

  // Magnetic call to action.
  const cta = $('.cta');
  const button = $('.cta a');
  if (cta && button) {
    const x = gsap.quickTo(button, 'x', { duration: 0.4, ease: 'power3.out' });
    const y = gsap.quickTo(button, 'y', { duration: 0.4, ease: 'power3.out' });
    cta.addEventListener('pointermove', (e) => {
      const r = button.getBoundingClientRect();
      const dx = e.clientX - (r.left + r.width / 2);
      const dy = e.clientY - (r.top + r.height / 2);
      const near = Math.hypot(dx, dy) < 140;
      x(near ? dx * 0.25 : 0);
      y(near ? dy * 0.25 : 0);
    });
    cta.addEventListener('pointerleave', () => {
      x(0);
      y(0);
    });
  }
}

/* --------------------------------------------------------------- reveals */

function reveals() {
  const groups = [
    { trigger: '#why .card-grid', items: '#why .card' },
    { trigger: '#start .card-grid', items: '#start .card' },
    { trigger: '#get-started', items: '#get-started' },
  ];
  const pending: Array<() => void> = [];
  for (const g of groups) {
    const trigger = $(g.trigger);
    const items = $$(g.items);
    if (!trigger || !items.length) continue;
    // Only hide what is still below the fold: never anything the visitor can already see.
    if (trigger.getBoundingClientRect().top < window.innerHeight * 0.9) continue;
    gsap.set(items, { opacity: 0, y: 28 });
    let tween: gsap.core.Tween | null = null;
    const play = () => {
      if (!tween) tween = gsap.to(items, { opacity: 1, y: 0, duration: 0.7, ease: 'power3.out', stagger: 0.1 });
    };
    ScrollTrigger.create({ trigger, start: 'top 88%', once: true, onEnter: play });
    // After an instant jump, whatever is in view must be fully shown, not still fading in.
    pending.push(() => {
      if (trigger.getBoundingClientRect().top >= window.innerHeight * 0.9) return;
      play();
      tween?.progress(1);
    });
  }
  landing.flushReveals = () => pending.forEach((f) => f());
}

/* ----------------------------------------------------------------- scene */

function pinnedScene(scene: HTMLElement) {
  const pinEl = $('[data-pin]', scene)!;
  const yaml = $('[data-yaml]', scene)!;
  const cv = $('[data-cv]', scene)!;
  const viewport = $('[data-viewport]', scene)!;
  const lines = $$('[data-line]', scene);
  const tinted = lines.filter((l) => !l.hasAttribute('data-theme-line'));
  const themeLine = $('[data-theme-line]', scene);
  const themeValue = $('[data-theme-name]', scene);
  const themeName = $('[data-theme-name]', scene);
  const previews = $$('[data-preview]', scene);
  const imgs = previews.map((p) => $('img', p)!);
  const veil = $('[data-veil]', scene)!;
  const scan = $('.scan', veil)!;
  const edge = $('[data-edge]', scene)!;
  const badges = $$('[data-badge]', scene);
  const titles = $$('[data-act-title]', scene);
  const chipsEl = $('[data-chips]', scene)!;
  const chips = $$('[data-chip]', scene);
  const tabs = $$('[data-tab]', scene);
  const bar = $('[data-progress]', scene)!;
  const skip = $('[data-skip]');
  const names = previews.map((p) => p.dataset.themeName ?? '');

  scene.dataset.mode = 'pin';
  chips.forEach((c, i) => (c.tabIndex = -1, c.setAttribute('aria-selected', String(i === 0))));
  previews.forEach((p) => (p.hidden = false));

  const scrollRange = (img: HTMLElement) => () => -Math.max(0, img.offsetHeight - viewport.offsetHeight);
  const gap = () => cv.offsetLeft - (yaml.offsetLeft + yaml.offsetWidth);

  const tl = gsap.timeline({ defaults: { ease: 'none' } });
  const at = (p: number) => p; // timeline position equals the 0..100 moment
  // Offsets below are layout-based (transform-free), so they hold whatever the timeline state is.

  // Initial states.
  gsap.set(cv, { opacity: 0, transformPerspective: 1400, rotationY: -26, x: 90, scale: 0.92 });
  gsap.set(yaml, { x: () => (cv.offsetWidth + gap()) / 2 });
  gsap.set(badges, { scale: 0, opacity: 0 });
  gsap.set(veil, { yPercent: 0 });
  gsap.set(scan, { opacity: 1 });
  previews.forEach((p, i) => gsap.set(p, { zIndex: i + 1, clipPath: i === 0 ? 'inset(0 0% 0 0)' : 'inset(0 100% 0 0)' }));
  gsap.set(edge, { xPercent: -100, opacity: 0 });
  gsap.set(titles.slice(1), { opacity: 0, y: 14 });
  gsap.set(tinted, { '--tint': 0 });
  if (themeLine) gsap.set(themeLine, { '--focus': 0 });

  // Act 1: write. Lines light up one after another.
  tl.fromTo(lines, { opacity: 0.22, x: -8 }, { opacity: 1, x: 0, duration: 4, stagger: 0.75, ease: 'power2.out' }, at(3));

  // Act titles cross-fade.
  const swap = (from: number, to: number, p: number) => {
    tl.to(titles[from], { opacity: 0, y: -14, duration: 2 }, at(p - 1)).to(titles[to], { opacity: 1, y: 0, duration: 2 }, at(p));
  };
  swap(0, 1, MOMENTS.generate.start);
  swap(1, 2, MOMENTS.theme.start);

  // Act 2: the YAML stays in place (lines take a resting tint, `theme:` is the focus) and the CV enters in 3D.
  tl.to(tinted, { '--tint': 1, duration: 3, stagger: 0.2, ease: 'power1.out' }, at(27));
  if (themeLine) tl.to(themeLine, { '--focus': 1, duration: 3, ease: 'power1.out' }, at(28));
  tl.to(yaml, { x: 0, duration: 8, ease: 'power2.inOut' }, at(29));
  tl.to(cv, { opacity: 1, rotationY: 0, x: 0, scale: 1, duration: 8, ease: 'power3.out' }, at(33));
  tl.to(veil, { yPercent: 100, duration: 10, ease: 'none' }, at(42));
  tl.to(scan, { opacity: 0, duration: 0.5 }, at(51.5));
  tl.to(badges, { scale: 1, opacity: 1, duration: 2, stagger: 0.5, ease: 'back.out(2)' }, at(52));
  tl.to(imgs[0], { y: scrollRange(imgs[0]), duration: 7, ease: 'power1.inOut' }, at(53.5));

  // Act 3: chips appear, scroll back to the top, then each theme change is a wipe.
  tl.to(chipsEl, { opacity: 1, duration: 3 }, at(MOMENTS.theme.start + 1));
  tl.to(imgs[0], { y: 0, duration: 4, ease: 'power1.inOut' }, at(63));
  const wipe = (i: number, center: number) => {
    tl.set(edge, { xPercent: -100, opacity: 1 }, at(center - 5))
      .to(previews[i], { clipPath: 'inset(0 0% 0 0)', duration: 10 }, at(center - 5))
      .to(edge, { xPercent: 0, duration: 10 }, at(center - 5))
      .to(edge, { opacity: 0, duration: 0.5 }, at(center + 5));
  };
  wipe(1, MOMENTS.wipe1);
  wipe(2, MOMENTS.wipe2);
  tl.to(imgs[1], { y: scrollRange(imgs[1]), duration: 5, ease: 'power1.inOut' }, at(78));
  tl.set({}, {}, at(100)); // pin the timeline length to 100

  // Derived UI: act tab, theme chip, `theme:` YAML line and progress follow the scroll position.
  let act = -1;
  let theme = -1;
  let flash = false; // no flash for the first sync (page load)
  const sync = (progress: number) => {
    const p = progress * 100;
    gsap.set(bar, { scaleX: progress });
    const a = actAt(p);
    if (a !== act) {
      act = a;
      tabs.forEach((t, i) => t.setAttribute('aria-pressed', String(i === a)));
    }
    const th = themeAt(p);
    if (th !== theme) {
      theme = th;
      chips.forEach((c, i) => c.setAttribute('aria-selected', String(i === th)));
      if (themeName) themeName.textContent = names[th] ?? '';
      if (flash && themeValue) {
        themeValue.classList.remove('flash');
        void themeValue.offsetWidth; // restart the CSS animation
        themeValue.classList.add('flash');
      }
    }
  };

  const st = ScrollTrigger.create({
    trigger: scene,
    pin: pinEl,
    start: () => `top top+=${Math.round(headerHeight())}`,
    end: () => `+=${Math.round(window.innerHeight * 4.4)}`,
    scrub: true,
    animation: tl,
    invalidateOnRefresh: true,
    anticipatePin: 1,
    onUpdate: (self) => sync(self.progress),
    onRefresh: (self) => {
      landing.pin = { start: self.start, end: self.end };
      sync(self.progress);
    },
    onToggle: (self) => skip?.setAttribute('data-active', String(self.isActive)),
  });
  landing.pin = { start: st.start, end: st.end };
  sync(st.progress);
  flash = true;

  const sceneAt = (moment: number) => {
    window.scrollTo({ top: st.start + (st.end - st.start) * (moment / 100), behavior: 'instant' });
    requestAnimationFrame(() => landing.flushReveals?.());
  };
  landing.sceneAt = sceneAt;

  // Act tabs and theme chips jump to the matching moment of the timeline.
  const acts = [MOMENTS.write.end, MOMENTS.generate.end, MOMENTS.horizonShown];
  const themeMoments = [MOMENTS.write.end, MOMENTS.wipe1 + 5, MOMENTS.wipe2 + 5];
  const onTab = (i: number) => () => sceneAt(acts[i]);
  const onChip = (i: number) => () => sceneAt(i === 0 ? MOMENTS.theme.start + 4 : i === 1 ? MOMENTS.horizonShown : themeMoments[i]);
  const tabHandlers = tabs.map((t, i) => [t, onTab(i)] as const);
  const chipHandlers = chips.map((c, i) => [c, onChip(i)] as const);
  [...tabHandlers, ...chipHandlers].forEach(([el, fn]) => el.addEventListener('click', fn));

  window.dispatchEvent(new Event(LAYOUT_EVENT));

  // Cleanup when the viewport drops below 900px.
  return () => {
    [...tabHandlers, ...chipHandlers].forEach(([el, fn]) => el.removeEventListener('click', fn));
    st.kill();
    tl.kill();
    gsap.set([cv, yaml, ...lines, ...(themeLine ? [themeLine] : []), ...badges, veil, scan, edge, chipsEl, ...titles, ...previews, ...imgs, bar], { clearProps: 'all' });
    gsap.set([...tinted, ...(themeLine ? [themeLine] : [])], { clearProps: '--tint,--focus' });
    themeValue?.classList.remove('flash');
    previews.forEach((p) => (p.hidden = false));
    skip?.setAttribute('data-active', 'false');
    delete scene.dataset.mode;
    landing.pin = null;
    landing.sceneAt = null;
  };
}
