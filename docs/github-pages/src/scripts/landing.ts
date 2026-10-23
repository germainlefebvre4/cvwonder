// Landing entry. Loaded only on the pages that render the Hero override. It stays tiny: the
// navigation and the no-motion tab set run right away, and the animation library is
// imported after first paint so the hero HTML is never held back. With reduced motion
// requested, nothing animates and the library is never fetched.
import { setupNav, applyHash } from './nav';
import { setupThemeTabs } from './tabs';

document.documentElement.dataset.cvwLanding = '';

const scene = document.querySelector<HTMLElement>('[data-scene]');
const reduced = matchMedia('(prefers-reduced-motion: reduce)').matches;
const wide = matchMedia('(min-width: 900px)');

setupNav();

// Static rendering: tabs for narrow screens and reduced motion (also the fallback if the
// animation library fails to load).
if (scene && (reduced || !wide.matches)) setupThemeTabs(scene);

if (!reduced) {
  const load = () =>
    import('./animate')
      .then((m) => m.start(scene, setupThemeTabs))
      .then(() => applyHash())
      .catch(() => {
        // Library failed to load: the static page is already complete.
        if (scene && !scene.dataset.mode) setupThemeTabs(scene);
      });
  if ('requestIdleCallback' in window) requestIdleCallback(load, { timeout: 1200 });
  else setTimeout(load, 200);
}
