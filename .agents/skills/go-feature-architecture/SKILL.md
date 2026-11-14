---
name: Go Feature Architecture
description: Build, extend or review Go backend code in this repository using its feature-modular Clean Architecture — composition root, internal/<feature> packages, repository ports, service rules, HTTP adapters, platform plumbing, SQLite transactions and test/coverage conventions. Use for any API, persistence, auth or staff-backoffice backend task.
---

# Go feature architecture (maulid.id)

This skill encodes the backend conventions of this repository. It adapts the Go
reference material listed in `references/sources.md` to a **monolith with
feature-closed packages**, as the project prompt requires — not the flat
`model/ repository/ usecase/ delivery/` layout some references show.

## When to use it

- adding or changing an API endpoint, domain rule, repository query or migration;
- adding a feature package or moving code between features;
- touching auth, sessions, CSRF, staff authorization or validation;
- writing Go unit, repository or HTTP tests; reviewing backend code.

For React work use `react-feature-architecture`; for releases use
`monolith-delivery`; for security review use `security-audit`.

## Non-negotiables

1. **One composition root.** `cmd/server` loads config, calls `app.New`, starts
   `http.Server` and handles graceful shutdown. No feature builds its own server.
2. **Feature-closed packages.** `internal/<feature>/` owns its domain types,
   repository port, service/use cases, SQLite adapter, HTTP adapter and tests.
3. **Dependencies point inward.** HTTP → service → repository port → adapter.
   Platform packages never import features; features may import other features only
   through an explicit public type or a narrow interface.
4. **The server enforces domain rules.** The UI may reflect a rule, but validation,
   roles and ownership are decided against trusted persistence on every request.
5. **Validated publishing.** `post.Service.Create` reloads the author's
   `validation_status` and rejects anything but `validated`; never trust a client
   claim.
6. **Small, explicit interfaces.** Define a port where it is consumed; accept
   interfaces, return concrete structs. No `init()` side effects.
7. **Tests and gates.** Unit/HTTP/repository tests with statement coverage > 80%,
   `go test -race`, `go vet`, `gofmt`.

## Repository layout

```text
cmd/server/main.go        # process wiring: config, app.New, timeouts, signals
cmd/staff-bootstrap/      # operator-only first staff account (no HTTP route)
internal/app/             # composition: DB, services, handlers, middleware, health
internal/account/         # accounts, passwords, sessions, CSRF issue, staff bootstrap
internal/post/            # feed, validated publishing, idempotent likes
internal/staff/           # staff account projection, validation changes, audit
internal/platform/
  authn/                  # principal, cookie/bearer resolution, Require/Optional, CSRF
  config/                 # env loading, validation, secure defaults
  database/               # SQLite open, pragmas, idempotent schema
  httpx/                  # JSON envelope, strict decode, headers, origin, logging
  ratelimit/              # bounded fixed-window limiter
internal/webui/           # embedded React build + SPA fallback
```

## Feature anatomy

A feature package typically contains:

| File | Responsibility |
| --- | --- |
| `<feature>.go` | domain types, sentinel errors, repository port, service rules |
| `repository_sqlite.go` | SQL adapter implementing the port; parameterized SQL only |
| `http.go` | request/response DTOs, `RegisterRoutes`, error→status mapping |
| `util.go` | IDs, timestamp formatting, LIKE escaping |
| `*_test.go` | service fakes, repository against `t.TempDir()` SQLite, `httptest` |

Example of the dependency style (from `internal/post`):

```go
// The post feature does not import the account repository. It declares the one
// fact it needs, and app wires the adapter.
type Authorizer interface {
    ValidationStatus(ctx context.Context, accountID string) (string, error)
}

func (s *Service) Create(ctx context.Context, authorID, body string) (Post, error) {
    if authorID == "" {
        return Post{}, ErrUnauthenticated
    }
    status, err := s.authorizer.ValidationStatus(ctx, authorID)
    if err != nil {
        return Post{}, ErrUnauthenticated
    }
    if status != "validated" {
        return Post{}, ErrNotValidated
    }
    // validate body, persist, return the stored post
}
```

Cross-feature access is allowed only when the type genuinely belongs to the other
feature's public contract — `staff` imports `account.ValidationStatus` because that is
the canonical vocabulary. It still does not touch `account` tables directly.

## HTTP adapter conventions

```go
func (h *Handler) RegisterRoutes(mux *http.ServeMux, require, optional func(http.Handler) http.Handler) {
    mux.Handle("GET /api/v1/posts", optional(http.HandlerFunc(h.list)))
    mux.Handle("POST /api/v1/posts", require(http.HandlerFunc(h.create)))
}
```

- Use Go 1.22+ method-and-path patterns (`GET /api/v1/...`, `{id}` + `r.PathValue`).
- `require` = authenticated + CSRF for cookie writes; `optional` = attach a principal
  when valid, never reject.
- Decode with `httpx.DecodeJSON` (JSON media type, 32 KiB, unknown fields and trailing
  data rejected) — never `json.NewDecoder(r.Body)` directly.
- Write with `httpx.JSON` / `httpx.Error`; client messages are safe Indonesian text;
  internal errors are logged, never returned.
- Map sentinel errors to status codes in one switch per feature:
  input → 400, unauthenticated → 401, forbidden/validation → 403, not found → 404,
  conflict → 409, else 500.
- Read the principal from `authn.FromContext`; never accept `role`, `account_id` or
  `validation_status` from a request body.

## Persistence rules

- Parameterize every value. `LIKE` search escapes `\ % _` and declares `ESCAPE`.
- Wrap invariants that span rows in one transaction (validation update + audit insert).
- Use composite keys for idempotency (`likes(post_id, account_id)`).
- Timestamps are fixed-width UTC (`2006-01-02T15:04:05.000000000Z`) so TEXT order matches
  time order; format through the feature's `formatTime`.
- SQLite opens with WAL, `busy_timeout`, foreign keys, one connection, `0700` directory
  and `0600` file. Schema creation is idempotent `CREATE TABLE IF NOT EXISTS`.
- Do not add a second migration mechanism or a second database driver without a
  recorded decision.

## Idioms used here

- `context.Context` is always the first parameter on I/O methods.
- Errors: sentinel values per feature for expected outcomes; wrap infrastructure
  errors with `fmt.Errorf("...: %w", err)`; compare with `errors.Is/As`.
- Constructors return `*Service` / `*Handler`; dependencies are interfaces or concrete
  pointers injected by `app.New`. No package-level mutable state.
- IDs and tokens come from `crypto/rand`; never `math/rand`.
- Logging uses `log/slog` with key/value fields; never log bodies, headers, tokens or
  passwords.
- Configuration is read once at startup, validated, and fails closed.

## Testing

```sh
go test ./...
go test -race ./...
go vet ./...
./scripts/check-go-coverage.sh   # >80% statements, internal/...
```

Patterns:

- **Service tests** use small fakes for the repository/authorizer; assert the rule and
  the sentinel error, not SQL.
- **Repository tests** open a real SQLite file in `t.TempDir()` via `database.Open`;
  exercise constraints, transactions, ordering, escaping and malformed rows.
- **HTTP tests** build the real handler stack with `httptest.Server`, real services and
  temp SQLite; log to `io.Discard`; assert status, envelope code and body.
- **App tests** (`internal/app/app_test.go`) cover wiring, middleware order, headers,
  health, CSRF and the full validated-publish path.
- Boundary cases worth keeping: anonymous vs unauthenticated, unvalidated vs validated,
  CSRF present/missing, expired/revoked sessions, unknown fields, oversized bodies,
  invalid canonical statuses, `LIKE` metacharacters, nonexistent IDs.

## Review checklist

- [ ] Rule enforced in the service, not only in the handler or UI.
- [ ] No cross-feature repository/table access; ports are narrow and consumer-owned.
- [ ] All SQL parameterized; no string-built statements; search escaped.
- [ ] JSON strict; limits respected; safe error envelopes.
- [ ] Principal identity and current validation reloaded from persistence.
- [ ] Transactions where invariants span rows; audit row written with the change.
- [ ] `gofmt -l`, `go vet`, `go test`, `go test -race`, coverage gate pass.

See `references/layering-and-idioms.md` for the reference-to-code mapping and
`references/testing.md` for concrete fixture patterns. Sources and access status are in
`references/sources.md`.
