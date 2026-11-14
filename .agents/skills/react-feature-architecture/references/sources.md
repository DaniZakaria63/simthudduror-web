# React references and how they were adapted

Read on 2026-11-14. Full access notes are in `docs/research/initial.sources.md` and
`docs/research/initial.reference-verification.md`.

| Source | Access | What was taken | What was not taken |
| --- | --- | --- | --- |
| `https://www.mdskills.ai/tags/react` | HTTP 200, index of 106+ React skills | Modular component patterns, separation of state management, styling practice, skill-driven development | Individual third-party skill implementations (not vendored) |
| `https://github.com/callstackincubator/agent-skills` (commit `59e9d17446a248b21d86f078fcd0bd2d36db1ca4`) | Full tree via GitHub API | Testing discipline, boundary/barrel-export management, keeping feature internals private | React Native-specific setup |
| `https://medium.com/@ashwathvia/react-architecture-a-practical-guide-for-senior-frontend-developers-b5d83bd0ac07` | Full text through a normal browser session (direct HTTP hit a Cloudflare challenge) | Feature slices, server/client/local state separation, logic in hooks/services, composition over inheritance, convention over configuration | Mandatory TanStack Query + Zustand; this slice uses feature hooks + one session context instead |
| `https://www.patterns.dev/ai/skills/catalog/` | HTTP 200, 58-pattern catalog | Module/provider/factory patterns, code-splitting/tree-shaking, component encapsulation | Patterns that need libraries or routing structures this monolith does not use |

## Concrete mapping

| Reference idea | Repository expression |
| --- | --- |
| Feature slice owns UI, state, services | `web/src/features/{landing,account,posts,backoffice}` each with `api.ts`, `types.ts`, hooks/session, components/pages, CSS |
| Public entry point / barrel boundary | `feature/index.ts`; `web/scripts/check-boundaries.mjs` rejects deep imports |
| Server state separated from client state | `useFeed` for posts; `AccountProvider` for session; local `useState` for forms |
| Business logic out of JSX | `api.ts`, `useFeed`, `format.ts`, page handlers |
| Composition over inheritance | small `shared/components` primitives composed by pages |
| Convention over configuration | feature layout, naming, trailing CSS import, enforced by lint |
| Behavior-oriented tests | Testing Library queries + `userEvent`; see `testing.md` |

## Deliberate deviations

- **No TanStack Query / Zustand.** The article recommends them for large apps. This
  repository has one feed endpoint and one session; a feature hook and a context are
  simpler and keep the dependency surface small. If a future slice needs caching,
  pagination or cross-feature client stores, adopt a library inside the owning feature
  and journal the decision.
- **Feature-modular Clean Architecture instead of flat layers.** The Go references
  recommend `model/`, `repository/`, `usecase/`, `delivery/` at the root. Here those
  layers live *inside* each `internal/<feature>/` package, because the prompt requires
  feature-closed modules.
- **The project prompt wins.** Where a reference suggests a different structure
  (e.g. global layers, mandatory query library), this repository follows
  `docs/plans/initial.prompt.txt` and `docs/plans/initial.plan.md`.
