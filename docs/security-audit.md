# Security audit

Run: Cloudflare `security-audit` skill, revision
`c1c8a8c1471069fb0e188eeaff69b8e8db6564a8`, profile `quick`, scope = all application
source (`cmd/`, `internal/`, `web/src/`, `web/scripts/`, manifests, build files).

The complete run artifacts live outside the repository at
`/tmp/opencode/security-audit/maulid.id/run-1/`:

- `run-metadata.json` — target, source identity, profile, scope, execution policy
- `architecture.md` — principals, trust boundaries, source-to-sink paths
- `coverage-ledger.json` — 12 coverage units, validated
- `findings.json` — 4 records, validated against the skill's schema
- `REPORT.md`, `FINDINGS-DETAIL.md`, `NEEDS-VALIDATION.md`

Both validators pass:

```sh
node <skill-dir>/validate-coverage-ledger.cjs <run-dir>/coverage-ledger.json  # PASS: 12 coverage units valid
node <skill-dir>/validate-findings.cjs <run-dir>/findings.json                # PASS: 4 findings valid
```

## Summary

| Record | Verdict | Severity | Subject |
| --- | --- | --- | --- |
| F-1 | confirmed | low | Registration conflict discloses whether an email has an account |
| NV-1 | needs validation | — | Auth limiter keys on the socket peer; a reverse proxy collapses clients |
| R-1 | rejected | — | Plaintext CSRF token at rest enables session takeover |
| R-2 | rejected | — | SQL injection through feed search |

### F-1 — Registration email enumeration (low, accepted)

`POST /api/v1/auth/register` returns `201` for a new address and
`409 email_unavailable` for one that already has an account
(`internal/account/http.go:178-179`, sentinel from
`internal/account/repository_sqlite.go:30`; asserted by
`internal/account/http_test.go:186`). An unauthenticated visitor can test addresses
one request at a time. Only existence is disclosed.

Remediation is a neutral acknowledgement for both outcomes, ideally with email
verification; that changes the product's registration contract and needs an email
flow, so it is recorded as an accepted limitation. The per-IP authentication limiter
is the compensating control. If the product later adds email verification, fix this
before exposing the endpoint broadly.

### NV-1 — Rate-limit identity behind a proxy (deployment fact missing)

`httpx.ClientIP` uses `RemoteAddr` and ignores forwarding headers
(`internal/platform/httpx/httpx.go:121`); the auth guard keys the limiter on that
address (`internal/account/http.go:131`). Behind a reverse proxy, all clients share
one 12-per-5-minute bucket: sign-in can be denied for everyone and per-client
brute-force protection disappears. Whether this applies depends on topology that is
not in the repository.

Deployment requirement until resolved: terminate at a known edge, enforce
authentication rate limits there, and record the assumption next to
`MAULID_TRUSTED_ORIGINS`. If the app must read a proxy-supplied client address, add an
explicit trusted-proxy allowlist — never trust `X-Forwarded-For` unconditionally.

## Controls reviewed and found sound

- **Passwords** — Argon2id (19 MiB, t=2, p=1), 16-byte salt, constant-time compare,
  dummy verification for unknown accounts.
- **Sessions** — 32-byte random tokens, only SHA-256 persisted as the lookup key,
  expiry/revocation checked per request, current account row reloaded each request,
  server-side logout.
- **CSRF** — session-bound token, constant-time comparison on every unsafe
  cookie-authenticated write; bearer clients are exempt by design.
- **Authorization** — publishing reloads `validation_status` from SQLite through a
  port; staff role enforced in middleware *and* service; validation change and audit
  insert share one transaction; no HTTP route grants staff.
- **Input** — strict JSON (media type, 32 KiB, unknown fields and trailing data
  rejected), parameterized SQL, escaped `LIKE`, bounded queries and bodies.
- **Output** — React text rendering, boundary lint forbids `dangerouslySetInnerHTML`
  and browser storage, CSP without inline script, `nosniff`, `DENY`,
  `frame-ancestors 'none'`.
- **Transport and operations** — Secure cookies enforced in production, HSTS in
  secure mode, `HttpOnly`, `SameSite=Lax`, no CORS, read/write/idle timeouts, graceful
  shutdown, private SQLite permissions, fail-closed config, logs without bodies or
  headers.

## Dependency results

Observed in this revision:

- `go run golang.org/x/vuln/cmd/govulncheck@v1.9.0 ./...` → `No vulnerabilities found.`
- `npm audit --audit-level=high` (web) → `found 0 vulnerabilities`

A dependency scan is not a source audit; the source review above is the primary
evidence.

## Coverage and limitations

One `quick` pass does not exhaust the target. Not reviewed: dependency provenance
beyond the scanners, hostile-load behaviour, SQLite file tampering, backups, edge
caching/SRI, and OS hardening. No prior audit ledger exists.

The skill's default topology delegates phases to independent agents. This run executed
the full quick-phase sequence in one reader with a separate re-read pass for candidate
validation and record verification, so hunter/verifier independence is weaker than
separate contexts. Re-check F-1 and NV-1 independently before treating them as final.

## Reproducing the audit posture

Source review needs no execution. Any target-controlled execution must run under an
OS sandbox with all namespaces unshared, no external network (isolated loopback only),
a read-only target and toolchain, scratch-only writes, an empty allowlisted
environment, and explicit CPU/memory/process/file-size/wall limits. The earlier
preflight in `.agents/journal/2026-11-14-principal-review.md` established that the
local `bwrap` setup enforces those controls; this run did not need to execute target
code because every record is source- or repository-test-grounded.
