# Implementation completion — 2026-11-14 (continuation)

## Session and hand-off

The previous principal session (Codex) stopped at its usage limit while the OpenCode
implementation worker (dispatch `ctx_ae0a40627a67`, DeepSeek V4.1 Flash `high`) was
still building the React surfaces. The worker had completed the Go monolith and most
of the frontend, then continued on its own until this continuation session interrupted
its turn at 15:21 to take over implementation directly (user instruction: “just
implement it”). No further agent was started; the touched worktree is the same
dedicated `feat/initialize` checkout.

State at takeover: Go packages built and all backend tests passed; the frontend had
`App.tsx`, landing, account, posts and backoffice features, but `tsc` failed on four
test typing issues, one backoffice test raced a re-render, and the build tooling
(`Makefile`, `.gitignore`, Go coverage gate, README, architecture/security docs) did
not exist.

## Code and verification changes in this session

- Fixed test typing without weakening assertions: typed the `Account` fixtures in the
  account/backoffice suites, removed an unused import/parameter, and replaced a
  narrowing-prone `Headers` capture with a mutable holder.
- Fixed the backoffice suite race: the test now waits for staff data before asserting
  the heading, instead of latching the loading-state heading that React then replaces.
- Fixed a real defect found by the runtime check: the audit entry returned by
  `PATCH /api/v1/staff/accounts/{id}/validation` had an empty `actor_name`, so the
  backoffice showed a blank actor. `staff.SQLRepository.SetValidation` now loads the
  actor name in the same transaction; repository and HTTP tests assert `Staf`.
- Fixed the landing-page contrast defect the design inspector measured: the
  `.contribution a { color: inherit }` rule outranked `.contribution__button`, making
  “Buat akun” white on white. The rule now applies to paragraph links only.
- Raised the smallest type step to 16–17 px (the inspector measured 14 px running
  text) and darkened `--muted` from `#5c5c5c` to `#4f4f4f` for ≥7:1 running text.
- Rebuilt the feed section as a two-column editorial reading layout (large heading and
  marginal note in a sticky left rail, feed column at up to 42 rem) instead of a
  narrow feed with a mostly empty page.
- Added staggered reveals to the footer columns; the inspector now counts 15 moving
  elements (warn threshold 16) with reduced-motion fallbacks.
- Added `Makefile`, `.gitignore`, `scripts/check-go-coverage.sh`, `README.md`,
  `docs/architecture.md`, `docs/security-audit.md`, and updated `AGENTS.md` with the
  real commands, architecture rules and verification expectations.

## Checks observed in this revision

Backend (Go 1.27.4 at
`/home/dani/.local/share/maulid-toolchains/go1.27.4/go/bin/go`):

- `gofmt -l .` → no files; `go vet ./...` → clean
- `go build ./...` → success
- `go test ./...` → all packages ok
- `go test -race ./...` → all packages ok
- `go test -coverprofile … ./internal/...` → **87.3%** statements
  (`./scripts/check-go-coverage.sh` passes the strict `>80%` gate)
- `go run golang.org/x/vuln/cmd/govulncheck@v1.9.0 ./...` → “No vulnerabilities found.”

Frontend (Node 26.5.0, npm 11.17.0, locked manifest):

- `npm run typecheck` → clean
- `npm test` → 13 files, **69 tests passed**
- `npm run coverage` → statements **96.31%**, branches **87.24%**, functions
  **97.24%**, lines **96.44%**; every metric passes the strict `>80%` gate
- `npm run lint` → “Feature boundaries OK”, “Source order OK”
- `npm run build` → production bundle built; `make build` embeds it and builds
  `bin/maulid` + `bin/staff-bootstrap`
- `npm audit --audit-level=high` → “found 0 vulnerabilities”

Runtime, one binary on `127.0.0.1:18080` with a disposable SQLite file:

- `/healthz` 200, `/readyz` 200, `/` 200 (embedded React), `/akun` and `/backoffice`
  client routes served by the SPA fallback
- staff bootstrap CLI created the first staff account from stdin; no HTTP path exists
- anonymous publish → 401; unvalidated member publish → 403 `validation_required`;
  staff validation without CSRF → 403; with CSRF → 200 + audit row; member session
  immediately shows `validated`; validated member publish → 201; like twice →
  idempotent `like_count`; nonexistent like → 404; logout → 204 then session → 401
- cookie observed `HttpOnly; SameSite=Lax` (Secure in production); security headers
  present (CSP, `X-Frame-Options: DENY`, `nosniff`, referrer/permissions/COOP)
- design inspector (real Chrome, 1440/768/390): **0 FAIL, 2 WARN** (no footage; 15
  moving elements) after the fixes above; screenshots reviewed at all three widths
- `node /tmp/opencode/verify-ui.mjs` → **15/15** checks: anonymous backoffice gate,
  real form login, staff status, account page, backoffice table/audit/filter, staff
  validation round trip, no horizontal overflow at 390 px on landing/backoffice/account,
  feed search filter, reduced-motion content visibility

## Security audit

Cloudflare `security-audit` revision `c1c8a8c1471069fb0e188eeaff69b8e8db6564a8`,
profile `quick`, all application surfaces. Ledger (12 units) and findings (4 records)
both pass the skill's validators. Summary in `docs/security-audit.md`; artifacts at
`/tmp/opencode/security-audit/maulid.id/run-1/`.

- F-1 confirmed low: registration conflict reveals whether an email is registered.
  Accepted limitation; a real fix needs a neutral acknowledgement/email verification,
  which would change the registration contract.
- NV-1 needs validation: the auth limiter keys on `RemoteAddr`; behind a proxy all
  clients share one bucket. Deployment topology is not in the repository.
- R-1, R-2 rejected with source evidence (CSRF-at-rest, feed SQL injection).
- Limitations: one quick pass; single-reader topology with a separate re-read
  verification pass; no prior ledger; deployment/edge controls not observable.

## Open items

- F-1 remains an accepted low risk; NV-1 remains a deployment prerequisite.
- Keep the worktree until the user confirms the branch is merged.

## Delivery

Commits on `feat/initialize` (feature-prefixed, smallest coherent slices):

```text
ce697c8 backend: add feature-modular Go monolith with SQLite accounts, posts and staff validation
392e244 web: add editorial landing page, account screens and staff backoffice
90ba42a tooling: add build, coverage gates, boundary checks and developer docs
931d2b6 docs: record plan evidence, architecture, research sources and security audit
```

`git push -u origin feat/initialize` succeeded; `git ls-remote origin
refs/heads/feat/initialize` returned `931d2b6e03056a0a0feff1edb32a3f575d53aee0`,
identical to the local HEAD, and the worktree was clean at that point. A final rebuild
of the pushed revision served `/`, `/backoffice` and the hashed CSS with 200 and a
staff login with 200 on a fresh database.


## Final verify gate

- Added a single root `make verify` target (lint + race + coverage + web build +
  govulncheck/npm audit + `go build ./...`) and documented it in `README.md` and
  `AGENTS.md`.
- Observed `make verify` on 2026-11-14 with Go 1.27.4 and Node 26.5.0: gofmt/vet OK;
  `go test -race ./...` OK; backend coverage 87.3% (>80%); Vitest 69/69 with
  statements 96.31%, branches 87.24%, functions 97.24%, lines 96.44%; typecheck and
  Vite production build OK; `govulncheck` "No vulnerabilities found."; `npm audit`
  0 vulnerabilities; final `go build ./...` OK.
- Committed as `6b4a7be tooling: add full verify gate` and pushed to
  `origin/feat/initialize`; local and remote heads matched.
