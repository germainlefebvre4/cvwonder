import {themes as prismThemes} from 'prism-react-renderer';
import type {Config} from '@docusaurus/types';
import type * as Preset from '@docusaurus/preset-classic';

const config: Config = {
  title: 'CV Wonder',
  tagline: 'Generate your CV with Wonder!',
  favicon: 'img/favicon.ico',

  url: 'https://cvwonder.fr',
  baseUrl: '/',

  organizationName: 'germainlefebvre4',
  projectName: 'cvwonder',
  onBrokenLinks: 'throw',
  onBrokenMarkdownLinks: 'warn',

  i18n: {
    defaultLocale: 'en',
    locales: ['en'],
  },

  presets: [
    [
      'classic',
      {
        docs: {
          sidebarPath: './sidebars.ts',
          editUrl:
            'https://github.com/germainlefebvre4/cvwonder/tree/main/docs/github-pages/',
        },
        theme: {
          customCss: './src/css/custom.css',
        },
      } satisfies Preset.Options,
    ],
  ],

  themeConfig: {
    image: 'img/social-card.png',
    metadata: [
      {name: 'keywords', content: 'cv, resume, generator, yaml, themes, professional'},
      {name: 'description', content: 'Generate beautiful, professional CVs from YAML in seconds with CV Wonder'},
      {name: 'og:type', content: 'website'},
      {name: 'og:title', content: 'CV Wonder | Professional CV Generator'},
      {name: 'og:description', content: 'Generate beautiful, professional CVs from YAML in seconds with CV Wonder'},
    ],
    navbar: {
      title: 'CV Wonder',
      logo: {
        alt: 'CV Wonder',
        src: 'img/logo.svg',
      },
      items: [
        // {
        //   type: 'docSidebar',
        //   sidebarId: 'documentationSidebar',
        //   position: 'left',
        //   label: 'Documentation',
        // },
        {
          to: '/docs/getting-started',
          label: 'Getting Started',
          position: 'left',
        },
        {
          to: '/docs/cli',
          label: 'CLI',
          position: 'left',
        },
        {
          to: '/docs/themes',
          label: 'Themes',
          position: 'left',
        },
        {
          to: '/docs/export',
          label: 'Export',
          position: 'left',
        },
        {
          href: 'https://github.com/germainlefebvre4/cvwonder',
          label: 'GitHub',
          position: 'right',
        },
      ],
    },
    footer: {
      style: 'dark',
      links: [
        {
          title: 'Docs',
          items: [
            {
              label: 'Getting Started',
              to: '/docs/getting-started',
            },
            {
              label: 'CV Format',
              to: '/docs/export',
            },
          ],
        },
        {
          title: 'Themes',
          items: [
            {
              label: 'Theme Gallery',
              to: '/docs/themes',
            },
            {
              label: 'Creating Themes',
              to: '/docs/themes/write-your-theme',
            },
          ],
        },
        {
          title: 'Resources',
          items: [
            {
              label: 'GitHub',
              href: 'https://github.com/germainlefebvre4/cvwonder',
            },
            {
              label: 'Releases',
              href: 'https://github.com/germainlefebvre4/cvwonder/releases',
            },
            {
              label: 'Report Issues',
              href: 'https://github.com/germainlefebvre4/cvwonder/issues',
            },
          ],
        },
      ],
      copyright: `Copyright © ${new Date().getFullYear()} Germain LEFEBVRE. Built with Docusaurus.`,
    },
    prism: {
      theme: prismThemes.github,
      darkTheme: prismThemes.dracula,
      additionalLanguages: ['bash', 'yaml', 'go'],
    },
    colorMode: {
      defaultMode: 'light',
      disableSwitch: false,
      respectPrefersColorScheme: true,
    },
  } satisfies Preset.ThemeConfig,
  markdown: {
    mermaid: true,
  },
  themes: ['@docusaurus/theme-mermaid'],
};

export default config;
