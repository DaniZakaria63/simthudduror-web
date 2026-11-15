# Component contracts — operational extract

Full anatomy: `docs/design/nextdoor-components.md`.

## Shell

- **Web:** left navigation + central feed + optional right rail. Feed is widest.
- **Mobile:** bottom navigation, ≤5 labelled destinations; one primary action.
- **Rail:** utility only; max 2–3 cards; collapses below the feed on narrow screens.

## Post card

```text
[avatar] Author · status · relative time                [overflow]
Body / headline
[image/event when supported]
[Like] [Comment] [Share]                               [quiet utility]
```

- Identity row first; body second; actions third.
- Name 14px/600; metadata 12–13px muted; body 15–16px/1.5–1.6.
- Card 12px radius, 16–20px padding, hairline border, no heavy shadow.
- States: default, hover, focus-visible, pressed-like, pending, error, loading.
- Never render an empty media frame or fake comment/share action.

## Identity row

- Avatar/initial 32–40px; name; `Sudah divalidasi` / `Belum divalidasi`; relative time.
- Staff: `Dikelola staf maulid.id`.
- No paid-check icon, follower count or trophy.

## Account status / validation

- Show the canonical state as text: `Belum divalidasi`, `Menunggu peninjauan staf`,
  `Sudah divalidasi`.
- One sentence explains what the state means and what happens next.
- If a checklist is used, every step must be real and resumable; no fake progress.
- After validation, state what is now possible (`Tulis kabar`).
- Do not use a gold check, shield, trophy, KYC funnel, address/map or location flow.

## Reaction bar

- Order: like → comment → share → overflow.
- Icon + label or icon + accessible name; count only if real.
- Pressed uses `aria-pressed` and primary colour.
- Hide unsupported actions.

## Feed header

- Section title; search; one filter row; one primary action.
- “Tulis kabar” for validated, “Minta validasi” for unvalidated.
- Sticky only when it does not cover the first card.

## Composer

- Author identity and status visible at top.
- Body field; rules line; cancel + publish.
- Unvalidated: publish disabled, explanation + link to account/contribution.
- Mobile: full-screen/sheet; keyboard-aware; no covered publish button.
- States: empty, typing, submitting, success, validation error, server error.

## Event card (future)

- Date block/gutter; title; `WIB` time; venue; organiser; one action.
- Real attendance only; no invented counts.
- Maulid agenda uses a date gutter.

## Publisher card (future)

- Publisher lockup; verification label; headline; optional 16:9 image; section label.
- Requires a source/featured model in the API.

## Overflow menu / moderation

- Context-sensitive items, each with a one-line consequence:
  `Sembunyikan — Hapus dari feed`, `Laporkan — Kirim untuk ditinjau`.
- Destructive items use the danger colour and require a confirm that names the effect.
- `Blokir` only if DMs/follows exist; they do not at launch.
- Show `Laporkan` only when reports are stored and staff can review them.
- Guidelines: publish as a readable page with four short rules; link from the composer.

## Menus and sheets

- Overflow menu items state consequences.
- Sheet: 16px top radius, handle, title, close, focus trap, Escape.
- Report/hide only when a moderation path exists.

## States

| State | Copy | Action |
| --- | --- | --- |
| Empty feed | “Belum ada kabar…” | write/validate |
| Filtered empty | “Belum ada kabar yang cocok…” | clear filter |
| Loading | “Memuat kabar…” | none |
| Error | “Kabar tidak dapat dimuat.” | “Coba lagi” |
| Unvalidated | staff-review explanation | account/validation link |

## Reuse

- Reuse `Button`, `Field`, `StatusMessage`, `LoadingState` before new primitives.
- Shared components live in `web/src/shared/components/`.
- New component → colocated CSS + behavior tests.
- No new icon library/font/colour token without a recorded decision.
