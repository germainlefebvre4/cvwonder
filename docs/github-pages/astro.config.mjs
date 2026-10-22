// @ts-check
import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';
import starlightLinksValidator from 'starlight-links-validator';
import mermaid from 'astro-mermaid';
import { fileURLToPath } from 'node:url';
import { redirects } from './redirects.ts';

// Exposed to the header component; `dev` when no release version is injected.
process.env.CVWONDER_VERSION ||= 'dev';

export default defineConfig({
  site: 'https://cvwonder.fr',
  trailingSlash: 'always',
  redirects,
  vite: {
    resolve: { alias: { '~': fileURLToPath(new URL('./src', import.meta.url)) } },
    define: { 'import.meta.env.CVWONDER_VERSION': JSON.stringify(process.env.CVWONDER_VERSION) },
  },
  integrations: [
    mermaid({ theme: 'default', autoTheme: true }),
    starlight({
      title: 'CV Wonder',
      description: 'Generate beautiful, professional CVs from YAML in seconds with CV Wonder.',
      favicon: '/favicon.svg',
      logo: {
        light: './public/img/logo-light.svg',
        dark: './public/img/logo-dark.svg',
      },
      defaultLocale: 'root',
      locales: {
        root: { label: 'English', lang: 'en' },
        fr: { label: 'Français', lang: 'fr' },
      },
      head: [
        { tag: 'link', attrs: { rel: 'preload', href: '/fonts/inter-latin.woff2', as: 'font', type: 'font/woff2', crossorigin: '' } },
        { tag: 'link', attrs: { rel: 'icon', href: '/favicon.ico', sizes: '48x48' } },
        { tag: 'meta', attrs: { property: 'og:image', content: 'https://cvwonder.fr/share-card.png' } },
        { tag: 'meta', attrs: { property: 'og:image:width', content: '1200' } },
        { tag: 'meta', attrs: { property: 'og:image:height', content: '630' } },
        { tag: 'meta', attrs: { name: 'twitter:card', content: 'summary_large_image' } },
        { tag: 'meta', attrs: { name: 'twitter:image', content: 'https://cvwonder.fr/share-card.png' } },
      ],
      social: [
        { icon: 'github', label: 'GitHub', href: 'https://github.com/germainlefebvre4/cvwonder' },
      ],
      editLink: {
        baseUrl: 'https://github.com/germainlefebvre4/cvwonder/edit/main/docs/github-pages/',
      },
      sidebar: [
        { label: 'Start', translations: { fr: 'Démarrer' }, items: [{ autogenerate: { directory: 'start' } }] },
        { label: 'Guides', translations: { fr: 'Guides' }, items: [{ autogenerate: { directory: 'guides' } }] },
        { label: 'Themes', translations: { fr: 'Thèmes' }, items: [{ autogenerate: { directory: 'themes' } }] },
        { label: 'Reference', translations: { fr: 'Référence' }, items: [{ autogenerate: { directory: 'reference' } }] },
        { label: 'Contribute', translations: { fr: 'Contribuer' }, items: [{ autogenerate: { directory: 'contribute' } }] },
      ],
      components: {
        SocialIcons: './src/components/SocialIcons.astro',
        Footer: './src/components/Footer.astro',
      },
      expressiveCode: { defaultProps: { frame: 'code' } },
      customCss: [
        './src/styles/theme.css',
      ],
      plugins: [starlightLinksValidator({ errorOnFallbackPages: false })],
    }),
  ],
});
