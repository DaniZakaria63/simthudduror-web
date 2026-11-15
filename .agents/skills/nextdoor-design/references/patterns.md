# Interaction patterns

Each pattern states the Nextdoor evidence, the underlying rule, and the maulid.id
adaptation.

## 1. Local feed framing

**Evidence:** Web and iOS headers place the neighbourhood before the feed; the feed is a
stack of posts from people in that place.

**Rule:** State the scope before the content. The reader should never wonder “whose feed is
this?”.

**maulid.id:** “Kabar komunitas” is the scope. A region filter can be added later; do not
invent a neighbourhood name.

## 2. Post identity row

**Evidence:** `[avatar] Name · Neighbourhood · time · audience icon` sits above the body.
Publisher content replaces the neighbourhood with a publisher lockup and verification
label.

**Rule:** Identity, provenance and time come before the message.

**maulid.id:** `[initial] Author · Sudah divalidasi · 2 jam lalu`. For staff: `Dikelola staf
maulid.id · 2 jam lalu`. The avatar falls back to an initial on the primary green.

## 3. Post body and attachment

**Evidence:** Text post → body; event post → body + event card with date block, title,
time, place, thumbnail; image post → full-width image with optional domain chip; sponsored
post → same geometry as editorial with a “Sponsored” label.

**Rule:** The attachment is part of the post, not a separate card floating in the feed.
One attachment per post is enough.

**maulid.id:** Text only today. Design the attachment slot for a future image or event, but
do not render an empty frame or a fake image.

## 4. Reactions and actions

**Evidence:** Heart, comment, share and overflow are outlined glyphs; “Post insights / View”
is a quiet utility row. Counts are small and secondary. Some cards use pills, some use
icon-only actions.

**Rule:** Actions are quiet, consistent and ordered by importance: like → comment → share →
overflow. Counts are real or omitted.

**maulid.id:** Like is the only interaction the API supports today. Comment/share stay
hidden until their APIs exist. Do not show a dead icon.

## 5. Composer gate

**Evidence:** The composer opens with the author identity and audience chip visible. Intent
cards reduce the blank-page problem. Rules are surfaced at the moment of posting.

**Rule:** Show who is posting, where it will go, and what the rules are before the publish
button.

**maulid.id:** Author + status at top; body field; rules line; disabled “Terbitkan” with an
explanation for unvalidated accounts. No audience selector until there are audiences.

## 6. Validation and trust

**Evidence:** Nextdoor has no “verified neighbour” badge. Ordinary neighbours show name +
neighbourhood + time. Trust comes from an address-gated signup (implicit, not displayed) and
from role labels on institutions: `Local publisher`, `Admin`, agency badge, business check,
`Sponsored`. Verification itself is a named, resumable checklist: invitation → explicit
“submitted” acknowledgement → ticked item, with no ambiguous pending state.

**Rule:** Use two registers, not one badge.

- Implicit: the account exists and has a staff-reviewed state.
- Explicit: the byline words `Sudah divalidasi` or `Dikelola staf maulid.id`.

**maulid.id:**

- Account states: `Belum divalidasi` → `Menunggu peninjauan staf` → `Sudah divalidasi`, each
  with one sentence and the next step.
- Composer gate explains what validation enables and links to the contribution section.
- After validation, say what is now possible (`Tulis kabar`), not a generic congratulations.
- No gold/blue check, no shield, no trophy, no KYC funnel, no address/map/location flow.
- Do not fake a progress bar; a checklist is acceptable only if the steps are real.

## 7. Events

**Evidence:** Event cards use a date block, title, time, place, organiser and one action
(Interested/Join). Luma and SCMP use a date-gutter agenda; Meetup uses avatar stacks and
“160 going”.

**Rule:** Time and place must be explicit, including timezone. One RSVP action per event.

**maulid.id:** Events are future work. When added: `WIB`, venue, organiser, date-gutter
agenda for Maulid season, real attendance only.

## 8. Search and filters

**Evidence:** Search field at the top; tabs (`All / Posts / Businesses / Pages`) and a
single row of dropdown filters; compact text-first results.

**Rule:** Search is a first-class action, not a hidden icon. One filter row; no dashboard.

**maulid.id:** Existing text search over body. Add a clear button and empty-state recovery;
do not add category filters until the API supports them.

## 9. Navigation

**Evidence:** Web: left sidebar with Home, For Sale & Free, Faves, Groups, Events, Post CTA,
Settings/Help/Invite. Mobile: five-item tab bar with a raised Post FAB.

**Rule:** A short, labelled destination list; one primary action; icon + label.

**maulid.id web:** Kabar, Kontribusi, Akun, Backoffice (staff). Mobile: Beranda, Cari,
Tulis, Akun. No notification tab until notifications exist.

## 10. Right rail

**Evidence:** Neighbourhood alerts, invite card, business CTA, ads, chat dock. The rail is
narrower and quieter than the feed.

**Rule:** The rail supports the feed; it is not a second feed.

**maulid.id:** Tentang komunitas, account validation status, staff guidance, seasonal
Maulid card. Maximum 2–3 cards. Collapse below the feed on narrow screens.

## 11. Empty, loading, error and permission states

**Evidence:** Short plain text, often one recovery action; some products add a topic
fallback. Illustration-only empty states are rare.

**Rule:** Say what happened and what to do next. No fake content.

**maulid.id:**
- Empty feed: “Belum ada kabar. Akun yang sudah divalidasi dapat membagikan kabar pertama.”
- Filtered empty: “Belum ada kabar yang cocok. Coba kata lain.”
- Error: “Kabar tidak dapat dimuat.” + “Coba lagi”.
- Unvalidated: “Akunmu belum divalidasi. Staf meninjau akun sebelum penerbitan dibuka.” +
  link to account/contribution.

## 12. Moderation and safety

**Evidence:** The overflow menu is the moderation surface, with consequence-first items:

```text
Mute    Hide posts from this neighbour in the feed
Block   Prevent this neighbour from seeing what you share or interacting with you
Report  Flag for review
```

Own-post menu: turn off notifications, bookmark, delete, close discussion, close poll.
Comment menu: edit, delete. Ad menu: report, hide, about ads. Block has a three-bullet
confirmation and a persistent “You have blocked…” state with Unblock. Guidelines appear once
in onboarding and again in context.

**Rule:** Rare/destructive actions live in overflow with explicit consequences. A report
button exists only when a review path exists.

**maulid.id:**

- Ship `Sembunyikan` (client-side hide) when useful; it needs no backend.
- Ship `Laporkan` only when reports are stored and staff can review them; until then, do
  not render the item or a dead menu.
- `Blokir` is out of scope at launch (no DMs/follows); do not copy it for symmetry.
- Publish `Aturan komunitas` as a readable page with four short rules and link it from the
  composer footer.
- Do not expose a report reason taxonomy unless staff can consume it.

## 13. Seasonal modules

**Evidence:** No Nextdoor seasonal/Ramadan module was found in the captured corpus. The
closest dated modules are “Today's local news”, “Get started” checklists and local-deals
explanation pages, none of which are seasonal greetings. The earlier Ramadan note was based
on a misread capture and has been withdrawn.

**Rule:** Seasonal content is a maulid.id editorial decision, not a Nextdoor-derived
pattern. Keep it calm and non-commercial.

**maulid.id:** a static Maulid card during the season: mark → greeting → one action
(“Bagikan doa” or “Tulis kabar”). No countdown, no discount, no ornate clip-art. Treat it as
an intentional addition and record the decision.

## 14. Notifications and messages

**Evidence:** Separate bell/message icons in the header; mobile tab bar may carry badges.
Preferences and privacy are in settings.

**Rule:** Notification channels are explicit and controllable.

**maulid.id:** Defer until the API has notifications. Do not render a bell with a fake badge.

## 15. Responsive adaptation

**Evidence:** Web three-column layout; mobile bottom tab bar and single column. The same
post grammar survives both.

**Rule:** Collapse the rail and sidebar before shrinking the reading column. Keep touch
targets and labels.

**maulid.id:** At 1024px, hide the right rail or move it below the feed; at 768px, collapse
the left nav into a top menu or bottom bar; at 390px, single column and bottom navigation.

## 16. Monetization separation (avoid)

**Evidence:** Nextdoor labels sponsored content but gives it the same geometry as news;
local deals and business CTAs occupy the rail.

**Rule:** If monetization exists, it must be structurally and visually distinct from
editorial content. maulid.id has no monetization at launch.

**maulid.id:** No ads, sponsored cards, paywalls, member-only gates or promo countdowns.
