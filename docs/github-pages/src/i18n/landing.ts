// Labels of the landing page, one dictionary per locale. Components read them
// through `landingText(Astro.locals.starlightRoute.lang)`; MDX pages only pass content.
export interface LandingText {
  terminal: string;
  heroCvAlt: string;
  seeHow: string;
  scene: {
    label: string;
    yamlFile: string;
    cvFile: string;
    acts: [string, string, string];
    actTitles: [string, string, string];
    actHints: [string, string, string];
    chips: string;
    tabs: string;
    preview: string;
    htmlBadge: string;
    pdfBadge: string;
  };
  rail: { label: string; intro: string; scene: string; why: string; start: string; getStarted: string };
  skip: string;
}

const en: LandingText = {
  terminal: 'Terminal',
  heroCvAlt: 'A CV generated with the default theme',
  seeHow: 'See how it works',
  scene: {
    label: 'From YAML to CV',
    yamlFile: 'cv.yml',
    cvFile: 'cv.html',
    acts: ['Write', 'Generate', 'Theme'],
    actTitles: ['Write your CV in YAML', 'Generate it in seconds', 'Choose a theme'],
    actHints: ['One readable file holds everything.', 'HTML and PDF, rendered from the same data.', 'Same data, a different look. One line changes.'],
    chips: 'Themes',
    tabs: 'Acts of the demo',
    preview: 'Theme preview',
    htmlBadge: 'cv.html',
    pdfBadge: 'cv.pdf',
  },
  rail: {
    label: 'Page sections',
    intro: 'Introduction',
    scene: 'How it works',
    why: 'Why CV Wonder',
    start: 'Where to start',
    getStarted: 'Get started',
  },
  skip: 'Skip the demo',
};

const fr: LandingText = {
  terminal: 'Terminal',
  heroCvAlt: 'Un CV généré avec le thème default',
  seeHow: 'Voir comment ça marche',
  scene: {
    label: 'Du YAML au CV',
    yamlFile: 'cv.yml',
    cvFile: 'cv.html',
    acts: ['Écrire', 'Générer', 'Thème'],
    actTitles: ['Écrivez votre CV en YAML', 'Générez-le en quelques secondes', 'Choisissez un thème'],
    actHints: ['Un seul fichier lisible contient tout.', 'HTML et PDF, générés depuis les mêmes données.', 'Mêmes données, un autre style. Une seule ligne change.'],
    chips: 'Thèmes',
    tabs: 'Étapes de la démo',
    preview: 'Aperçu du thème',
    htmlBadge: 'cv.html',
    pdfBadge: 'cv.pdf',
  },
  rail: {
    label: 'Sections de la page',
    intro: 'Introduction',
    scene: 'Comment ça marche',
    why: 'Pourquoi CV Wonder',
    start: 'Par où commencer',
    getStarted: 'Commencer',
  },
  skip: 'Passer la démo',
};

export function landingText(lang: string | undefined): LandingText {
  return lang?.startsWith('fr') ? fr : en;
}
