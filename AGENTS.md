# AGENTS.md

## Project

maulid.id is a community news-sharing platform for Maulid Rasulullah events in Indonesia.

It is intentionally a **monolith** with three surfaces sharing one codebase and deployment:
- Public landing page (React)
- Mobile API (`/api/v1`)
- Backend + backoffice for staff (React + Go)

Do not split these into separate apps/services without an explicit decision from the user.

## Domain rules

- Any user can interact with posts.
- Only **validated users** may create posts. Enforce this in the backend/API, not only in the UI.
- Validation status is a core account concept; keep its naming and behavior consistent across the API, backoffice, and landing page. The canonical values are `unvalidated` and `validated`, surfaced as “Belum divalidasi” and “Sudah divalidasi”.
- Staff accounts are created only through `cmd/staff-bootstrap`; never add an HTTP route that grants the staff role.

## Project skills

The project's reference material — the prompt sources and the Nextdoor design reference — is
extracted into loadable project skills under `.agents/skills/` (index: `docs/skills.md`).
Load the matching skill before the work:

| Skill | Use for |
| --- | --- |
| `react-feature-architecture` | React screens, components, hooks, state, styling, tests |
| `go-feature-architecture` | Go API, services, repositories, auth, staff code |
| `security-audit` | security review/audit; auth, session, validation changes |
| `monolith-delivery` | prompt→plan, branch/worktree, commits, journal, push |
| `nextdoor-design` | **mandatory** before any UI, layout, component, visual or design-review work |

They are project-specific and cite their sources and adaptations in each skill's
`references/`. They never override this file or the prompt/plan.

## Commands

Go 1.27.4+ and Node.js 22+ are required. The full command reference is in `README.md`.

```sh
make deps          # go mod download + npm ci
make build         # web build → embed → bin/maulid, bin/staff-bootstrap
make test          # go test ./... + npm test
make coverage      # backend statements and frontend metrics must be >80%
make lint          # gofmt check, go vet, feature-boundary and source-order checks
make race          # go test -race ./...
make security      # govulncheck ./... + npm audit
make run           # development server on :8080
make run-dev       # Vite dev server on :5173 (proxies /api to :8080)
make staff-bootstrap EMAIL=... NAME=...   # operator-only first staff account
```

A non-zero exit from any check means the task is not done. Do not commit a change
that lowers coverage to 80% or below, breaks `make lint`, or leaves `go test -race`
failing.

## Architecture rules

Read `docs/architecture.md` before structural changes; it is the source of truth for
module boundaries. In short:

- Backend features live in `internal/<feature>` and own their rules, repository
  port, SQLite adapter, HTTP adapter, and tests. Cross-feature access goes through an
  explicit interface or the canonical `account` types — never through another
  feature's repository or tables.
- Platform plumbing belongs in `internal/platform/*` and stays transport-level.
- Frontend features live in `web/src/features/<feature>` and are imported only
  through their `index.ts`. Anything used by more than one feature belongs in
  `web/src/shared/`.
- React components keep the requested order: imports, then variables/state/logic,
  UI output, and exactly one trailing stylesheet import. `make lint` enforces it.
- Never use `dangerouslySetInnerHTML` or browser storage for credentials; the CSRF
  token stays in memory.

## Design reference

All public, account and backoffice UI follows the Nextdoor-derived design reference in
`docs/design/` and the `nextdoor-design` skill. Before designing, building, restyling or
reviewing any screen, component or visual token:

1. load `nextdoor-design` and read the reference that matches the task;
2. apply `docs/design/ui-review-checklist.md`;
3. keep the post grammar `author · validation status · time → content → actions`;
4. attach desktop/tablet/phone screenshots and record deliberate deviations.

The reference is structural, not a clone. Never copy Nextdoor's logo, brand assets,
illustrations, photography, copy, ads or exact compositions. Do not add advertising,
sponsored cards, paywalls, points, streaks, leaderboards, paid verification ticks or
fintech-style identity checks. If the API cannot support a reference pattern, document the
schema work instead of faking it in the UI.

## Verification expectations

`docs/plans/initial.plan.md` maps every requirement to its evidence. Before claiming
work complete, run the full gate and record the actual numbers in
`.agents/journal/`:

- backend: formatting, vet, unit/HTTP tests, race, statement coverage >80%,
  `govulncheck`;
- frontend: typecheck, boundary/source-order lint, component/service tests,
  statements/branches/functions/lines >80%, `npm audit`, production build;
- runtime: all three surfaces served by one binary, the validation invariant
  exercised through the API, and the public page reviewed in a real browser;
- UI: `docs/design/ui-review-checklist.md` completed, with desktop/tablet/phone
  screenshots and keyboard, focus, 200% zoom and reduced-motion checks.

## Docs and journals

- `docs/plans/*.prompt.txt` are preserved verbatim; each gets a neighbouring
  `*.plan.md` before structural implementation.
- `docs/research/` holds the source matrix, raw retrievals, and image provenance.
- `docs/security-audit.md` records the security-audit run, findings, and remaining
  limitations.
- `docs/design/` records the mandatory Nextdoor-derived UI reference and review checklist.
- `.agents/journal/` records agent work, decisions, checks, failures, and fixes.
