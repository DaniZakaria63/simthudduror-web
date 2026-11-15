# feat/designing — Mobbin research setup — 2026-11-15

## Session

- Agent: OpenCode (DeepSeek V4.1 Flash 1017), principal/executor for this slice.
- Request: install Mobbin MCP, search Islamic community news-sharing references, create
  `feat/designing` from `origin/development` in its own worktree, use subagents.
- Worktree: `/home/dani/WebstormProjects/maulid.id-feat-designing`
- Branch: `feat/designing`, created from `origin/development` (`74ae972`, merge of
  `feat/initialize`).

## Mobbin MCP

Commands actually run:

```sh
opencode mcp add mobbin --global --url https://api.mobbin.com/mcp
opencode mcp list
```

Observed:

- `opencode mcp add` wrote the server to `~/.config/opencode/opencode.json`:
  `{"type":"remote","url":"https://api.mobbin.com/mcp"}`.
- First `mcp list` reported `⚠ mobbin needs authentication`; the subsequent list reported
  `✓ mobbin connected` after the MCP client completed its OAuth connection.
- A live `search_screens` call returned four screens (Meetup, Nextdoor, Insight Timer,
  Formula 1), confirming the tools work from this session.
- Mobbin tools exposed: `search_screens`, `search_flows`, `search_sections`,
  `list_collections`, `get_collection_contents`, `create_collection`, `save_content`,
  `add_comment`.

## Worktree creation

```sh
git fetch origin --prune
git worktree add -b feat/designing \
  /home/dani/WebstormProjects/maulid.id-feat-designing origin/development
```

Observed: `Preparing worktree (new branch 'feat/designing')`, branch tracking
`origin/development`, HEAD `74ae972`. The OpenCode session was moved to the worktree with
`tools.opencode.session_move`, so all subsequent paths are relative to it.

## Research workers

Two subagents ran in parallel with non-overlapping briefs:

| Session | Brief | Result |
| --- | --- | --- |
| `ses_e205af0a6ffdwBOFkfI2NCyFmC` | Islamic/Muslim apps, Indonesian community apps, community-news references | 12 ranked references, Islamic-adjacent modules, confirmation that Mobbin has no dedicated Islamic app |
| `ses_e205af09effe83a7IA61G31WZU` | Mobile feed/detail/publishing patterns, web landing/sections, trust and moderation | 9 mobile references, 8 web/section references, publishing/trust flow findings, weak-query list |

Both were instructed to inspect previews, preserve Mobbin URLs, separate observation from
interpretation, and not save collections or edit files.

### Accepted findings

- Mobbin has **no dedicated Islamic app** and no major Indonesian news app. Repeated named
  and category probes returned unrelated wellness/community products; Indonesian-language
  queries are not indexed.
- The closest Islamic community-content pattern is **Nextdoor's Ramadan Mubarak card**
  (mark → greeting → one verb, no commerce). Islamic-adjacent modules exist in Binance
  (Sharia), Snoonu (halal certification), Grab/GoPay (Indonesian Ramadan/THR language) and
  Careem (donation).
- Strong structural references: Nextdoor (local feed + news strip + composer), Lex
  (place-based post feed), Reddit (source-first row), Yahoo News Singapore (SEA news row),
  SCMP (latest stories, News Agenda), Luma (upcoming/past events), NYT/Quartz/Wise/PayPal/
  Ramp (editorial web), Dropbox Community (honest size signals), Nextdoor publisher sheet
  (plain-language trust label).
- Trust findings: show validation as plain language (`Local publisher`, `Admin · 2w`) rather
  than paid-check iconography; put identity and rules inside the composer; publish rules as
  a readable numbered page; put moderation in one overflow menu with consequences.
- Empty states: Gojek, SCMP and Deepstash show recovery actions, not illustration-only
  screens.

### Rejected or narrowed

- Fintech identity verification (ID photo/selfie) is disproportionate and does not fit a
  religious community; staff validation is the product rule.
- Reddit/Yahoo vote and read counts are not recommended at current scale; they can reward
  sensation and look implausible.
- Sponsored cards in the editorial stream, paywalls, streaks/points, autoplay and
  Ramadan commerce tropes were treated as anti-patterns.
- Title/image cards, category pills, event modules and regional filters were kept as
  design directions but marked as requiring API/schema work; they are not implementable in
  the current `Post` model.

The final report is `docs/research/mobbin-design-references.md` (86 unique Mobbin URLs,
grouped by theme, with a source index and a “fits today / needs first” table).

## Availability result

Exact Islamic products were queried by name (`Umma`, `Muslim Pro`, `Quran.com`, `Muzz`) and
by category. Mobbin did not return them; it returned adjacent meditation, community and
social products. The research therefore uses community-news and neighbourhood references,
labels the gap explicitly, and does not claim an Islamic precedent that Mobbin does not
contain.

## Files added

- `docs/plans/designing.prompt.txt` — request preserved verbatim.
- `docs/plans/designing.plan.md` — scope, requirements, verification, limitations.
- `docs/research/mobbin-design-references.md` — reference report.
- `.agents/journal/2026-11-15-designing.md` — this record.

## Checks

- `opencode mcp list` → `✓ mobbin connected`.
- `git fetch origin --prune` → `origin/development` advanced to `74ae972`.
- `git worktree add -b feat/designing ... origin/development` → worktree created.
- `git diff --check` → exit 0 after removing Markdown hard-break trailing spaces.
- `git diff --cached --stat` → four docs files, 1,065 insertions; no Go/React source.
- Docs-only change; no Go/React source touched, so `make` gates are not run for this slice.

## Limitations and next steps

- Mobbin catalogue coverage and plan access limit the references; absence is not proof that
  a product does not exist.
- The current `Post` model has no title, image, category or event fields, so title-led cards
  and event modules need a separate API/data plan before implementation.
- Next: turn the reference report into a design direction for the public landing page and
  the future mobile app, then implement behind the normal verification gate.
