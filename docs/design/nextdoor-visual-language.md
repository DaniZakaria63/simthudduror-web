# Nextdoor visual language — observed reference and maulid.id adaptation

**Evidence:** Mobbin web capture `https://mobbin.com/screens/03d5f316-fcd3-4389-9cc0-a5406574beca`
(1920×1320) and the additional screens cited in the skill references. Values marked
**observed** were sampled or measured from the capture; values marked **estimated** are
read from the image because Mobbin is not a token export. Do not present estimates as
Nextdoor's official design tokens.

## 1. Colour roles

| Role | Nextdoor observation | maulid.id adaptation |
| --- | --- | --- |
| Canvas | Light grey, sampled `#F8F9FA` | `--ui-canvas: #F7F8F7` |
| Surface | White cards and panels, `#FFFFFF` | `--ui-surface: #FFFFFF` |
| Ink | Charcoal body/headline text, sampled `#2F2F2F` | `--ui-ink: #222222` |
| Muted | Grey metadata, estimated `#6B6B6B` | `--ui-muted: #626B66` |
| Border | Hairline grey, sampled `#EBEDED` | `--ui-border: #E3E7E4` |
| Primary | Deep community green, sampled `#188050` | `--ui-primary: #1F6F4A` |
| Primary strong | Darker green for pressed/hover, sampled `#107040` | `--ui-primary-strong: #17553A` |
| Strong CTA | Dark navy/charcoal button, sampled `#222E46` | `--ui-strong: #1F2A44` |
| Danger | Not visible in the sampled capture | `--ui-danger: #B42318` |
| Focus | Assumed browser/default or green ring | `--ui-focus: #1F6F4A` at 3px |

Rules:

- Green is a **role**, not a decoration. It marks the primary action, the active
  navigation item, the selected tab or the verified/positive state.
- Nextdoor runs **two primaries with two jobs**: green for the social commit
  (`Post`, `Continue`, FAB), navy for the operational commit (`Join`, `Interested`, `Save`,
  `Next`). maulid.id maps green to publishing/positive and `--ui-strong` to operational
  confirmation. Destructive actions use `--ui-danger`.
- Do not reuse Nextdoor's exact brand green. maulid.id uses its own deeper green so the
  product does not read as a Nextdoor clone.
- The canvas is never pure white when cards are white; the grey canvas is what separates
  the feed from the page.
- Borders carry structure. Avoid stacking drop shadows on every card.
- Red is reserved for destructive actions and errors, not for Maulid or seasonal accents.

### iOS observed ranges (estimated)

| Role | Light | Dark |
| --- | --- | --- |
| Canvas | `#F2F3F4` | `#0F1011` |
| Surface | `#FFFFFF` | `#1C1D1F` |
| Ink | `#202124` | `#F3F4F5` |
| Muted | `#6B7280` | `#A1A6AD` |
| Hairline | `#E5E7E8` | `#2A2C2E` |
| Social green | `#1B6B47` | same |
| Operational navy | `#1F2B47` | same |
| Danger | `#D92D20` | same |
| Toast | `#E7E9F5` | `#2A2D3A` |
| Illustration cream | `#F6F0E2` | n/a |
| Publisher check | `#1D6FE0` | same |

These are ranges read from captures, not official tokens. Use the maulid.id `--ui-*` mapping
above, not the raw Nextdoor values.

## 2. Typography

Nextdoor's UI is sans-serif throughout; its hierarchy comes from size and weight, not from
a second typeface.

| Role | Nextdoor observation (estimated) | maulid.id adaptation |
| --- | --- | --- |
| Wordmark | 24–26px, heavy, tight | Keep EB Garamond wordmark for maulid.id identity |
| Page/section title | 20–24px, 600 | EB Garamond display for editorial sections |
| Card name | 14px, 600 | `--font-ui` 14px / 600 |
| Body | 15–16px, 400, 1.4–1.5 line height | `--font-ui` 15–16px / 1.55 |
| Metadata | 12–13px, 400, muted | `--font-ui` 12–13px / 400 |
| Navigation | 15–16px, 500 | `--font-ui` 15px / 500 |
| Button | 14–15px, 600 | `--font-ui` 14px / 600 |
| Composer placeholder | 15–16px, muted | `--font-ui` 15px / 400 |

Adaptation decision:

- Keep **EB Garamond** for the public editorial voice (masthead, section display
  headings, long-form article titles).
- Introduce a **UI sans** for navigation, buttons, metadata, forms, cards and the
  backoffice. Use a system stack first to avoid adding a font asset:
  `-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif`.
- Do not set long body copy in the UI sans if the existing editorial body voice is working;
  choose per surface and record the decision.
- Never use all-caps for long labels. Small caps only for short kickers.
- Numerals in timestamps and counts use tabular figures where alignment matters.

## 3. Spacing and layout rhythm

**Observed web proportions** (1920px capture, measured):

- Left navigation: ~17% of the viewport (~243px at 1440).
- Central reading column: ~38–39% (~553px at 1440).
- Right utility rail: ~19% (~271px at 1440).
- The column-to-rail ratio is about 2:1; the remaining width is margin, not more content.
- The feed is the widest content zone; the rail is visibly quieter.
- Cards are separated by ~12px and have 20px padding; the column does **not** widen to
  fill the viewport when the rail is absent.
- The active navigation item is a light grey fill with a filled icon; it is not a green
  block.

**Spacing scale for maulid.id:**

```text
4 / 8 / 12 / 16 / 20 / 24 / 32 / 40 / 48 / 64
```

- Card padding: 16px mobile, 20px desktop.
- Gap between feed cards: 12px.
- Gap between sections: 32–48px.
- Sidebar item height: 40–44px; horizontal padding 12–16px.
- Form control height: 40–44px; primary button 44–48px.
- Right-rail card padding: 16px; gap 12px.
- Page horizontal padding: 16px mobile, 24px tablet, 32px desktop.

## 4. Shape, borders and elevation

| Element | Nextdoor observation | maulid.id rule |
| --- | --- | --- |
| Feed card | ~12px radius, white, hairline border | `--radius-card: 12px` |
| Right-rail card | ~12px radius, white, hairline border | 12px |
| Primary button | Pill / fully rounded | `--radius-control: 999px` |
| Search field | Pill | 999px |
| Composer field | Rounded rectangle, ~12px | 12px |
| Avatar | Circle | 50% |
| Media | Top corners follow card radius | inherit card radius |
| Sheet/modal | ~16–20px top radius | 16px |
| Shadow | Subtle, only for floating layers | `0 8px 24px rgb(0 0 0 / 0.08)` |

- Do not put a shadow on every feed card. Use the canvas/border contrast.
- Hover on a card may change the border or background, not lift the card dramatically.
- Focus rings use the primary colour with a 3px offset and must remain visible on white
  and on photos.

## 5. Iconography

- Outline icons, 20–24px, consistent stroke (approximately 1.75–2px).
- Navigation always pairs icon with label. Icon-only actions (overflow, close, share) have
  accessible names.
- Icons are not decorative confetti; each one maps to an action.
- Use the project's existing icon source or inline SVG components. Do not import a new
  icon library for a single screen without a recorded decision.

## 6. Imagery

- Nextdoor relies on member-uploaded photography: houses, pets, streets, people, events.
  Images are ordinary and local, not art-directed stock.
- Feed images are full-width, roughly 4:3 or 1:1, cropped to the card width.
- Event thumbnails are square or 4:3 and small.
- Avatars are circular and use a fallback initial.
- For maulid.id: use real Indonesian mosque/community photographs with credit; a
  deterministic monogram fallback is better than a generic stock image.
- Do not use generic Middle-Eastern mosque silhouettes, lanterns or "Arabian nights"
  ornament as a stand-in for Indonesian Islam.

## 7. Motion

- Nextdoor's visible motion is restrained: hover/pressed states, sheets, small transitions.
- Duration: 150–200ms for controls, 250–300ms for sheets.
- Easing: standard ease-out; no bounce, no spring, no parallax.
- Respect `prefers-reduced-motion`: remove transform-based entrance animations and smooth
  scrolling.
- No autoplay video. No infinite marquee. No countdown animation.

## 8. Density and scannability

- One feed item per ~300–400px of vertical space on mobile when it has an image; 2–3
  text-only items in the same space.
- Metadata is small but never below 12px equivalent.
- The eye should travel: name/status → headline/body → action; the action row is the
  quietest part.
- Long Indonesian headlines wrap to two lines; do not truncate the headline before the
  dek.
- Right-rail cards are short: one idea, one action, no nested cards.

## 9. Accessibility notes from the reference

- Nextdoor pairs navigation icons with labels; this is the right default for maulid.id.
- Touch targets in the captures look comfortable (roughly 40–48px); do not shrink them to
  fit more controls.
- White cards on a light grey canvas provide a subtle but visible boundary; borders must
  remain at 3:1 against the canvas for non-text UI.
- Feed images have no visible alt text in a screenshot; maulid.id must provide meaningful
  alt text or mark decorative images correctly.
- The green primary on white must be verified at 4.5:1 for text; if it fails, use the
  darker green for text and reserve the lighter green for fills.

## 10. Dark mode

Nextdoor captures include a dark theme: black page canvas, dark grey cards, white text,
muted grey metadata, and the green active/FAB state preserved.

- Define separate dark tokens for canvas, surface, border, ink, muted and primary; do not
  invert the light palette mechanically.
- Ads remain bright in the reference and become visually loud; maulid.id has no ads.
- Verify contrast independently in both themes.
- Dark mode is optional for maulid.id; do not add it until the light theme passes review.

## 11. What not to copy visually

- Nextdoor's logo, wordmark, green brand equity, illustrations, icon set, photography,
  advertising units and exact card compositions.
- The "neighborhood" visual language (map pins, house illustrations, local-deal badges) is
  tied to a different product model.
- Dense ad rails, "For Sale & Free", "Faves" and business-promotion surfaces are outside
  maulid.id's domain.
- The US-centric content voice and neighbourhood terminology.

The goal is to reproduce **clarity, locality and calm utility**, not Nextdoor's identity.
