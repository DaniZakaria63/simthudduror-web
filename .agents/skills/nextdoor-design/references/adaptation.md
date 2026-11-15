# Adaptation to maulid.id

## Product mapping

| Nextdoor | maulid.id |
| --- | --- |
| Neighbourhood | Maulid community / Indonesian Muslim readers |
| Home feed | Kabar komunitas |
| Post | Kabar |
| Neighbour (address-gated) | Reader / account |
| Address-gated membership (implicit) | Akun tervalidasi staf (explicit text) |
| Local publisher | Akun tervalidasi / staf |
| Event | Kegiatan Maulid (future; needs an event domain) |
| Group | Future community grouping; not in the first release |
| For Sale & Free / Faves | Not applicable |
| Business page / local deals / ads | Not applicable; do not add |
| Invite neighbours | Future account invitation; not in the first release |
| Notifications / messages | Future; needs API and storage |
| Report / hide / mute | Report/hide only when a moderation path exists |

## What to adopt now (fits the current API)

1. **Web shell:** left navigation, central feed, optional right utility rail.
2. **Feed card:** author, validation status, relative time, body, like action, overflow
   placeholder only if a real action exists.
3. **Feed header:** section title, search field, one filter row.
4. **Composer gate:** a visible “Tulis kabar” path; unvalidated accounts get an explanation
   and a link to account status.
5. **Trust strip:** “Semua akun dapat membaca dan menyukai. Penerbitan hanya untuk akun
   yang divalidasi staf.”
6. **Empty/loading/error states:** Indonesian, one recovery action.
7. **Right rail:** about the community, account status, staff guidance, and a static
   seasonal Maulid card when appropriate.
8. **Mobile shell:** bottom navigation with Home, Cari, Tulis, Akun; Tulis opens the gate.
9. **Typography:** keep EB Garamond for editorial display; add a UI sans for controls and
   metadata.
10. **Tokens:** canvas/surface/border/primary/muted/radius/spacing from the visual language
    reference.

## What to prepare but not ship yet (needs schema/API)

| Feature | Why it is blocked |
| --- | --- |
| Post title, image, category | `Post` has only `body`; no columns or API fields |
| Comments | No comment domain or API |
| Notifications | No notification storage or delivery |
| Groups | No membership/group domain |
| Events / Maulid agenda | No event domain, timezone or RSVP model |
| Regional filter | No location field |
| “Kabar pilihan”/featured | No featured flag or staff curation model |
| Pinned “Pengumuman” | No announcement type |
| Moderation actions | No report/moderation backend; a report button without a review path is dishonest |

Do not fake these with static content or client-only state. Document the schema work in a
separate plan.

## Screen-by-screen adaptation

### Public landing page

Current: `Masthead` → `Kabar komunitas` feed → `ContributionSection`.

Recommended composition:

```text
Top bar: brand · search · account action
Left nav: Kabar · Kontribusi · Akun · (Backoffice for staff)
Centre:  masthead/compact heading · filter row · feed
Right:   about community · validation status · Maulid seasonal card
```

- Keep the masthead's editorial identity, but reduce its height on the feed-first view.
- The first screen must show the search/filter and at least the start of the feed.
- The right rail is not a second feed; no “popular posts” invented from nothing.

### Feed card

```text
[avatar] Author name · Sudah divalidasi · 2 jam lalu        [overflow]
Body (title when the API supports it)
[image/event attachment when supported]
[Like] [Comment when supported] [Share when supported]       [quiet metadata]
```

- No fake comment count; hide the action until comments exist.
- No engagement score. Like count is real or omitted.
- The author row is the trust unit.

### Composer

```text
Author identity + validation status
[Body field]
Rules line
[Cancel] [Terbitkan]   (disabled + explanation when unvalidated)
```

- For an unvalidated account, “Terbitkan” is disabled and the explanation links to
  `/#kontribusi` or `/akun`.
- The composer does not need Nextdoor's “Sell / Event / Poll” intent cards; maulid.id has
  one content type today.
- Keep the rules line short and link to the full guidelines.

### Account screen

- Show `Belum divalidasi` / `Sudah divalidasi` as a plain status row with an explanation.
- Show what staff validation enables: publishing kabar.
- Do not show a fake progress bar or KYC checklist.

### Backoffice

- Use the same tokens but a denser table layout; Nextdoor's consumer card grammar does not
  apply to staff review queues.
- Keep one primary action per review row (`Validasi`, `Tolak`) and state consequences.
- The public design reference must not force editorial whitespace onto the backoffice.

### Mobile app (future)

- Five destinations maximum: `Beranda`, `Cari`, `Tulis`, `Akun` (and one more only when
  the API supports it).
- `Tulis` opens the composer gate.
- No notification tab until notifications exist.
- No maps, neighbourhood verification, classifieds or business surfaces.

## Copy adaptation

| Nextdoor | maulid.id |
| --- | --- |
| “Hello, neighbors!” | “Assalamu’alaikum, warga Maulid” only if the community uses it; otherwise “Kabar komunitas” |
| “Post” | “Tulis kabar” / “Terbitkan” |
| “Reply” | “Balas” (when comments exist) |
| “Local publisher” / “Admin” | “Akun tervalidasi staf” / “Dikelola staf maulid.id” |
| “See more” | “Selengkapnya” |
| “Send invitation” | “Undang akun” (future) |
| “Sponsored” | Never |

Do not translate Nextdoor idioms literally. Indonesian copy should be shorter and calmer.

## Anti-patterns for maulid.id

- Neighbourhood address verification and maps.
- House/neighbour illustrations and “neighbor” social mechanics.
- Classifieds, local deals, business pages and ads.
- Sponsored content in the kabar stream.
- Paid verification ticks.
- Upvote/streak/points/leaderboards.
- US “Lower Allston” local framing; maulid.id can use Indonesian regions later but not a
  fake neighbourhood identity.
- Dense multi-column newsroom grids.
- Autoplay, heavy carousels and countdowns.
- English-first UI copy.

## Implementation sequence

1. **Design tokens and shell** — canvas, surface, border, primary, UI sans, radii, spacing;
   web three-zone shell; mobile bottom nav.
2. **Feed card and feed header** — identity/status/time, body, real like action, search and
   filter.
3. **Composer gate** — account status, explanation, disabled publish.
4. **Right rail** — about/status/guidance, seasonal card.
5. **Empty/loading/error states** — Indonesian, recovery action.
6. **Responsive and accessibility pass** — 1440/1024/768/390, keyboard, zoom, reduced
   motion.
7. **Future schema work** — title/image/category, comments, notifications, events, groups,
   moderation, in separate plans.

Each step must pass `docs/design/ui-review-checklist.md` and the normal frontend gates.
