# Landing-page critique — round 2 (final)

Screens: `.taste-skill/inspect/desktop-fold.jpg`, `desktop-page-01.jpg`,
`desktop-page-02.jpg`, `tablet-fold.jpg`, `mobile-fold.jpg`, `mobile-page-01.jpg`.
Inspector state at this round: **0 FAIL, 2 WARN** (no footage; 15 moving elements).

## Look pass notes

1. **One second** — “a printed community bulletin”: the band of mosque photographs over
   the large Garamond masthead reads as one thing.
2. **Squint** — three masses: photographic band, masthead field, blue contribution
   field. Clear size difference; nothing competes with the headline.
3. **Cover the logo** — the structure (photo band → masthead → reading column →
   validation sequence) is an editorial community-news page, not a SaaS page. A
   neighbouring community-news site could use the same frame; the specific content
   (Maulid, validation-before-publishing) is what fills it.
4. **Edges** — three deliberate left edges: page margin, the indented contribution
   steps, and the feed's rail. Consistent.
5. **Space** — section gaps are clearly larger than internal gaps; the feed rail's
   white below the note is a rest, not a hole, because the reading column holds the
   weight.
6. **Type** — one family; headline/body jump is >4×; longest headline is two lines at
   all three widths; no orphan; body measure stays under ~75 characters.
7. **Colour** — white ground, black ink, one blue field and one red/black accent block
   in the photo band; blue returns only on the primary action and the contribution
   field.
8. **Images** — four real, credited photographs of Indonesian mosque/Maulid settings,
   one grade (black and white), square crops; no moving picture (no honest clip of
   these communities was available).
9. **Sameness** — photo band / masthead split / reading column / staggered blue steps /
   black footer close: no repeated layout.
10. **Phone** — its own composition: band crops differently, headline keeps two lines at
    50.7 px, search and post stack, contribution steps indent instead of staggering
    across columns.
11. **Fullness** — every desktop screen carries content; the feed rail is the only
    deliberate rest.
12. **Substance** — the page states what the site is, shows the real feed, explains who
    may publish and how validation works, and links the account and staff surfaces.
13. **Collisions** — none observed at 1440/768/390; no horizontal overflow on landing,
    account, or backoffice at 390 px.
14. **Competitor + live thing** — the live thing is the API-backed feed with visitor
    search and signed-in likes. The competitor test weakens creativity: a community
    news bulletin is a known frame, and the validation idea appears in words rather
    than in the composition.

## Score

| Criterion | Weight | Score | Reason |
| --- | --- | --- | --- |
| Design | 40% | 7.0 | Coherent editorial system; capped at 7 by the still-page warning and limited motion. |
| Usability | 30% | 8.0 | Purpose and action clear in three seconds; comfortable phone layout; the desktop feed column is narrower than it could be. |
| Creativity | 20% | 7.0 | Clear idea, but the editorial-bulletin frame survives the competitor test; capped at 7. |
| Content | 10% | 8.0 | Specific Indonesian copy, real feed, honest image credits; the validation sequence is explained rather than dramatised. |

**Overall 7.4** — a confident editorial page that would not embarrass a designer, short
of award level because the idea lives mostly in the type and the copy.

## Fixes applied across the two rounds

1. **Contrast FAIL** — “Buat akun” was white on white because `.contribution a`
   outranked the button rule; scoped the inherit rule to paragraph links.
2. **Body-size FAIL** — running text measured 14 px; raised the small step to 16–17 px
   and reviewed the result at all three widths.
3. **Grey-body WARN** — muted ink darkened to 7:1+.
4. **Fullness** — rebuilt the feed as a two-column reading section with a sticky rail;
   the previous layout left most of the page empty beside a narrow feed.
5. **Motion** — added staggered footer reveals; 15 elements now move, still under the
   16 warn threshold, all transform/opacity with reduced-motion fallbacks.

## What still falls short

- No moving picture: no honest clip of these specific communities/mosques was
  available, and a generic stock clip would misrepresent the subject.
- Motion is entrances only; the signature masthead move is subtle at first load.
- The desktop feed column could take more width; it is capped at 42 rem to keep the
  measure readable.
- The photography is contextual (mosques and a 2022 Maulid gathering), not maulid.id
  events; the page says so in the footer.
