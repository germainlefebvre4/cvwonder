# Spec Delta

## MODIFIED Requirements

### Requirement: Single accent color with accessible contrast
The site SHALL use a single accent color, vermilion, `#c2410c` in light mode and `#fb923c` in dark mode, for text, links and buttons, and every use of the accent for text, links or button backgrounds SHALL meet a contrast ratio of at least 4.5:1 against its background. Purely decorative, non-text effects, such as the landing page background, MAY use additional hues. Any text painted with a gradient SHALL meet the same ratio at every point of the gradient.

#### Scenario: Light mode contrast
- **WHEN** the accent is used for text or links on the light background
- **THEN** its contrast ratio against that background SHALL be at least 4.5:1

#### Scenario: Dark mode contrast
- **WHEN** the accent is used for text or links on the dark background
- **THEN** its contrast ratio against that background SHALL be at least 4.5:1

#### Scenario: Warning remains distinct from the accent
- **WHEN** a caution admonition and a note admonition are displayed together
- **THEN** they SHALL be visually distinguishable in both light and dark modes

#### Scenario: Decorative hues stay out of text and controls
- **WHEN** the landing page background uses hues other than the accent
- **THEN** no text, link or button SHALL use those hues as its color

#### Scenario: Gradient text contrast
- **WHEN** a word is painted with a color gradient on the landing page
- **THEN** every color of the gradient SHALL have a contrast ratio of at least 4.5:1 against the background in both light and dark modes
