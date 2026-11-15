# Nextdoor component contracts

These are the component shapes observed in the Mobbin captures, translated into contracts
for maulid.id. A contract describes what the component must contain and how it behaves; it is
not a copy of Nextdoor's markup.

## App shell — web

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

- Left nav ~17% of width, feed ~38–39%, right rail ~19% in the 1920 capture; the rest is
  margin. Feed:rail ≈2:1.
- Top bar ~82px; search is centred in the viewport, not in the column.
- Feed cards: 20px padding, 12px gap, 1px hairline, 12px radius, no resting shadow.
- Feed is the widest content zone; the rail is visually quieter.
- At medium widths, collapse the rail first, then the left nav.
- The shell owns page background, max-width and sticky behaviour.

**maulid.id:** `Kabar`, `Kontribusi`, `Akun`, `Backoffice` (staff). The right rail carries
community context, validation status and guidance, not ads.

## App shell — mobile

```text
Header: [compass] community ⌄ · bell · chat · avatar
Feed:   onboarding module (optional) · posts · local-news rail
Tab bar: Home · Search · For Sale · Faves · raised green Post FAB
```

- Bottom bar 64–72px plus safe area; active item green, inactive grey.
- The Post FAB is a raised circle at the right; it is the composer entry.
- Header is sticky; the feed scrolls under it.
- Dark mode swaps canvas/surface tokens but keeps the green role.
- maulid.id replaces the commerce tabs with `Beranda · Cari · Tulis · Akun`.

## Top bar

- Brand at left, search centred in the viewport, notification/message/avatar at right.
- Search is a pill input with a leading icon and a visible placeholder.
- Account actions are icon+avatar; badges are small and real.
- Height is comfortable (~64px web), with a hairline bottom border or canvas contrast.

**States:** default, search focus, notification badge, signed-out (Masuk).

**maulid.id:** brand + “Cari kabar” + account/status. Do not show notification/message icons
until those features exist.

## Left navigation

- White surface, hairline right border.
- Items: icon + label, ~40–44px tall, 12–16px horizontal padding, rounded active block.
- Active item: light grey fill and a filled icon; not a green block and not a heavy pill.
- Primary action: filled green pill with a plus/pen icon, full width of the nav content.
- Secondary links at the bottom: Settings, Help Center, Invite neighbors, in a quieter style.
- The nav is sticky and independently scrollable if it exceeds the viewport.

**maulid.id:** Home/Kabar, Kontribusi, Akun. Add `Backoffice` only for staff. Primary action
“Tulis kabar” opens the composer gate. No For Sale/Faves/Groups/Events until those domains
exist.

## Right rail

- Stack of white cards, 12px radius, 16px padding, 12px gap.
- One idea per card; one action per card.
- First card can be the community/identity card: name, short context, “See all” link.
- Second card can be an invitation/status card with an explanation and a dark CTA.
- Third card can be a contextual promotion; maulid.id replaces this with guidance or a
  seasonal card.
- The chat dock is a floating utility, not part of the rail stack.

**maulid.id:** “Tentang komunitas” (purpose, validation rule), “Status akun” (Belum/Sudah
divalidasi), “Panduan” or a Maulid seasonal card. Maximum 3 cards.

## Feed header

- Left: “Filter by” dropdown or a small set of tabs.
- Right: primary “Post” action.
- Optional section title above the row.
- Sticky on scroll when the feed is long; must not cover the first card.

**maulid.id:** section title “Kabar komunitas”, search field, one filter row (when the API
supports it), primary “Tulis kabar”/“Minta validasi”.

## Post card

Anatomy:

```text
[avatar]  Author name · status · relative time · audience   [overflow]
Body / headline
[optional image, event card, or link preview]
[Like] [Comment] [Share]                                  [quiet utility]
```

- Identity row is the trust unit. Avatar 32–40px, name 14px/600, metadata 12–13px muted.
- Body 15–16px, line height 1.5–1.6; max readable width.
- Action row is quiet: outlined/ghost buttons or icon+label; counts are secondary.
- Overflow at the far right; menu items state consequences.
- Card padding 16–20px; hairline border; 12px radius; no heavy shadow.
- On mobile, the action row can be icon-only with accessible names.

**Variants:**

| Variant | Additional content | maulid.id status |
| --- | --- | --- |
| Text | body only | supported |
| Image | full-width image after body, optional domain chip | future schema |
| Event | date block + title + time + place + thumbnail | future domain |
| Poll | question + full-width options + real vote percentages | future domain |
| Place/listing | square thumbnail + title + price/status + location | not applicable |
| Publisher/news | publisher lockup + headline + image + reactions | future/source model |
| Sponsored | same geometry + “Sponsored” label | never |

**States:** default, hover (border/background shift), focus-visible, like-pending,
like-active, error, loading skeleton.

**maulid.id:** text + like only today. The author row uses the initial avatar, name,
`Sudah divalidasi`, relative time. Do not render an empty image frame.

## Identity row

- Avatar/initial, name, status/role, relative time, optional audience icon.
- Name is the only strong text in the row; status and time are muted.
- Publisher variant: logo + publisher name + verification label + time.
- Staff variant: “Dikelola staf maulid.id” instead of a person's name.

**Rules:** status is plain text; no gold/blue check; no follower count unless it is real and
useful.

## Account status / validation

- Canonical states as text: `Belum divalidasi` → `Menunggu peninjauan staf` →
  `Sudah divalidasi`.
- Each state has one sentence and the next step; no fake progress bar.
- After validation, state what is now possible (`Tulis kabar`).
- No gold check, shield, trophy, KYC funnel, address/map or location permission.
- If a checklist is used, every item must be real and resumable.

## Reaction bar

- Ordered: like → comment → share → overflow.
- Icon + optional count; no vote arrows; no emoji reaction picker unless specified.
- Pressed state uses primary colour and `aria-pressed`.
- Pending state disables only the pressed action; no whole-card spinner.

**maulid.id:** Like only. Comment/share hidden until APIs exist.

## Post utility footer

- Low-contrast line for real metadata: “Post insights / View” in Nextdoor.
- maulid.id: omit unless there is real data (e.g. edited timestamp, moderation state).

## Publisher/news card

- Publisher lockup: logo, name, verification label, relative time.
- Headline (2–3 lines), optional image (16:9), optional short summary.
- Section label (“Today's local news”, “Kabar pilihan”) and “See more”.
- Reactions/comments can appear below the item when the news is social.

**maulid.id:** a future “Kabar pilihan”/staff-curated module. It needs a featured flag or
editorial collection in the API; do not fake it with a hard-coded list.

## Event card / event row

- Date block (day + month) or date gutter.
- Title, time (`WIB`), venue, organiser, optional thumbnail.
- One action: Interested / Join / Lihat kegiatan.
- Optional real attendee count/avatar stack.

**maulid.id:** future event domain. The Maulid-season agenda is the strongest use: date
gutter, title, venue, WIB, organiser.

## Poll

- Question as body; 2–4 full-width option rows.
- Options have a border/background state; after voting, show real percentages and the vote
  count for the poll, not per-option fake data.
- One vote per account; the selected option is visually distinct and accessible.
- **maulid.id:** future domain; requires a poll model, vote storage and uniqueness rules.

## Listing / classified

- Square thumbnail, bold title, price/status, location line, overflow.
- **maulid.id:** not applicable. Do not add marketplace surfaces.

## Local-news module

- Section label + “See more” pill; horizontal rail on mobile, stacked or rail on web.
- Publisher lockup: logo, name, verification check, `Local publisher · time`.
- Headline 2–3 lines; 16:9 image; domain chip; reactions.
- Optional neighbour reaction/comment under the news card.
- **maulid.id:** future “Kabar pilihan”. Needs a source/publisher field and a featured
  mechanism. Do not hard-code a publisher list.

## Onboarding / growth module

- Optional module above the feed: horizontal cards (“Stay up to date”, business), progress
  line, dismiss control.
- **maulid.id:** replace with a validation/onboarding explanation only if the account state
  needs it. No 7-step growth funnel, no notification priming.

## Composer

Anatomy:

```text
Author identity · validation status
[body textarea]
[rules line]
[Cancel] [Terbitkan]
```

- Identity and status at the top.
- Body is the primary field; no title until the API supports one.
- Publish is the only filled action; cancel is quiet.
- Unvalidated: publish disabled, explanation visible, link to account/contribution.
- Rules line is short and links to the full guidelines.
- On mobile, a full-screen route or sheet; keyboard must not cover the publish action.

**States:** empty, typing, submitting, success, validation error, server error, unvalidated.

## Comment row and composer (future)

- Row: avatar, name, status/role, time, body, reply action, overflow.
- Composer pinned above the keyboard; one text field + send.
- Authority comments (staff/author) can be pinned/separated, not just badged.
- Report/hide only when the moderation backend exists.

## Search and filters

- Search input: visible label or accessible name, leading icon, clear button, 40–44px height.
- Tabs/chips: plain text with an underline or a soft active block; not a wall of pills.
- Filters: one row, max 2–3 controls; a “Filter by” dropdown is enough.
- Results: compact text-first rows; show result count and a recovery action when empty.

## Buttons

| Type | Shape | Use |
| --- | --- | --- |
| Primary | filled green, pill, 44–48px | one per surface |
| Secondary | outlined or dark filled | invite/send/confirm |
| Quiet | text/ghost | cancel, see more, secondary |
| Icon | 40–44px hit area, accessible name | overflow, close, share |

- Hover/pressed states are required.
- Disabled state keeps readable contrast and explains why when the reason is not obvious.
- Loading state preserves width and shows a spinner/label.

## Inputs and selects

- Visible label above the field; placeholder is not a label.
- 40–44px height, 12px radius or pill for search, hairline border, focus ring.
- Inline error message below the field, linked by `aria-describedby`.
- Selects use a native control or an accessible listbox with keyboard support.

## Tabs and chips

- Tabs switch views; chips filter or tag. Do not use the same style for both.
- Active state: underline or a light grey block/filled icon; text weight change is enough.
- Horizontal scrolling is allowed for chips; the page must not scroll horizontally.

## Avatars

- Circle; image or initial fallback on primary green.
- Sizes: 24px inline, 32–40px card, 48–64px profile.
- Decorative avatars use empty alt; meaningful avatars use the person's name.

## Badges and status labels

- Status is a text label next to the name, not a competing trophy.
- `Sudah divalidasi` uses a small check glyph + text at most.
- No gold/blue paid-check styling.
- Role labels (`Staf`) are plain text.

## Menus, sheets and modals

- Overflow menu: anchored, keyboard navigable, items with consequences.
- Sheet: 16px top radius, drag handle, title, close; focus trapped; Escape closes.
- Modal: only for a blocking decision; otherwise prefer a sheet or route.
- Toasts: short, non-blocking, `role="status"`; no toast for destructive success without an
  undo path.

## Empty, loading, error and permission states

| State | Content | Action |
| --- | --- | --- |
| Empty feed | one sentence explaining why | “Tulis kabar” if validated, otherwise validation link |
| Filtered empty | “Belum ada kabar yang cocok.” | clear filter |
| Loading | skeleton or “Memuat kabar…” | none |
| Error | “Kabar tidak dapat dimuat.” | “Coba lagi” |
| Unvalidated | explain staff review | “Lihat status akun” / “Cara validasi” |
| Offline/network | plain message | retry |

No fake content, no illustration-only dead end.

## Web measurements (estimated)

| Element | Value |
| --- | --- |
| Top bar | ~82px |
| Left nav | ~243px at 1440 (~17%) |
| Feed column | ~553px at 1440 (~38–39%) |
| Right rail | ~271px at 1440 (~19%) |
| Card padding | 20px |
| Card gap | 12px |
| Card radius | 12px |
| Search field | 44–52px, pill |
| Primary button | 40–44px, pill |
| Events grid | 3-up |
| For Sale grid | 6-up |
| Post detail | centred modal (~16px radius); publisher variant split-pane |
| Menus/popovers | ~8px radius + soft shadow |

Nextdoor web was captured at ~1920px only; tablet/phone breakpoints are maulid.id decisions.

## Mobile measurements (estimated)

| Element | Value |
| --- | --- |
| Card radius | 12–14 |
| Card padding | 12–16 |
| Card gap | 8–10 |
| Author avatar | 36–40 |
| Name | 14–15 / 600 |
| Metadata | 12 / 400 |
| Body | 15–16 / 21 |
| Reaction hit area | 36–40 |
| Tab bar | 49 + safe area |
| FAB | 56 |
| Sheet radius | 18–20 top |
| Sheet row | 56 |
| Form field | 44–48 |
| Primary button | 48–50 |
| Comment field | 36–40, send 32 |
| Toast | 48 high, 8–10 radius |
| Profile cover | 16:9 |
| Feed photo | 4:3 or 16:9, inset 8, radius 8–10 |

These are read from the 390pt captures and labelled estimated. They are not official tokens.

## Responsive rules

| Width | Web behaviour |
| --- | --- |
| ≥1280 | three zones, sticky left nav, right rail visible |
| 1024–1279 | left nav + feed; rail below feed or hidden |
| 768–1023 | top navigation or collapsible nav; single feed; no rail |
| <768 | single column; bottom navigation; composer full-screen |

Mobile specifics: 16px page padding, 12px card gap, 44px targets, bottom bar 64–72px,
safe-area insets, keyboard-aware composer.

## Component reuse rules

- Reuse `shared/components` primitives (Button, Field, StatusMessage, LoadingState) before
  creating new ones.
- A component used by two features moves to `shared/`, not copied.
- Every new component gets a colocated CSS file and behavior tests.
- No new icon library, font or colour token without a recorded decision.
