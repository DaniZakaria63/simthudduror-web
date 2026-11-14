---
name: Security Audit
description: Run or extend the Cloudflare security-audit workflow for this monolith — trust-boundary reconnaissance, coverage ledger, hunting, candidate validation, validated findings and reporting — with this repository's known findings, sandbox rules and stack-specific checks. Use for security reviews, pre-release audits, or changes to auth, sessions, validation or staff access.
---

# Security audit (maulid.id)

This skill adapts the Cloudflare `security-audit` skill (MIT) to this repository. The
complete upstream procedure is vendored under `references/upstream/` at the revision
this project pinned, so the audit is reproducible offline.

- Upstream: <https://github.com/cloudflare/security-audit-skill>
- Revision: `c1c8a8c1471069fb0e188eeaff69b8e8db6564a8`
- Vendored: `references/upstream/SKILL.md` plus its companions, schema, validators and
  `LICENSE`.
- Project record: `docs/security-audit.md`; run artifacts live outside the repository.

## When to use it

- before a release or a production deployment change;
- when auth, sessions, CSRF, authorization, validation, SQL or file serving change;
- when the user asks for a review, audit, pentest or security report;
- when triaging a vulnerability claim.

For a one-line question (“is this cookie flag right?”) use guidance mode: read the
relevant upstream companion and answer, without creating a run directory.

## Modes

| Mode | When | Output |
| --- | --- | --- |
| Guidance | a specific question or review | findings in the conversation, no artifacts |
| Full audit | explicit audit/pentest/report request | run directory, ledger, findings JSON, reports |

This repository's prompt (`docs/plans/initial.prompt.txt`) explicitly requires the
codebase to be checked with the Cloudflare skill, so the initial delivery used **full
audit mode** with profile `quick`.

## Full audit, condensed

Read `references/upstream/SKILL.md` before running; it is the authority. The phases:

1. **Reconnaissance** — map entry surfaces, principals, trust boundaries and controls;
   seed a deterministic coverage ledger. Relevant upstream: `RECONNAISSANCE.md`.
2. **Coverage-led hunting** — one hunter assignment per ledger unit; structured
   candidates with concrete file:line traces. Relevant: `HUNTING.md`,
   `ATTACK-CLASSES.md`, `CLIENT-SIDE.md`, `WEB-PROTOCOL-AND-AUTH.md`.
3. **Candidate validation** — a fresh verifier attempts to refute each candidate from
   source (and bounded local execution when allowed).
4. **Structured output** — `findings.json` with `confirmed` / `needs_validation` /
   `rejected`; validate it and the ledger:

   ```sh
   node references/upstream/validate-coverage-ledger.cjs <run>/coverage-ledger.json
   node references/upstream/validate-findings.cjs <run>/findings.json
   ```

5. **Record verification** — re-check every confirmed record's source trace and
   conditions.
6. **Report** — derive `REPORT.md`, `FINDINGS-DETAIL.md`, `NEEDS-VALIDATION.md`; no
   live-probe instructions.

A run ends only when both validators pass and the reports exist, or when
`run_status: "incomplete"` plus the exact reason and gap are recorded.

### Run directory

Create it **outside** the target, e.g. `~/security-audit-skill/maulid.id/run-<N>`, with
`run-metadata.json`, `architecture.md`, `coverage-ledger.json`, `findings.json`,
`REPORT.md`, `FINDINGS-DETAIL.md`, `NEEDS-VALIDATION.md`. Never commit audit artifacts,
scratch files, agent directories or promoted evidence into the repository.

### Execution safety

Source review is read-only. Any target-controlled execution (build, tests, server,
browser) must run inside an OS-enforced sandbox with:

- no external network (isolated loopback only);
- an empty, allowlisted environment and scratch-local HOME/TMPDIR/caches;
- read-only target and toolchain mounts, scratch-only writes;
- explicit CPU, memory, process, file-size and wall-clock limits.

If those controls are unavailable, do not execute target code: record the gap as
`needs_validation` with a safe validation plan. Use dummy principals and fixtures; never
probe deployed endpoints, real users or shared services.

## This application's threat model

Strongest controls and where to verify them:

| Boundary | Control | Where |
| --- | --- | --- |
| Visitor → registration | forced `member`/`unvalidated`; strict JSON; Argon2id | `internal/account/service.go`, `http.go` |
| Browser ↔ API | `HttpOnly; SameSite=Lax` cookie, session-bound CSRF, origin check on auth writes, no CORS | `internal/platform/authn/authn.go`, `internal/account/http.go` |
| Mobile ↔ API | explicit bearer token, SHA-256 lookup, expiry/revocation | `internal/account/service.go`, `internal/account/repository_sqlite.go` |
| Member → publisher | current `validation_status` reloaded from SQLite | `internal/post/post.go` + `internal/app/app.go` |
| Member → staff | role in middleware **and** service; CLI-only bootstrap | `internal/staff/*`, `cmd/staff-bootstrap` |
| Input → SQL | parameterized statements + escaped `LIKE` | `internal/post/repository_sqlite.go`, `internal/staff/repository_sqlite.go` |
| Text → HTML | React text rendering; boundary lint bans `dangerouslySetInnerHTML` | `web/src`, `web/scripts/check-boundaries.mjs` |
| Static files | cleaned paths inside `embed.FS`; GET/HEAD only | `internal/webui/webui.go` |
| Operator → production | fail-closed config, Secure cookies, HSTS, timeouts | `internal/platform/config/config.go`, `cmd/server/main.go` |

Stack-specific checks that have produced value:

- client-supplied `role`/`validation_status` must never be read;
- CSRF must be required for every unsafe cookie write, including likes and staff PATCH;
- session reload must pick up validation changes immediately;
- `DecodeJSON` limits and unknown-field rejection must cover every write handler;
- SQL search must escape `\ % _` and bind the pattern;
- error envelopes must not leak internal errors; logs must omit bodies/headers/secrets;
- rate limiting must key on a trustworthy client identity (see NV-1 below);
- embedded-file serving must not escape the FS or serve directories.

## Current audit state (run 1, 2026-11-14)

Profile `quick`, all application surfaces, 12 coverage units, 4 records. Summary in
`docs/security-audit.md`.

| Record | Verdict | Note |
| --- | --- | --- |
| `maulid/account/registration-email-enumeration` | confirmed, low | 409 `email_unavailable` reveals an existing address; accepted limitation until neutral registration + email verification exist |
| `maulid/account/rate-limit-client-identity` | needs validation | limiter keys on `RemoteAddr`; behind a proxy all clients share a bucket — confirm topology and edge limiting |
| `maulid/session/plaintext-csrf-at-rest` | rejected | session key is stored hashed; CSRF alone does not authenticate |
| `maulid/post/feed-search-sql-injection` | rejected | bound pattern + escaped metacharacters |

Do not re-report a rejected claim unless the source changed; do not treat
`needs_validation` as a confirmed vulnerability or assign it severity.

## Known limitations of run 1

- One `quick` pass, no prior ledger; it does not exhaust the target.
- The skill's multi-agent topology was executed by one sequential reader with a
  separate re-read verification pass; independence is weaker than separate contexts.
- Not reviewed: dependency provenance beyond `govulncheck`/`npm audit`, hostile-load
  behavior, SQLite file tampering, backups, edge caching/SRI, OS hardening.
- Deployment topology (proxy/CDN/TLS) is not in the repository, which is why NV-1
  remains open.

## Definition of done for a new audit

- [ ] Scope/profile and source ref recorded in `run-metadata.json`.
- [ ] Ledger covers every in-scope surface; `validate-coverage-ledger.cjs` passes.
- [ ] Every candidate has a trace, evidence and a verdict; `validate-findings.cjs` passes.
- [ ] Confirmed records have bounded evidence and the smallest effective fix.
- [ ] `docs/security-audit.md` updated with the new run, findings and limitations.
- [ ] No audit artifacts, secrets or live-probe instructions committed.

See `references/sources.md` for the upstream file map and this repository's audit
decisions.
