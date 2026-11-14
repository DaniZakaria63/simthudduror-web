# Security-audit sources

## Upstream skill

- Repository: <https://github.com/cloudflare/security-audit-skill>
- Commit: `c1c8a8c1471069fb0e188eeaff69b8e8db6564a8`
- License: MIT, Copyright (c) 2025–2026 Cloudflare, Inc. (full text in
  `upstream/LICENSE`)
- Vendored files: `SKILL.md`, `RECONNAISSANCE.md`, `HUNTING.md`,
  `ATTACK-CLASSES.md`, `CLIENT-SIDE.md`, `WEB-PROTOCOL-AND-AUTH.md`,
  `VALIDATION-AND-REPORTING.md`, `report-schema.json`,
  `validate-coverage-ledger.cjs`, `validate-findings.cjs`.
- Not vendored: companions for domains this monolith does not contain (AI/LLM, desktop/
  mobile IPC, memory safety/binary, protocols/RPC, cloud/IAM, data isolation, supply
  chain). Fetch them from the same revision if the product grows into those domains.
- The upstream repository also ships validator tests; they are not required to run an
  audit and are not vendored here.

## How this repository applied it

| Upstream requirement | Repository evidence |
| --- | --- |
| Six-phase full audit | `docs/security-audit.md` (run 1, quick profile) |
| Pinned revision and local copy | this directory; `docs/research/raw/cloudflare_skill.md` |
| Deterministic coverage ledger | run 1 ledger validated with `validate-coverage-ledger.cjs` (12 units) |
| Schema-valid findings, independent verification | `findings.json` validated with `validate-findings.cjs` (4 records); re-read verification pass |
| Bounded local evidence, dummy principals | source- and repository-test-grounded records; no live probing |
| Sandbox before target execution | preflight recorded in `.agents/journal/2026-11-14-principal-review.md`; run 1 needed no target execution |
| Target-neutral report | `REPORT.md`, `FINDINGS-DETAIL.md`, `NEEDS-VALIDATION.md` derived from findings |

## Decisions and deviations

- **Profile `quick`.** The target is small and the prompt required a check, not an
  exhaustive pentest. The report states partial coverage explicitly; a `standard` run is
  justified before a public launch.
- **Accepted low finding.** Registration email enumeration was not silently “fixed” by
  changing the UX; it is documented with a compensating control because the proper fix
  needs email verification.
- **Needs-validation allocation.** The proxy-identity issue is a deployment fact, not a
  missing source fact; it is recorded with a local and deployment validation plan rather
  than assigned severity.
- **Topology.** One sequential reader performed the phases with a separate re-read
  verification pass; this is disclosed as a limitation instead of claiming independent
  agents.

## Related project material

- `docs/research/initial.sources.md` — full source matrix and access status
- `docs/research/initial.reference-verification.md` — Medium access, toolchain and image
  provenance verification
- `docs/security-audit.md` — durable run summary, controls reviewed, limitations
- `.agents/journal/2026-11-14-principal-review.md` — sandbox preflight evidence
- `.agents/journal/2026-11-14-implementation.md` — final checks and delivery
