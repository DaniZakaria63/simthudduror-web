# Reference Verification & Provenance Report

**Date**: 2026-11-14  
**Author**: Antigravity Research Worker (Supervised by Codex)  
**Task ID**: `task_ce6807c04b69` / Dispatch `ctx_1d64f9328130`  
**Target Repository**: `maulid.id`  
**Branch**: `feat/initialize`  
**Artifact Path**: `docs/research/initial.reference-verification.md`

---

## 1. Executive Summary

This report establishes verifiable factual findings across three priority research areas requested by the principal:
1. **Named Medium Sources**: Exact access behavior, observed full text, and corroborated architectural concepts for the Ashwathvia and Frederich Blessy articles.
2. **Frontend Toolchain & Testing Documentation**: Primary source requirements from official React, Vite, Vitest, and Testing Library documentation (Node support, SPA setup, coverage thresholds/semantics, and behavior-oriented testing).
3. **Selected Image Provenance**: Authoritative metadata queried directly from Wikimedia Commons `imageinfo` API for the four selected project assets in `docs/research/image-provenance.json`.

---

## 2. Priority Medium Sources Verification

### 2.1 Access Status & Transport Behavior

| Source | Direct HTTP (curl / automated client) | Full Browser Session (Orca Headless Browser) | Observed Outcome |
|---|---|---|---|
| **Ashwathvia**: React Architecture | **HTTP 403 Forbidden** (Cloudflare bot challenge / WAF mitigation) | **HTTP 200 OK** (Interactive DOM rendered via `orca tab create`) | Full article readable without credentials or paywall bypass. |
| **Frederich Blessy**: Golang Architecture | **HTTP 403 Forbidden** (Cloudflare bot challenge / WAF mitigation) | **HTTP 200 OK** (Interactive DOM rendered via `orca tab create`) | Full article readable without credentials or paywall bypass. |

**Observation**: Standard command-line clients (`curl`) trigger Cloudflare bot mitigation challenges (`cf-mitigated: challenge`, `content-type: text/html`, status 403). However, when accessed via an ordinary browser instance (Orca browser runtime), Cloudflare mitigation resolves naturally without any login, paywall barrier, or authentication bypass. The full text of both articles is publicly readable.

---

### 2.2 Ashwathvia: *React Architecture: A Practical Guide for Senior Frontend Developers*

- **Canonical URL**: `https://medium.com/@ashwathvia/react-architecture-a-practical-guide-for-senior-frontend-developers-b5d83bd0ac07`
- **Author**: Ashwathvia (`https://medium.com/@ashwathvia`)
- **Published Date**: November 29, 2025 (4 min read)
- **Observed Content & Architectural Breakdown**:
  - **Core Principles**:
    1. *Separation of Concerns*: UI, state, domain logic, networking, and side-effects must reside in different layers.
    2. *Single Responsibility*: Each component or module should do one thing well.
    3. *Predictable Data Flow*: Deterministic state flow and event-driven patterns.
    4. *Composition over Inheritance*: Small, reusable, composable component primitives.
    5. *Convention over Configuration*: Naming conventions and folder structures to eliminate decision fatigue.
  - **Recommended Directory Structure & Feature-Based Architecture**:
    - Recommends feature slicing where each feature folder owns its presentation, state, hooks, and services (e.g., `features/invoices/components/`, `features/invoices/hooks/`, `features/invoices/services/`).
    - Cross-cutting or shared primitives live in `components/ui/` or `shared/`.
  - **State Segregation**:
    - Categorizes state into three distinct layers:
      1. *Server State*: Data fetched from APIs (caching, revalidation). Recommended: TanStack Query (React Query).
      2. *Client Global State*: Cross-feature UI states (modals, auth sessions, themes). Recommended: Zustand.
      3. *Local Component State*: Ephemeral form or toggle state (`useState`, `useReducer`).
  - **Logic Encapsulation**:
    - Business rules and data transformation must not leak into React JSX components. Use custom hooks (`useInvoices.ts`) and pure utility functions (`transformInvoiceResponse.ts`).
  - **Architectural Flow**:
    `UI Component → Custom Hook → API Service → Transformation Layer → Global/Server State → UI`

---

### 2.3 Frederich Blessy: *Breakdown Golang Architecture*

- **Canonical URL**: `https://medium.com/@frederich/breakdown-golang-architecture-b71bbcc74b3a`
- **Author**: Frederich Blessy (`https://medium.com/@frederich`)
- **Published Date**: March 25, 2023 (3 min read)
- **Observed Content & Architectural Breakdown**:
  - **Evaluated Architecture Archetypes**:
    The author evaluates several Go community architectures:
    1. `katzien/go-structure-examples` (`https://github.com/katzien/go-structure-examples`)
    2. `disintegration/bebop` (`https://github.com/disintegration/bebop`)
    3. `netlify/gocommerce` (`https://github.com/netlify/gocommerce`)
    4. `bxcodec/go-clean-arch` (`https://github.com/bxcodec/go-clean-arch`)
  - **Recommended Layout (Clean Architecture + DDD)**:
    - Recommends adapting `bxcodec/go-clean-arch` with 4 concentric domain layers:
      1. *Model Layer*: Independent domain entities and core business rules.
      2. *Repository Layer*: Interfaces and persistence adapters abstracting storage engines.
      3. *Usecase Layer*: Business logic and transaction coordination across domain entities.
      4. *Delivery Layer*: Transport boundary (e.g. `delivery/http/` handlers and routing).
    - Top-level directory layout:
      ```text
      go-service/
      ├── config/
      ├── console/
      ├── db/
      ├── delivery/
      │   └── http/
      ├── model/
      ├── repository/
      ├── usecase/
      ├── utils/
      ├── go.mod
      ├── go.sum
      └── main.go
      ```
  - **Core Principles**:
    - Strict adherence to the Single Responsibility Principle (SRP) in Go functions and methods.
    - Testability achieved by decoupling use cases from concrete databases via Go interfaces.
  - **Scaffolding Tool**: Mentions author's generator package `github.com/notblessy/bless` (`bless tolong-generate`).

---

### 2.4 Corroboration & Project Applicability

| Author Recommendation | Maulid.id Project Adoption / Alignment | Architectural Classification |
|---|---|---|
| Feature-based modularity (`features/<name>/`) | Approved for implementation (`web/src/features/<feature>/` and `internal/<feature>/`); acceptance requires completed source review. | **Approved architecture** (prompt & plan requirement). |
| State segregation (TanStack Query + Zustand) | Permitted but optional: plan establishes that minimal React hooks suffice for the initial foundation slice. | **Advisory / Optional** (not a hard gate). |
| 4-layer Clean Architecture (`model`, `repository`, `usecase`, `delivery`) | Adopted inside feature packages (`internal/account/`, `internal/post/`) rather than flat root layers. | **Adapted monolith pattern** (feature-modular Clean Architecture). |

---

## 3. Official React, Vite, and Vitest Documentation Verification

### 3.1 Node.js Runtime Requirements
- **Vite 6 / 7 (Official Docs `vite.dev`)**:
  - Explicit Node requirement: **Node.js version 20.19+ or 22.12+**.
  - Current workspace runtime: Node `v26.5.0` (exceeds minimum prerequisite).
- **React 19 (Official Docs `react.dev`)**:
  - This application renders React in the browser. The selected build/test packages determine the development Node requirement; do not infer a universal Node minimum for browser React itself. The final locked manifests and their engine requirements are authoritative for the project.

### 3.2 Vite React SPA Configuration
- Standard single-page application setup uses `@vitejs/plugin-react` (or `@vitejs/plugin-react-swc`).
- In a Go monolith deployment, the compiled SPA assets (`web/dist/`) are served statically or embedded directly via Go `//go:embed` into the server binary, routing `/api/v1/*` to backend handlers and falling back to `index.html` for client-side routing.

### 3.3 Vitest Coverage Configuration & Threshold Semantics
Official Vitest documentation (`vitest.dev/config/#coverage`) defines coverage options powered by `@vitest/coverage-v8` or `@vitest/coverage-istanbul`:

- **Threshold Specification**:
  ```ts
  // vite.config.ts or vitest.config.ts
  import { defineConfig } from 'vitest/config'

  export default defineConfig({
    test: {
      coverage: {
        provider: 'v8',
        reporter: ['text', 'json', 'html'],
        include: ['src/**/*.{ts,tsx}'],
        exclude: ['src/**/*.d.ts', 'src/**/types.ts', 'src/main.tsx'],
        thresholds: {
          statements: 80.01, // or >80 gate
          branches: 80.01,
          functions: 80.01,
          lines: 80.01,
          // perFile: true // Optional: enforce per individual file
        }
      }
    }
  })
  ```
- **Coverage Inclusion & Exclusion Semantics**:
  - `coverage.include`: By default, Vitest measures only files that are imported during test execution. To prevent untested/dead code from escaping measurement, explicit file globs (e.g. `['src/**/*.{ts,tsx}']`) must be provided.
  - `coverage.exclude`: Patterns specifying files omitted from calculation (e.g. type definitions, test setup files, build scripts).
  - `test.exclude` vs `test.coverage.exclude`: `test.exclude` controls test file discovery, while `test.coverage.exclude` controls source instrumentation.
  - **Threshold Gate Behavior**: Setting a threshold (e.g. `80`) causes Vitest to fail the build process with a non-zero exit code if measured coverage is less than the configured percentage.

### 3.4 Testing Library Behavior-Oriented Guidance
Official Testing Library documentation (`testing-library.com` Guiding Principles):
- **Core Axiom**: *"The more your tests resemble the way your software is used, the more confidence they can give you."*
- **Behavior over Implementation**:
  - Tests must interact with the rendered DOM nodes via accessible queries (`getByRole`, `getByLabelText`, `getByText`) rather than querying component internal state, props, instance methods, or CSS selectors.
  - Event dispatch should simulate real user behavior via `@testing-library/user-event`.
  - Refactoring internal component implementation without changing external behavior must not break the test suite.

---

## 4. Image Provenance Verification (Selected Set)

Authoritative metadata retrieved from the Wikimedia Commons MediaWiki API (`prop=imageinfo&iiprop=extmetadata|url|dimensions`):

| File Path in Repo | Wikimedia File / Page ID | Title / Object Name | Creator / Artist | Date | License & Code | Direct Wikimedia URL | Attribution & Licensing Obligations |
|---|---|---|---|---|---|---|---|
| `assets/masjid-istiqlal-jakarta.jpg` | [File:Masjid Istiqlal Jakarta.jpg](https://commons.wikimedia.org/wiki/File:Masjid_Istiqlal_Jakarta.jpg)<br>(ID: `181378467`) | Masjid Istiqlal Jakarta | **AR El Syarif** | 2026-01-07 | **CC BY-SA 4.0** | `https://upload.wikimedia.org/wikipedia/commons/f/f8/Masjid_Istiqlal_Jakarta.jpg` | **Attribution Required**; **Share-Alike Required**: Derivative works or modifications must be licensed under CC BY-SA 4.0. |
| `assets/masjid-baiturrahman-aceh.jpg` | [File:Masjid Baiturrahman Aceh.jpg](https://commons.wikimedia.org/wiki/File:Masjid_Baiturrahman_Aceh.jpg)<br>(ID: `166419886`) | Masjid Baiturrahman Aceh | **ROFIQOHHASAN** | 2025-05-31 | **CC BY-SA 4.0** | `https://upload.wikimedia.org/wikipedia/commons/3/3b/Masjid_Baiturrahman_Aceh.jpg` | **Attribution Required**; **Share-Alike Required**: Derivative works or modifications must be licensed under CC BY-SA 4.0. |
| `assets/peringatan-maulid-nabi-2022-di-masjid-nurul.jpg` | [File:Peringatan Maulid Nabi 2022 di Masjid Nurul Huda Pulo Lor-13.jpg](https://commons.wikimedia.org/wiki/File:Peringatan_Maulid_Nabi_2022_di_Masjid_Nurul_Huda_Pulo_Lor-13.jpg)<br>(ID: `155309401`) | Peringatan Maulid Nabi 2022 di Masjid Nurul Huda Pulo Lor-13 | **Indonesiagood** | 2022-10-22 | **CC BY 4.0** | `https://upload.wikimedia.org/wikipedia/commons/3/35/Peringatan_Maulid_Nabi_2022_di_Masjid_Nurul_Huda_Pulo_Lor-13.jpg` | **Attribution Required**: Must credit creator (`Indonesiagood`), link to license, and indicate changes. No share-alike restriction. |
| `assets/masjid-menara-kudus.jpg` | [File:Masjid Menara Kudus.jpg](https://commons.wikimedia.org/wiki/File:Masjid_Menara_Kudus.jpg)<br>(ID: `15611448`) | Masjid Menara Kudus | **PL09Puryono** | 2011-06-20 | **CC0 1.0** (Public Domain Dedication) | `https://upload.wikimedia.org/wikipedia/commons/8/85/Masjid_Menara_Kudus.jpg` | **Public Domain**: Attribution not legally required, though credit is good practice. |

**Summary of Image Obligations**:
1. Two assets (`Masjid Istiqlal Jakarta` and `Masjid Baiturrahman Aceh`) carry **CC BY-SA 4.0**. If cropped, compressed, or transformed into WebP, the transformed assets must remain under CC BY-SA 4.0, and visible attribution to AR El Syarif and ROFIQOHHASAN must be maintained on the landing page or credits modal.
2. One asset (`Peringatan Maulid Nabi 2022...`) carries **CC BY 4.0**. Requires clear attribution to Indonesiagood.
3. One asset (`Masjid Menara Kudus`) is dedicated to **CC0 1.0 Public Domain**. No legal attribution required.
4. All four records in `docs/research/image-provenance.json` match the live Wikimedia API metadata exactly.

---

## 5. Remaining Limitations

1. **Medium Transport Constraints**: Direct automated HTTP client access (e.g. via cURL, Go HTTP clients, or CI scripts) will encounter Cloudflare mitigation HTTP 403. Full-text verification requires browser-mediated execution.
2. **Third-Party Article Status**: The two Medium articles represent personal engineering blogs reflecting established industry patterns rather than canonical formal standards. Project architectural requirements are governed by `initial.prompt.txt` and `initial.plan.md`.
