# Journal: Reference Verification & Provenance

**Date**: 2026-11-14  
**Agent**: Antigravity Research Worker  
**Supervisor**: Codex (Principal)  
**Task ID**: `task_ce6807c04b69`  
**Dispatch ID**: `ctx_1d64f9328130`  

---

## 1. Context & Objectives

The principal dispatched this research verification task to resolve three outstanding source gaps:
1. Verification of the two Medium architecture articles previously marked as HTTP 403 in `initial.sources.md`.
2. Authoritative official documentation verification for React, Vite, Vitest, and Testing Library (Node requirements, SPA monolith integration, coverage threshold semantics, and behavior-driven testing).
3. Provenance verification of the actual four selected image assets in `docs/research/image-provenance.json` using live API queries.

Constraints:
- Work strictly inside `feat/initialize` branch.
- No editing of application code.
- No subagents or helper workers.
- Only modify `docs/research/initial.reference-verification.md` and this journal entry.
- No paywall bypass or fake claims; distinguish observed facts from advice.

---

## 2. Actions & Findings

### Medium Source Verification
1. **Direct HTTP probe**:
   - `curl -sI` against both Medium URLs yielded `HTTP/2 403` with header `cf-mitigated: challenge`.
2. **Orca Browser Probe**:
   - Launched Orca headless browser tab via `/home/dani/.config/orca/linux-orca-cli-shim/orca tab create`.
   - Navigated to `https://medium.com/@ashwathvia/react-architecture-a-practical-guide-for-senior-frontend-developers-b5d83bd0ac07`.
   - Evaluated DOM tree and accessibility snapshot: rendered cleanly with no paywall or bot block. Extracted full outline, 14 architectural principles, folder structure (`features/<feature>/`), and state layers (React Query + Zustand).
   - Navigated to `https://medium.com/@frederich/breakdown-golang-architecture-b71bbcc74b3a`.
   - Evaluated DOM tree and accessibility snapshot: rendered cleanly without paywall or bot block. Extracted author's analysis of 4 Go architectures, recommending `bxcodec/go-clean-arch` with 4 layers (`model/`, `repository/`, `usecase/`, `delivery/http/`).
   - Closed browser tab.

### Official Toolchain & Testing Documentation
1. **Node.js**: Vite documentation requires Node.js 20.19+ or 22.12+. Node v26.5.0 installed on the host satisfies and exceeds requirements.
2. **Vitest Coverage**:
   - `test.coverage.thresholds`: enforces strict gate (lines, statements, functions, branches); exits non-zero if coverage is below threshold.
   - `coverage.include` vs `coverage.exclude`: default only tracks imported files; explicit include glob `src/**/*.{ts,tsx}` ensures all files are counted.
3. **Testing Library**:
   - Guiding principle: test behavior, not implementation details.
   - Deal with rendered DOM nodes via accessible roles (`getByRole`, etc.) rather than component internals.

### Image Provenance Verification
1. Queried Wikimedia Commons MediaWiki API (`prop=imageinfo&iiprop=extmetadata|url|dimensions`) for the four titles:
   - `Masjid Istiqlal Jakarta.jpg` (Page ID `181378467`): AR El Syarif, CC BY-SA 4.0.
   - `Masjid Baiturrahman Aceh.jpg` (Page ID `166419886`): ROFIQOHHASAN, CC BY-SA 4.0.
   - `Peringatan Maulid Nabi 2022 di Masjid Nurul Huda Pulo Lor-13.jpg` (Page ID `155309401`): Indonesiagood, CC BY 4.0.
   - `Masjid Menara Kudus.jpg` (Page ID `15611448`): PL09Puryono, CC0.
2. Verified that all metadata in `docs/research/image-provenance.json` exactly aligns with live Wikimedia Commons API data. Noted CC BY-SA 4.0 share-alike and attribution requirements.

---

## 3. Deliverables

- Detailed research document written to [docs/research/initial.reference-verification.md](file:///home/dani/WebstormProjects/maulid.id/docs/research/initial.reference-verification.md).
- Journal entry recorded in `.agents/journal/2026-11-14-reference-verification.md`.
- No application code or other files modified.
