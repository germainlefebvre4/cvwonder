// Old URL (under /docs/) -> new URL. One entry per page listed in
// redirects.fixture.json, plus the /docs/ root. GitHub Pages cannot send a 301,
// so Astro emits meta-refresh pages for these.
export const redirects: Record<string, string> = {
  '/docs/': '/start/',

  // Start
  '/docs/installation': '/start/install/',
  '/docs/installation/binary': '/start/install/',
  '/docs/installation/docker': '/start/install/',
  '/docs/getting-started': '/start/first-cv/',
  '/docs/getting-started/generate-your-cv': '/start/first-cv/',

  // Guides
  '/docs/getting-started/write-your-cv': '/guides/write-your-cv/',
  '/docs/guides': '/guides/',
  '/docs/guides/bulk-mode': '/guides/bulk-mode/',
  '/docs/guides/theme-customization': '/guides/theme-customization/',
  '/docs/validation': '/guides/validate-your-cv/',
  '/docs/validation/validate': '/guides/validate-your-cv/',
  '/docs/export': '/guides/export-pdf-and-html/',
  '/docs/export/html': '/guides/export-pdf-and-html/',
  '/docs/export/pdf': '/guides/export-pdf-and-html/',
  '/docs/export/missing-format': '/guides/export-pdf-and-html/',

  // Themes
  '/docs/themes': '/themes/',
  '/docs/themes/library': '/themes/gallery/',
  '/docs/themes/use-a-theme': '/themes/use/',
  '/docs/themes/install-remote-theme': '/themes/install-remote/',
  '/docs/themes/write-your-theme': '/themes/write/',
  '/docs/themes/publish-your-theme': '/themes/publish/',
  '/docs/themes/troubleshooting': '/themes/troubleshooting/',

  // Reference
  '/docs/themes/theme-functions': '/reference/theme-functions/',
  '/docs/cli': '/reference/cli/',
  '/docs/cli/commands': '/reference/cli/commands/',
  '/docs/cli/options': '/reference/cli/options/',
  '/docs/validation/schema': '/reference/yaml-schema/',
  '/docs/validation/schema-reference': '/reference/yaml-schema/',
  '/docs/validation/common-errors': '/reference/common-errors/',

  // Contribute
  '/docs/contributing': '/contribute/',
  '/docs/contributing/architecture': '/contribute/architecture/',
  '/docs/contributing/how-to-contribute': '/contribute/',
  '/docs/contributing/how-to-contribute/issues': '/contribute/issues/',
  '/docs/contributing/how-to-contribute/pull-requests': '/contribute/pull-requests/',
  '/docs/contributing/how-to-contribute/development': '/contribute/development/',
  '/docs/contributing/how-to-contribute/code-of-conduct': '/contribute/code-of-conduct/',
  '/docs/contributing/maintainer': '/contribute/',
  '/docs/contributing/maintainer/git-workflow': '/contribute/workflow/',
  '/docs/contributing/maintainer/release-flow': '/contribute/release/',

  // Footer
  '/docs/license': '/license/',
};
