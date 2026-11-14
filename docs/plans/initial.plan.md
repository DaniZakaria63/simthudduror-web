# Initial monolith foundation

Status: principal interpretation drafted; external-reference research is in progress. Implementation acceptance requires the evidence below, not this document alone.

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

Research is currently active as task `task_6a00049242a6`, dispatch `ctx_960a7eaec30e`. The principal must review its findings and resolve remaining uncertainty before dispatching the final implementation scope. Use native Orca mailbox supervision, never launch a second principal or a polling scheduler.

## Decisions still requiring evidence

- Resolve actual access to every named reference and use only available content; unavailable links remain clearly identified.
- Confirm a supported reproducible Go/Node toolchain and record the exact versions and local commands.
- Finalize the auth/session/CSRF and storage details against the researched source/security guidance.
- Choose and record the public landing-page direction before implementation.
- At completion replace this draft status and fill in measured acceptance evidence; any unmet item keeps the goal active.
