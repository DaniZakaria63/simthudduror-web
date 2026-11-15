# Nextdoor reference for maulid.id

**Selected:** 2026-11-15 by the project owner.
**Evidence base:** Mobbin captures of Nextdoor's web and iOS products. The full URL list is
in `.agents/skills/nextdoor-design/references/sources.md`.

## 1. Why Nextdoor

maulid.id is a community news product, not a general social network, a marketplace or a
publishing CMS. Nextdoor is the closest reference found in Mobbin because it demonstrates,
in one product:

- a local/community feed that people return to;
- a post card with identity, place and time before the message;
- an embedded event module;
- a composer with visible identity and audience;
- a local publisher/news treatment with a plain-language verification label;
- a web three-zone shell and a mobile tab bar;
- moderation and community rules that are readable.

It is **not** a like-for-like Islamic product. Mobbin has no dedicated Islamic community
news app. Nextdoor's structure is the reference; its subject matter, US neighbourhood
language, ads and business surfaces are not.

## 2. Product comparison

| Axis | Nextdoor | maulid.id |
| --- | --- | --- |
| Subject | Neighbourhood life | Maulid Rasulullah community news |
| Geography | Hyperlocal US neighbourhood | Indonesia, community-wide with optional regions later |
| Who can post | Verified residents and businesses | Staff-validated accounts only |
| Trust signal | Address-gated membership; role labels for publishers/agencies | `Sudah divalidasi` by staff |
| Feed unit | Post with optional event/photo | Kabar with body; image/title/event are future |
| Interaction | Reactions, comments, shares, messages | Like today; comments/share later |
| Monetisation | Ads, local deals, business pages | None |
| Seasonality | No seasonal module in the captured corpus | Maulid greeting/doa module (maulid.id decision) |
| Platform | Web + iOS | Public web today; mobile API for a future app |

## 3. Information architecture

### Web (observed)

```text
Top bar: brand · search · notifications · messages · account
Left navigation: Home · For Sale & Free · Faves · Groups · Events · Post CTA
                 Settings · Help Center · Invite neighbors
Centre: filter row · post feed
Right rail: neighbourhood card · invite card · business CTA · ads · chat dock
```

- Left navigation ~243px at 1440 (~17%); feed ~553px (~38–39%); right rail ~271px (~19%).
- The left nav is the product's spine; the feed is the central reading column; the rest is
  margin. The column does not widen to fill when the rail is absent.
- The top bar is fixed (~82px); search is centred in the viewport, not in the column.
- Active nav = light grey fill + filled icon, not green.
- The right rail is utility and monetisation, not a second feed.
- Search is in the top bar; the feed has a separate “Filter by” control.
- The primary action is the green **Post** pill in the left navigation and at the top of
  the feed.

Sources: [web home](https://mobbin.com/screens/03d5f316-fcd3-4389-9cc0-a5406574beca),
[web home variant](https://mobbin.com/screens/50649f88-3142-465e-ab6d-037207fd3f78),
[community rail](https://mobbin.com/screens/efe11897-4db8-46ec-b5f1-e6d8fa164e53).

### Mobile (observed)

```text
Status/header: [compass] Lower Allston ⌄ · bell · chat · avatar
Feed: onboarding module (optional) · post cards · Today's local news rail
Tab bar: Home · Search · For Sale · Faves · raised green Post FAB
```

- The tab bar is the spine; the composer is a raised green action at the right.
- Post detail, comments and composer are pushed screens or sheets.
- The feed keeps the same identity grammar as web.
- Dark mode is supported: black canvas, dark cards, green active state.

Sources: [iOS home](https://mobbin.com/screens/251e08e3-4688-4ea7-be88-557ec9871142),
[For Sale & Free + local news](https://mobbin.com/screens/54cefa73-8bc8-4c81-add3-ed34b1f570fe),
[dark mode](https://mobbin.com/screens/d178d1a1-d1fe-4cc4-b19e-24dc5f6fd108),
[group feed](https://mobbin.com/screens/9fa1c528-5530-4f11-bc24-0244ea7279fe),
[search](https://mobbin.com/screens/6d9b9a65-a2bd-4f67-b0e4-4e7ba3f0068a),
[composer flow](https://mobbin.com/flows/e258bc98-834d-4cce-b62a-cbc96c477653?tab=screens).

## 4. Screen composition

The complete captured screen list is in `nextdoor-screen-inventory.md`.

### 4.1 Web home

- Fixed top bar: brand at left, search pill centred in the viewport, then bell, chat and
  avatar at right. Height ~82px.
- Left navigation (~243px at 1440): icon + label rows, 40–44px tall. The active item uses a
  **light grey fill with a filled icon**, not a green block. A green **Post** pill sits in
  the middle group; Settings/Help/Invite are plain text at the bottom.
- Central reading column ~553px at 1440 (38–39% of the viewport): white cards on a light
  canvas, 12px radius, 1px hairline, no resting shadow, 20px padding, 12px gap.
- Feed header: “Filter by” pill on the left, green **+ Post** pill on the right.
- Post cards: identity row → body → attachment → actions → quiet utility footer.
- Right rail ~271px (19%): 2–3 short utility cards, each one idea + one action; plus the
  chat dock/ads in the reference.
- Sponsored cards reuse the editorial geometry and add a small “Sponsored” line. This is
  exactly the ambiguity maulid.id must avoid.
- The column does not widen to fill the viewport when the rail is absent; the layout keeps
  the reading measure.

### 4.2 Web post detail

**Capture status:** captured as a modal overlay; the publisher version is a split-pane
overlay.

- Clicking a post opens a centred white modal (~16px radius) over the dimmed feed. It holds
  identity, body, attachment, action row, comments and a pinned comment composer.
- The empty comment state is a two-line block (“No comments” / “Be the first to comment!”).
- The publisher version is a split pane: media viewer left (~58%), 380–400px reading pane
  right with publisher identity, headline, standfirst, navy “Read article” CTA, reactions
  and comments.
- The close × floats outside the panel at the top right.
- This keeps the origin page visible and gives back-navigation for free.

### 4.3 Web groups and events

**Capture status:** captured.

- Groups directory: two columns — “Groups near you” list + “Your groups” rail with Create.
  Rows are square avatar, name, Join/Request/Joined pill, right-aligned member count,
  hairline separators.
- Group page: banner, name + overflow, description clamp + See more, Invite split-button,
  member count; optional Group location card; “Welcome”/“Finish setup” horizontal carousel;
  progress bar; inline composer; pinned post.
- Events: wide container (~838px at 1440), 3-up card grid. Card = 4:3 image with a date
  block overlaid bottom-left, 2-line title, time · venue, full-width outline RSVP pill that
  expands to Going/Interested.
- Your events: Host / Going / Interested / Past / Invited underline tabs.
- Empty: line illustration + headline + one muted line + one dark pill CTA.

### 4.4 Web search

**Capture status:** captured.

- Back chevron + horizontally scrollable tabs: All, Posts, Businesses, Pages, For Sale &
  Free, Event, Neighbors.
- “All” is sectioned by content type: rows for events/people, a 3-up grid for marketplace.
- A single-type tab is one long list of the same row.
- Result rows are full-width white cards with hairline borders and ~10px gaps.
- Empty results were not captured; use an Indonesian recovery state.

### 4.5 Web composer

**Capture status:** captured (composer entry, popover, listing attachment, AI suggestion).

- Composer card: author row, prompt, growing textarea, attachment zone, tool row.
- Audience selector pill (globe/lock/people + label + chevron) and Post pill.
- “Create something” popover: Sell or give away, Create an event, Poll your neighbors.
- The listing/event/poll attachments are removable inline cards.
- For maulid.id: one content type, author status visible, body, rules, publish. No
  marketplace/poll intent cards until those domains exist.

### 4.6 Mobile home

- Header: `[compass] Lower Allston ⌄ · bell · chat · avatar`.
- Feed: one column; post card identity row, body, actions, event/image/poll/listing
  attachment.
- Optional “Get started on Nextdoor” module above the feed: horizontal growth cards and a
  progress line (“1 of 7 steps completed”).
- A “Today's local news” rail with publisher lockups, a “See more” pill and a peeking next
  card.
- Bottom tab bar with labels `Home`, `Search`, `For Sale`, `Faves` and a raised green
  `Post` FAB.
- The first screen shows the feed, not a dashboard.
- Dark mode uses a black canvas and dark cards with the same grammar.

**Post variants observed:** text, image, event (large date block), poll (options with vote
percentages), For Sale listing, local-news publisher card, sponsored card.

### 4.7 Mobile post detail and comments

**Capture status:** captured (iOS post detail, comments, reaction tray, report).

- Pushed screen with a back control; post first, then comments.
- Comment rows repeat author, time and body; the composer is pinned above the keyboard.
- Authority comments can be separated/pinned (staff/author) rather than mixed into the
  general thread.
- Reaction tray is six emojis above the composer; the reaction summary is emoji chips plus
  a count, not a single like total.

### 4.8 Mobile composer

- Author identity and audience chip at the top.
- Intent cards can reduce the blank-page problem, but maulid.id has only one content type.
- Rules line near the publish button.
- Unvalidated accounts see an explanation instead of a dead publish button.

### 4.9 Mobile groups/events

**Capture status:** captured (groups directory, group page, events index, event detail,
RSVP, create event).

- Groups: identity header, member count, feed, events; joining swaps Join for
  Invite/video/compose and inserts a welcome checklist.
- Events: date/time/place, organiser, attendee stack and one RSVP action; RSVP is a state
  selector (`Interested` → `Going`).
- Maulid-season agenda: date gutter, event title, venue, WIB time, organiser.

### 4.10 Profiles, settings, notifications

**Capture status:** iOS notifications, settings and profile are captured; web equivalents are
not.

- Profile: cover, overlapping avatar, name/location, Edit profile, Dashboard (self only),
  Faves, Groups, Posts; new-neighbour system post with Add an intro.
- Settings: Account / Notifications / Privacy tabs; read-only identity, password, data
  download, deactivate.
- Notifications: Neighborhood / My activity / Alerts; Today / Last 7 days; activity rows.
- Notification preferences: category list → Allow / Email / Push toggles.
- Privacy: four-level audience ladder.
- maulid.id defers notifications, messages, groups and profiles until the API supports
  them. Do not render empty shells.

## 5. Component model

The detailed component contracts are in `nextdoor-components.md`. The core set:

- app shell (web three zones / mobile tab bar);
- top bar (brand, search, account actions);
- left navigation (destinations, primary action, secondary links);
- feed header (title, search, filter, primary action);
- post card (identity, body, attachment, actions, utility footer);
- publisher/news card (publisher lockup, headline, image, reaction/comment);
- event card (date block, title, time, place, action);
- composer (identity, status, body, rules, publish);
- account status / validation checklist;
- comment row and comment composer;
- group card/row; event row; search result row;
- right-rail card;
- buttons, inputs, selects, tabs/chips, avatars, badges, menus, sheets, toasts;
- empty/loading/error/permission states.

## 6. Design principles

The 21 principles are in `nextdoor-principles.md`. The short version:

1. Local context before content.
2. The feed is the product.
3. Identity before content.
4. Trust is plain language, not decoration.
5. One primary action per surface.
6. Stack, do not grid.
7. Progressive disclosure.
8. The right rail is utility.
9. Honest states.
10. Calm utility, not engagement mechanics.
11. Same grammar across web and mobile.
12. Accessibility is structural.
13. Monetisation is separate or absent.
14. Indonesian by default.

## 7. Visual language

The full reference is in `nextdoor-visual-language.md`. Key observations from the 1920×1320
web capture:

- Canvas `#F8F9FA`; white surfaces; ink `#2F2F2F`; hairline border `#EBEDED`; primary green
  `#188050`; dark CTA `#222E46`.
- Sans-serif UI with hierarchy from size/weight; 12–16px metadata/body; no decorative
  serif in the app chrome.
- 12px card radius, pill controls, circle avatars.
- Three-zone web layout: left nav ~17%, feed ~38–39%, right rail ~19%; the rest is margin.
- Small card gaps and modest shadows; structure comes from borders and canvas contrast.

maulid.id keeps **EB Garamond** for its editorial voice and introduces a **UI sans** for
controls/metadata. It uses its own deeper green, not Nextdoor's exact brand green.

## 8. Trust, validation and moderation

### The identity grammar

Every post, web and iOS, follows:

```text
avatar → name (+ role label when applicable) → place/scope · time · audience → body
```

Nothing sits between the reader and the author's identity, and nothing sits between identity
and the message.

### How Nextdoor carries trust

1. **Implicit, universal trust:** signup requires an address. The address verifies the person
   to the neighbourhood, but the resulting profile does **not** display a “verified” badge.
   Ordinary neighbours show only name, neighbourhood and time.
2. **Explicit, role-based trust:** `Local publisher`, `Admin`, public-agency badge, business
   check + rating, `Sponsored`. These appear only on the content types that need them.
3. **Behavioural proof:** real follower/post/member counts, real places, join announcements,
   business awards. These are visible but not a reputation score for people.
4. **Provenance on derived content:** an AI summary is labelled and lists its sources.

This is a deliberate split. A badge on every neighbour would be noise; a badge on nobody
would make publisher content indistinguishable. Nextdoor chooses an address gate for
everyone and visible role labels for institutions.

### Validation is a resumable checklist, not a mystery

The onboarding flow runs: account → code → name → address → verified payoff → interests →
recommendations → guidelines → first post. Business verification is a separate, named,
resumable task:

- pre-state: “Build trust. Verify your business” with the value proposition;
- in-flight: “Thank you for submitting verification request.” plus what happens next;
- post-state: the checklist item gets a tick and the banner changes;
- no ambiguous “pending” limbo.

The payoff screen says “You're verified!” and immediately shows what is now available.

### maulid.id mapping

- **Implicit trust:** the account exists and has a staff-reviewed state.
- **Explicit trust:** the byline string `Sudah divalidasi` or `Dikelola staf maulid.id`.
- **No badge:** do not render a green check, trophy, shield or “verified” seal. Nextdoor's
  strongest signal is a role marker, not a quality claim; copying its shape would import a
  meaning maulid.id cannot back.
- **Account states:** `Belum divalidasi` → `Menunggu peninjauan staf` → `Sudah divalidasi`,
  each with one sentence and the next step. Do not show a fake progress bar or KYC checklist.
- **Payoff:** after validation, say what is now possible (`Tulis kabar`), not a generic
  congratulations.
- **Do not add:** country/street fields, map pickers, address autocomplete, location
  permission or “at home” checks. maulid.id has no geospatial model.

### Moderation

Nextdoor's moderation surface is the overflow menu, with consequence-first items:

```text
Mute    Hide posts from this neighbour in the feed
Block   Prevent this neighbour from seeing what you share or interacting with you
Report  Flag for review
```

Own-post menu: turn off notifications, bookmark, delete, close discussion, close poll.
Comment menu: edit, delete. Ad menu: report, hide, about ads. Block has a three-bullet
confirmation and a persistent “You have blocked…” state with Unblock. Guidelines appear
once in onboarding and again in context.

**maulid.id:**

- Ship `Sembunyikan` (client-side hide) when useful; it needs no backend.
- Ship `Laporkan` only when reports are stored and staff can review them. Until then, do not
  render the item or a dead menu.
- `Blokir` is out of scope at launch (no DMs/follows); do not copy it for symmetry.
- Publish `Aturan komunitas` as a readable page with four short rules and link it from the
  composer footer.
- Do not expose a report reason taxonomy unless the moderation workflow can consume it.

## 9. Monetisation (observed, rejected)

Nextdoor monetises the same reading column. Observed formats:

| Format | How it appears | CTA |
| --- | --- | --- |
| In-feed business post | business logo + name, `Sponsored` where the neighbourhood would be, same card geometry | `Schedule Service Now` / `Contact us` |
| In-feed display ad | brand + `Sponsored`, headline, body, creative, footer bar | `Shop Now` |
| Utility/safety ad | `PG&E Corp.` + `Sponsored`, safety copy, video still | `Stay safe` / `Learn more` |
| Full-column ad | creative occupies the whole reading column | implicit |
| Rail ad | small logo tile + `Sponsored` + copy; tiny `Ad` marker at the card edge | `Book now` / `Learn more` |
| Carousel ad | two creatives with peek, brand header, `Sponsored` | `Learn more` |
| Video ad | autoplay with pause control and progress bar | `Learn more` |
| Sponsored testimonial | first-person “Hi neighbour!” copy that mimics a neighbour post | `Learn more` / `Switch Now` |

The label is honest, but the sponsored card occupies the same geometry, the same author
line and sometimes the same voice as a neighbour post. The `Ad` marker on rail units is easy
to miss. For a trust-first community, this is the most dangerous pattern to import: the harm
is not the ad, it is that readers learn to distrust the author line.

**maulid.id decision:** no ads, sponsored cards, local deals, business pages, paywalls or
promo countdowns at launch. If monetisation is ever revisited, it requires a separately named
surface, separate geometry, no author byline, an explicit `Berbayar` label at the top of the
card, and a published policy before the first unit ships.

## 10. Adaptation and boundaries

The detailed mapping is in `../.agents/skills/nextdoor-design/references/adaptation.md`.

Adopt now: web three-zone shell, feed card grammar, feed header/search, composer gate,
right-rail utility, mobile tab bar, honest states, UI sans + EB Garamond pairing.

Prepare later: title/image/category, comments, notifications, groups, events, regional
filter, featured kabar, pinned announcements, moderation — all require schema/API work.

Never: ads, sponsored cards, paywalls, points/streaks/leaderboards, paid verification,
fintech KYC, autoplay, neighbourhood/map mechanics, classifieds, business pages, copied
Nextdoor assets or copy.

## 11. Limitations

- Mobbin is a capture library, not a design-token export. Spacing, font sizes and colours
  are estimated from images unless sampled; sampled values are labelled.
- The web captures are desktop-only (~1920px); tablet/phone web breakpoints are not
  evidenced and must be treated as maulid.id decisions.
- There are no Nextdoor marketing-site sections in Mobbin; do not infer a public website
  palette or layout from these captures.
- Captures include promoted and account-specific states. No Nextdoor seasonal module was
  found; the Maulid module is a maulid.id editorial decision, not a Nextdoor precedent.
- The reference shows what shipped; it does not prove a metric.
- Nextdoor's brand and assets are proprietary. Only patterns and measurements may be
  adapted.
- maulid.id's current `Post` model supports body/author/time/likes only. The richer
  components are documented for future work, not implemented in this slice.

## 12. Where to go next

- Components: `nextdoor-components.md`
- Screen inventory: `nextdoor-screen-inventory.md`
- Content and tone: `nextdoor-content.md`
- Principles: `nextdoor-principles.md`
- Visual language: `nextdoor-visual-language.md`
- Review gate: `ui-review-checklist.md`
- Agent skill: `.agents/skills/nextdoor-design/SKILL.md`
- Raw sources: `.agents/skills/nextdoor-design/references/sources.md`
