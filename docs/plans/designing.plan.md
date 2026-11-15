# feat/designing — Mobbin reference research plan

## Interpretation

The request has four concrete outcomes:

1. Install the Mobbin MCP server in this OpenCode installation.
2. Search Mobbin for Islamic community / news-sharing products and adjacent community-news
   products that can inform maulid.id's design.
3. Move the working branch to a dedicated worktree created from `origin/development` as
   `feat/designing`.
4. Use subagents for the research work, then consolidate their findings in the repository.

This slice is **research and branch setup**. It does not change production UI code. Any
resulting redesign gets its own plan and verification evidence before structural
implementation.

## Requirements and evidence

| # | Requirement | Evidence |
| --- | --- | --- |
| R1 | Mobbin MCP is configured and usable | `opencode mcp list` reports the `mobbin` server; a `search_screens` call returns screens |
| R2 | Worktree exists on `feat/designing` from `origin/development` | `git worktree list`, `git status -sb`, `git log -1` |
| R3 | Research covers Islamic/Muslim products and adjacent community-news products | research report with product names, Mobbin URLs and observations |
| R4 | Research distinguishes visible observation from interpretation | report sections and wording |
| R5 | Subagents contribute and their findings are reconciled | journal records the subagent sessions and what was accepted/rejected |
| R6 | No production UI changed in this slice | `git diff --stat` contains docs only |

### Status (2026-11-15)

| # | Status | Evidence |
| --- | --- | --- |
| R1 | Done | `opencode mcp list` → `✓ mobbin connected`; live searches returned screens |
| R2 | Done | worktree at `/home/dani/WebstormProjects/maulid.id-feat-designing`, branch `feat/designing`, base `74ae972` |
| R3 | Done | `docs/research/mobbin-design-references.md` covers Islamic-adjacent, community-news and editorial references |
| R4 | Done | each section separates visible observations from “Borrow”/interpretation |
| R5 | Done | two subagent sessions recorded in the journal; findings reconciled and deduplicated |
| R6 | Done | docs-only diff; no Go/React source changed |

## Scope

- Mobbin queries: iOS/Android screens, flows and web sections.
- Reference themes: Islamic apps, Indonesian community apps, local/community news feeds,
  post publishing and validation cues, editorial web landing pages.
- Output: `docs/research/mobbin-design-references.md` and a journal entry.
- Out of scope: implementing a redesign, changing copy, committing generated images,
  pushing new product behaviour.

## Worktree

```text
branch:   feat/designing
worktree: /home/dani/WebstormProjects/maulid.id-feat-designing
base:     origin/development @ 74ae972
```

## Research method

1. Search screens, flows and sections separately; keep each query to one concept.
2. Inspect the returned previews, not just app names and tags.
3. Prefer several relevant products over one exemplar.
4. Preserve the product name and canonical Mobbin URL for every reference.
5. Record what Mobbin does **not** contain, especially exact Islamic community-news apps.
6. Translate findings into adaptations for maulid.id, not copies.

## Verification

- `opencode mcp list` and one live Mobbin search.
- `git worktree list` and branch/base check.
- `git diff --check` and docs-only diff review.
- No `make` gate is required because no Go/React source changes. If any code is added
  later, the full gate in `AGENTS.md` applies.

## Known limitations

- Mobbin access depends on the user's Mobbin plan and OAuth session.
- Mobbin is a commercial reference library; absence of an app does not mean the app does
  not exist, only that it was not returned by the searches.
- Reference screens show that a product used a pattern; they do not prove the pattern
  improves outcomes.
- Designs, branding, copy and assets must be adapted, never copied.

## Next step (not in this slice)

Turn the reference findings into a design direction and an implementation plan for the
public landing page and the future mobile app, then implement it behind the normal test,
coverage, race, security and browser gates.
