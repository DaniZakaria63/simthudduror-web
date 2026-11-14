# Architecture

maulid.id is one deployable monolith with three surfaces that share a domain and a
database:

| Surface | Served by | Root |
| --- | --- | --- |
| Public landing page | embedded React build | `/` |
| Backoffice | embedded React build | `/backoffice`, `/masuk`, `/daftar`, `/akun` |
| Mobile/JSON API | Go handlers | `/api/v1/*`, `/healthz`, `/readyz` |

Production is a single process and a single artefact: React is built once, embedded
with `go:embed`, and served by the same `http.Server` that owns the API. There is no
second origin, separate frontend service, or runtime proxy. That keeps cookies,
CSP, and CSRF behaviour identical for the two browser surfaces and preserves one
release unit.

## Composition

`cmd/server` is a small composition root:

1. `config.Load()` reads and validates the environment and fails closed.
2. `app.New()` opens SQLite, builds each feature's service and handler, mounts
   routes, and wraps the mux in recoverer → access log → security headers.
3. `main` starts `http.Server` with read/write/idle timeouts, header limits, and a
   10-second graceful shutdown on `SIGINT`/`SIGTERM`.

`cmd/staff-bootstrap` is the only way to create a staff account. It opens the same
database, writes the account directly through `account.Service.BootstrapStaff`, and
never exposes an HTTP route.

## Feature modules

Backend features live under `internal/<feature>` and own their domain rules,
persistence port, SQLite adapter, HTTP adapter, and tests:

- `internal/account` — accounts, passwords, sessions, CSRF token issue, login rate
  limiting, and the operator-only staff bootstrap. Exports `Service`, `Handler`,
  `Repository`, and the canonical `Role`/`ValidationStatus` types.
- `internal/post` — public feed, publishing, and likes. It does not import
  `internal/account`; it accepts a narrow `Authorizer` port whose
  `ValidationStatus(ctx, accountID)` is implemented in `internal/app` over the
  account repository. That keeps the validation rule owned by the post use case
  while the account table stays owned by the account feature.
- `internal/staff` — the staff-only account projection, validation changes and the
  audit log. It imports `internal/account` only for the canonical status type.

Shared infrastructure is deliberately small and transport-level:

- `internal/platform/authn` — the request principal, cookie/bearer extraction, and
  `Require`/`Optional` middleware. `Require` writes 401/403; `Optional` never writes.
- `internal/platform/httpx` — JSON envelope, strict body decoding, `Problem` errors,
  origin checks, client IP, security headers, recoverer and access log.
- `internal/platform/database` — SQLite open, pragmas and the idempotent schema.
- `internal/platform/config`, `internal/platform/ratelimit` — configuration and the
  bounded authentication limiter.
- `internal/webui` — the embedded build, immutable asset caching, and SPA fallback for
  client routes.

No feature reads another feature's tables through its repository. Cross-feature
behaviour goes through an explicit interface or the canonical types, which is why
the validation lookup is a one-method port instead of a direct account-repository
call.

## Frontend modules

`web/src/features/<feature>/` owns everything the feature needs: `api.ts` (typed
fetch calls), `types.ts`, local state/context, components, pages, and its CSS. A
feature exposes a public entry point (`index.ts`); other features and `App.tsx`
import only that entry point. `web/src/shared/` holds cross-feature primitives: the
API client, buttons, fields, status/loading states, header/footer, formatting,
reveal hook, and design tokens/fonts.

Feature boundaries are executable, not conventional:

- `web/scripts/check-boundaries.mjs` rejects deep imports into another feature,
  `dangerouslySetInnerHTML`, and browser storage APIs.
- `web/scripts/check-source-order.mjs` enforces the requested component shape:
  imports first, then logic/state, UI output next, and exactly one trailing
  stylesheet import.

Components import their stylesheet at the bottom (`import "./X.css"`), so a screen
reads top-to-bottom as import → variables/logic → output → style. Styles live in the
colocated stylesheet instead of JSX.

## Request lifecycle

1. `SecurityHeaders` sets CSP (`default-src 'self'`, no inline script, no framing),
   `nosniff`, `DENY`, referrer/permissions/COOP policies, HSTS in production, and
   `no-store` for `/api/` and health responses.
2. `authn.Require` or `authn.Optional` resolves a credential. Both call
   `account.Service.Authenticate`, which hashes the presented token with SHA-256,
   joins `sessions` to `accounts`, checks expiry/revocation, and returns a principal
   built from the **current** row. A role or validation change therefore applies on
   the next request.
3. Cookie-authenticated unsafe methods must send `X-CSRF-Token`; the middleware
   compares its SHA-256 to the session's stored hash in constant time. Bearer
   clients skip CSRF because their credential is explicit.
4. Feature handlers decode JSON strictly (`application/json`, 32 KiB cap, unknown
   fields and trailing data rejected), apply domain rules, and write either a JSON
   result or the common error envelope:
   `{"error":{"code":"...","message":"..."}}`.
5. Panics become a generic 500; the access log records method, path, status, and
   duration but never bodies, headers, or credentials.

## Data model

SQLite runs with WAL, `busy_timeout=5000`, foreign keys on, one connection, and
private `0700` directory / `0600` file permissions.

| Table | Purpose | Notes |
| --- | --- | --- |
| `accounts` | identity, role, validation | `role IN ('member','staff')`, `validation_status IN ('unvalidated','validated')`, unique case-insensitive email |
| `sessions` | opaque session lookup | primary key is the SHA-256 of the token; stores CSRF token, expiry and revocation time |
| `posts` | public feed | author FK, body, timestamp; descending created-at index |
| `likes` | post interaction | composite primary key `(post_id, account_id)` makes like/unlike idempotent |
| `audit_logs` | staff actions | actor, target, action, from/to validation status, timestamp |

Timestamps are written as fixed-width UTC strings with nanoseconds
(`2006-01-02T15:04:05.000000000Z`), so lexicographic SQL comparisons agree with
chronological order. Schema creation is idempotent (`CREATE TABLE IF NOT EXISTS`) and
runs on every start; the initial deployment has no destructive migration step.

## Domain invariant

Publishing is checked in `post.Service.Create`:

```text
authorID empty                     -> 401 unauthenticated
authorizer lookup fails            -> 401 unauthenticated
current validation_status != validated -> 403 validation_required
```

The authorizer reads the account row inside the same request; a client-supplied
status is never accepted. The HTTP adapter maps those sentinel errors to status
codes, and the React composer only appears for validated accounts — the UI reflects
the rule, it does not enforce it.

Interaction is intentionally broader: any authenticated account can like/unlike a
post regardless of validation. Likes are recorded with `INSERT OR IGNORE` semantics
on a composite key, so repeated requests converge on the same state and counts.

## Frontend state

`AccountProvider` owns the session: it loads `/api/v1/auth/session` on mount, keeps
the CSRF token in memory (never `localStorage`), and exposes `signIn`, `signUp`,
`signOut`, and `reload`. The feed hook owns list/search state with abortable
requests and optimistic-safe updates after a like. Every screen renders explicit
loading, empty, success, and error states, and the backoffice gates anonymous and
non-staff accounts before showing data.
