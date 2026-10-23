# documentation-landing Specification

## Purpose

Defines the observable behavior of the animated landing page of the CV Wonder documentation site: its entrance animation, the scroll scene that shows a YAML file becoming a CV, how visitors move past that scene, and how the page degrades without animation.

## Requirements

### Requirement: Hero entrance animation
The landing page hero SHALL animate in on load: the title revealing word by word, the supporting text and actions fading in, then the terminal and the CV preview appearing. The hero text SHALL be readable without waiting for the animation.

#### Scenario: Animated entrance
- **WHEN** a visitor opens the landing page with animations enabled
- **THEN** the title, tagline, actions, terminal and CV preview SHALL appear in sequence

#### Scenario: Text is not held back
- **WHEN** the animation code has not loaded or has failed
- **THEN** the title, tagline and actions SHALL be fully visible

---

### Requirement: Hero terminal types the three commands
The hero terminal SHALL type the commands `cvwonder init`, `cvwonder serve` and `cvwonder generate` in that order, each followed by its output line.

#### Scenario: Commands typed in order
- **WHEN** the hero finishes appearing
- **THEN** the three commands SHALL be typed one after the other, each followed by its output

#### Scenario: Pointer interaction
- **WHEN** a visitor with a pointing device moves the pointer over the hero terminal area
- **THEN** the terminal and the CV preview SHALL tilt and shift subtly in response, and return to rest when the pointer leaves

---

### Requirement: Scroll scene shows YAML becoming a CV
The landing page SHALL include a scene, driven by the visitor's scrolling, in three acts: writing the YAML, generating the CV, and choosing a theme. The scene SHALL stay on screen while it plays and the page SHALL move on to the next section once it ends.

#### Scenario: Act one, write
- **WHEN** the scene starts
- **THEN** the lines of an example YAML file SHALL light up one after another

#### Scenario: Act two, generate
- **WHEN** the visitor scrolls into the second act
- **THEN** the YAML lines SHALL fly away except the theme line, and the CV SHALL render from top to bottom under a visible scan line, followed by indicators of the HTML and PDF outputs

#### Scenario: Act three, choose a theme
- **WHEN** the visitor scrolls into the third act
- **THEN** the CV SHALL be re-rendered in the other themes in turn through a visible wipe, the theme line of the YAML SHALL show the current theme name, and the matching theme chip SHALL be highlighted

#### Scenario: CV legibility
- **WHEN** the CV is displayed in the scene on a screen at least 1280 pixels wide
- **THEN** the CV preview SHALL be displayed at about its native screenshot width so that headings and section titles are readable without zooming

---

### Requirement: The scene follows native scrolling
The scene SHALL be a direct function of the page scroll position, in both directions. The page SHALL NOT capture the wheel, snap to positions, or lock scrolling while the scene plays.

#### Scenario: Scrolling back replays the scene in reverse
- **WHEN** the visitor scrolls up through the scene
- **THEN** the scene SHALL play backwards and return to its starting state

#### Scenario: Scrolling is never intercepted
- **WHEN** a visitor scrolls through the scene at any speed with wheel, keyboard, touch or scrollbar
- **THEN** the page SHALL scroll normally without being held, snapped or redirected

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
- **WHEN** a visitor activates the tab of an act while the scene is on screen
- **THEN** the scene SHALL show that act fully played, and the active tab SHALL follow the scroll position

#### Scenario: Hero link to the scene
- **WHEN** a visitor activates the "See how it works" link in the hero
- **THEN** the page SHALL show the beginning of the scene

#### Scenario: Deep link to a section
- **WHEN** a visitor opens the landing page with a section anchor in the URL
- **THEN** the page SHALL show that section with its content visible

#### Scenario: Return to the scene
- **WHEN** a visitor who is below the scene activates the rail link for the scene
- **THEN** the scene SHALL show its beginning

---

### Requirement: Section rail shows the current position
The landing page SHALL display a vertical rail of section markers for the intro, the scene, "Why CV Wonder", "Where to start" and the closing call to action. The marker of the section being viewed SHALL be highlighted and each marker SHALL have an accessible name.

#### Scenario: Current section highlighted
- **WHEN** the visitor scrolls from one section to the next
- **THEN** the highlighted marker SHALL change to the section being viewed

#### Scenario: Marker labels
- **WHEN** a visitor hovers or focuses a marker
- **THEN** the name of its section SHALL be displayed

#### Scenario: Rail on narrow screens
- **WHEN** the viewport is narrower than 1180 pixels
- **THEN** the rail MAY be hidden

---

### Requirement: Content stays complete without animation
When the visitor prefers reduced motion, or JavaScript is unavailable, the landing page SHALL show all its content statically: the hero, the YAML example, the CV, every section and the call to action. No content SHALL remain hidden waiting for an animation, and every theme preview SHALL be viewable.

#### Scenario: Reduced motion preference
- **WHEN** a visitor whose system requests reduced motion opens the landing page
- **THEN** no animation SHALL play and all sections SHALL be fully visible

#### Scenario: No JavaScript
- **WHEN** the landing page is displayed without JavaScript
- **THEN** the hero, the YAML example, the CV previews, the sections and the call to action SHALL all be visible

#### Scenario: Themes viewable without animation
- **WHEN** a visitor on the static rendering wants to see another theme
- **THEN** they SHALL be able to view each of the three theme previews

---

### Requirement: Small screens get a simplified scene
On viewports narrower than 900 pixels, the scene SHALL NOT stay on screen while scrolling. Its content SHALL be laid out in a single column, and the other animations MAY remain.

#### Scenario: Phone-sized viewport
- **WHEN** a visitor scrolls the landing page on a viewport narrower than 900 pixels
- **THEN** the page SHALL scroll continuously with no section held in place, and the YAML, the CV and the theme previews SHALL all be reachable

---

### Requirement: Landing page supports light and dark modes
The animated landing page SHALL render correctly in both light and dark modes, following the site's mode. The aurora background SHALL have a variant for each mode, and the CV frames SHALL stay legible in both.

#### Scenario: Dark mode
- **WHEN** the site is in dark mode
- **THEN** the landing page SHALL show the aurora on the dark background

#### Scenario: Light mode
- **WHEN** the site is in light mode
- **THEN** the landing page SHALL show a light aurora variant, with all text meeting the contrast requirement of the site

#### Scenario: Switching mode on the page
- **WHEN** a visitor switches mode while on the landing page
- **THEN** the landing page SHALL update without reloading

---

### Requirement: Landing page is available in English and French
All texts of the landing page, including the hero, the acts of the scene, the tab, chip, rail and skip labels, SHALL be available in English and French, following the language being browsed.

#### Scenario: French landing page
- **WHEN** a visitor opens the French landing page
- **THEN** every label of the scene, the rail and the skip control SHALL be in French

#### Scenario: Language-independent navigation
- **WHEN** a visitor uses the rail or the skip control on the French landing page
- **THEN** the visitor SHALL stay on the French page

---

### Requirement: Animations do not slow down the rest of the site
Animation code SHALL be delivered only with the landing page and SHALL add no more than 60 KB, gzip-compressed, of script to it. Other pages SHALL receive none of it. The landing page SHALL NOT shift its layout when animation starts.

#### Scenario: Other pages unaffected
- **WHEN** a documentation page other than the landing page is built
- **THEN** it SHALL NOT include the animation code

#### Scenario: Script budget
- **WHEN** the built landing page is inspected
- **THEN** the animation script it adds SHALL NOT exceed 60 KB, gzip-compressed

#### Scenario: No layout shift
- **WHEN** the hero animation and the scene start
- **THEN** surrounding content SHALL NOT move
