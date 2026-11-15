# Visual language — operational reference

Full evidence and measurements: `docs/design/nextdoor-visual-language.md`.

## Colour

| Token | Value | Use |
| --- | --- | --- |
| `--ui-canvas` | `#F7F8F7` | page background |
| `--ui-surface` | `#FFFFFF` | cards, panels, sheets |
| `--ui-ink` | `#222222` | headings and body |
| `--ui-muted` | `#626B66` | metadata, secondary copy |
| `--ui-border` | `#E3E7E4` | cards, dividers, inputs |
| `--ui-primary` | `#1F6F4A` | primary action, active nav, selection |
| `--ui-primary-strong` | `#17553A` | pressed/hover primary |
| `--ui-strong` | `#1F2A44` | dark secondary CTA (invite/send) |
| `--ui-danger` | `#B42318` | destructive actions and errors |
| `--ui-focus` | `#1F6F4A` | 3px focus ring, 2px offset |

iOS estimates for reference: canvas `#F2F3F4` light / `#0F1011` dark; surface `#FFFFFF` /
`#1C1D1F`; ink `#202124` / `#F3F4F5`; muted `#6B7280` / `#A1A6AD`; hairline `#E5E7E8` /
`#2A2C2E`; social green `#1B6B47`; operational navy `#1F2B47`; danger `#D92D20`; toast
`#E7E9F5`; illustration cream `#F6F0E2`; publisher check `#1D6FE0`. Use the maulid.id tokens
above, not these raw values.

Rules:

- Green is a role, not decoration. Do not use it for every icon or heading.
- Nextdoor uses two primaries: green for the social commit (`Post`, `Continue`, FAB) and
  navy for the operational commit (`Join`, `Interested`, `Save`, `Next`). maulid.id maps
  green to publishing/positive and `--ui-strong` to operational confirmation.
- Do not reuse Nextdoor's exact brand green; `#1F6F4A` is maulid.id's own.
- Never rely on colour alone to communicate validation, error or selection.
- Verify AA contrast before shipping. If the light green fails for text, use
  `--ui-primary-strong` for text and reserve `--ui-primary` for fills.

## Type

| Role | Size / weight | Face |
| --- | --- | --- |
| Editorial display | 32–64px / 400 | EB Garamond |
| Section title | 20–28px / 400–500 | EB Garamond or UI sans |
| Card name | 14px / 600 | UI sans |
| Body | 15–16px / 400, 1.5–1.6 | UI sans or EB Garamond for long form |
| Metadata | 12–13px / 400 | UI sans |
| Navigation | 15px / 500 | UI sans |
| Button | 14px / 600 | UI sans |
| Kicker | 11–12px / 600, small caps | UI sans |

```css
--font-ui: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto,
  "Helvetica Neue", Arial, sans-serif;
--font-display: "EB Garamond", Georgia, "Times New Roman", serif;
```

- Keep the existing EB Garamond editorial identity.
- Use the UI sans for navigation, controls, metadata, forms and the backoffice.
- Do not load a new font without a recorded decision.
- Use tabular numerals for timestamps and counts.

## Spacing and shape

```text
spacing: 4 / 8 / 12 / 16 / 20 / 24 / 32 / 40 / 48 / 64
card radius: 12px
control radius: 999px
sheet radius: 16px top
avatar: 50%
card gap: 12px
control height: 40–44px
primary button height: 44–48px
```

- Card padding 16px mobile / 20px desktop.
- Page padding 16px mobile / 24px tablet / 32px desktop.
- Use borders and canvas contrast; avoid shadow on every card.
- Floating layers only: `0 8px 24px rgb(0 0 0 / 0.08)`.

## Icons and imagery

- Outline icons, 20–24px, consistent stroke; label navigation items.
- Icon-only actions need accessible names.
- Use honest Indonesian mosque/community photography with credit.
- Fallback: monogram tile, not generic stock.
- No Middle-Eastern silhouette/lantern clip-art as a stand-in for Indonesian Islam.
- Image ratios: feed 4:3 or 1:1; event thumb 1:1; avatar circle.

## Motion

- 150–200ms controls; 250–300ms sheets; ease-out.
- No bounce, parallax, autoplay or countdown animation.
- Respect `prefers-reduced-motion`.
- Animate only transform/opacity; never animate layout properties for decoration.

## Density

- One image post per 300–400px of mobile feed.
- 2–3 text-only posts in the same space.
- Metadata never below 12px equivalent.
- Headline max two lines on cards; do not truncate before the dek.
- Right-rail cards: one idea, one action, no nested cards.

## Dark mode

- Optional, not required for the first implementation.
- If added, define separate canvas/surface/border/ink/muted tokens; do not invert the light
  palette mechanically.
- Preserve the green active/FAB role and verify contrast in both themes.
- Do not let images or ads become the only source of contrast.
