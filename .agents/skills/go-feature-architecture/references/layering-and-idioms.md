# Go layering and idioms — worked examples

## Request flow

```text
main → app.New
  ├─ database.Open(path)                 SQLite + schema
  ├─ account.NewRepository(db) → NewService(repo, ttl)
  ├─ post.NewRepository(db)  → NewService(repo, validationReader{accountRepo})
  ├─ staff.NewRepository(db) → NewService(repo)
  └─ handlers.RegisterRoutes(mux, require, optional)
main → http.Server{Handler: SecurityHeaders(AccessLog(Recoverer(mux)))}
```

A request then goes: middleware (resolve principal, CSRF) → handler (decode,
authorize) → service (rules) → repository port → SQLite adapter.

## Rule placement

| Rule | Where | Why |
| --- | --- | --- |
| JSON shape, media type, body size | `httpx.DecodeJSON` | one implementation for every feature |
| CSRF for cookie writes | `authn.Require` | transport concern; every unsafe route inherits it |
| “only staff” | `staff.Service` **and** middleware `require` | defense in depth; service is the last trusted decision point |
| “only validated may publish” | `post.Service.Create` | needs current persistence, not a token claim |
| “any account may like” | `post.Service.Like` | domain rule |
| Audit row with status change | `staff.SQLRepository.SetValidation` transaction | invariant spans two rows |
| Role/validation reload | `account.Service.Authenticate` | revocation takes effect on the next request |

## Interfaces at the consumer

```go
// internal/post/post.go
type Authorizer interface {
    ValidationStatus(ctx context.Context, accountID string) (string, error)
}

// internal/app/app.go
type validationReader struct{ repo account.Repository }

func (v validationReader) ValidationStatus(ctx context.Context, id string) (string, error) {
    status, err := v.repo.ValidationStatus(ctx, id)
    return string(status), err
}
```

The post feature never imports `account`, so it cannot depend on account internals. The
adapter lives in the composition root where both features are already known.

## Errors

```go
var (
    ErrUnauthenticated = errors.New("unauthenticated")
    ErrNotValidated    = errors.New("account not validated")
    ErrPostNotFound    = errors.New("post not found")
)

// infrastructure errors carry context and remain inspectable
return Post{}, fmt.Errorf("insert post: %w", err)
```

Handlers translate sentinels once; unexpected errors become a logged 500 with a generic
message. Never return `err.Error()` to the client.

## Transactions

```go
tx, err := r.db.BeginTx(ctx, nil)
defer tx.Rollback()

// read current row
// no-op if status unchanged (no audit row)
// update + insert audit in the same tx
// commit
```

Rule of thumb: if two writes must both happen or neither, wrap them. Use the transaction
to read the “before” value too, so the audit's `from_status` cannot race.

## Context and time

- `ctx` first on every I/O method, passed straight through to `QueryContext`/`ExecContext`.
- Use `time.Now().UTC()` in services via an injectable `now func() time.Time` for tests.
- Format timestamps with the feature's fixed-width helper before writing SQLite.

## Concurrency and lifecycle

- `http.Server` has read-header/read/write/idle timeouts; `main` waits for
  `SIGINT`/`SIGTERM` and calls `Shutdown` with a 10 s budget.
- The rate limiter is mutex-guarded and bounded; no unbounded goroutines or maps.
- SQLite uses one connection so writes serialize and transactions behave predictably.

## Things not to do

- Do not open a database or start a goroutine from a package `init()`.
- Do not add a global logger singleton inside a feature; accept `*slog.Logger`.
- Do not read `os.Getenv` in a feature; config is injected.
- Do not construct SQL with `fmt.Sprintf`, even for ORDER BY; use fixed clauses.
- Do not let an HTTP handler call another feature's repository; add a port/service call.
