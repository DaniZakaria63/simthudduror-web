# Initial Monolith Research & Source Matrix (Revised)

**Date**: 2026-11-14  
**Author**: Antigravity Research Worker (Supervised by Codex)  
**Target Repository**: maulid.id  
**Scope**: Minimal working Go + React feature-modular monolith architecture, security audit methodology, canonical domain vocabularies, toolchain recommendations, and contextual image provenance.

**Follow-up status**: This matrix records the initial HTTP retrievals. `initial.reference-verification.md` supersedes its Medium-access limitations: ordinary public-browser access rendered both complete articles, and both Antigravity and the principal read them without credentials or bypass. The approved plan selects only the two validation states `unvalidated` and `validated`; the example `pending`/`rejected` values below are not implementation requirements. Actual selected image provenance is in `image-provenance.json` and the follow-up report, rather than the original candidate list below.

---

## 1. Source Matrix

| # | Reference URL | Access Status & Date / Revision | Retrieval Mode & Artifact | Pertinent Guidance & Architectural Takeaways | Limitations / Caveats |
|---|---|---|---|---|---|
| 1 | `https://www.mdskills.ai/tags/react` | **HTTP 200** (2026-11-14) | Direct scrape (HTML parsed) | Directory of 106+ React agent skills (e.g., UI/UX Pro Max, shadcn/ui v4 MCP, Frontend Developer skill). Highlights modular component patterns, state management separation, and styling best practices. | Curated aggregator/index; individual skill implementations require downstream inspection. |
| 2 | `https://github.com/callstackincubator/agent-skills` | **HTTP 200** (2026-11-14, Commit `59e9d17446a248b21d86f078fcd0bd2d36db1ca4`) | Full Git tree via GitHub API | Agent skills for React/React Native testing (`plugins/vendored/.agents/skills/react-native-testing`), bundle optimization, barrel export management, and test automation. | Primarily tailored to React Native; testing conventions and modular boundaries adapt to React SPA/web. |
| 3 | `https://medium.com/@ashwathvia/react-architecture-a-practical-guide-for-senior-frontend-developers-b5d83bd0ac07` | **HTTP 403 Forbidden** (2026-11-14; Medium WAF/Cloudflare bot challenge) | Public Search Index & Mirror Snippets | Guides senior frontend architecture: organizing code into self-contained feature slices (`components/`, `hooks/`, `services/`, `state/`); strictly isolating server state (TanStack Query) from client state (Zustand); separating business logic from presentation. | Direct full-text blocked by Medium HTTP 403. Findings derived from indexed summaries and technical excerpts. |
| 4 | `https://www.patterns.dev/ai/skills/catalog/` | **HTTP 200** (2026-11-14) | Direct scrape (HTML parsed) | 58 design and performance skills across JavaScript/React (e.g., Module pattern, Provider pattern, Factory pattern, `/vite-bundle-optimization`, `/react-data-fetching`). Highlights code-splitting, tree-shaking, and clean component encapsulation. | High-level pattern catalog; patterns must be implemented within our single-page Vite architecture. |
| 5 | `https://dev.to/leapcell/clean-architecture-in-go-a-practical-guide-with-go-clean-arch-51h7` | **HTTP 200** (2026-11-14) | Full text saved: `docs/research/raw/leapcell_clean_arch.html` | Practical implementation of Uncle Bob's Clean Architecture in Go based on `bxcodec/go-clean-arch`: layers for domain entities/interfaces, use cases (business logic), repository (database adapters), and delivery (HTTP handlers). Enforces inward dependency rule. | Demonstrates layered architecture by technical layer across whole repo; must be adapted to feature-modular packaging (`internal/<feature>/`). |
| 6 | `https://medium.com/@frederich/breakdown-golang-architecture-b71bbcc74b3a` | **HTTP 403 Forbidden** (2026-11-14; Medium WAF/Cloudflare bot challenge) | Public Search Index & Mirror Snippets | Architectural breakdown of Go services adapting Clean Architecture and DDD: directory layout (`model/`, `repository/`, `usecase/`, `delivery/http/`), dependency injection via interfaces, and isolation of domain logic from database drivers. | Direct full-text blocked by Medium HTTP 403. Findings derived from public search index and technical summaries. |
| 7 | `https://reintech.io/blog/go-project-structure-2026-clean-architecture-best-practices` | **HTTP 200** (2026-11-14) | Full text saved: `docs/research/raw/reintech_go_clean_arch.html` | Modern (2026) Go project layout standard: `cmd/api/main.go` for composition root, `internal/` to encapsulate packages, decoupling layers via interfaces, accepting interfaces and returning concrete structs, avoiding `init()` side effects, and structuring for testability. | Covers multi-binary/service layouts; for our monolith, single root under `cmd/server/` or `cmd/api/` is recommended. |
| 8 | `https://developers.cyberagent.co.jp/blog/archives/59647/` | **HTTP 200** (2026-11-14) | Direct scrape (HTML parsed) | Practical Clean Architecture in Go from CyberAgent/ABEMA engineering. Emphasizes pragmatic, non-dogmatic layer separation (domain, usecase, interface, infrastructure) and introducing automated dependency linter (`go-arch-lint`) to enforce boundaries. | Evaluated in large monorepo microservice context; boundary enforcement principles directly inform our monolith package rules. |
| 9 | `https://mcpmarket.com/tools/skills/idiomatic-go-patterns-9` | **HTTP 200** (2026-11-14) | Full text saved: `docs/research/raw/mcpmarket_idiomatic_go.html` | Idiomatic Go design patterns: accept interfaces, return structs; context propagation for timeouts/cancellation; standardized error wrapping (`fmt.Errorf("%w", ...)`, `errors.Is`, `errors.As`); zero-values useful; avoiding goroutine leaks and `init()` side effects. | Focuses on language-level idioms rather than directory organization. |
| 10 | `https://github.com/cloudflare/security-audit-skill` | **HTTP 200** (2026-11-14, Commit `c1c8a8c1471069fb0e188eeaff69b8e8db6564a8`) | Full file saved: `docs/research/raw/cloudflare_skill.md` | Security guidance and audit methodology for codebases/APIs. Located at `skills/security-audit/SKILL.md` (and companion docs `CLIENT-SIDE.md`, `WEB-PROTOCOL-AND-AUTH.md`). Defines the exact 6-phase audit workflow, strict trust boundary discipline, candidate validation, and report schema validation. | Complete 6-phase run is reserved for formal audit mode; in guidance mode, apply the threat model and hunting disciplines directly to design and code review. |
| 11 | `https://docs.typesafe.ai/llms.txt` (and HTTP API / confidence docs) | **HTTP 200** (2026-11-14) | Full files saved: `docs/research/raw/typesafe_llms.txt`, `typesafe_confidence.md`, `typesafe_coding_agents.md` | TypeSafe's Jev is a **System One** decision model (not a chat or code-completion LLM). It evaluates state and returns typed answers (`Choice`, `Score`) with calibrated probabilities and confidence ratings. **Confidence is advisory only** and never acts as authority or permission. System behavior and security controls remain strictly in application code. | Must never be used as a security gate, permission bypass, or code generator. |

---

## 2. Direct Observations vs. Advice

### Direct Observations (Verified System & Source Facts)
- **Local Toolchain State**:
  - `node`: v26.5.0, `npm`: 11.17.0, `pnpm`: 12.16.1, `git`: /usr/bin/git are present on PATH.
  - `go` is **not installed** on PATH.
  - Official current Go releases (from `https://go.dev/dl/?mode=json`):
    - **Go 1.27.4** (latest stable release)
    - **Go 1.26.11** (current supported maintenance release)
    - Distro package candidate on Ubuntu 26.04.1 LTS: `golang-go 2:1.26~1`.
    - *(Note: Go 1.22 and Go 1.23 are obsolete/unsupported upstream).*
  - Repository git status is clean on branch `feat/initialize`, containing only untracked `.agents/` and `docs/`.
- **Domain Constraints & Canonical Vocabulary (AGENTS.md & initial.plan.md)**:
  - Single monolith with 3 surfaces: Landing Page (React), Mobile API (Go REST JSON), Backoffice (React).
  - Validation status concept: **`validation_status`** is the canonical field and concept across the API, database, backoffice, and public site (e.g. `pending`, `validated`, `rejected`).
  - Access rule: Every user/account may view and interact with posts (e.g., likes). Only accounts with `validation_status == "validated"` may create posts. This invariant must be strictly enforced on the server.
  - Component code order convention: Top is imports & variables, middle is UI output (JSX), bottom is styling.
  - Test coverage gate: strictly > 80% statement/branch/line coverage on both backend and frontend.

### Architectural Recommendations
- **Monolith Structure**:
  - Place Go backend in `cmd/server/main.go` as the single composition root.
  - Feature packages under `internal/` (e.g. `internal/post`, `internal/account`, `internal/auth`, `internal/staff`).
  - Each feature owns its domain rules, service/use cases, repository interfaces, and HTTP handlers.
- **Single Deployment**:
  - The React build output (`web/dist` or `frontend/dist`) is embedded into the Go binary using `//go:embed`.
  - The Go server serves the embedded SPA files for web routes and routes `/api/v1/*` to the JSON API.
- **Frontend Architecture**:
  - Feature-sliced structure under `web/src/features/<feature>/` (e.g. `features/posts`, `features/accounts`, `features/backoffice`).
  - Common/global UI components under `web/src/shared/components/`.
  - Server state via TanStack Query; client UI state via Zustand.

---

## 3. Cloudflare Security Audit Guidance (Exact Methodology)

From `skills/security-audit/SKILL.md` (commit `c1c8a8c1471069fb0e188eeaff69b8e8db6564a8`):

### Exact Six Phases (Verbatim Headings & Flow)
1. **Reconnaissance** — map the source, trust boundaries, local build paths, companion selections, prior evidence, and initial deterministic coverage ledger with `RECONNAISSANCE.md`.
2. **Coverage-led hunting waves** — assign isolated hunters from the ledger and collect structured candidate results with `HUNTING.md`, `ATTACK-CLASSES.md`, and the selected domain companions (`CLIENT-SIDE.md`, `WEB-PROTOCOL-AND-AUTH.md`).
3. **Candidate validation** — consolidate fingerprints and give every candidate to a fresh source verifier as defined in `VALIDATION-AND-REPORTING.md`.
4. **Structured output** — write all final `confirmed`, `needs_validation`, and `rejected` records to `findings.json`; validate it with `report-schema.json` and `validate-findings.cjs`, and validate the coverage claim with `validate-coverage-ledger.cjs`.
5. **Independent record verification** — use fresh agents to verify final source claims and reconcile corrections or state changes.
6. **Target-neutral report** — derive `REPORT.md`, `FINDINGS-DETAIL.md`, and `NEEDS-VALIDATION.md` from the final records, with no live-probe instructions.

### Application to maulid.id
- **Trust Boundary Verification**: Ensure the server does not rely on client-supplied headers or UI state for `validation_status`. The check must be performed against trusted server persistence during post creation.
- **Client-Side Sinks**: Ensure React renders event text safely without `dangerouslySetInnerHTML`.
- **Authentication & CSRF**: Protect API endpoints with session tokens / Bearer tokens and enforce SameSite cookies on browser-driven requests.

---

## 4. Advisory Role of TypeSafe Jev

From `docs.typesafe.ai` documentation:
- Jev is a **System One** model that performs fast, typed decisions (`Choice`, `Score`) based on a structured application state.
- **Non-Authority Principle**: Jev outputs probabilities and confidence scores (0.0 to 1.0). **Confidence is never authorization or permission**.
- Jev cannot generate code or text responses.
- In `maulid.id`, Jev may provide advisory scoring (e.g., content categorization or advisory spam scoring), but all security-critical checks (enforcing `validation_status`, role authorization, and data mutations) must be deterministic code running in the Go server.

---

## 5. Toolchain & Installation Guidance

- **Frontend**: Node `v26.5.0`, npm `11.17.0`, and pnpm `12.16.1` are available. React 19 / Vite / Vitest are supported.
- **Backend (Go)**:
  - Official supported versions: **Go 1.27.4** (stable) or **Go 1.26.11** (maintenance).
  - Official installation approach:
    1. Download official archive from `https://go.dev/dl/go1.27.4.linux-amd64.tar.gz`.
    2. Extract to `/usr/local/go` (or `~/.local/go`).
    3. Add `/usr/local/go/bin` to `PATH` in environment configuration.
  - Distribution package alternative on Ubuntu 26.04: `sudo apt-get update && sudo apt-get install -y golang-go` (installs `2:1.26~1`).
  - *(Per instructions, no packages or software have been installed by the research worker).*

---

## 6. Contextual Indonesian Mosque & Cultural Imagery Provenance

For the public landing page (purely contextual cultural photography; not claimed as actual platform events):

1. **Masjid Raya Baiturrahman (Banda Aceh)**
   - **Source**: Wikimedia Commons ([File:Baiturrahman_Mosque.jpg](https://commons.wikimedia.org/wiki/File:Baiturrahman_Mosque.jpg))
   - **Direct Image URL**: `https://upload.wikimedia.org/wikipedia/commons/c/c5/Baiturrahman_Mosque.jpg`
   - **Creator / Artist**: `durentiga` (Flickr / Wikimedia Commons)
   - **License**: Creative Commons Attribution 2.0 Generic (CC BY 2.0)
   - **Usage / Context**: Hero background or cultural showcase of historical Indonesian Islamic architecture.

2. **Masjid Istiqlal (Jakarta)**
   - **Source**: Wikimedia Commons ([File:Masjid_Istiqlal_Jakarta.jpg](https://commons.wikimedia.org/wiki/File:Masjid_Istiqlal_Jakarta.jpg))
   - **Direct Image URL**: `https://upload.wikimedia.org/wikipedia/commons/f/f8/Masjid_Istiqlal_Jakarta.jpg`
   - **Creator / Artist**: AR El Syarif
   - **License**: Creative Commons Attribution-ShareAlike 4.0 International (CC BY-SA 4.0)
   - **Usage / Context**: Iconic national mosque architecture highlighting modern Indonesian community spaces.

3. **Masjid Menara Kudus (Central Java)**
   - **Source**: Wikimedia Commons ([File:Masjid_Menara_Kudus.jpg](https://commons.wikimedia.org/wiki/File:Masjid_Menara_Kudus.jpg))
   - **Direct Image URL**: `https://upload.wikimedia.org/wikipedia/commons/8/85/Masjid_Menara_Kudus.jpg`
   - **Creator / Artist**: PL09Puryono
   - **License**: Creative Commons CC0 1.0 Universal Public Domain Dedication
   - **Usage / Context**: Historical cultural syncretism representing traditional Javanese Islamic heritage.

4. **Masjid Agung Demak (Central Java)**
   - **Source**: Wikimedia Commons ([File:Masjid_Agung_Demak.jpg](https://commons.wikimedia.org/wiki/File:Masjid_Agung_Demak.jpg))
   - **Direct Image URL**: `https://upload.wikimedia.org/wikipedia/commons/c/ce/Masjid_Agung_Demak.jpg`
   - **Creator / Artist**: Hastosuprayogo
   - **License**: Creative Commons Attribution-ShareAlike 4.0 International (CC BY-SA 4.0)
   - **Usage / Context**: One of the oldest mosques in Indonesia, symbolizing traditional Nusantara wooden tiered-roof architecture.
