# React testing patterns

## Stack

- Vitest 5, jsdom, Testing Library (`@testing-library/react`, `user-event`),
  `@testing-library/jest-dom` matchers.
- Config: `web/vite.config.ts` (`test` block), setup: `web/src/test/setup.ts`.
- Fixtures: `web/src/test/fixtures.ts` (`mockFetch`, `jsonResponse`,
  `accountFixture`, `staffFixture`, `postFixture`).

## What to test

For every screen:

| State | What proves it |
| --- | --- |
| loading | the labelled loading state is visible before the request resolves |
| empty | the feature's empty message renders when the API returns `[]` |
| success | real content appears via accessible queries |
| request error | the server message or the safe fallback renders in an alert |
| permission gate | anonymous / non-staff see the gate and a working link |
| validation | unvalidated vs validated accounts see the right affordances |

For every mutation: pending/disabled control, success feedback, failure feedback,
and the state the server returned (not an invented count).

## Mocking at the network boundary

```ts
import { jsonResponse, mockFetch, postFixture } from "@/test/fixtures";

mockFetch(({ url, method }) => {
  if (url.includes("/auth/session")) {
    return jsonResponse({ account: staffFixture, csrf_token: "csrf", expires_at: "…" });
  }
  if (url.includes("/staff/accounts") && method === "PATCH") {
    return jsonResponse({ account: { ...member, validation_status: "validated" }, audit: null });
  }
  return jsonResponse({}, 404);
});
```

Why: the feature's `api.ts`, error mapping and hook transitions all stay in the path.
Stubbing the hook instead would skip the behaviour users depend on.

## Querying like a user

```tsx
const user = userEvent.setup();
render(<AccountPage />, { wrapper: MemoryRouter });

expect(await screen.findByText("Belum divalidasi")).toBeInTheDocument();
await user.click(screen.getByRole("button", { name: "Suka (0)" }));
expect(await screen.findByText("Kabar disukai.")).toBeInTheDocument();
```

- Prefer `getByRole` with an accessible name, then `getByLabelText`, then text.
- Use `findBy*` for async arrivals; do not wrap in manual `waitFor` unless asserting
  a negative or a sequence.
- Do not assert internal state, props, CSS classes, or mock call counts when a visible
  outcome exists.
- One behavioural assertion set per test; avoid snapshotting whole pages.

## Coverage gate

```sh
cd web
npm run coverage
```

- V8 provider, `include: src/**/*.{ts,tsx}`, threshold 80 for statements, branches,
  functions and lines.
- `node scripts/check-coverage.mjs` re-checks the summary and **fails at 80 or
  below** (the plan requires strictly greater).
- Only `types.ts`, `main.tsx`, `*.d.ts` and `src/test/**` are excluded; do not exclude
  real logic to reach the number.
- Current observed values (2026-11-14): statements 96.31%, branches 87.24%,
  functions 97.24%, lines 96.44%; 69 tests in 13 files.

## Common failure modes

| Symptom | Cause | Fix |
| --- | --- | --- |
| `findByRole` resolves but `toBeInTheDocument` fails | the node was replaced by a re-render (e.g. loading → staff) | wait for a stable staff-data element first, then assert the heading |
| “not wrapped in act” | state update after the test ends | await the user event and the visible result; cleanup is automatic |
| fetch is undefined | test stubbed `global.fetch` without restoring | use `mockFetch`; `setup.ts` calls `vi.unstubAllGlobals()` |
| coverage passes locally but fails CI | test.exclude vs coverage.exclude confusion | remember `test.exclude` only hides tests; only `coverage.exclude` hides source |
