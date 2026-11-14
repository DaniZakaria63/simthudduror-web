# Go testing patterns

## Three layers of tests

| Layer | File pattern | What it proves |
| --- | --- | --- |
| Service/rule | `service_test.go`, `helpers_test.go` | domain rules and sentinel errors with fakes |
| Repository | `repository_test.go` | real SQL against a temp SQLite file |
| HTTP | `http_test.go` | status codes, envelopes, cookies, CSRF, auth gates |
| Wiring | `internal/app/app_test.go` | middleware order, headers, full publish path, health |

## Service tests with fakes

```go
type fakeStore struct {
    accounts []AccountSummary
    audit    []AuditEntry
    gotActor string
}
func (f *fakeStore) ListAccounts(context.Context, ListFilter, int) ([]AccountSummary, error) {
    return f.accounts, nil
}

func TestServiceEnforcesStaffRole(t *testing.T) {
    service := NewService(&fakeStore{})
    if _, err := service.List(context.Background(), "member", ListFilter{}); !errors.Is(err, ErrForbidden) {
        t.Fatalf("member list error = %v", err)
    }
}
```

Keep fakes tiny and assert the rule, not the query.

## Repository tests on real SQLite

```go
func newStaffDB(t *testing.T) *sql.DB {
    t.Helper()
    db, err := database.Open(filepath.Join(t.TempDir(), "staff.db"))
    if err != nil { t.Fatalf("open database: %v", err) }
    t.Cleanup(func() { db.Close() })
    return db
}
```

Seed rows with explicit IDs and fixed timestamps. Exercise:

- unique constraints (duplicate email → `ErrEmailTaken`);
- transaction behaviour (validation update + audit insert; unchanged status writes no
  audit);
- escaping (`Dua_` must not match `DuaX`), ordering and limits;
- malformed persisted rows (e.g. invalid timestamp) → error, not panic.

## HTTP tests

Build the real stack in a fixture:

```go
db := newStaffDB(t)
accountService := account.NewService(account.NewRepository(db), time.Hour)
staffService := NewService(NewRepository(db))
mux := http.NewServeMux()
require := authn.Require(accountService)
optional := authn.Optional(accountService)
account.NewHandler(accountService, ratelimit.New(1000, time.Minute), discardLogger(), false, nil).
    RegisterRoutes(mux, require, optional)
NewHandler(staffService, discardLogger()).RegisterRoutes(mux, require)
server := httptest.NewServer(mux)
```

`httptest.Server` gives real cookies and real `Set-Cookie` handling. Assert:

- anonymous 401, member 403, staff 200 for staff routes;
- CSRF missing → 403, valid → 200 + audit;
- register/login → cookie + `csrf_token`; logout → 204 then session 401;
- unvalidated publish → 403 `validation_required`; validated publish → 201;
- like twice → same like count (idempotent), unknown post → 404;
- unknown JSON field / oversized body / wrong content type → 4xx, not 500;
- the response after a like reflects `liked_by_me`.

Use `slog.New(slog.NewTextHandler(io.Discard, nil))` for silent logs, and a helper
`request(method, path, body, cookie, csrf, bearer)` to keep tests readable.

## App-level test

`internal/app/app_test.go` should cover what unit tests cannot: the actual middleware
order (security headers on every response), `/healthz` and `/readyz`, SPA fallback,
the validated-publish invariant end to end, and that the `post` feature sees validation
changes made through `staff` without re-login.

## Commands

```sh
gofmt -l .                       # must print nothing
go vet ./...
go test ./...
go test -race ./...
./scripts/check-go-coverage.sh   # total statements over internal/... must be >80%
```

Current observed values (2026-11-14): all packages pass, race passes, total statements
87.2–87.3% (the `database` package is the lowest at 76.7% because its remaining branches
are filesystem/ping failures that need privileged faults; the gate is the total, and the
per-package number is reported honestly rather than gamed with exclusions).

## Table-driven tests

Use a table when the same rule has many inputs (statuses, limits, malformed values),
but keep one intent per row and name rows by behaviour:

```go
tests := []struct{ name, status string; wantErr bool }{
    {"canonical unvalidated", "unvalidated", false},
    {"canonical validated", "validated", false},
    {"legacy pending", "pending", true},
}
```
