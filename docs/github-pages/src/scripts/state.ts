// State shared between the light navigation code and the lazily loaded animation chunk.
export interface Pin {
  start: number;
  end: number;
}

export const landing: {
  /** Scroll range of the pinned scene; null when the scene is not pinned. */
  pin: Pin | null;
  /** Scrolls the page to a position inside the scene, expressed as 0..100 of its timeline. */
  sceneAt: ((moment: number) => void) | null;
  /** Finishes every reveal that is already in view (used after instant jumps). */
  flushReveals: (() => void) | null;
} = { pin: null, sceneAt: null, flushReveals: null };

export const LAYOUT_EVENT = 'cvw:layout';
