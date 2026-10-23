// Theme chips as a tab set, for narrow screens and reduced motion: one preview at a time.
export function setupThemeTabs(scene: HTMLElement) {
  const chips = [...scene.querySelectorAll<HTMLButtonElement>('[data-chip]')];
  const previews = [...scene.querySelectorAll<HTMLElement>('[data-preview]')];
  const themeName = scene.querySelector<HTMLElement>('[data-theme-name]');
  const names = previews.map((p) => p.dataset.themeName);

  scene.dataset.mode = 'tabs';
  const select = (index: number, focus = false) => {
    if (scene.dataset.mode !== 'tabs') return;
    chips.forEach((chip, i) => {
      chip.setAttribute('aria-selected', String(i === index));
      chip.tabIndex = i === index ? 0 : -1;
      if (i === index && focus) chip.focus();
    });
    previews.forEach((p, i) => (p.hidden = i !== index));
    if (themeName && names[index]) themeName.textContent = names[index] ?? '';
  };
  chips.forEach((chip, i) => {
    chip.addEventListener('click', () => select(i));
    chip.addEventListener('keydown', (e) => {
      const step = e.key === 'ArrowRight' ? 1 : e.key === 'ArrowLeft' ? -1 : 0;
      if (step) {
        e.preventDefault();
        select((i + step + chips.length) % chips.length, true);
      }
    });
  });
  select(0);
}
