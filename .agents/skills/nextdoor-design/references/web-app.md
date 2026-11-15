# Nextdoor web app — screen and layout reference

**Captured evidence:** the four web screens in `sources.md` (web home ×2, community rail,
local publisher sheet). Items marked **estimated** are read from the image; items marked
**not captured** are product surfaces that Mobbin did not return for this slice and must not
be described as if verified.

## 1. Shell

```text
┌──────────────────────────────────────────────────────────────┐
│ Top bar: brand        search             notifications  account│
├───────────────┬──────────────────────────────┬───────────────┤
│ Left nav      │ Feed header                  │ Right rail    │
│ destinations  │ filter · primary action      │ utility cards │
│ primary CTA   ├──────────────────────────────┤               │
│ secondary     │ Post card                    │               │
│ links         │ Post card                    │               │
│               │ ...                          │               │
└───────────────┴──────────────────────────────┴───────────────┘
```

**Observed proportions** at 1920px: left nav ~17% (~243px at 1440), central feed ~38–39%
(~553px), right rail ~19% (~271px). The column-to-rail ratio is about 2:1; the remaining
width is margin. The column does not widen to fill when the rail is absent.

**Observed surfaces:** white cards on a light grey canvas (`#F8F9FA`), hairline borders
(`#EBEDED`), 12px card radius, 20px card padding, 12px gap, no resting shadow.

**Top bar:** fixed, ~82px; brand left, search pill centred in the viewport, bell/chat/avatar
right. Business accounts drop the bell/chat. Faves renders without a search field.

## 2. Top bar

- Brand at the far left, above the left navigation.
- Search field centred in the viewport, with a leading search icon and a contextual
  placeholder such as “Search for Landscaping” / “Search for House Cleaner”.
- Notifications bell, messages icon and avatar at the far right.
- White background; the feed starts below it.
- The top bar does not carry the primary Post action; that sits in the left nav and the feed
  header.

**maulid.id:** brand, “Cari kabar”, account status/avatar. Omit notification/message icons
until those APIs exist.

## 3. Left navigation

Observed order:

```text
Home
For Sale & Free
Faves
Groups
Events
[ Post ]  ← filled green pill
Settings
Help Center
Invite neighbors  ← secondary links at the bottom
```

- Each destination is an icon + label row, about 40–44px tall.
- The active item uses a **light grey fill and a filled icon**, not a green block.
- The Post pill is full-width inside the nav content, with a plus/pen icon and white text.
- Secondary links are smaller, quieter and anchored near the bottom.
- The sidebar is white with a hairline right border and its own scroll if needed.
- Business accounts replace the middle group with Dashboard / Ads / Notifications / Chats /
  Profile; the shell itself stays constant.

**maulid.id:** `Kabar`, `Kontribusi`, `Akun`, `Backoffice` (staff only), primary “Tulis
kabar” / “Minta validasi”. No For Sale/Faves/Groups/Events until those domains exist.

## 4. Feed header

- A “Filter by” dropdown on the left.
- A green **+ Post** button on the right.
- Sits above the first card with a small vertical gap.
- The filter is compact; it does not become a toolbar of chips.

**maulid.id:** section title + search + one filter row when supported + primary action.

## 5. Feed and post cards

Observed post anatomy:

```text
[initial avatar]  Name
                   Neighbourhood · time · audience icon        [overflow]
Body text (first person, short paragraphs)
[optional image or attached place/event card]
[heart] [comment] [share]                                      [...]
Post insights                                            View
```

- The identity row is compact: 36–40px avatar, 14–15px name, 12–13px muted metadata.
- Body text is the dominant content, 15–16px, ~1.45 line-height, clamped at ~4 lines with
  “… see more”.
- Card padding ~20px; card gap ~12px; image attachments are full-bleed to the card edges.
- Attached place/event cards use a square thumbnail, bold title, green open/status badge,
  address and a small action.
- Reaction glyphs are outlined and quiet; the overflow sits at the far right.
- “Post insights / View” is a low-contrast utility footer, not a primary action.
- Sponsored cards reuse the editorial card geometry and add a small “Sponsored” label.

**Variants observed:** text, text + image, text + poll, text + event chip, text + listing,
quoted post, long-form, sponsored.

**maulid.id:** text + like today; author/status/time; no image frame, no fake actions.

## 6. Right rail

Observed cards (top to bottom):

1. **Neighbourhood card** — map thumbnail, neighbourhood name, “See all alerts”, chevron.
2. **Neighbourhood stats** — “2,262 neighbors”, progress bar, “49% of 2,835 households”.
3. **Invite neighbours** — explanation, address checkbox list with distances, dark
   “Send invitation” button.
4. **Business CTA** — photo, “Own a local business?”, explanation, “Create page”, chevron.
5. **Ads** — 1–2 sponsored units.
6. **Chats dock** — floating bottom utility.

- Rail ~271px at 1440 (19%), 16px card padding, 12px gap.
- Cards are white, 12px radius, 1px hairline, no resting shadow.
- One idea and one action per card; density drops as you go down.
- The dark CTA (`#222E46`) is used for the invite action; the green primary is reserved for
  Post.

**maulid.id:** replace with “Tentang komunitas”, “Status akun”, “Panduan” or a seasonal
Maulid card. No ads, business CTA, invite address or chat dock until the domain exists.

## 7. Local publisher sheet

Observed on web:

- Split layout: large photo left, white panel right.
- Right panel: publisher avatar, publisher name, “Local publisher”, relative time, title,
  two-line summary, dark “Read article” pill.
- Below: reaction count, comments count, named commenters, comment input pinned at the
  bottom.
- The trust label is plain language, not a badge icon.

**maulid.id:** a future “Kabar pilihan”/linked-news treatment. It needs a source/publisher
field and a featured mechanism; do not hard-code it.

## 8. Search

**Captured.** Back chevron + horizontally scrollable tabs: `All`, `Posts`, `Businesses`,
`Pages`, `For Sale & Free`, `Event`, `Neighbors`.

- “All” is sectioned by content type: rows for events/people, a 3-up grid for marketplace,
  each section with a “See all” link.
- A single-type tab is one long list of the same row.
- Result rows are full-width white cards with hairline borders and ~10px gaps.
- No empty-results capture; use an Indonesian recovery state.

**maulid.id:** current text search over `body`; add a clear button, result count and an
Indonesian recovery state. One filter row only when the API supports more filters.

## 9. Composer

**Captured.** Composer card: author row, prompt, growing textarea, attachment zone, tool row;
audience selector pill and Post pill. The “Create something” popover offers Sell/Event/Poll.
Attachments (listing/event/poll) are removable inline cards; an AI suggestion card can be
inserted above the text.

- For maulid.id, implement the gate in `docs/design/nextdoor-components.md`; do not copy
  the marketplace/poll intent cards, because maulid.id has one content type today.

## 10. Groups and events

**Captured.**

- Groups directory: two columns — “Groups near you” list + “Your groups” rail with Create.
  Rows = square avatar, name, Join/Request/Joined pill, right-aligned member count.
- Group page: banner, name + overflow, description, Invite split-button, member count;
  optional location card; “Welcome/Finish setup” carousel; progress bar; inline composer;
  pinned post.
- Events: wide container, 3-up card grid; 4:3 image, overlaid date block, 2-line title,
  time · venue, full-width outline RSVP pill; “Your events” underline tabs; event detail is a
  narrow column with a 16:9 cover and an RSVP state selector.

## 11. Notifications, messages, profile, settings

**Captured on web.**

- Notifications/activity: back-chevron page with a bordered panel of activity rows (icon,
  action sentence, relative time) and a privacy footnote; neighbourhood activity uses tabs
  and 28px avatars.
- Chats: two-pane list + thread; empty state with illustration and “Send a message”.
- Profile: banner, overlapping avatar, name/location, dark Message pill, overflow, Posts,
  skeleton rows while loading.
- Settings: Account / Notifications / Privacy tabs; two-column name fields, email copy,
  password/mobile, native selects, Save at the card foot.
- Notification settings: five groups × Allow / Email / Push toggles.
- Privacy: profile visibility controls with plain-language explanations.

**maulid.id:** defer these until the API exists. Do not render empty shells.

## 12. Moderation

**Captured on web.**

- Post overflow includes Report post, Hide (with consequence), About ads; own post adds
  Edit/Delete/Close discussion; ad overflow adds Review ad/About ads.
- Confirm modals state the consequence and use a red destructive button.
- The ad report form has a category taxonomy and an optional 500-character note.
- Community guidelines appear as four icon rows with one acknowledgement.

**maulid.id:** do not show report/hide until staff can review reports. When added, use the
same consequence-first menu and publish the rules. Do not build ad-report plumbing.

## 13. Responsive behaviour

**Nextdoor evidence is desktop-only (~1920).** No tablet or phone web breakpoint is captured;
the table below is a maulid.id decision, not a Nextdoor observation.

| Width | maulid.id behaviour |
| --- | --- |
| ≥1280px | three zones; sticky left nav; right rail visible |
| 1024–1279px | left nav + feed; rail below the feed or hidden |
| 768–1023px | collapsible/top navigation; single feed; no rail |
| <768px | single column; bottom navigation; full-screen composer |

- Collapse the rail before shrinking the reading column.
- Keep the primary action visible at every width.
- Do not let the rail become a horizontal carousel on mobile.
- Do not attribute these breakpoints to Nextdoor.

## 14. Implementation contract for maulid.id

1. Build the shell first: top bar, left nav, reading column, optional rail.
2. Use the visual tokens in `visual-language.md`; do not inline raw colours.
3. Use the post card contract in `components.md`.
4. Use the composer gate; show account status and the validation path.
5. Use honest empty/loading/error states.
6. Test at 1440, 1024, 768 and 390.
7. Complete `docs/design/ui-review-checklist.md`.
