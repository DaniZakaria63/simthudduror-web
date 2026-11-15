# Mobbin design references for maulid.id

**Date:** 2026-11-15
**Branch:** `feat/designing`
**Worktree:** `/home/dani/WebstormProjects/maulid.id-feat-designing`
**Tool:** Mobbin MCP (`https://api.mobbin.com/mcp`) through OpenCode
**Task intent sent with every query:** “Research design references for an Islamic community news-sharing app”

## 1. Method and scope

Two research workers ran Mobbin searches in parallel, then the results were reviewed and
reconciled here.

- Screens: iOS and web.
- Flows: publishing, onboarding, verification, moderation.
- Sections: editorial heroes, latest-story lists, community join blocks, footers.
- Every reference below was inspected as an image, not selected from metadata alone.
- Product names and Mobbin URLs are preserved so each source can be reopened.
- Visible observations are separated from interpretation; the “Borrow” lines are
  adaptations for maulid.id, not copies.

The product context used throughout: an Indonesian community news platform for Maulid
Rasulullah; any account can read and interact; only accounts validated by staff can publish;
the public web page and a future mobile app share one API; staff work from a backoffice.

## 2. Availability finding: no direct Islamic or Indonesian product in Mobbin

This negative result shapes the whole report.

| Query | Platform | What came back |
| --- | --- | --- |
| `Umma Muslim community app` | iOS | Insight Timer, WhatsApp, Strava |
| `Muslim Pro prayer app` | iOS | Happier Meditation, MyDyson, Waking Up |
| `Quran.com app home` | iOS | Waking Up, Babbel |
| `Muzz Muslim social app` | iOS | Clubhouse, Numo |
| `Hallow Catholic prayer and community app` | iOS | Open, Insight Timer, QUITTR, Ahead |
| `mosque community website homepage` | Web sections | Shopify Community, Sprout Social, Mural |
| `Indonesian news app home` | iOS | FotMob, Opera |
| `detikcom Indonesian news app` | iOS | Apple News, Finimize, Google |
| `Kompas news app Indonesia` | iOS | South China Morning Post, Bloomberg, Yahoo News |
| `Tribun News app Indonesia` | iOS | Digg, Artsy, FotMob |

Repeated probes for prayer, Quran, mosque, hijri calendar, zakat and named Islamic products
returned zero or unrelated screens. Mobbin also does not carry the major Indonesian news
apps, and Indonesian-language queries are not indexed.

What **does** exist is Islamic content or compliance inside secular products:

- **Nextdoor — Ramadan Mubarak module.** The closest Islamic community-content pattern
  found anywhere in Mobbin.
- **Binance — Sharia Products.** A plain compliance claim with a crescent mark.
- **Snoonu (Qatar) — Halal Certified.** Halal as a filterable attribute, not a theme.
- **Grab / GoPay — Indonesian Ramadan and THR campaigns.** Local seasonal visual language.
- **Careem — Gaza donation flow.** A restrained charity pattern.

**Consequence:** do not claim a like-for-like Islamic app precedent. Use community-news,
local-news and editorial references for structure, and the Islamic-adjacent modules above
for tone and seasonal treatment. The design argument must come from maulid.id's own
material — Maulid news, Indonesian mosque communities, staff validation — not from a
borrowed religious aesthetic.

## 3. Reference observations

### 3.1 Community feeds and local news

**Nextdoor — home feed with local news**
<https://mobbin.com/screens/251e08e3-4688-4ea7-be88-557ec9871142>

- Location is the first control (“Lower Allston”), making the feed feel local before the
  first post.
- A community post can embed a structured event card (date, time, place, thumbnail) inside
  the feed instead of forcing a separate events tab.
- “Today's local news” is a distinct module: publisher logo, name, verification check,
  “Following · 8h”, headline and a 16:9 photo, followed by a “See more” pill.
- The green **Post** action is a raised FAB overlapping the tab bar. Publishing is always
  one tap away, even for an account that cannot publish yet.
- Reaction row uses outlined circular buttons; a low-contrast “Post insights / View”
  utility row sits underneath.

**Borrow:** a provenance strip on every kabar (`author · status · relative time`), a
staff-curated “Kabar pilihan” rail, and a persistent but gated contribution action.

**Nextdoor — news item carrying neighbour commentary**
<https://mobbin.com/flows/f3c0a105-1443-4497-9f1c-20aa376762aa?tab=screens>

- A publisher lockup (`WCVB Channel 5 Boston` with a check, “Local publisher · 23h”), a
  headline and a photo with a domain chip are followed by a divider and a neighbour's
  reaction block.
- A “Sponsored” card reuses the same geometry as editorial content. That is a warning, not
  a pattern to copy.
- The publisher → headline → photo → **community response** stack is the strongest
  structural idea found for making news social without turning it into a comment thread.

**Lex — place-based community feed**
<https://mobbin.com/screens/36bffdb3-8d25-413f-9698-596d80c9c31a>

- The header is a place (“San Francisco, C…”) with search, notifications, messages and
  avatar; the scope is known before the first post.
- Cards carry `name · place · relative time`, a two-line headline, body, small topical chips
  and a right-aligned Chat action.
- Reactions are glyphs without counts; this reads calmer than a score at community scale.
- An attached event block with its own thumbnail and “Join now” shows how to mix
  announcements into a text feed.

**Reddit — source line above the headline**
<https://mobbin.com/screens/5ea420f9-81bb-440c-bd5c-3ee75378c040>

- Each row starts with `subreddit · time · source domain`, then a two-line headline, then a
  small right-side thumbnail.
- Action pills carry large counts, which suit Reddit's scale and would look implausible in
  a young community.
- The source-first micro-line is the cheapest credibility cue for linked or aggregated news.

**Nextdoor — search over community posts**
<https://mobbin.com/screens/6d9b9a65-a2bd-4f67-b0e4-4e7ba3f0068a>

- Search field with tabs (`All / Posts / Businesses / Pages`) and one row of filters
  (`Most Recent`, `Distance`, `All Time`).
- Results are compact text-first cards: avatar, name, place, relative time, body, optional
  document thumbnail, secondary “See all N comments”.
- This is the closest analogue to maulid.id's text search. One filter row is enough; no
  dashboard of controls.

**WhatsApp — Communities**
<https://mobbin.com/screens/919e16e9-a197-4dc0-8f5c-c2411ee16521>

- Communities group channels under a name; each row has an avatar, channel name, last
  message, time and a pin on the announcement row.
- Unread state is a small dot, not a red badge. Hierarchy is only two levels: group →
  channel.
- **Borrow:** treat staff announcements as a distinct pinned “Pengumuman” channel, separate
  from open kabar. It expresses “only validated accounts publish” without a badge farm.

**LinkedIn — organisation post with event card**
<https://mobbin.com/screens/22dc6a41-322d-4dbd-b7c8-c4067d13740b>

- Org name, follower count, age and an event card with cover, date range, explicit `WIB`
  timezone, “Online” line and a “View event” pill.
- Authority is stated in text (`Admin · 2w`) rather than iconography.
- **Borrow:** “Dikelola staf maulid.id · 2 hari” as a byline, and `WIB` on every event time.

### 3.2 Editorial and news feeds

**Apple News — Top Stories and multi-rail home**
<https://mobbin.com/screens/2450aa3b-98e8-45ed-a7bc-afb38e478b73>
<https://mobbin.com/screens/8c36d7b2-c9df-4841-a027-96c19c0791e1>

- Section chips sit under the masthead; the lead story is a large image with publisher
  logo, headline, relative time and a short summary.
- Secondary stories collapse to a thumbnail row with the same source/time grammar.
- Colour-coded rail labels (`Top Stories`, `For You`, `Today's Podcast`) repeat the same
  four-slot rhythm, so density stays scannable.

**Borrow:** one lead story + compact rows, with source-first metadata and one accent colour
per rail.

**Yahoo News Singapore — SEA news list**
<https://mobbin.com/screens/ccc01d9b-1faa-479f-b37e-4de3f2aff929>

- The only Southeast Asian news list Mobbin surfaced.
- Row order: 16px source logo + name + age (`CNN 2d`) → bold two-line headline → grey
  two-line snippet → right thumbnail → overflow.
- Region title is centred in the header; bottom bar is Home / Top stories / Profile.
- **Borrow:** source → headline → dek → thumbnail. Keep Indonesian headlines at two lines
  and truncate the dek, never the headline.

**South China Morning Post — Latest Stories**
<https://mobbin.com/screens/de37e4cb-68ec-4421-89de-bbb62ce49fb1>

- A plain `LATEST STORIES` label, then time and section **above** each headline. The
  headline is the largest element; the image sits between stories.
- Alternating text-only and image-backed stories creates rhythm without a uniform grid.
- This is the closest web-like reading experience found on mobile and fits a text-first feed.

**Bloomberg — Top News**
<https://mobbin.com/screens/cb43f713-ac4a-4534-a687-4d6a8bacbc1d>

- A section kicker (`MARKETS`, `POLITICS`) precedes the headline; time and read time follow.
- Lead story full-width; secondary stories form a two-column block, then a compact row.
- Audio and bookmark actions are icon-only and secondary.

**Yahoo News — Top stories**
<https://mobbin.com/screens/d255761c-9756-4152-a3d8-3d6048663b88>

- “Top story” label, headline, then engagement metadata (“34,520 reads · 91 articles”).
- Following rows use a thumbnail, headline and the same metadata pattern.
- Read counts can reward sensation; if maulid.id shows counts, they must describe real
  activity, not a score.

**FotMob — dense sports news**
<https://mobbin.com/screens/bfd98a17-3d48-4413-92e8-d7c68067079e>

- `News` title, `For you / Latest / Transfers / Leagues` tabs, “Trending” label.
- Lead story is a full-width image; following stories are thumbnail rows.
- A source tag and relative time sit under each headline.
- **Borrow:** the lead + list rhythm is the most reusable feed shape for maulid.id.

**Opera — news tab**
<https://mobbin.com/screens/07272a4c-bc23-4b70-986b-632a7fbe26bf>

- Plain text tabs (`For you`, `Headlines`, `Politics`, `Sports`) with an underline active
  state.
- Large story cards show source and relative time under the headline.
- A section-level “See More” keeps the first screen editorial rather than endless.

**Particle News — card-deck home**
<https://mobbin.com/screens/a44dd637-6e5c-4477-9d5f-76223aa7847d>

- Oversized greeting over a deep-blue hero card: “Popular” pill, `13 ARTICLES · 4 HR AGO`,
  three-line headline, one-line deck, “2 questions” footer.
- The next rail is a labelled row with a two-line deck, so hierarchy is one hero + compact
  rows, not a uniform grid.
- **Borrow:** allow one kabar per day to become the lead card with a status line
  (“12 kabar · 4 jam lalu”), the rest compact.

**Formula 1 — news list and latest video**
<https://mobbin.com/screens/46a4c1f9-86e8-4b38-a68f-971091077477>

- Category label above the headline, thumbnail on the left, timestamp below.
- A centred “View latest news” pill separates the lead story from the archive.
- Articles and video share one row anatomy — useful if kabar eventually carries recordings.

**Digg — feed layout choice**
<https://mobbin.com/screens/1b39801a-55cd-4115-afd9-57fe07af037a>

- A sheet offers `Comfy` (larger images, more detail) and `Compact` (more posts, smaller
  images). Density is a user preference.
- `My Feed / All Digg` separates personal and global views.
- **Borrow:** “Semua kabar” now; a following view later if the API grows one.

**Finimize — news card with author and read time**
<https://mobbin.com/screens/e26e811a-a2a7-412e-bd9c-7765498ef92e>

- `NEWS` kicker, large image, headline, author avatar/name, relative time, “2 min read”.
- Bookmark and play are the only row actions; no like/comment clutter.
- **Borrow:** the author block is the trust unit; a reading-time estimate is cheap and
  useful.

### 3.3 Article detail and reading

**Yahoo News — article detail**
<https://mobbin.com/screens/8dbd4068-1bdb-4da4-8813-0ae057d80a67>

- Order: publisher lockup → three-line headline → byline → absolute date and time zone →
  read time; bookmark and overflow top right.
- Share sheet includes a preview card with thumbnail, truncated title and domain. That
  preview is how the article is re-seen in WhatsApp groups.
- **Borrow:** design the share preview deliberately; title and image must survive cropping.

**HYPE — compact editorial header**
<https://mobbin.com/screens/c059b3ff-8855-47f3-9dcc-1b74d20f4858>

- Four-line header: small-caps kicker → three-line headline → grey dek → source row with
  logo, author, relative time.
- Full-width 4:5 image, then reaction, comment, share, bookmark.
- **Borrow:** kicker becomes category; source row becomes `author · Sudah divalidasi · time`.

**theScore — byline and image credit**
<https://mobbin.com/screens/b7676dd5-d5a3-4bf5-a723-f0728ef6d976>

- Headline, circular avatar, author, relative time, read time, right-aligned image credit.
- In a product with contributed photos, visible credit is both editorial duty and trust.

**The Atlantic — mobile article**
<https://mobbin.com/screens/0a960744-cc1d-44fb-a07f-a5ccf42e0d12>

- Small-caps kicker, serif headline, tiny grey photo credit directly under the image.
- **Borrow:** the credit line; it is missing from most community products and matters when
  contributors upload photographs.

**Artsy — minimal news list**
<https://mobbin.com/screens/287e2490-24ca-4ebe-81ae-1745748b91b5>

- Image, headline, author and date; no excerpt, tags or card borders.
- Useful when updates are short rather than long articles.

**Citizen — developing-story timeline** (secondary)
<https://mobbin.com/screens/97df1830-c5e3-411a-a6dc-3cf22c4bf411>

- Full-bleed photo, then a timestamped list of updates with text and optional audio.
- Only the chronological update list is transferable for multi-day Maulid coverage
  (“Pembaruan 14:20”). The red/black incident styling is not.

### 3.4 Events and Maulid-season programming

**Meetup — group and event home**
<https://mobbin.com/screens/6bd6f052-9f86-4111-b6cf-59e35c3085fb>
<https://mobbin.com/screens/b99188cb-7a5a-4d08-8bb5-fe86cd4cba50>

- “New event” / “New group” chips under the greeting, “Incoming” invitation cards, then
  groups.
- Event cards use a compact image, date/time, place and overflow action.
- The events list shows avatar stacks, “160 going” and an “Online” chip; RSVP and share sit
  at the right edge.
- **Borrow:** “N hadir” with an avatar stack under each kegiatan row. Social proof about
  attendance, not popularity.

**Eventbrite — group page**
<https://mobbin.com/screens/582b5426-fac7-41eb-86bc-05d178af3c10>

- `Events / Collections / About` tabs separate the group's identity from its activity.
- Event rows keep a square image, title, date/time and save/share; past events are grouped.

**Luma — upcoming/past and event web hero**
<https://mobbin.com/screens/d57a0760-5c3a-4d4c-93fe-5037f0d7b071>
<https://mobbin.com/screens/73073983-e4eb-4274-a6b6-63305d3873ca>
<https://mobbin.com/screens/ff217eba-443e-4ed5-b4d8-fae280bbfe6b>

- Segmented `Upcoming / Past` control and date-grouped rows keep a long history readable.
- Date headers in grey caps, square thumbnail, organizer avatar/name, right-aligned price,
  then title, clock line and pin line in a muted column.
- Web hero states the local time range (“5–10 PM GMT+8”) and starts the events list
  immediately below.
- The blank area below a short list is honest: no filler cards, no invented recommendations.
- **Borrow:** “Waktu Indonesia Barat” and the venue on every event row; precision is trust.

**Discord — community events**
<https://mobbin.com/screens/f88c4c9d-3dbf-4a89-9adf-51f3c16dad5f>

- Event cards show time, host, attendee count and one “Interested” button.
- Community identity (server name, member count) sits above the event list.

**South China Morning Post — News Agenda**
<https://mobbin.com/screens/c9ca1bcb-aac9-478b-935a-908beccd15fe>

- Cream panel, “Key events” checkbox, two date pickers, Apply button.
- Huge day numerals in a left gutter, event headline, subline and date; chevrons instead of
  cards; a “Show more events” pill closes the list.
- **Borrow:** the date-gutter agenda is the best pattern for Maulid-season programming —
  printable, scannable, and able to hold pengajian, lomba and community events.

**Strava — group events** (supporting)
<https://mobbin.com/screens/cb844353-6828-40e9-9aff-bf37ffe1b302>

- Date-chip agenda rows with a clear RSVP path.

### 3.5 Editorial web landing pages and sections

**The New York Times — homepage**
<https://mobbin.com/screens/04f507e7-fbef-41e8-9caa-aee8a86ba505>

- Utility row above the wordmark: search, date line, editions, account.
- Centered masthead with a single-line section nav.
- Three-column grid: lead headline + dek + read time, a large photo, and a right rail with
  a human-scale community story plus smaller items.
- **Borrow:** centered wordmark + date + short nav reads as a publication rather than a
  startup. Scale the masthead down for repeat visits; keep the identity.

**Quartz — dense editorial grid with pill kickers**
<https://mobbin.com/screens/f4c1e482-190b-41f3-ac26-832048c1d524>
<https://mobbin.com/screens/d906bb5c-edc3-4242-96cb-b6759df2ca07>

- Centered wordmark, search, account, filled Subscribe button.
- Columns mix tall brand images and pill kickers (`Business News`, `#Pharma`) above
  headlines, with a tiny “2 min read” on every item.
- **Borrow:** pill kickers for Maulid categories (Sejarah, Tradisi, Komunitas, Lomba) make a
  mixed feed feel organised without tabs.

**Wise Newsroom — Latest news**
<https://mobbin.com/sites/sections/457c75c9-46f5-4c9b-9ab7-b3bd761ac372>

- “Latest news” heading left, search right, one underlined topic row, three cards.
- Card: image, dated meta line, headline, two-line excerpt.
- **Borrow:** this is a complete spec for the public feed header, and it degrades to one
  column on mobile.

**PayPal Newsroom — image-free latest news**
<https://mobbin.com/sites/sections/ef719f12-8489-4ee5-ac8f-819dfc01839d>

- Navy band, white heading, three cards with date first, headline, three-line excerpt and
  “Read more”.
- **Borrow:** the deterministic image-free fallback for contributor posts without a photo.
  Date-first ordering signals recency honestly.

**Ramp — text columns with source attribution**
<https://mobbin.com/sites/sections/37fef8e2-7c40-4c08-9415-af3bbb95d288>

- Four text-only columns: headline, excerpt, muted attribution (`Bloomberg, November 2025`).
- **Borrow:** `outlet, month year` is the right treatment for linked third-party news; it
  sits beside original posts without pretending to be original reporting.

**Dropbox Community — join block with real size signals**
<https://mobbin.com/sites/sections/97a24690-777e-4c78-8e0e-4b11596e4fba>

- Split hero with a short explanation and an illustration; below it a search input, a stat
  row (`1,263,475 members · 534,707 posts · 1,166 online`) and a “Start a discussion”
  button.
- **Borrow:** the honest version of social proof — say how many contributors and posts
  exist. If the numbers are small, still show them; it reads as candour.

**Shopify Community — welcome, search, categories, announcements**
<https://mobbin.com/sites/sections/4543e419-1343-437c-b8d6-8dc9baef424f>

- Welcome statement explains the community's purpose before sign-in.
- Search is the primary action next to “Ask a Question”; four category cards give the
  community a visible information architecture.
- An announcements carousel carries current, staff-curated material.

**Sprout Social Arboretum — welcome and trending**
<https://mobbin.com/sites/sections/cebd58a5-a64e-4118-a052-bb46d88bab2d>

- Hero is a welcome plus search/join prompt, not a marketing slogan.
- “Trending Topics” is a two-column list of dated threads; the gradient is decorative.
- **Borrow:** the useful content starts immediately below the welcome.

**Mural Community — recent discussions and sidebar**
<https://mobbin.com/sites/sections/f5db1f78-fe82-49cf-a7ef-6186b29df74f>

- “Recent Discussions” carries category, replies, views and last activity; the sidebar has
  a featured template, quick links and tags.
- Restrained visual language suits a small staff team better than a heavily branded forum.

**The New Yorker — editorial footer**
<https://mobbin.com/sites/sections/84724b3c-961d-4245-a46a-f176fe4bec67>

- Black background, small white type, generous columns, a circular monogram.
- Two tiers: editorial links (sections, archive, about) then utility/legal links in smaller
  type.
- **Borrow:** two-tier footer with the masthead mark. Skip the long legal paragraph;
  Indonesian products need Privacy, Terms, Contact and little else.

**sweetgreen — mission statement** (supporting)
<https://mobbin.com/sites/sections/1e44b0a4-1609-4ecc-baf0-7296db75c194>

- Sage half-panel with a giant serif statement opposite a photo.
- **Borrow:** the best template for “why contribute” — a statement, not a feature list.

**Duna — compact news list** (supporting)
<https://mobbin.com/sites/sections/33709e60-a59e-4f6f-97ef-9395e466f005>

- 56px square thumbnails, title, `Conversations · Sep 24, 2025`.

**Patreon News hero** (supporting)
<https://mobbin.com/sites/sections/3b0dd719-f5ea-4284-a497-50df22f5ae66>

- Full-bleed photo, small overline, three-line serif headline, small “Read Article” pill.

**beehiiv newsletter section** (supporting)
<https://mobbin.com/sites/sections/53487269-e180-458b-b298-736592e9c836>

- Centered “Join the newsletter”, one line of copy, input + Subscribe, thin rules and a small
  ornament. Only relevant if a newsletter exists; there is no backend for one today.

**Nextdoor web — community sidebar** (supporting)
<https://mobbin.com/screens/efe11897-4db8-46ec-b5f1-e6d8fa164e53>

- Right rail: `Downtown Menlo Park`, `2,262 neighbors`, `49% of 2,835 households`.
- **Borrow:** a small “Komunitas” rail with real totals, not invented statistics.

### 3.6 Islamic-adjacent and Indonesian seasonal patterns

**Nextdoor — Ramadan Mubarak module**
<https://mobbin.com/screens/feb2cdda-699f-4be5-a11a-244612039997>

- The closest Islamic community-content pattern found in Mobbin.
- Centered crescent illustration, bold “Ramadan Mubarak”, one explanatory sentence
  (“Celebrate the communal spirit of Ramadan with your Muslim neighbours…”), one green pill
  “Share good wishes”.
- It is a full-width in-feed card above “Local Deals”: no price, no product, no countdown —
  a calm seasonal notice inside an otherwise commercial feed.
- **Borrow:** a seasonal Maulid band with the same three-part anatomy (mark → greeting → one
  verb pill), in Indonesian, inviting a greeting or doa rather than a purchase.

**Binance — Sharia Products**
<https://mobbin.com/screens/39875189-194c-467c-9f84-4d0f41183c29>

- Green-tinted midnight header, “Earn yields that are certified Sharia Compliant”, crescent
  mark, and the compliance claim in one sentence above a plain list.
- **Borrow:** state a religious assurance plainly; do not decorate it.

**Snoonu (Qatar) — Halal Certified**
<https://mobbin.com/screens/17404f2c-db93-46c7-8489-74aedd0fded6>

- A “Halal Certified” icon row above category chips, plus an English/Arabic switch.
- Halal is a filterable attribute, not the theme.
- **Borrow:** if categories arrive, make them functional attributes, not ornaments.

**Grab — Ramadan gifting and halal search**
<https://mobbin.com/screens/406422a4-9ab8-4388-8dc6-b017ccf3bc72>
<https://mobbin.com/screens/edbfd20c-6206-4109-bf5f-367bf8a825de>

**GoPay — THR theme picker**
<https://mobbin.com/screens/2c69d338-1e68-4ebe-88d9-a8f88063398e>

- Indonesian Ramadan/THR visual language: dark green and gold, batik-adjacent ornament,
  Indonesian headlines.
- **Borrow:** the palette and local idiom, not the commerce. Avoid generic Middle-Eastern
  mosque silhouettes; Indonesian masjid/musala architecture and batik-adjacent geometry are
  the honest vocabulary.

**Careem — Gaza donation flow** (supporting)
<https://mobbin.com/screens/6cec8e7c-0a09-4775-8b45-ceb02a07a9ee>

- Cause title, “Total amount to be donated”, single dark-green payment action.
- Only relevant if infaq/contribution is added later; it is not in the current domain.

**Insight Timer — gratitude wall and community**
<https://mobbin.com/screens/bd16fe51-7b47-4f35-9840-b1e13c28c75e>
<https://mobbin.com/screens/1a977741-4c5a-4c0f-888a-9494086d091a>

- A simple prompt (“What are you grateful for today?”) produces a chronological, named,
  location-tagged feed.
- The community section mixes listening activity with people; the page feels alive without a
  separate social graph.
- A persistent five-label bottom bar; no icon-only navigation.
- **Borrow:** the prompt-plus-chronological-list model for a future doa or reflection
  section, not for the news feed itself.

**QUITTR — community forum inside a habit app**
<https://mobbin.com/screens/1087c060-dab5-4a37-8d41-0c3258a6391c>

- `Forum / Teams / Friends` tabs and search keep a support community navigable.
- Posts show author, relative time, body and a single upvote; a floating compose button is
  the only accent action.

**Open — daily practice and community**
<https://mobbin.com/screens/44b88cf9-7a95-4ead-a63e-a52aeb11c847>

- A daily programme banner establishes a shared rhythm; the community list is a stream of
  small repeated actions with names and timestamps.
- Dark and atmospheric, but the list stays text-first.

### 3.7 Publishing, validation and moderation

The product rule is “only staff-validated users publish”. Trust must be legible **before**
the composer opens, and the composer must explain the rule at the moment of posting.

**Identity and audience inside the composer**

- **Nextdoor:** composer opens with the author's avatar, name and neighbourhood under the
  close button, and an audience chip “Anyone”. Three intent cards (“Sell or give away”,
  “Create an event”, “Poll your neighbors”) replace a blank box.
  <https://mobbin.com/flows/e258bc98-834d-4cce-b62a-cbc96c477653?tab=screens>
- **Runna:** scope as a dropdown chip (`GENERAL`) plus two fields, “Subject title” and “What
  are your thoughts?”.
  <https://mobbin.com/flows/2d97de76-1809-444d-be9b-fc4f02fe0505?tab=screens>
- **Handshake:** audience labelled inline (“People at your school”) and confirmed with a
  “Posted” toast on the feed.
  <https://mobbin.com/flows/5c80547c-d194-4d5d-bddc-ee7dc9205467>
- **X:** audience sheet lists communities with member counts; after selection the composer
  shows the community chip, an “Also share with followers” checkbox, and a rules line with a
  link.
  <https://mobbin.com/flows/99a52f6e-d598-4720-bfc7-4fdcb64d4d0a?tab=screens>
- **Instagram:** explains consequences inline (“People you share this content with can see
  the location…”) rather than hiding them in settings.
  <https://mobbin.com/flows/fe19b5d5-dcd7-4378-960d-e31353d6279f>

**Validation shown as facts, not flair**

- **Airtasker:** a “BADGES” panel lists Police Check, Digital ID, Payment Method, Mobile,
  LinkedIn, each with “Verified on Wednesday, 3 May 2023”.
  <https://mobbin.com/screens/aecc1557-5bb9-40ba-b1a1-69801acf3314>
- **Airbnb:** one linked row, “Identity verified”, under the profile header.
  <https://mobbin.com/screens/0673af7e-b7a2-4c55-a1c9-a02ae004f278>
- **X:** a single “Get verified” pill beside the display name; this is a paid-check
  semantic, not staff validation.
  <https://mobbin.com/screens/04adbd4e-9a4b-4bf4-9027-0ee4b817b4f3>
- **Nextdoor publisher sheet (web):** plain-language “Local publisher” label beside the
  publisher name, with title, summary, reactions and comments.
  <https://mobbin.com/screens/d748f96c-1850-44c9-b0c3-2f38894f97b5>

**Rules published where people can read them**

- **X community About:** numbered rules and a moderator list with roles.
  <https://mobbin.com/screens/e5ffe316-9faa-48c9-bed1-0775a3df5f06>
- **Reddit Rules:** numbered rules with one-line explanations and an “Auto Enforced” tag on
  machine-checked ones.
  <https://mobbin.com/screens/abd49111-29b0-406e-bb80-7645871fe99c>
- **Nextdoor Community Guidelines:** warm opening, filter chips, numbered sections.
  <https://mobbin.com/screens/7546ad42-8bdd-4c1c-894c-5ec39388788c>
- **Azar / Vestiaire:** consent screens with examples and an explicit consequence callout
  (“…can result in a permanent ban”).
  <https://mobbin.com/screens/8d757f53-7677-4f6c-a2df-0d82f3423c16>
  <https://mobbin.com/screens/3873f821-6783-4b6c-8ee3-a5bb8db3c379>

**Moderation and thread authority**

- **Nextdoor overflow:** Hide, Mute, Report, Turn off notifications, Share, Bookmark — each
  with a one-line consequence.
  <https://mobbin.com/flows/b278d4f1-98e7-4915-91c2-264f86cff9c6?tab=screens>
- **Weverse:** “Artist's comments” pinned above “All comments”; authority is separated, not
  just badged.
  <https://mobbin.com/screens/c9514ae9-bf1e-43e4-be2d-d6141035ed3f>
- **Medium:** “Responses (834)” sheet with a sort control and a shield icon.
  <https://mobbin.com/screens/b757ef66-ab07-4e8f-a687-89d07e5b9041>

**Interpretation for maulid.id**

1. Put the account state in the persistent header or account chip: `Belum divalidasi` /
   `Sudah divalidasi`, with the same wording as the API.
2. On every post, show `author · Sudah divalidasi · 2 jam lalu` in plain language. No gold
   check, no trophy, no “verified” blue tick.
3. Open the composer with the author identity and status visible. For an unvalidated user,
   show a disabled publish button plus an explanation and a link to the contribution
   section — not an empty error.
4. Add a short rules line at the bottom of the composer (“Kabar harus benar, sopan, dan
   terkait Maulid…”), matching a readable rules page.
5. Keep moderation in one overflow menu with consequences spelled out; do not expose a
   report button without a review path in the backend.
6. Avoid fintech KYC styling (ID photo, selfie, “2 minutes”). Staff validation by a human is
   the product rule; show the outcome, not a compliance funnel.

## 4. Constraints from the current codebase

The references point toward title-led cards, images, categories and event modules. The
current implementation does not support all of that yet:

- `Post` has `id`, `author_id`, `author_name`, `body`, `created_at`, `like_count` and
  `liked_by_me` only (`web/src/features/posts/types.ts`, `internal/post/post.go`).
- The public feed is a text search over `body`; there is no category, region, media or event
  field (`internal/post/repository_sqlite.go`).
- `PostCard` renders author, date, body paragraphs and a like control; it has no title or
  image slot (`web/src/features/posts/components/PostCard.tsx`).
- The landing page is `Masthead` → `Kabar komunitas` feed → `ContributionSection`
  (`web/src/features/landing/pages/LandingPage.tsx`).
- Only validated accounts can create posts; the API returns `validation_required` otherwise.

| Design move | Fits today? | Needs first |
| --- | --- | --- |
| Lead story + compact rows | Yes, with body excerpts | Nothing |
| `author · status · time` metadata | Yes | Nothing |
| Relative Indonesian time | Yes | A formatter change |
| Reading-time estimate | Yes | Word-count helper |
| Image-free fallback card | Yes | Nothing |
| Search/filter row | Partly | More filter parameters |
| Images, titles, categories | No | Post schema + API + migration |
| Event cards / agenda | No | Event domain + API |
| Regional filter | No | Location field + API |
| “Kabar pilihan” staff curation | No | Featured flag or staff collection |
| Pinned “Pengumuman” channel | No | Announcement type |

Therefore the immediately safe design moves are typographic and structural: lead/latest
rhythm, clearer section labels, better feed metadata, a visible contribution path and an
honest trust explanation. Do not simulate unsupported features with fake content.

## 5. Patterns worth adapting

1. **Two card grammars.** Publication-first for curated/linked news
   (`source · time → headline → dek → thumbnail`); social-first for contributor posts
   (`avatar · author · status · time → body → actions`). Make the difference obvious in the
   meta row, not with colour.
2. **Lead story + compact list.** One full-width lead item, then thumbnail or text rows. It
   gives an editorial page an entry point without a uniform card grid.
3. **Provenance on every item.** Author, plain-language validation state and relative time.
   This is the trust unit for a platform where only validated accounts publish.
4. **Source before headline.** For linked or aggregated news, print the publisher before the
   headline. It is the cheapest credibility cue and costs one line.
5. **Plain text tabs, restrained active states.** `Kabar`, `Kontribusi`, `Akun` do not need
   pill filters; an underline or weight change is enough.
6. **A visible but gated contribution path.** A persistent “Tulis kabar” action that leads
   unvalidated users to an explanation and account status, never a dead end.
7. **Honest empty and error states.** Indonesian copy, one recovery action, optional topic
   suggestions. Gojek's “Let's try another word?” and SCMP's topic fallback are good models:
   <https://mobbin.com/screens/6f0b1c31-fd99-474c-92f8-e1a05ef3eb48>
   <https://mobbin.com/screens/b517ee46-0024-40af-b036-a84002ba3d33>
8. **Deterministic image fallback.** Use an initial-monogram tile or a category colour, not
   a stock photo. PayPal and Ramp show that image-free cards can look finished.
9. **Date-first event rows for Maulid season.** Day gutter, WIB time, venue, organizer and
   attendance count. The SCMP agenda and Luma rows are the strongest models.
10. **A calm seasonal module.** Nextdoor's Ramadan card — mark, greeting, one verb — is the
    right register. No countdown, no discount, no commerce.
11. **Two-tier editorial footer.** Categories and about links first, legal links in smaller
    type, with the masthead mark.
12. **Text-first density for Indonesian mobile data.** Headlines and metadata carry the
    page; images support. Avoid autoplay and carousel-heavy heroes.

## 6. What not to copy

- **Hyper-local US framing** (garage sales, “Lower Allston”, identity-classifieds tags). The
  card anatomy transfers; the subject matter does not.
- **Neighbourhood-social mechanics** (follow requests, DMs, hyperlocal ads) are larger than
  this product's scope and would change the domain.
- **Fintech-grade identity verification** (ID photo, selfie, “2 minutes”). Disproportionate,
  privacy-heavy and culturally off for a religious community. Staff validation is the rule.
- **Blue-check semantics.** A paid verification tick is not “validated by staff”; keep the
  status as text.
- **Vote counts, streaks, points, leaderboards.** Scoring religious participation reads as
  insincere and looks implausible at community scale.
- **Paywalls, member-only chips and subscription interstitials.** A Maulid announcement must
  never be interrupted by a sales sheet.
- **Advertising inside the editorial stream.** Nextdoor mixes sponsored cards into news
  geometry; at launch this would undercut the credibility the product is earning.
- **Ramadan commercial tropes** (crescent + lanterns + discount badges, generic
  Middle-Eastern mosque silhouettes). Use Indonesian masjid/musala architecture and
  batik-adjacent geometry if ornament is needed.
- **Emergency/crime aesthetics.** Citizen's black-and-red incident styling implies danger;
  a Maulid community needs warmth.
- **Mega-navigation and three-column complexity.** NYT-scale information architecture needs
  NYT-scale staff. Keep 3–4 destinations and one reading column.
- **Autoplay video and heavy hero carousels.** Expensive on Indonesian mobile data and
  hostile to low-end Android.
- **English-first microcopy.** All empty, error, status and action strings must be Indonesian
  (“Belum ada kabar”, “Akun belum divalidasi”, “2 jam lalu”).
- **Copied branding, copy, protected layouts or imagery.** Adapt the information architecture
  and interaction pattern; never reproduce a product's visual identity.

## 7. Search quality notes

Broad category queries were unreliable; named-product and behaviour queries returned far
higher signal.

Weak queries not worth repeating:

| Query | Why it failed |
| --- | --- |
| `local events app home feed with nearby events` | Ticketing/commerce apps (DICE, StubHub, Partiful), no community-news feed |
| `verify identity with selfie and document` | 100% fintech KYC; no community-relevant verification |
| `community platform homepage with posts feed` (web) | Mostly SaaS collaboration tools |
| `editorial footer section` | Generic startup footers; only The New Yorker and SSENSE had character |
| `community membership join section` | SaaS marketing CTAs; only Dropbox Community carried real proof |
| `search results empty state in news app` | Half unrelated empty states; only Gojek, SCMP and Deepstash were useful |
| `community app verification badge profile screen` | Dating/social apps dominated; only Airtasker, Airbnb and X transferred |
| `create a community post` (flow) | Reddit/Flo/QUITTR composers, no newsroom-style composer |
| `news app onboarding` (flow) | Logo-then-paywall screens; thin for a read-first product |
| Indonesian-language queries | Not indexed (`berita Indonesia aplikasi beranda` returned 0) |

Worth repeating: `local news website homepage with latest stories`, `news site latest
stories list section`, `publishing a post in a neighborhood app`, `community guidelines and
moderation rules screen`, and named-app queries.

## 8. Open questions for the design direction

- Should the public page remain an editorial landing page, or become a feed-first product
  with the masthead as a compact header?
- Does “Kabar” need categories (kegiatan, kajian, dokumentasi, pengumuman) or is a single
  chronological feed enough for the first release?
- Where should validation status live on the public page: header, contribution section, or
  both?
- Is a Maulid-season agenda a real product requirement, or a later addition? It needs an
  event domain and API before it can be designed honestly.
- Should the mobile API expose a regional filter before the mobile app is designed?
- Which references should become a saved Mobbin collection for the team? (Not saved in this
  research slice.)

## 9. Source index

### Community and local news

| Product | Platform | Mobbin |
| --- | --- | --- |
| Nextdoor home | iOS | <https://mobbin.com/screens/251e08e3-4688-4ea7-be88-557ec9871142> |
| Nextdoor news + commentary | iOS flow | <https://mobbin.com/flows/f3c0a105-1443-4497-9f1c-20aa376762aa> |
| Nextdoor search | iOS | <https://mobbin.com/screens/6d9b9a65-a2bd-4f67-b0e4-4e7ba3f0068a> |
| Nextdoor publisher sheet | Web | <https://mobbin.com/screens/d748f96c-1850-44c9-b0c3-2f38894f97b5> |
| Nextdoor community rail | Web | <https://mobbin.com/screens/efe11897-4db8-46ec-b5f1-e6d8fa164e53> |
| Lex | iOS | <https://mobbin.com/screens/36bffdb3-8d25-413f-9698-596d80c9c31a> |
| Reddit News | iOS | <https://mobbin.com/screens/5ea420f9-81bb-440c-bd5c-3ee75378c040> |
| WhatsApp Communities | iOS | <https://mobbin.com/screens/919e16e9-a197-4dc0-8f5c-c2411ee16521> |
| LinkedIn org post | iOS | <https://mobbin.com/screens/22dc6a41-322d-4dbd-b7c8-c4067d13740b> |

### Editorial feeds and article detail

| Product | Platform | Mobbin |
| --- | --- | --- |
| Apple News Top Stories | iOS | <https://mobbin.com/screens/2450aa3b-98e8-45ed-a7bc-afb38e478b73> |
| Apple News rails | iOS | <https://mobbin.com/screens/8c36d7b2-c9df-4841-a027-96c19c0791e1> |
| Yahoo News Singapore | iOS | <https://mobbin.com/screens/ccc01d9b-1faa-479f-b37e-4de3f2aff929> |
| South China Morning Post | iOS | <https://mobbin.com/screens/de37e4cb-68ec-4421-89de-bbb62ce49fb1> |
| Bloomberg | iOS | <https://mobbin.com/screens/cb43f713-ac4a-4534-a687-4d6a8bacbc1d> |
| Yahoo News top stories | iOS | <https://mobbin.com/screens/d255761c-9756-4152-a3d8-3d6048663b88> |
| Yahoo News article | iOS | <https://mobbin.com/screens/8dbd4068-1bdb-4da4-8813-0ae057d80a67> |
| FotMob news | iOS | <https://mobbin.com/screens/bfd98a17-3d48-4413-92e8-d7c68067079e> |
| Opera news | iOS | <https://mobbin.com/screens/07272a4c-bc23-4b70-986b-632a7fbe26bf> |
| Particle News | iOS | <https://mobbin.com/screens/a44dd637-6e5c-4477-9d5f-76223aa7847d> |
| Formula 1 news | iOS | <https://mobbin.com/screens/46a4c1f9-86e8-4b38-a68f-971091077477> |
| Digg feed | iOS | <https://mobbin.com/screens/1b39801a-55cd-4115-afd9-57fe07af037a> |
| Finimize | iOS | <https://mobbin.com/screens/e26e811a-a2a7-412e-bd9c-7765498ef92e> |
| HYPE article | iOS | <https://mobbin.com/screens/c059b3ff-8855-47f3-9dcc-1b74d20f4858> |
| theScore article | iOS | <https://mobbin.com/screens/b7676dd5-d5a3-4bf5-a723-f0728ef6d976> |
| The Atlantic article | iOS | <https://mobbin.com/screens/0a960744-cc1d-44fb-a07f-a5ccf42e0d12> |
| Artsy news | iOS | <https://mobbin.com/screens/287e2490-24ca-4ebe-81ae-1745748b91b5> |
| Citizen timeline | iOS | <https://mobbin.com/screens/97df1830-c5e3-411a-a6dc-3cf22c4bf411> |

### Events

| Product | Platform | Mobbin |
| --- | --- | --- |
| Meetup home | iOS | <https://mobbin.com/screens/6bd6f052-9f86-4111-b6cf-59e35c3085fb> |
| Meetup events | iOS | <https://mobbin.com/screens/b99188cb-7a5a-4d08-8bb5-fe86cd4cba50> |
| Eventbrite group | iOS | <https://mobbin.com/screens/582b5426-fac7-41eb-86bc-05d178af3c10> |
| Luma upcoming/past | iOS | <https://mobbin.com/screens/d57a0760-5c3a-4d4c-93fe-5037f0d7b071> |
| Luma Discover | iOS | <https://mobbin.com/screens/73073983-e4eb-4274-a6b6-63305d3873ca> |
| Luma web hero | Web | <https://mobbin.com/screens/ff217eba-443e-4ed5-b4d8-fae280bbfe6b> |
| Discord events | iOS | <https://mobbin.com/screens/f88c4c9d-3dbf-4a89-9adf-51f3c16dad5f> |
| SCMP News Agenda | iOS | <https://mobbin.com/screens/c9ca1bcb-aac9-478b-935a-908beccd15fe> |
| Strava group events | iOS | <https://mobbin.com/screens/cb844353-6828-40e9-9aff-bf37ffe1b302> |

### Web and sections

| Product | Platform | Mobbin |
| --- | --- | --- |
| The New York Times | Web | <https://mobbin.com/screens/04f507e7-fbef-41e8-9caa-aee8a86ba505> |
| Quartz | Web | <https://mobbin.com/screens/f4c1e482-190b-41f3-ac26-832048c1d524> |
| Quartz dark home | Web | <https://mobbin.com/screens/d906bb5c-edc3-4242-96cb-b6759df2ca07> |
| Wise Latest news | Section | <https://mobbin.com/sites/sections/457c75c9-46f5-4c9b-9ab7-b3bd761ac372> |
| PayPal Latest News | Section | <https://mobbin.com/sites/sections/ef719f12-8489-4ee5-ac8f-819dfc01839d> |
| Ramp Latest news | Section | <https://mobbin.com/sites/sections/37fef8e2-7c40-4c08-9415-af3bbb95d288> |
| Dropbox Community join | Section | <https://mobbin.com/sites/sections/97a24690-777e-4c78-8e0e-4b11596e4fba> |
| Shopify Community | Section | <https://mobbin.com/sites/sections/4543e419-1343-437c-b8d6-8dc9baef424f> |
| Sprout Social Arboretum | Section | <https://mobbin.com/sites/sections/cebd58a5-a64e-4118-a052-bb46d88bab2d> |
| Mural Community | Section | <https://mobbin.com/sites/sections/f5db1f78-fe82-49cf-a7ef-6186b29df74f> |
| The New Yorker footer | Section | <https://mobbin.com/sites/sections/84724b3c-961d-4245-a46a-f176fe4bec67> |
| sweetgreen mission | Section | <https://mobbin.com/sites/sections/1e44b0a4-1609-4ecc-baf0-7296db75c194> |
| Duna news list | Section | <https://mobbin.com/sites/sections/33709e60-a59e-4f6f-97ef-9395e466f005> |
| Patreon News hero | Section | <https://mobbin.com/sites/sections/3b0dd719-f5ea-4284-a497-50df22f5ae66> |
| beehiiv newsletter | Section | <https://mobbin.com/sites/sections/53487269-e180-458b-b298-736592e9c836> |

### Islamic-adjacent / Indonesian seasonal

| Product | Platform | Mobbin |
| --- | --- | --- |
| Nextdoor Ramadan | iOS | <https://mobbin.com/screens/feb2cdda-699f-4be5-a11a-244612039997> |
| Binance Sharia | iOS | <https://mobbin.com/screens/39875189-194c-467c-9f84-4d0f41183c29> |
| Snoonu halal | iOS | <https://mobbin.com/screens/17404f2c-db93-46c7-8489-74aedd0fded6> |
| Grab Ramadan | iOS | <https://mobbin.com/screens/406422a4-9ab8-4388-8dc6-b017ccf3bc72> |
| Grab halal search | iOS | <https://mobbin.com/screens/edbfd20c-6206-4109-bf5f-367bf8a825de> |
| GoPay THR | iOS | <https://mobbin.com/screens/2c69d338-1e68-4ebe-88d9-a8f88063398e> |
| Careem donation | iOS | <https://mobbin.com/screens/6cec8e7c-0a09-4775-8b45-ceb02a07a9ee> |
| Insight Timer | iOS | <https://mobbin.com/screens/bd16fe51-7b47-4f35-9840-b1e13c28c75e> |
| Insight Timer member warmth | iOS | <https://mobbin.com/screens/1a977741-4c5a-4c0f-888a-9494086d091a> |
| QUITTR | iOS | <https://mobbin.com/screens/1087c060-dab5-4a37-8d41-0c3258a6391c> |
| Open | iOS | <https://mobbin.com/screens/44b88cf9-7a95-4ead-a63e-a52aeb11c847> |

### Publishing, validation and moderation

| Product | Platform | Mobbin |
| --- | --- | --- |
| Nextdoor composer | iOS flow | <https://mobbin.com/flows/e258bc98-834d-4cce-b62a-cbc96c477653> |
| Runna composer | iOS flow | <https://mobbin.com/flows/2d97de76-1809-444d-be9b-fc4f02fe0505> |
| Handshake composer | iOS flow | <https://mobbin.com/flows/5c80547c-d194-4d5d-bddc-ee7dc9205467> |
| X community composer | iOS flow | <https://mobbin.com/flows/99a52f6e-d598-4720-bfc7-4fdcb64d4d0a> |
| Instagram composer | iOS flow | <https://mobbin.com/flows/fe19b5d5-dcd7-4378-960d-e31353d6279f> |
| Airtasker badges | iOS | <https://mobbin.com/screens/aecc1557-5bb9-40ba-b1a1-69801acf3314> |
| Airbnb identity | iOS | <https://mobbin.com/screens/0673af7e-b7a2-4c55-a1c9-a02ae004f278> |
| X verified | iOS | <https://mobbin.com/screens/04adbd4e-9a4b-4bf4-9027-0ee4b817b4f3> |
| X community About | iOS | <https://mobbin.com/screens/e5ffe316-9faa-48c9-bed1-0775a3df5f06> |
| Reddit rules | iOS | <https://mobbin.com/screens/abd49111-29b0-406e-bb80-7645871fe99c> |
| Nextdoor guidelines | iOS | <https://mobbin.com/screens/7546ad42-8bdd-4c1c-894c-5ec39388788c> |
| Azar rules | iOS | <https://mobbin.com/screens/8d757f53-7677-4f6c-a2df-0d82f3423c16> |
| Vestiaire rules | iOS | <https://mobbin.com/screens/3873f821-6783-4b6c-8ee3-a5bb8db3c379> |
| Nextdoor overflow | iOS flow | <https://mobbin.com/flows/b278d4f1-98e7-4915-91c2-264f86cff9c6> |
| Weverse comments | iOS | <https://mobbin.com/screens/c9514ae9-bf1e-43e4-be2d-d6141035ed3f> |
| Medium responses | iOS | <https://mobbin.com/screens/b757ef66-ab07-4e8f-a687-89d07e5b9041> |

### Empty states

| Product | Platform | Mobbin |
| --- | --- | --- |
| Gojek no results | iOS | <https://mobbin.com/screens/6f0b1c31-fd99-474c-92f8-e1a05ef3eb48> |
| SCMP no results | iOS | <https://mobbin.com/screens/b517ee46-0024-40af-b036-a84002ba3d33> |
| Deepstash no results | iOS | <https://mobbin.com/screens/b537fc59-dd61-4293-bf46-70257fbac79e> |

Mobbin captures show that these products used a pattern. They do not prove the pattern
improved retention, conversion or trust. Accessibility and performance must be verified in
maulid.id's own implementation, and no protected design, branding, copy or asset may be
reproduced.
