---
name: Nextdoor Design Reference
description: Mandatory design reference and UI review rules for maulid.id, derived from Mobbin captures of Nextdoor's web and iOS products. Use before designing, building, restyling or reviewing any public page, account screen, feed, post, composer, event, empty state, navigation or backoffice surface.
---

# Nextdoor design reference (maulid.id)

This skill turns the Nextdoor reference selected by the project owner into operating rules
for maulid.id. It is a **structural and interaction reference**, not a visual clone.

Read `docs/design/README.md` for the human-readable index. The long-form evidence lives in
`docs/design/nextdoor-reference.md`, `nextdoor-components.md`,
`nextdoor-principles.md`, `nextdoor-visual-language.md` and `ui-review-checklist.md`.

## When to use it — mandatory

Load this skill before any of these:

- designing or building a new page, screen, section or component;
- changing layout, typography, colour, spacing, iconography, imagery or motion;
- adding or changing a feed card, post card, composer, event card, empty state, navigation
  or sheet;
- reviewing a UI change for design quality or accessibility;
- writing a plan that includes visual work.

If the task is pure Go/API/security work with no visual surface, this skill is not required.

## Non-negotiables

1. **Local context first.** Say what the reader is looking at (community, feed, scope)
   before the first item.
2. **The feed is the product.** One reading column; no dashboard of modules around it.
3. **Identity before content.** Every post shows author, canonical validation status and
   relative time before the body.
4. **Plain-language trust.** Use `Belum divalidasi` / `Sudah divalidasi`; never a paid-check
   tick, trophy or vague badge.
5. **One primary action per surface.** One filled action; everything else is secondary or
   in overflow.
6. **Stack, not grid.** The feed is a vertical card/row stack. No masonry, no equal-weight
   tile grid for posts.
7. **Progressive disclosure.** Rare, destructive and administrative actions live in an
   overflow menu or sheet, with consequences stated.
8. **The right rail is utility.** It carries context, account status or guidance — never a
   second feed.
9. **Honest states.** Empty/loading/error/permission states use Indonesian copy and one next
   action; no fake posts or fake counts.
10. **Calm utility.** No points, streaks, leaderboards, engagement scores, urgency timers or
    paywalls.
11. **Same grammar across platforms.** Web and mobile keep the same post information order;
    only the layout adapts.
12. **Accessibility is structural.** Focus-visible, labels, landmarks, 44px targets,
    reduced motion and AA contrast are part of the component.
13. **Indonesian by default.** All labels, dates, empty states and errors are Indonesian.
14. **No copied identity.** Do not reuse Nextdoor's logo, green brand equity, illustrations,
    icon set, photography, copy, ad formats or exact compositions.
15. **No unsupported features.** If the API does not have title, image, category or event
    fields, do not simulate them with fake content. Document the schema work instead.
16. **Media is optional.** A text-only post is complete; never render an empty media frame.
17. **Two primaries, two jobs.** Green for publishing/positive, `--ui-strong` for operational
    confirmation, red only for destructive actions and errors.

## Layout contracts

### Web

- Three zones where the page allows it: left navigation (~17%), central reading column
  (~38–39%), right utility rail (~19%) at desktop. The column does not widen to fill when
  the rail is absent.
- Reading column is the widest content zone; right rail is quieter and optional.
- Top bar: product mark, search/filter, account actions; search centred in the viewport.
- Left navigation: a short list of destinations and one primary “Tulis kabar” action; the
  active item is a light grey fill + filled icon, not a green block.
- Feed: filter row, then a single column of cards separated by 8–12px.
- Right rail: at most 2–3 short cards; collapses below the feed or disappears on narrow
  widths.
- A page that omits a zone must say why in the plan.

### Mobile

- Root shell: at most four labelled destinations plus one raised composer action.
- Notifications, messages and profile are header utilities until they are real surfaces.
- Feed is one column with 12–16px page padding.
- Composer and filters use sheets or full-screen routes, not cramped inline dialogs.
- Sticky bars must not cover content, the keyboard or the safe area.

## Component contracts

Use the detailed anatomy in `references/components.md`. The minimum rules:

- **Post card:** identity row → body/headline → attachment (if any) → actions → quiet
  utility footer. Actions: like, comment (when available), share (when available),
  overflow. No fake counts.
- **Composer:** author identity and validation status visible at the top; disabled publish
  for unvalidated accounts with an explanation and a path to validation; rules line at the
  bottom; one primary publish action.
- **Event card:** date/time with `WIB`, venue, organiser, one RSVP/interested action,
  optional attendee count only when it is real.
- **Publisher/trust label:** plain text next to the name, not a competing badge.
- **Empty state:** one sentence, one recovery action, optional topic suggestions.
- **Overflow menu:** Hide/Mute/Report/Share/Bookmark as applicable, each with a one-line
  consequence.
- **Buttons:** one filled primary, outlined/quiet secondary, text tertiary; pill radius for
  primary controls.
- **Forms:** visible labels, 40–44px control height, inline errors, no placeholder-only
  labels.

## Visual tokens

The full mapping is in `references/visual-language.md` and
`docs/design/nextdoor-visual-language.md`. Summary:

```text
canvas      #F7F8F7
surface     #FFFFFF
ink         #222222
muted       #626B66
border      #E3E7E4
primary     #1F6F4A   (maulid.id green, not Nextdoor's brand green)
primary-2   #17553A
strong      #1F2A44
danger      #B42318
radius      12px cards, 999px controls, 50% avatars
spacing     4 / 8 / 12 / 16 / 20 / 24 / 32 / 40 / 48 / 64
font-ui     system sans stack
font-display EB Garamond (existing editorial voice)
```

- Green marks the primary action, active navigation and positive/selected state. It is not
  decoration.
- White cards sit on the light canvas; use borders, not shadows, for structure.
- Keep EB Garamond for the public editorial voice; introduce the UI sans only for
  controls, metadata and forms. Record the decision.
- No new font, icon library or colour token without a recorded design decision.

## Required references

Read in this order for the task:

| Task | Read |
| --- | --- |
| Web page, layout or navigation | `references/web-app.md` |
| Mobile screen or navigation | `references/mobile-app.md` |
| Full captured screen list | `docs/design/nextdoor-screen-inventory.md` |
| Post, composer, event, sheet, menu | `references/components.md` |
| Feed, trust, moderation, states, seasonal | `references/patterns.md` |
| Copy, tone, dates, empty states | `references/content.md` |
| Colour, type, spacing, imagery, motion | `references/visual-language.md` |
| Adapting to maulid.id, anti-patterns, roadmap | `references/adaptation.md` |
| Source URLs and capture limitations | `references/sources.md` |

## Review gate

Before calling UI work complete:

1. Apply `docs/design/ui-review-checklist.md`.
2. Attach desktop (1440), tablet (1024/768) and phone (390) screenshots.
3. Verify keyboard, focus-visible, 200% zoom and reduced motion.
4. Run the ordinary frontend gates: lint, tests, coverage >80%, build.
5. Record the commands, results and any limitation in `.agents/journal/`.
6. Record deliberate deviations in the plan or journal; silent drift is a defect.

A UI change without this evidence is not complete.

## Boundaries

Do not add advertising, sponsored cards, paywalls, member-only gates, promo countdowns,
streaks, points, badges, leaderboards, fintech identity verification, autoplay video,
heavy carousels, US hyperlocal framing or copied Nextdoor assets/copy.

Do not use `dangerouslySetInnerHTML` or browser storage for credentials.

## Maintenance

- If a design convention changes, update this skill, the matching `docs/design/` file and
  the checklist in the same change.
- If the skill and `docs/design/` disagree, `docs/design/` wins and the skill is corrected.
- New Mobbin evidence must include the product, platform and Mobbin URL in
  `references/sources.md`.
- Mobbin captures show what shipped; they do not prove a metric. Do not present them as
  performance evidence.
