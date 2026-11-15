---
name: React Feature Architecture
description: Build, extend or review React code in this repository using its feature-sliced architecture — feature ownership, public entry points, component reading order, state layers, API services, styling, and behavior tests above 80% coverage. Use for any landing-page, account, posts or backoffice UI task.
---

# React feature architecture (maulid.id)

This skill encodes how React is written in this repository. It adapts the reference
material listed in `references/sources.md` to this project's monolith: one Vite SPA,
served by the Go binary, with three surfaces (public landing page, account flow,
staff backoffice).

## When to use it

Load this skill before any of these:

- adding or changing a screen, page, component, hook, context or stylesheet;
- adding a feature or moving code between features;
- writing or repairing Vitest / Testing Library tests;
- reviewing React code for architecture, accessibility or state ownership.

Do not use it for Go work (see `go-feature-architecture`) or for the release
workflow (see `monolith-delivery`).

Before any visual, layout or component-design decision, also load
`nextdoor-design` and apply `docs/design/ui-review-checklist.md`. This skill owns code
structure; the design skill owns composition, hierarchy, tokens and the review gate.

## Non-negotiables

1. **Feature-sliced, closed modules.** `web/src/features/<feature>/` owns its
   components, pages, local state, API calls, types and styles. Another feature may
   import it only through `web/src/features/<feature>/index.ts`.
2. **Shared only when shared.** A component used by two or more features lives in
   `web/src/shared/`. Never move a domain-specific component to `shared` “for
   convenience”.
3. **One-way dependencies.** `features/* → shared/*`, never
   `shared → features`, and never feature → feature internals. `App.tsx` may import
   feature entry points only.
4. **Component reading order.** Imports first; then types/constants/state/handlers;
   then the UI output; then exactly one trailing stylesheet import at the bottom.
5. **Logic is not JSX.** Data transformation and request logic live in `api.ts`,
   pure helpers, or a feature hook — not inline in render.
6. **Behavior tests, >80% coverage.** Tests use accessible queries and user events;
   the gate fails at 80% or below.
7. **No unsafe rendering or storage.** `dangerouslySetInnerHTML`/`innerHTML` and
   `localStorage`/`sessionStorage` are rejected by the boundary lint. CSRF tokens stay
   in memory (`shared/api/client.ts`).

## Repository layout

```text
web/src/
  main.tsx                     # createRoot + BrowserRouter; imports tokens/base CSS
  App.tsx                      # routes, shell, session provider, reveal hook
  features/
    landing/                   # Masthead, feed section, ContributionSection
    account/                   # session context, login/register/account pages
    posts/                     # feed hook, PostCard, PostFeed, composer, api
    backoffice/                # staff page, staff api + types
  shared/
    api/client.ts              # typed fetch wrapper, ApiError, in-memory CSRF
    components/                # Button, Field, LoadingState, StatusMessage, header/footer
    hooks/useReveal.ts         # scroll-reveal lifecycle
    styles/tokens.css          # fonts, colour/type/space tokens
    styles/base.css            # reset, base typography, motion helpers
    format.ts                  # date/initial formatting
  test/                        # setup, fetch fixtures
```

Each feature normally contains:

| File | Owns |
| --- | --- |
| `index.ts` | the only public surface: components, hooks, API functions, types |
| `api.ts` | typed calls through `shared/api/client.ts`; no React |
| `types.ts` | domain types shared inside the feature (excluded from coverage) |
| `use<X>.ts` | server-state fetching, aborting, error/loading transitions |
| `session.tsx` | only where a feature owns a cross-component context |
| `components/` | presentational units plus colocated `.css` |
| `pages/` | route-level screens composed from components and hooks |

## Component shape

The exact order the prompt asks for:

```tsx
import { useState } from "react";
import { Link } from "react-router-dom";
import { useSession } from "@/features/account";
import type { Post } from "./types";
import { PostCard } from "./components/PostCard";

type FeedProps = {
  posts: Post[];
};

export function Feed({ posts }: FeedProps) {
  const { account } = useSession();
  const [pendingId, setPendingId] = useState<string | null>(null);

  async function handleLike(post: Post) {
    setPendingId(post.id);
    // …
  }

  return (
    <section aria-labelledby="feed-title">
      <h2 id="feed-title">Kabar komunitas</h2>
      {posts.map((post) => (
        <PostCard key={post.id} post={post} onToggleLike={handleLike} />
      ))}
      <Link to="/akun">Kelola akunmu</Link>
    </section>
  );
}

import "./Feed.css";
```

Rules that go with it:

- The stylesheet import is the last line and there is exactly one per component
  (`web/scripts/check-source-order.mjs` enforces it).
- No inline `style={{}}`; values come from CSS variables in `tokens.css`.
- Class names are BEM-ish and feature-scoped (`post-card__meta`, `masthead__band`).
- Props are named `XProps`; `type` over `interface` is the house style.
- Handlers are named for the action (`handleSubmit`, `toggleLike`), and async work
  sets a pending state and renders feedback.

## Data and state ownership

Layer the state instead of centralising everything:

| Layer | This repository | Rule |
| --- | --- | --- |
| Server state | `useFeed` (posts), `AccountProvider` session load | Fetch in the owning feature's hook/context; use `AbortController`, expose `status`, `error`, `reload` |
| Cross-feature client state | `AccountProvider` / `useSession` | Only genuinely shared state; context lives in the feature that owns the domain |
| Local UI state | `useState` in pages/forms | Form fields, pending flags, filter text |
| Derived state | `useDeferredValue` (search) | Do not duplicate server data into another store |

Data flow:

```text
page → feature hook/context → feature api.ts → shared/api/client.ts → /api/v1
```

- Components never call `fetch` directly.
- `api.ts` returns typed domain objects and throws `ApiError`; message text comes from
  the server envelope or a safe Indonesian fallback.
- Every screen renders explicit `loading`, `empty`, `error`, `success` states; the
  backoffice additionally gates `anonymous` and `non-staff`.
- After a mutation, update the owned server state (`feed.updatePost`) rather than
  refetching everything unless correctness requires it.

The referenced articles recommend TanStack Query and Zustand. They are **not**
required here and are not installed: the current slice is small enough for feature
hooks plus one session context. If a future slice needs caching or cross-feature
stores, introduce them inside the owning feature, document the decision in
`.agents/journal/`, and do not create a second source of truth for session or feed
state.

## Forms, links and accessibility

- Every input has a `label` (via `shared/components/Field.tsx`), plus hint/error ids
  wired through `aria-describedby`.
- Use real `<button type="submit">`, real `<Link>` for navigation, and visible focus.
- Errors use `role="alert"`; loading states have a visible label.
- The public copy is Indonesian; validation vocabulary is exactly
  `Belum divalidasi` / `Sudah divalidasi` in the UI and `unvalidated` / `validated`
  in code.
- Pages must remain usable at 390 px without horizontal overflow.

## Styling

- Tokens first: `--bg`, `--ink`, `--muted`, `--accent`, `--field-*`,
  `--step-body`, `--step-head`, `--margin`, `--section-gap`.
- One type family (EB Garamond, locally hosted); two sizes carry most of the page.
- Flat colour, no gradients, no shadows at rest, no decorative rules; section change
  comes from tone/colour fields.
- Responsive rules are written per component; test 1440 / 768 / 390.

## Testing

- Runner: Vitest + jsdom + Testing Library; setup in `web/src/test/setup.ts`.
- Mock at the network boundary with `mockFetch` from `web/src/test/fixtures.ts`,
  not by stubbing feature modules.
- Query as a user would: `getByRole`, `getByLabelText`, `findByText`; use
  `userEvent`; assert visible outcomes, not internal state.
- Cover the states that matter: loading, empty, error, success, permission gates,
  validation, pending/disabled controls.
- Commands:

```sh
cd web
npm run typecheck
npm test
npm run coverage   # fails when statements/branches/functions/lines are <= 80%
npm run lint       # feature boundaries + component source order
npm run build
```

When you add a feature, add its tests in the same feature folder (`*.test.tsx` next
to the unit under test) and keep the coverage gate green without excluding real
logic. `types.ts`, `main.tsx` and `test/**` are the only exclusions.

## Review checklist

- [ ] No deep import into another feature; the public `index.ts` is enough.
- [ ] No domain component in `shared/` unless it is genuinely cross-feature.
- [ ] Imports → logic → UI → stylesheet, one stylesheet import.
- [ ] Fetching lives in a hook/service; components stay declarative.
- [ ] Loading, empty, error and success states exist and are reachable in tests.
- [ ] Labels, roles, focus and 390 px layout verified.
- [ ] `npm run typecheck && npm test && npm run coverage && npm run lint` pass.

See `references/component-shape.md` for a longer annotated example and
`references/testing.md` for the fixture and test patterns. Sources and the exact
adaptations are recorded in `references/sources.md`.
