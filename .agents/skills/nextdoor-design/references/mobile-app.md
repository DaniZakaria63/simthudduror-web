# Nextdoor iOS app — screen and layout reference

**Captured evidence:** the iOS screens and flows listed in `sources.md`, plus the additional
home-feed captures below. Values marked **estimated** are read from images.

## 1. Header

```text
[compass/logo] Lower Allston ⌄        [bell] [chat] [avatar]
```

- The neighbourhood selector is the first control; the compass/logo mark sits to its left.
- The bell carries a small red dot for unread notifications.
- Chat is a separate icon; the avatar is a circle with an initial or photo.
- The header is white (or dark in dark mode), with the feed starting immediately below.
- A sticky header may gain a subtle border on scroll.

**maulid.id:** brand/community label + account avatar. No bell/chat until notifications and
messages exist.

## 2. Navigation model

**Root shell (observed):** four labelled consumer destinations — `Home`, `Search`,
`For Sale`, `Faves` — plus a raised green **Post** FAB at the right. Notifications, messages
and profile are utilities in the Home header, not tabs.

```text
Home     Search     For Sale     Faves     [green Post FAB]
```

- Five visual slots: four text+icon destinations and one raised FAB.
- Active destination uses a filled green icon/label; inactive items are grey.
- The FAB is a 56pt green circle with a white pen glyph; it overlaps the bar.
- The bar is white, with a hairline top border, 49pt + safe area.

**Observed variant:** a floating capsule tab bar (four icons + labels in a white pill, no
FAB) appears on Search results, Events, Your events, Group page and Faves. Treat it as a
variant to resolve, not two systems.

**maulid.id:** `Beranda`, `Cari`, `Tulis` (FAB), `Akun`. No For Sale/Faves. Keep at most
five slots; notifications/messages are future header utilities, not tabs.

## 3. Home feed composition

Top-to-bottom order observed across captures:

1. Header (neighbourhood, bell, chat, avatar).
2. Optional onboarding module (“Get started on Nextdoor”) with cards and progress.
3. Post cards.
4. “Today's local news” module (a horizontal rail with a “See more” pill).
5. More post cards and/or sponsored cards.
6. Bottom navigation.

- The feed is a single column; cards are white on a light grey canvas (or dark on black).
- Card gap is small (~8–12px); cards are edge-to-edge with 12–16px page padding.
- The reading order is always identity → body → attachment → actions.
- Sponsored cards reuse the post geometry and add a small “Sponsored” label.

## 4. Post card anatomy

Observed:

```text
[avatar] Name
         Neighbourhood · time · audience icon        [overflow]
Body text
[optional attachment: event, place, poll, image, listing]
[heart] [comment] [share]                            [...]
[Post insights]                                      [View]
```

- Avatar 32–40px; name 14px/600; neighbourhood/time 12–13px muted.
- Audience icon (globe) sits after the time; it indicates who can see the post.
- Body 15–16px, line height 1.4–1.5.
- Action row uses outlined circular buttons or icon-only glyphs; counts can appear inside
  the pill (`heart 3`, `comment 4`).
- “Post insights / View” is a low-contrast utility row.
- Overflow sits at the far right of the action row.

**Variants observed:**

| Variant | Visible structure |
| --- | --- |
| Text | body + reaction row |
| Place/listing | square thumbnail, bold title, price/status, location chip |
| Event | large date block (`24 Jul`), title, time/address, thumbnail |
| Poll | question + 2–4 full-width options + vote percentages |
| Image | full-width photo, optional overlay controls |
| For Sale & Free | category label in identity line, listing card |
| Sponsored | “Sponsored” label, same geometry, ad creative + CTA |
| Local news | publisher lockup + headline + image + optional neighbour reaction |

## 5. Today's local news

- Section label “Today's local news” with a white “See more” pill at the right.
- A horizontal rail of publisher cards; the next card peeks at the right edge.
- Card: publisher logo + name + verification check + “Local publisher · 1h”, headline
  (2–3 lines), 16:9 image, domain chip (`masslive.com`), then reactions.
- Some news cards include a neighbour's reaction/comment underneath.
- The module mixes editorial and community response without turning the news into a normal
  post.

**maulid.id:** future “Kabar pilihan” module. Needs a source/publisher field and a featured
mechanism; do not fake it.

## 6. Onboarding module

Observed “Get started on Nextdoor”:

- Section title with an overflow menu.
- Horizontal cards: “Stay up to date” with a bell illustration and “Notify me” dark button;
  a business card with a “Get...” action.
- Progress line: “1 of 7 steps completed” with a progress bar.
- The module sits above the feed and can be dismissed.

**maulid.id:** a future account/onboarding module could explain validation and the feed.
Do not copy the 7-step growth funnel or notification priming.

## 7. Group feed

Observed:

- Back chevron, invite icon, settings icon.
- “Groups” section with “See all”; rows show an avatar, group name and member count
  (`Babysitting 223 members`).
- “Posts” section with “See activity”.
- Posts below; poll and listing attachments appear in the group feed too.

**maulid.id:** groups are future work. The section grammar (title + See all + rows) is
useful for any future collection.

## 8. Search

Observed:

- Search field with the query; tabs `All / Posts / Businesses / Pages`.
- A row of filter dropdowns (`Most Recent`, `Distance`, `All Time`).
- Result rows: initial avatar, name, neighbourhood, time, body preview, optional document
  thumbnail, secondary “See all N comments”.
- The result is text-first; the search does not become a grid.

**maulid.id:** current text search; add clear, result count and an Indonesian empty state.

## 9. Post detail and comments

- **Header:** 44pt, circular back, centred “Post details”; share/overflow at right.
- **Content:** article/post body, reaction row, emoji summary (`3 reactions`), comments list.
- **Composer:** pinned above the home indicator; six-emoji tray sits directly above it;
  send changes from pale to dark when text exists.
- **Comment row:** 28–32pt avatar, body 14–15/20, actions Like / Reply / overflow at 12–13pt;
  new/reply comments show `name · time · neighbourhood` above the body.
- **Empty:** `No comments` / `Be the first to comment!` — inline, no illustration.
- **States:** one/two/no comments, reaction tray open, keyboard open, report, hide/mute/delete
  confirms, toasts.
- **maulid.id:** comments are future work. Do not show a comment action until the API and
  moderation path exist.

## 10. Events

- **Index:** back + “Events” + compose; search; “Your events” chip; large cards with cover,
  date, title, venue, attendee line, full-width `Interested` button.
- **Detail:** 16:9 hero, date, title 24–26, venue, host row, description, attendee line,
  `Interested`/`Going` state selector, reaction row, Upcoming rail, comment composer.
- **RSVP sheet:** Going / Interested / Not Interested.
- **Create:** cover, title, start/end date-time pills, location, video link, details; Next
  enabled when valid; date picker is a sheet.
- **Empty:** `No events yet` + Create an event.
- **maulid.id:** future event domain. Use `WIB`, venue, organiser and a date-gutter agenda
  for Maulid season. Real attendance only.

## 11. Groups

- **Directory:** rows with 56pt rounded-square illustration, name, member count, navy Join.
- **Group page:** 16:9 cover, title 22, member count, full-width Join; About rows; posts use
  the same card grammar.
- **Joined:** Invite pill, video circle, green compose; “Welcome to the group” checklist.
- **Coachmark:** “Introduce yourself / What brings you to the group?” pointing at compose.
- **maulid.id:** groups are future work; the section grammar (title + See all + rows) is
  useful for any collection.

## 12. Notifications, settings and profile

- **Notifications:** back + title; `Neighborhood / My activity / Alerts` underline tabs;
  `Today / Last 7 days` groups; avatar rows with activity text and relative time.
- **Notification preferences:** category list → per-category `Allow / Email / Push` toggles.
- **Privacy:** four-level audience ladder (`Anyone on Nextdoor` / `Nearby neighborhoods` /
  `Your neighborhood` / `No one`) with icon, description and check.
- **Settings:** Account / Notifications / Privacy tabs; read-only identity, password, data
  download, deactivate.
- **Profile:** cover, overlapping avatar, name/location, Edit profile, Dashboard (self only),
  Faves, Groups, Posts; new-neighbour system post with Add an intro.
- **maulid.id:** defer notifications/messages/groups/profile until the API exists. No empty
  shells. The privacy ladder is not applicable without geodata.

## 13. Composer

Observed in the flow capture:

- Author avatar, name and neighbourhood at the top.
- Audience chip (`Anyone ⌄`) at the top right.
- Three intent cards: “Sell or give away”, “Create an event”, “Poll your neighbors”.
- Body field and media controls.
- Post action at the bottom.

**maulid.id:** one content type. Author identity + validation status, body, rules, publish.
No intent cards, no audience selector until there are audiences.

## 14. Post overflow and moderation

Observed menu:

```text
Hide — Remove post from your feed
Mute — Hide all posts from MassLive
Report — Flag for review
Turn off notifications — Don't show updates for this post
Share
Bookmark — Save post for later
```

- Each item has a consequence line.
- The menu is a sheet/action list anchored to the overflow.
- Community guidelines are a readable page with filter chips and numbered rules.

**maulid.id:** do not show report/hide until staff can review reports. When added, keep the
consequence-first menu and publish the rules.

## 15. Seasonal module

**No Nextdoor seasonal capture exists in this corpus.** An earlier draft cited a Ramadan
card; that was a misread of the Local Deals explanation screen. Do not present a Nextdoor
seasonal precedent.

**maulid.id:** the Maulid card is a maulid.id editorial decision: mark → greeting → one
action, no commerce, no countdown. Place it in the right rail on web or after the second
post on mobile, and let it be dismissed for the season.

## 16. Dark mode

- Captures show black background, dark grey cards (estimated `#1C1C1E`), white text and
  muted grey metadata; the green FAB and active nav remain green.
- The same card grammar and spacing apply.
- Ad creative remains bright and visually loud, which is a warning about contrast and
  attention.

**maulid.id:** if dark mode is added, define tokens for canvas/surface/border/text and
verify contrast; do not simply invert.

## 17. Type, spacing and shape (estimated)

### Type

| Role | Size / line height | Weight |
| --- | --- | --- |
| Display | 28–30 / 34 | 700 |
| Page title | 22–24 / 28 | 700 |
| Section | 18–20 / 24 | 600–700 |
| Card title | 16–17 / 21 | 600 |
| Body | 15–16 / 21 | 400 |
| Metadata | 12–13 / 16 | 400 |
| Micro | 10–11 / 13 | 400–500 |
| Button | 14–15 / 18 | 600 |

- One humanist geometric sans family; no serif.
- Sentence case; headings tight, metadata smaller and greyer.
- maulid.id keeps EB Garamond for editorial display and uses the UI sans for controls/meta.

### Spacing

```text
4 / 8 / 12 / 16 / 20 / 24 / 32 / 40
page margin 8–12
card padding 12–16
card gap 8–10
section gap 20–24
form gap 12
sticky footer 16 + safe area
```

### Shape

| Element | Estimate |
| --- | --- |
| Card radius | 12–14 |
| Inner media | 8–10 |
| Fields | 10–12 |
| Controls | 999 (pill) |
| Sheets | 18–20 top |
| Avatar | 50% |
| Card border | 1pt hairline, no shadow |
| Floating shadow | `0 8 20 rgba(0,0,0,.12)` |

### Element sizes

| Element | Estimate |
| --- | --- |
| Avatar | 24 / 28 / 32 / 36 / 40 / 56 / 64 / 80 |
| Reaction control | 36–40 hit area, 20–22 icon |
| Bottom bar | 49pt + safe area |
| FAB | 56pt |
| Sheet row | 56pt |
| Form field | 44–48pt |
| Primary button | 48–50pt |
| Comment composer | 36–40pt field, 32pt send |

### Colour ranges (iOS estimates)

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

These are ranges, not official tokens. maulid.id uses its own `--ui-*` values in the
visual-language reference.

## 18. Motion and interaction

- Sheets rise from the bottom with a grabber, scrim and rounded top corners; the feed behind
  is dimmed.
- The account drawer slides from the right at ~88% width with rounded left corners.
- Creation flows are full-screen modals with a fixed header commit action.
- Toasts overlay the tab bar with a dismiss X; they carry status, not navigation.
- The reaction tray and comment composer move with the keyboard.
- The home header scrolls away; the tab bar and FAB remain.
- Selecting an audience/filter closes the sheet and updates the chip/feed immediately; no
  Apply button.
- RSVP is a state selector, not a one-shot.
- No parallax, scroll-linked header collapse, spring or animated counters are evidenced.

**maulid.id:** keep motion to 150–200ms controls and 250–300ms sheets; respect
`prefers-reduced-motion`; no autoplay.

## 19. Adaptation for maulid.id

Adopt:

- Header with community context and account.
- Bottom navigation with `Beranda`, `Cari`, `Tulis`, `Akun`.
- One-column feed; identity → body → actions.
- Real reaction counts only; no fake comments/shares.
- Search result grammar.
- Consequence-first moderation menu when the backend exists.
- Calm seasonal card.

Reject:

- For Sale, Faves, business pages, local deals, ads, sponsored cards, ad carousels and
  autoplay video.
- For you / Nearby / Trending ranking framed as feed identity; prefer explicit, explainable
  scopes.
- Engagement counts as default proof, reaction pickers, red-dot activity pressure, virtual
  gifts, streaks and “1 of 7 steps” growth funnels.
- Neighbourhood/address verification, maps, invite address lists and `<1 mi away`.
- Business KYC as a posting gate; paid/role verification ticks.
- Marketplace grids, listing creation, Ads Manager, AI summaries of neighbours.
- US hyperlocal framing; maulid.id uses Indonesian community scope, not home addresses.
- Dense chrome with four tabs + FAB + bell + chat + avatar + filter sheet; keep the shell
  quieter and let the feed carry the product.

Defer:

- Groups, events, notifications, messages, profiles, comments, images, polls.
- All require API/schema work; document them in a separate plan.

## 20. Sources

See `sources.md` for the full list. Key home-feed captures:
[garage-sale feed](https://mobbin.com/screens/251e08e3-4688-4ea7-be88-557ec9871142),
[For Sale & Free + local news](https://mobbin.com/screens/54cefa73-8bc8-4c81-add3-ed34b1f570fe),
[local news + neighbour reaction](https://mobbin.com/screens/ed239926-f476-46ca-87ec-e897ff20028b),
[poll + event](https://mobbin.com/screens/1aae4c0b-e83c-4e67-9c2c-e5a11cff8698),
[dark mode](https://mobbin.com/screens/d178d1a1-d1fe-4cc4-b19e-24dc5f6fd108),
[group feed](https://mobbin.com/screens/9fa1c528-5530-4f11-bc24-0244ea7279fe).
