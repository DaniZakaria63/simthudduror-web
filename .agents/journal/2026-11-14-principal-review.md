# Principal verification and audit preparation

## Requested implementation model

`orca-team verify-session --session ses_orcahc14b3cdd7983b407 --json` returned exit 0 and status `high`: observed assistant-message provider `opencode-go`, model `deepseek-v4.1-flash`, variant `high` match the requested configuration. This proves the executing model/variant rather than only launcher input. The implementation dispatch `ctx_ae0a40627a67` remains active; it has not reported completion and must not be duplicated.

## Cloudflare audit preparation

The principal downloaded and read the pinned skill's `RECONNAISSANCE.md`, `HUNTING.md`, `VALIDATION-AND-REPORTING.md`, and ordinary attack classes. Relevant companion files, schemas, and validator sources are locally available at `/tmp/maulid-security-skill`, revision `c1c8a8c1471069fb0e188eeaff69b8e8db6564a8`. This is preparation, not a completed audit of the future application.

The full audit requires four independent reconnaissance assignments, deterministic source-boundary coverage units, focused hunting, a post-wave critic and a distinct final-clean critic for the standard profile, independent candidate/record validation, and schema-validated findings/coverage plus reports. Formal shared audit artifacts must be outside the target and principal-owned. Findings require concrete source traces; confirmations additionally need bounded local observed effects, not checklist deviations. Missing facts are recorded honestly.

`bwrap` and `unshare` are installed. A trusted `/usr/bin/true` probe succeeded in a fresh all-namespace Bubblewrap environment with read-only system mounts, an empty environment except PATH, and no external network. This establishes that namespace/mount isolation is usable; it does not by itself establish the entire required sandbox contract. Any audit target execution must additionally use scratch-only writes, explicit resource/disk/time limits, allowlisted safe environment, no dependency fetching, and the required parent-controlled evidence promotion. Ordinary worker development checks and screenshots are separate evidence and are not represented as sandboxed audit executions.

Application source and final source ref do not yet exist, so reconnaissance/hunting have not started. Final acceptance still depends on the requested implementation, >80% coverage gates, runtime/browser evidence, source audit/remediation, reviewable commits and final origin push.

## Independent image-source review

Queried Wikimedia Commons' `imageinfo` API for all four research-proposed file titles in one read-only request. All four exist; the API confirms the exact original image paths, creators and licenses in the note: `Baiturrahman_Mosque.jpg` (durentiga, CC BY 2.0), `Masjid_Istiqlal_Jakarta.jpg` (AR El Syarif, CC BY-SA 4.0), `Masjid_Menara_Kudus.jpg` (PL09Puryono, CC0), and `Masjid_Agung_Demak.jpg` (Hastosuprayogo, CC BY-SA 4.0). The worker was informed. This verifies provenance metadata only; successful local downloads, actual image inspection, contextual presentation and visible attribution remain implementation evidence to review.

The worker subsequently selected its own verified four-file set, recorded in `docs/research/image-provenance.json`: Istiqlal (AR El Syarif), a different Baiturrahman photograph (ROFIQOHHASAN), a historical 2022 Maulid gathering (Indonesiagood), and Menara Kudus (PL09Puryono). Principal opened and visually inspected all four downloaded WebP files under `web/public/images/`. They show the documented contextual subjects. The provenance/credits for the actual selected set must be used; the earlier research candidates are not interchangeable attribution.

## Audit sandbox preflight

A trusted Python preflight (not application code) succeeded under Bubblewrap with all namespaces unshared, capabilities dropped, target and system/tool mounts read-only, `/proc` read-only, a 32 MiB tmpfs assigned to `/scratch`, and no host home mounted. Environment was cleared and populated only with safe `PATH`, scratch-local `HOME`, `TMPDIR`, `PWD`, and `LC_ALL=C.UTF-8`. A 10-second wall timeout and inner `prlimit` imposed CPU 10 seconds, address space 512 MiB, processes 64, output file size 1 MiB and file descriptors 64.

Observed assertions: scratch write succeeded; a target write was prevented; host SSH directory was absent; a socket could not connect through the isolated network; a trusted child process could still spawn; and the configured environment names/resource limits matched exactly. No target application code ran, no data-store/process was probed, and no target files were changed.

Two initial preflight attempts were corrected: applying the low process limit before namespace creation prevented the host-side timeout from spawning (exit 126), and Bubblewrap adds a safe PWD variable that needed explicit allowlisting (assertion exit 1). Applying limits inside the namespace and explicitly setting PWD/locale produced exit 0. Future audit checks must keep limits inside the namespace, adapt memory/disk limits only as needed for the bounded check, and use trusted parent-side artifact promotion before retaining target-produced evidence. Namespace-local scratch files from this preflight were discarded automatically.

## Early implementation source review

Read the initial platform HTTP/authentication/config/rate-limiter/database packages and account model/password/session repository source while the worker continues implementing. This was an in-progress correctness review, not the formal completed-source Cloudflare audit. No application code was executed or edited by principal.

Sent mailbox `msg_645efca09936` to the worker with a definite functional defect: `authn.Optional` invokes the same response-writing `authenticate` helper as required authentication, so missing/invalid credentials emit 401 before the public handler. Requested silent optional resolution and HTTP response assertions for anonymous/invalid-cookie requests.

Also requested regression-backed corrections for variable-width RFC3339Nano timestamps compared as SQL TEXT, an authentication limiter whose live key map could exceed its stated bound, production insecure-cookie overrides, same-origin scheme comparison and explicit private SQLite file/directory permissions. These are review requests against evolving code, not independently validated security findings. Final source and tests must demonstrate their resolution before acceptance.

Read the subsequently added account service: registration explicitly constructs an unvalidated member, creates an account/session transaction, hashes opaque credentials before storage, and session authentication joins trusted current account data. These are source observations, not yet passing runtime evidence. HTTP adapters, domain posting checks, staff authorization and completed unit coverage remain to be reviewed.

## Parent-side audit evidence helper

Created trusted audit tooling outside the application at `/tmp/maulid-security-skill/promote.py`, with checks in `test_promote.py`. This is principal-owned evidence tooling, not application code or part of the app coverage percentage. The helper retains non-inheritable directory descriptors opened by no-follow component traversal before any future target execution; only allowlisted relative regular single-link files can be promoted within per-file/cumulative byte bounds. It rechecks source identity/metadata after bounded reads, walks destination parents by descriptor without following links, and creates destination files exclusively. It does not run target code or decide that processes have terminated; the principal must establish termination before calling it.

Observed `python3 -m unittest -v test_promote.py` exit 0, seven tests passed. Cases cover exact-byte copies/exclusive creation; absolute/traversal/empty/unlisted rejection; source symlink/hardlink/directory/FIFO rejection; symlinked source and destination parents; per-file/cumulative limits; real nested destinations; and descriptor closure. No application source or data was modified by these checks. Formal audit artifacts and findings still await a reviewable application revision.

## Targeted implementation regression check

Current source resolves optional authentication silently before the handler and leaves required authentication responsible for 401 responses. The timestamp formatter now uses fixed-width UTC nanoseconds, production refuses a false Secure-cookie override, database creation uses private file/directory permissions, and the limiter map has an explicit cap. Final tests must still prove the boundary cases and the completed application behavior.

To verify the identified optional-auth/HTTP concern after the edit, principal ran `go test -p 1 ./internal/platform/authn ./internal/platform/httpx` inside the established Bubblewrap isolation. Observed exit 0: authn passed (0.002s) and httpx passed (0.003s). This is a targeted in-progress regression check, not complete coverage, race testing or the formal final audit. Output was observed directly; no target-generated files were retained as formal audit artifacts.

Execution policy: all namespaces unshared, capabilities dropped, repository mounted read-only at `/target`, Go 1.27.4 toolchain read-only at `/toolchain`, existing dependency cache read-only at `/modcache`, `/proc` read-only, scratch-only tmpfs (512 MiB), no external network or dependency downloads. Safe explicit environment used scratch HOME/TMPDIR/GOCACHE/GOPATH, read-only GOMODCACHE, GOTOOLCHAIN=local, GOPROXY=off, GOSUMDB=off, GOTELEMETRY=off, GOENV=off, GOFLAGS=-mod=readonly, CGO_ENABLED=0, GOMAXPROCS=2, GOMEMLIMIT=512MiB, safe PATH/PWD/locale. Inner limits: CPU 60 seconds, address space 4 GiB, process count 128, file size 64 MiB, descriptors 256; outer wall timeout 120 seconds. The namespace and scratch were discarded on exit; the host source/database/cache were not writable.

## Account adapter and delivery preparation

The current-session adapter now uses the persisted session's actual expiry and CSRF token instead of inventing an expiry or rotating CSRF on a read. Staff bootstrap also has a separate account-create persistence path. Principal identified one remaining DTO consistency case in the evolving code: current-session `View` omits account CreatedAt and therefore serializes a zero date. Requested metadata consistency tests across registration/login/session and sent the passing targeted-check results to the worker in mailbox `msg_a777ee5bed0b`.

Verified read-only Git remote access with `GIT_SSH_COMMAND='ssh -o BatchMode=yes -o ConnectTimeout=10' git ls-remote --heads origin`, exit 0. Remote main and development both point to initial commit `902b3a08669d42d706c23372e4cff0a87a29fbe7`; feat/initialize does not yet exist remotely. No push occurred. Successful remote read is not proof of a completed push or write permission; final delivery must still push the reviewed head and verify the remote branch matches it.

## Remaining reference evidence

Dispatched an independent Antigravity follow-up task `task_ce6807c04b69`, dispatch `ctx_1d64f9328130`, receipt `1794640948866-8d7fbc75df672bdc-research`, with observed turn start. It may write only `docs/research/initial.reference-verification.md` and its isolated journal. It cannot edit application code or disturb the active OpenCode worker.

The remaining factual gaps are actual public access/content versus snippets for the two named Medium articles, precise primary-source React/Vite/Vitest toolchain/testing guidance, and independent source metadata for the worker's actual selected image set. The spec prohibits paywall/authentication bypass, invented full-text reading, unchanged retry loops and wholesale third-party HTML copies. An accurately established access limitation is a valid research result and does not silently become a claim of reading an unavailable article. No final acceptance or source-audit claim follows from this dispatch.

## Reference verification and post-service review

Antigravity follow-up completed successfully and provided `initial.reference-verification.md`. Principal independently created one Orca browser tab, opened both original Medium URLs and read each complete article with `exec ... read`; both results were untruncated. No login or paywall bypass was used. Closed only that principal-created tab afterward. The source-access gap is now resolved; the plan and original source matrix were updated to distinguish initial HTTP403 observations from current full-browser reading. Corrected an overbroad React Node-version statement and an implementation-completion claim in the research note: locked build-tool engines and completed source are the relevant evidence.

Read the new post service/repository/HTTP source and service/HTTP tests. Observed server-side validation lookup on publishing, public feed without private email, parameterized SQL/filter escaping, bounded feed queries, composite-key idempotent likes, and separate unvalidated-reader interaction cases. Actual completed integration and coverage results remain pending. Account/post test files cover browser cookie/CSRF and mobile bearer flows; they are not yet assumed to pass.

Sent mailbox `msg_cc29d230b40d` with remaining test-quality corrections: the account unknown-field fixture needs an actual unknown field; an oversized-JSON limit test should use syntactically valid oversized JSON rather than a first-byte syntax error; password-hiding tests should check serialization; and the persisted account PasswordHash field should defensively omit JSON serialization consistently with its comment. These are in-progress correctness/privacy review requests, not formal audit findings. The worker retains all application/test edits and must run the full suites and gates.
