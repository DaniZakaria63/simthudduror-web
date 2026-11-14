# Go references and how they were adapted

Read on 2026-11-14. Access status and raw artifacts are recorded in
`docs/research/initial.sources.md` (`docs/research/raw/`).

| Source | Access | What was taken | Adaptation |
| --- | --- | --- | --- |
| `https://dev.to/leapcell/clean-architecture-in-go-a-practical-guide-with-go-clean-arch-51h7` | Full text saved | Inward dependency rule; domain entities/interfaces, use cases, repository adapters, delivery layer; testability through interfaces | Layers live inside `internal/<feature>/` instead of flat root folders |
| `https://medium.com/@frederich/breakdown-golang-architecture-b71bbcc74b3a` | Full text through a normal browser session | Four concentric layers (model, repository, usecase, delivery), SRP, dependency injection via interfaces | Same layers, per feature; no generator tool (`bless`) adopted |
| `https://reintech.io/blog/go-project-structure-2026-clean-architecture-best-practices` | Full text saved | `cmd/` composition root, `internal/` encapsulation, interfaces at boundaries, no `init()` side effects, testability | Single binary, SQLite, no multi-service layout |
| `https://developers.cyberagent.co.jp/blog/archives/59647/` | Scraped | Pragmatic layer separation; automated dependency/boundary linting; avoid dogmatic purity | Boundary rules enforced by package layout + review (no `go-arch-lint` dependency added) |
| `https://mcpmarket.com/tools/skills/idiomatic-go-patterns-9` | Full text saved | Accept interfaces/return structs; context propagation; `%w` wrapping; `errors.Is/As`; zero-value usefulness; no goroutine leaks; no `init()` | Applied throughout `internal/*` |

## Reference-to-code map

| Reference idea | Code |
| --- | --- |
| Composition root | `cmd/server/main.go`, `internal/app/app.go` |
| Domain entities + rules | `internal/account/account.go`, `internal/post/post.go`, `internal/staff/staff.go` |
| Repository ports | `account.Repository`, `post.Repository`, `staff.Store` |
| Adapters | `*_sqlite.go` files |
| Delivery layer | `http.go` per feature |
| Dependency injection | `app.New` constructs repositories → services → handlers |
| Inward dependency rule | `post.Authorizer` is declared by the consumer and implemented in `app` |
| Idiomatic errors | sentinel errors per feature + `fmt.Errorf("%w")` for infrastructure |
| Context propagation | every repository/service method takes `ctx` first |
| Boundary enforcement | package layout + `go vet` + code review (no extra linter dependency) |

## Deliberate deviations

- **Not flat layers.** The two architecture articles show root-level `model/`,
  `repository/`, `usecase/`, `delivery/`. The prompt requires feature-closed modules,
  so each feature owns those layers. The dependency direction is the same; only the
  packaging differs.
- **No `go-arch-lint`.** The CyberAgent article uses it to enforce boundaries. This
  repository keeps the dependency rule in the package layout and tests; adding a tool
  would be considered only if the package graph grows.
- **SQLite, not a service database.** The persistence adapter is `modernc.org/sqlite`
  so the single binary needs no cgo and no external service.
- **One HTTP server.** No separate API gateway or internal RPC between features; the
  monolith calls in-process ports.

## Prompt precedence

If a reference conflicts with `docs/plans/initial.prompt.txt` or
`docs/plans/initial.plan.md` (for example by recommending microservices or flat
layers), the prompt and plan win. Record the decision in `.agents/journal/` when it is
structural.
