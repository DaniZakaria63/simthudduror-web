# Initial monolith foundation

Status: implemented and verified in the `feat/initialize` checkout (2026-11-14). Acceptance evidence below was observed on the final revision; two limitations remain explicit (deployment rate-limit identity and the low-severity registration enumeration).

## Prompt and repository analysis

Source: `initial.prompt.txt`, preserved verbatim. The repository starts at `902b3a0` with only `AGENTS.md` tracked and no application, toolchain, or tests. The existing `feat/initialize` branch and checkout are dedicated to this initial prompt; they must not be removed before the user confirms a merge.

The prompt is an architecture and engineering-workflow specification. It names three surfaces, a language for each, modularity conventions, tests, documentation, reference reading, security review, agent roles, and Git delivery. It does not specify a complete production product backlog. Implement a runnable foundation with one useful end-to-end domain slice, rather than inventing unrelated product features.

The project instructions add the core domain: community news about Maulid Rasulullah events in Indonesia. Every account may interact with posts. Only **validated users** may create posts, with this rule enforced on the server. Validation status must have one canonical name and vocabulary across the API, public site, and staff backoffice.

## Required behavior

1. A Go process exposes `/api/v1`, health endpoints, the React public landing page, and React backoffice routes from one deployment. Development may use a frontend development server proxy, but production has one application service and release artifact. No separate frontend service, microservices, mobile client, or application split.
2. A visitor can read the public explanation and list published posts. The page makes posting eligibility clear and links to the account flow. Any example event data is explicitly identified as sample content; do not invent confirmed events or community statistics.
3. Accounts can authenticate. The backend establishes identity and loads current authorization and validation status from trusted persistence. New account input cannot assign staff privileges or validated status. Staff bootstrapping is an explicit operator action, without a hardcoded password.
4. An authenticated account can interact with a post regardless of validation status. Implement a small concrete interaction such as liking/unliking a post, with deterministic and idempotent behavior.
5. A validated account can create a post. An unauthenticated request is rejected as unauthenticated; an authenticated unvalidated account is rejected as forbidden. The same checks apply to direct mobile/API requests and to browser UI calls. An account losing validation loses posting eligibility on subsequent requests.
6. Staff can inspect accounts and change their validation status in a protected backoffice. Unauthenticated visitors and ordinary accounts cannot read staff-only account data or change validation. Administrative operations leave a useful audit record.
7. Application data survives a normal server restart. Prefer SQLite for the initial single-service deployment, with migration/setup instructions, transactions where invariants require them, and a durable configured path. Do not silently substitute an in-memory production store.
8. Both React surfaces handle loading, empty, success, validation, and request-error states and remain usable on a narrow mobile viewport. Forms have labels, keyboard focus, and actionable feedback. The public site is in Indonesian; validation terminology is consistent.

## Feature boundaries and code conventions

- Backend: `cmd/server` is a small composition root. Feature packages under `internal/` own their domain rules, service/use cases, repository interfaces and HTTP adapters. Accounts/authentication, posts/interactions, and staff validation have explicit public entry points; no cross-feature access to private state. Shared HTTP, database and configuration plumbing belong in small platform packages.
- Frontend: one React + TypeScript application, organized as `web/src/features/<feature>/`. Each feature owns its components, service/API functions, local state, types, and styling. Import another feature only through an explicit public entry point. Global UI primitives and infrastructure used by multiple features live under `web/src/shared/`; domain-specific components do not become global for convenience.
- Screens and components follow the requested reading order: imports first, constants/types/state/variables next, UI output in the middle, styling at the bottom via a trailing stylesheet import or a clearly documented consistent colocated equivalent. Avoid mixing large style definitions into component logic. Apply equivalent predictable ordering to other files.
- Enforce feature dependency rules with an executable lint check where feasible. Keep dependencies purposeful and lock exact install resolution.
- Record generated code, decisions, checks, failures, and remediation under `.agents/journal/`. Every future `*.prompt.txt` gets a neighboring `*.plan.md` before structural implementation; retain the source prompt.

## Security acceptance

Use the actual Cloudflare security-audit skill at a recorded upstream revision. Read its instructions and applicable references; execute its audit procedure against the final source, document findings with file/line evidence, remediate material findings, and record any justified remaining limitations. A dependency scan alone does not prove source security.

Minimum implementation concerns: password hashing, cryptographically random expiring sessions, protected cookie/bearer handling suitable for browser/mobile clients, CSRF protection for cookie-authenticated writes, backend authorization and current validation lookup, bounded/strict JSON input, parameterized persistence, escaped React output, bounded login attempts, safe error responses, request/body/header/time limits, graceful shutdown, security headers, no committed credentials, and deliberate production configuration. Define and test the chosen mechanisms; do not claim exhaustive security.

## Verification and delivery

| Prompt requirement | Artifact | Evidence required before acceptance |
| --- | --- | --- |
| Read existing code and all references | `docs/research/initial.sources.md` | File inventory and source matrix with access limitations; no false claim of reading unavailable pages |
| Go API/backend; React public/backoffice | Go module, frontend manifest, server routing/build | Successful backend/frontend build and runtime requests to all three surfaces |
| One feature-modular monolith | `docs/architecture.md`, feature source layout | Source review, import-boundary check, one production start command/artifact |
| Feature-owned component/service/state/UI | Feature folders, shared public entry points | Source and boundary-lint evidence |
| Imports/variables/UI/styles ordering | React components, developer instructions | Component review and consistent ordering |
| Agent code journal | `.agents/journal/` | Actual execution/model/review/test evidence, not only intentions |
| Unit coverage strictly greater than 80% | Backend/frontend unit suites and gates | Measured backend statement coverage and frontend statements/branches/functions/lines all >80%; checks must fail at <=80% |
| Validated-user posting invariant | API service and HTTP tests | Anonymous 401, unvalidated 403, validated success, revocation enforcement, forged privilege rejection |
| Any account may interact | Interaction service/API and UI | Unvalidated account interaction succeeds; idempotency and nonexistent-post handling verified |
| Protected staff validation | Staff API/backoffice | Nonstaff rejected, authorized status updates succeed, next post request uses new status |
| Secure code reviewed with Cloudflare skill | `docs/security-audit.md` | Recorded revision, source audit, dependency results, addressed findings and explicit limitations |
| Specific branch/worktree and small commits | Git branch/worktree and history | A dedicated prompt branch, feature-prefixed small commits, no unrelated files |
| Push completed work to origin | Remote branch | Successful push and remote head matching reviewed local head |
| Keep worktree until human confirms merge | Existing checkout | No cleanup/removal without the user's merge confirmation |
| Preserve prompt and write interpretation | `initial.prompt.txt`, `initial.plan.md` | Prompt unchanged; final plan reflects actual implementation and evidence |
| Requested agent roles | Journal and Orca receipts | Antigravity research, verified OpenCode Deepseek V4.1 Flash high session, Jev advice with principal decisions; no recursive principal launch |

Backend checks include formatting, vetting, meaningful unit/HTTP tests, race detection, coverage and relevant vulnerability scanning. Frontend checks include type checking, lint/boundary checks, behavioral component/service tests and coverage, dependency audit, and production build. Commands must be derived from the created manifests and documented in README and AGENTS.md.

Runtime checks cover persistence across restart, direct API authorization and validation changes, browser navigation and form behavior, and public/backoffice rendering at desktop and mobile widths. Follow the loaded frontend design skill for the public landing page, including a recorded direction, actual imagery provenance where used, working domain-specific interaction, browser inspection and screenshot critique. Never claim a check that did not run.

## Execution ownership

Codex is principal: architecture, plan, delegation, review, and acceptance. The single OpenCode worker owns application changes and fixes, using `opencode-go/deepseek-v4.1-flash#high`; verify the actual assistant-session model after execution. Antigravity owns external research and may write only its isolated research note/journal. Jev is used for typed advisory judgments around routing and review; its confidence is not authorization and never replaces source or test evidence.

Research completed as task `task_6a00049242a6`, dispatch `ctx_960a7eaec30e`, with corrections reported under dispatch `ctx_85a2813d079a`. The principal reviewed `docs/research/initial.sources.md`, obtained and read the actual Cloudflare skill at `c1c8a8c1471069fb0e188eeaff69b8e8db6564a8`, and verified the official Go release manifest. Follow-up task `task_ce6807c04b69` resolved the two Medium access gaps through normal public-browser access; both Antigravity and principal read the full articles without credentials or bypass. Details are in `docs/research/initial.reference-verification.md`; the earlier direct-HTTP 403 observations remain historical transport evidence, not a current full-text limitation. Blog recommendations inform the adaptation but do not override this project's monolith and feature-boundary requirements. Use native Orca mailbox supervision, never launch a second principal or a polling scheduler.

## Approved implementation decisions

- Use the existing dedicated `feat/initialize` branch and `/home/dani/WebstormProjects/maulid.id` worktree for this single prompt. No new checkout is necessary because there is no competing application work. Keep it after delivery until the user confirms merge.
- Go **1.27.4**, verified against `https://go.dev/dl/?mode=json`, is installed at `/home/dani/.local/share/maulid-toolchains/go1.27.4/go/bin/go`. The official Linux amd64 archive SHA-256 is `7226736c2339719cda28fbf3a54d71aa63bcefe376e98cef00034268dab24840`. Do not modify global shell configuration; use the configured toolchain in commands and document portable prerequisites.
- Use SQLite persistence, Go's standard HTTP routing, one embedded React/TypeScript/Vite build, and feature-owned services/repositories/HTTP handlers. Small local React hooks are sufficient for this slice; TanStack Query and Zustand are optional, not requirements inferred from the linked articles.
- The canonical account field is `validation_status`, with exactly **`unvalidated`** and **`validated`** for this slice. Roles are separate (`member`, `staff`). The research note's example `pending`/`rejected` states are not adopted. Never accept either role or validation changes from self-registration/profile payloads.
- Browser sessions use secure HTTP-only cookies in production and explicit CSRF protection for state changes; mobile clients use opaque bearer credentials. Sessions expire and are revocable. Read trusted current account status when making authorization decisions, rather than stale client/token claims.
- Staff bootstrap is an operator-only command or explicit bootstrap mode using supplied credentials, never a public API or a default password. No production demo accounts or fabricated event seed data.
- The public landing direction is documented in `initial.design.md` and `.taste-skill/direction.json`: editorial-report seed 979061848. Verify actual image source/license metadata and inspect downloaded images before use; research-proposed URLs are candidates, not proven assets.
- Source-backed Cloudflare review is performed against the completed revision; the principal owns formal audit artifacts. Ordinary development checks are distinct from security-audit execution. The skill requires OS-enforced isolation for any target execution during the audit; use source evidence when execution controls are unavailable and disclose the exact gap.

## Implemented evidence

All commands were run from the repository root with Go 1.27.4
(`/home/dani/.local/share/maulid-toolchains/go1.27.4/go/bin/go`) and Node 26.5.0.
Full details and timestamps are in `.agents/journal/2026-11-14-implementation.md`.

| Requirement | Observed evidence |
| --- | --- |
| Read existing code and references | `docs/research/initial.sources.md` + `initial.reference-verification.md`; both Medium articles read through public browser access |
| Go API/backend; React public/backoffice | `make build` produced `bin/maulid`; the binary served `/`, `/akun`, `/backoffice` (200) and the API on one port |
| One feature-modular monolith | `docs/architecture.md`; `internal/{account,post,staff,platform,app,webui}`; frontend boundary lint reports “Feature boundaries OK” |
| Feature-owned component/service/state/UI | `web/src/features/*` with public `index.ts`; deep imports rejected by `check-boundaries.mjs` |
| Imports/variables/UI/styles ordering | `check-source-order.mjs` reports “Source order OK” |
| Agent code journal | `.agents/journal/2026-11-14-initial.md`, `-research`, `-reference-verification`, `-principal-review`, `-implementation` |
| Unit coverage strictly > 80% | backend `87.3%` statements; frontend statements `96.31%`, branches `87.24%`, functions `97.24%`, lines `96.44%`; both gates fail at `<= 80%` |
| Validated-user posting invariant | API: anonymous 401, unvalidated 403 `validation_required`, staff-validated 201, session reflects the change immediately; forged status impossible (server reads SQLite) |
| Any account may interact | unvalidated member like succeeded and repeated idempotently; nonexistent post 404 |
| Protected staff validation | anonymous 401, member 403, CSRF missing 403, staff update 200 + audit row; no HTTP staff creation path |
| Cloudflare security audit | `docs/security-audit.md`; 12-unit ledger and 4-record findings set both pass the skill validators; run artifacts under `/tmp/opencode/security-audit/maulid.id/run-1/` |
| Branch/worktree and small commits | `feat/initialize` checkout, feature-prefixed commits (see `git log`) |
| Push to origin | performed after the final gate; remote head reported in the delivery message |
| Keep worktree until merge | checkout retained; no cleanup performed |
| Preserve prompt and interpretation | `initial.prompt.txt` unchanged; this plan updated with observed evidence |
| Requested agent roles | Antigravity research dispatches, verified OpenCode DeepSeek V4.1 Flash `high` session, advisory Jev request recorded in the journal; no recursive principal |

Runtime and browser checks: health/ready 200; three surfaces served from the single
binary; cookie `HttpOnly; SameSite=Lax` (Secure in production); CSP and hardening
headers present; design inspector `0 FAIL, 2 WARN` at 1440/768/390 with screenshots
reviewed; `verify-ui.mjs` 15/15 checks including staff login through the real form,
backoffice validation round trip, 390 px overflow checks, feed search and reduced
motion. `govulncheck` and `npm audit --audit-level=high` report no vulnerabilities.

## Remaining limitations

- Registration returns a distinct conflict for an existing email (audit finding F-1,
  low). Fixing it requires a neutral acknowledgement and email verification, which the
  product does not have yet; the per-IP limiter is the compensating control.
- Authentication rate limiting keys on the socket peer address. Behind a reverse proxy
  all clients would share one bucket; deployment topology is not in the repository
  (audit record NV-1, needs validation). Enforce auth rate limits at the edge and
  document the deployment assumption.
- The public page has no moving picture: no honest clip of these communities was
  available, so none was used rather than substituting stock footage.
