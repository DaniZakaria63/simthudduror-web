# maulid.id

Community news about Maulid Rasulullah events in Indonesia: a public reading site,
a JSON API for the mobile client, and a protected backoffice for staff.

One Go process serves all three surfaces. The React public site and backoffice are
built into the same binary, so production runs one artefact against one SQLite file.

## Stack

| Layer | Choice |
| --- | --- |
| Server | Go 1.27.4, standard `net/http` router, `log/slog` |
| Storage | SQLite (`modernc.org/sqlite`, pure Go, no cgo) |
| Web | React 19 + TypeScript, Vite 8, React Router 7 |
| Tests | `go test` (with `-race`), Vitest + Testing Library, V8 coverage |
| Security | HTML-escaped React output, CSP and hardening headers, cookie + CSRF and bearer sessions, `govulncheck` |

## Quick start

Requirements: Go 1.27.4+, Node.js 22+ (the web manifest pins `node >=22`), `make`.

```sh
make deps            # go mod download + npm ci (locked)
make build           # builds web/dist, embeds it, builds bin/maulid and bin/staff-bootstrap
make staff-bootstrap EMAIL=staf@example.com NAME="Staf Maulid"   # prompts for a password
make run             # http://127.0.0.1:8080
```

`make run` uses `MAULID_ENV=development` and `./data/maulid.db`. The staff
bootstrap is an operator-only CLI action; there is deliberately no HTTP route that
creates staff accounts, so a fresh deployment cannot be claimed by a visitor.

Frontend development uses Vite with a proxy to the Go server:

```sh
make run       # terminal one: API on :8080
make run-dev   # terminal two: Vite on :5173
```

## Verification

```sh
make test        # go test ./... + npm test
make coverage    # backend statements >80% and frontend statements/branches/functions/lines >80%
make lint        # gofmt check, go vet, feature-boundary and source-order checks
make race        # go test -race ./...
make security    # govulncheck ./... + npm audit --audit-level=high
```

Coverage gates fail when the measured value is 80% or lower, as required by the
project plan. The latest measured values are recorded in
`.agents/journal/2026-11-14-implementation.md`.

## Configuration

| Variable | Default | Notes |
| --- | --- | --- |
| `MAULID_ADDR` | `:8080` | Listen address |
| `MAULID_DB_PATH` | `./data/maulid.db` | SQLite file; parent directories are created `0700`, the file `0600` |
| `MAULID_ENV` | `production` | `development` allows plain-HTTP localhost cookies |
| `MAULID_SESSION_TTL` | `168h` | Between 5m and 720h |
| `MAULID_TRUSTED_ORIGINS` | empty | Comma-separated extra browser origins for cookie writes |
| `MAULID_COOKIE_SECURE` | derived | Cannot be disabled in production |
| `MAULID_RATE_LIMIT_MAX` | `12` | Login/registration attempts per window per IP |
| `MAULID_RATE_LIMIT_WINDOW` | `5m` | Between 1s and 1h |

Production configuration fails closed: an invalid environment, an out-of-range
duration, or an attempt to disable `Secure` cookies in production stops startup.

## Repository layout

```text
cmd/server/              composition root: config, app, HTTP server, graceful shutdown
cmd/staff-bootstrap/     operator-only first-staff CLI
internal/app/            wires features, health checks and the embedded frontend
internal/account/        registration, login, sessions, staff bootstrap
internal/post/           public feed, validated publishing, idempotent likes
internal/staff/          account listing, validation changes, audit trail
internal/platform/       authn, config, database, httpx, ratelimit (shared plumbing)
internal/webui/          embedded React build and SPA fallback handler
web/src/features/        landing, account, posts, backoffice (each owns api/state/ui/styles)
web/src/shared/          cross-feature primitives, API client, tokens, fonts
web/scripts/             feature-boundary and source-order checks, coverage gate
scripts/                 backend coverage gate
docs/                    plan, architecture, research and security audit
.agents/journal/         dated record of agent work, checks and decisions
```

## API

`/api/v1` is the same API the mobile client uses. Browser clients authenticate with
the `maulid_session` cookie plus the `X-CSRF-Token` header for unsafe methods;
mobile clients send `Authorization: Bearer <token>` and receive the token only when
they request it (`X-Client: mobile`), never as a cookie.

| Method | Path | Access |
| --- | --- | --- |
| `GET` | `/healthz`, `/readyz` | public |
| `POST` | `/api/v1/auth/register` | public; always creates an unvalidated member |
| `POST` | `/api/v1/auth/login` | public, rate-limited |
| `POST` | `/api/v1/auth/logout` | authenticated |
| `GET` | `/api/v1/auth/session` | optional; 401 without a session |
| `GET` | `/api/v1/posts?q=` | public feed, optional viewer for `liked_by_me` |
| `POST` | `/api/v1/posts` | authenticated **and** currently validated |
| `PUT`/`DELETE` | `/api/v1/posts/{id}/like` | any authenticated account, idempotent |
| `GET` | `/api/v1/staff/accounts` | staff only |
| `PATCH` | `/api/v1/staff/accounts/{id}/validation` | staff only, CSRF, audited |
| `GET` | `/api/v1/staff/audit` | staff only |

Every authenticated request reloads the account from SQLite, so losing validation
removes publishing rights on the next request without a re-login.

## Documentation

- `docs/plans/initial.prompt.txt` — the original prompt, preserved verbatim
- `docs/plans/initial.plan.md` — interpretation, requirements and acceptance evidence
- `docs/plans/initial.design.md` — public page art direction
- `docs/architecture.md` — module boundaries, request lifecycle and data model
- `docs/security-audit.md` — Cloudflare security-audit run, findings and limits
- `docs/design/` — mandatory Nextdoor-derived UI reference, components and review checklist
- `docs/skills.md` — project agent skills extracted from the prompt's references
- `docs/research/` — reference matrix, raw retrievals and image provenance
