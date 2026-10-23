# Spec Delta

## MODIFIED Requirements

### Requirement: Scroll scene shows YAML becoming a CV
The landing page SHALL include a scene, driven by the visitor's scrolling, in three acts: writing the YAML, generating the CV, and choosing a theme. The scene SHALL stay on screen while it plays and the page SHALL move on to the next section once it ends. The YAML SHALL remain readable throughout, so that the CV is visibly rendered from it.

#### Scenario: Act one, write
- **WHEN** the scene starts
- **THEN** the lines of an example YAML file SHALL light up one after another, and all of them SHALL be lit before the second act begins

#### Scenario: Act two, generate
- **WHEN** the visitor scrolls into the second act
- **THEN** the YAML lines SHALL stay in place, take a resting tint, and the `theme:` line SHALL be visually emphasized; the CV SHALL render from top to bottom under a visible scan line, followed by indicators of the HTML and PDF outputs

#### Scenario: YAML stays visible once the CV is rendered
- **WHEN** the CV has finished rendering
- **THEN** every line of the YAML SHALL still be visible and readable, distinguished from the CV by its tint, and no line SHALL be hidden

#### Scenario: Act three, choose a theme
- **WHEN** the visitor scrolls into the third act
- **THEN** the CV SHALL be re-rendered in the other themes in turn through a visible wipe, the theme line of the YAML SHALL show the current theme name and keep its emphasis, and the matching theme chip SHALL be highlighted

#### Scenario: Theme name change is noticeable
- **WHEN** the theme shown in the CV changes during the third act
- **THEN** the theme name on the `theme:` line SHALL change with a short highlight, so the visitor sees which line drives the render

#### Scenario: YAML example size
- **WHEN** the scene is displayed in a viewport at least 1280 pixels wide and 650 pixels tall
- **THEN** the YAML window SHALL show an example of about 18 lines covering several sections of a CV (for example person, abstract, skills, career, education, languages) and the `theme:` line, entirely visible without scrolling inside the window

#### Scenario: CV legibility
- **WHEN** the CV is displayed in the scene on a screen at least 1280 pixels wide
- **THEN** the CV preview SHALL be displayed at about its native screenshot width so that headings and section titles are readable without zooming, whatever the size of the YAML window

#### Scenario: Theme line emphasis without the pinned scene
- **WHEN** the scene is displayed without the pinned animation (narrow screen, reduced motion or no JavaScript)
- **THEN** the `theme:` line SHALL still be visually distinguished from the other lines of the YAML example

---

### Requirement: Sections are reachable without playing the scene
A visitor SHALL be able to reach any section of the landing page directly without scrolling through the scene, by a side rail of section links, a skip control, tabs for each act and a hero link to the scene. Navigation to a section SHALL be immediate.

#### Scenario: Jump from the hero to a later section
- **WHEN** a visitor activates the rail link for "Where to start" while at the top of the page
- **THEN** the page SHALL show that section immediately, without traversing the scene, with its content visible

#### Scenario: Skip while the scene plays
- **WHEN** the scene is on screen
- **THEN** a "Skip the demo" control SHALL be visible, and activating it SHALL move the visitor to the "Why CV Wonder" section

#### Scenario: Skip control hidden outside the scene
- **WHEN** the scene is not on screen
- **THEN** the skip control SHALL NOT be visible or focusable

#### Scenario: Jump between acts
- **WHEN** a visitor activates the tab of the first or second act while the scene is on screen
- **THEN** the scene SHALL show that act fully played, and the active tab SHALL follow the scroll position

#### Scenario: Jump to the theme act
- **WHEN** a visitor activates the tab of the theme act while the scene is on screen
- **THEN** the scene SHALL show the horizon-timeline theme in the CV, with its capture displayed from the top, the `theme:` line SHALL show `horizon-timeline`, and the theme tab SHALL be the active tab

#### Scenario: Hero link to the scene
- **WHEN** a visitor activates the "See how it works" link in the hero
- **THEN** the page SHALL show the beginning of the scene

#### Scenario: Deep link to a section
- **WHEN** a visitor opens the landing page with a section anchor in the URL
- **THEN** the page SHALL show that section with its content visible

#### Scenario: Return to the scene
- **WHEN** a visitor who is below the scene activates the rail link for the scene
- **THEN** the scene SHALL show its beginning
