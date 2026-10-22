# Spec Delta

## Purpose

Defines the observable behavior of the CV Wonder documentation website (cvwonder.fr): URL structure, languages, navigation, search, visual identity, link integrity and publication. It is independent of the framework that builds the site.

## ADDED Requirements

### Requirement: Documentation is served without a /docs prefix
The site SHALL serve English content from the site root and French content under `/fr/`, with no `/docs/` path segment.

#### Scenario: English page URL
- **WHEN** a reader opens the installation page in English
- **THEN** its URL SHALL be of the form `https://cvwonder.fr/start/install/` with no `/docs/` segment

#### Scenario: French page URL
- **WHEN** a reader opens the French version of a translated page
- **THEN** its URL SHALL be the English URL prefixed with `/fr/`

---

### Requirement: Legacy documentation URLs redirect
Every page published under the previous `/docs/...` URL structure SHALL redirect to its new equivalent, or to the closest equivalent when pages were merged.

#### Scenario: Deep link from the README keeps working
- **WHEN** a reader follows `https://cvwonder.fr/docs/themes/install-remote-theme`
- **THEN** the reader SHALL land on the page that now documents remote theme installation

#### Scenario: Merged page redirects to its successor
- **WHEN** a reader follows an old URL whose page was folded into another page
- **THEN** the reader SHALL land on the page that now contains that content

#### Scenario: Redirect coverage is complete
- **WHEN** the redirect table is checked against the list of previously published pages
- **THEN** every old page URL SHALL have a redirect target that exists in the built site

---

### Requirement: English is the default language and French is available
The site SHALL offer English as the default language and French as a second language, with a language switcher visible on every page.

#### Scenario: Switching language on a translated page
- **WHEN** a reader switches from English to French on a page that has a French translation
- **THEN** the reader SHALL land on the French version of that same page

#### Scenario: Language switcher on every page
- **WHEN** any page is displayed
- **THEN** a language switcher offering English and French SHALL be visible

---

### Requirement: Untranslated French pages fall back to English
When a page has no French translation, the French site SHALL show the English content with a visible notice that the page is not yet translated, instead of a missing page.

#### Scenario: Untranslated page under /fr/
- **WHEN** a reader opens the French URL of a page that has no translation
- **THEN** the English content SHALL be displayed together with a notice saying the page is not available in French

#### Scenario: French navigation is complete
- **WHEN** a reader browses the French site
- **THEN** every navigation entry of the English site SHALL be present and SHALL lead to a page, translated or in fallback

---

### Requirement: Launch scope of the French translation
At launch, the French site SHALL provide translated content for the homepage, the navigation labels, the installation pages, the first-CV walkthrough and the guides section.

#### Scenario: Core path fully in French
- **WHEN** a French-speaking reader goes from the homepage to installation, then to the first-CV walkthrough
- **THEN** every page on that path SHALL be displayed in French with no untranslated-page notice

#### Scenario: Technical reference may remain in English
- **WHEN** a reader opens a page outside the launch scope, such as the theme functions reference
- **THEN** the page MAY be displayed in English under the fallback notice

---

### Requirement: Navigation is organised by reader intention
The primary navigation SHALL group content into Start, Guides, Themes, Reference and Contribute sections. The license SHALL be reachable from the footer.

#### Scenario: Primary navigation entries
- **WHEN** any page is displayed
- **THEN** the primary navigation SHALL list Start, Guides, Themes, Reference and Contribute

#### Scenario: Section content is auto-discovered
- **WHEN** a page is added to a section directory with its ordering metadata
- **THEN** it SHALL appear in that section's sidebar without a manual sidebar entry

#### Scenario: Former export and validation content remains reachable
- **WHEN** a reader looks for export formats or YAML validation
- **THEN** that content SHALL be found under Guides or Reference

---

### Requirement: Homepage presents the audiences and the product
The homepage SHALL present a hero with a short terminal demonstration of `cvwonder init`, `serve` and `generate`, a primary call to action toward the first-CV walkthrough, and three entry points: starting a first CV, generating CVs in bulk, and creating a theme.

#### Scenario: Primary call to action
- **WHEN** a reader opens the homepage
- **THEN** a prominent action SHALL lead to the first-CV walkthrough

#### Scenario: Three audience entry points
- **WHEN** a reader scrolls the homepage
- **THEN** three entry points SHALL be shown, leading respectively to the first-CV walkthrough, to the bulk mode guide, and to writing a theme

#### Scenario: Theme screenshots shown as previews
- **WHEN** the homepage shows theme previews
- **THEN** screenshots SHALL be cropped to a consistent preview area rather than displayed at full page height

---

### Requirement: Header shows the release version
The header SHALL display the documented release version and link to the corresponding GitHub release, with `dev` displayed when no version is provided at build time.

#### Scenario: Version supplied at build
- **WHEN** the site is built with a version value
- **THEN** the header SHALL show that version, linking to its GitHub release page

#### Scenario: No version supplied
- **WHEN** the site is built without a version value
- **THEN** the header SHALL show `dev`

---

### Requirement: Full-text search per language
The site SHALL provide full-text search over page content, restricted to the language currently being browsed, available from every page.

#### Scenario: Searching in English
- **WHEN** a reader searches for a term while browsing the English site
- **THEN** matching English pages SHALL be listed

#### Scenario: Searching in French
- **WHEN** a reader searches for a term while browsing the French site
- **THEN** matching French pages SHALL be listed, including pages shown in fallback

---

### Requirement: Broken internal links fail the build
The site build SHALL fail when an internal link or an internal anchor does not resolve to a page of the built site.

#### Scenario: Link to a removed page
- **WHEN** a page links to an internal page that does not exist
- **THEN** the build SHALL fail and report the page and the broken target

#### Scenario: Valid links
- **WHEN** all internal links resolve
- **THEN** the build SHALL succeed

---

### Requirement: Single accent color with accessible contrast
The site SHALL use a single accent color, vermilion, `#c2410c` in light mode and `#fb923c` in dark mode, and every use of the accent for text, links or button backgrounds SHALL meet a contrast ratio of at least 4.5:1 against its background.

#### Scenario: Light mode contrast
- **WHEN** the accent is used for text or links on the light background
- **THEN** its contrast ratio against that background SHALL be at least 4.5:1

#### Scenario: Dark mode contrast
- **WHEN** the accent is used for text or links on the dark background
- **THEN** its contrast ratio against that background SHALL be at least 4.5:1

#### Scenario: Warning remains distinct from the accent
- **WHEN** a caution admonition and a note admonition are displayed together
- **THEN** they SHALL be visually distinguishable in both light and dark modes

---

### Requirement: Light and dark modes with a visible logo
The site SHALL support light and dark modes, following the reader's system preference by default and allowing a manual switch. The logo SHALL remain visible in both modes.

#### Scenario: System preference respected
- **WHEN** a reader with a dark system preference opens the site for the first time
- **THEN** the site SHALL be displayed in dark mode

#### Scenario: Logo in dark mode
- **WHEN** the site is displayed in dark mode
- **THEN** the logo SHALL be rendered in a color that contrasts with the dark background

---

### Requirement: Social share metadata
Every page SHALL declare a title, a description and a social share image, so that links to the site display a preview card.

#### Scenario: Share image exists
- **WHEN** the built site is inspected
- **THEN** the share image referenced by the page metadata SHALL exist in the output

---

### Requirement: CLI reference reflects the current binary
The CLI reference SHALL display the help output of the `cvwonder` binary as generated from the source, so that it cannot drift from the CLI.

#### Scenario: Help output shown verbatim
- **WHEN** a reader opens the CLI reference
- **THEN** the page SHALL show the current `cvwonder --help` output

---

### Requirement: Site is published as static files on release
The documentation SHALL build to static files and be published to cvwonder.fr on each release tag, with the release version injected at build time.

#### Scenario: Release publishes the site
- **WHEN** a release tag is pushed
- **THEN** the site SHALL be built with that tag as its version and published to cvwonder.fr
