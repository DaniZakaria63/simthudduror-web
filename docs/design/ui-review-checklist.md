# UI review checklist

Every UI change must satisfy this checklist **and** attach evidence. A checked box without
the referenced evidence is not accepted.

## 1. Reference and scope

- [ ] `nextdoor-design` skill loaded and the relevant `docs/design/` sections read.
- [ ] The change names the Nextdoor-derived pattern it uses: feed, post card, composer,
      event, search, trust row, empty state, navigation or sheet.
- [ ] Any deliberate deviation is recorded with a reason in the plan or journal.
- [ ] The feature is achievable with the current API/data model, or the required
      schema/API work is listed as out of scope.

## 2. Composition

- [ ] Web: three-zone composition is respected where the page allows it — left
      navigation (~17%), central reading column (~38–39%), right utility rail (~19%) at
      desktop. The column does not widen to fill when the rail is absent.
- [ ] Web: active navigation uses a light grey fill and filled icon, not a green block.
- [ ] The reading column stays the widest and most prominent zone.
- [ ] The right rail carries utility or context, never the primary feed.
- [ ] Mobile: root shell uses at most four labelled destinations plus one composer action;
      notifications/messages/profile are header utilities, not tabs, until they exist as
      real surfaces.
- [ ] Mobile: bottom navigation has at most five labelled destinations.
- [ ] A text-only post is treated as complete; no empty media frame is rendered.
- [ ] The primary action is visible on the first screen and is singular.
- [ ] Cards are stacked in one column; the page does not become a grid of equal-weight
      boxes.
- [ ] Sticky regions (navigation, composer, filters, bottom bar) do not cover content or
      the keyboard.

## 3. Hierarchy and card anatomy

- [ ] Feed items answer, in order: who / status / when → what → actions.
- [ ] The headline or first body line is the largest text in the item.
- [ ] Metadata is secondary and legible at 12–14px equivalent, not lower than 12px.
- [ ] Relative time is Indonesian (`2 jam lalu`, `kemarin`), with an absolute date in the
      article/detail view.
- [ ] Every post shows the author and canonical validation status
      (`Belum divalidasi` / `Sudah divalidasi`).
- [ ] Validation states are complete: `Belum divalidasi` → `Menunggu peninjauan staf` →
      `Sudah divalidasi`, each with one sentence and the next step.
- [ ] Actions are consistent across cards: like, comment (when available), share (when
      available), overflow.
- [ ] No fake counts, fake authors, fake images or invented activity.

## 4. Components and states

- [ ] The change reuses existing `shared/components` primitives before adding a new one.
- [ ] Every interactive component defines default, hover, focus-visible, active, disabled,
      loading and error states where applicable.
- [ ] The composer opens with the author identity and validation status visible.
- [ ] An unvalidated account sees an explanation and the path to validation, not a dead
      button or a generic error.
- [ ] Empty, loading, error and permission states exist, use Indonesian copy, and offer one
      recovery action.
- [ ] Destructive or moderation actions state their consequence.

## 5. Visual language

- [ ] Colour roles come from tokens; no raw hex values in components.
- [ ] Primary green is used for the primary action and selection, not for decoration.
- [ ] Two primaries, two jobs: green for publishing/positive, `--ui-strong` for operational
      confirmation, red only for destructive actions and errors.
- [ ] Surfaces are white on a light canvas; borders are subtle; shadows are not used to
      fake depth on every card.
- [ ] Corner radii follow the token scale (card vs control vs avatar).
- [ ] Typography uses the project display face for editorial voice and a UI face for
      controls/metadata; no new font is introduced without a recorded decision.
- [ ] Icon+label pairing is used for navigation and primary actions; icon-only controls
      have accessible names.
- [ ] Imagery uses honest community/Indonesian mosque material with credit; no generic
      Middle-Eastern stock imagery.
- [ ] Motion is subtle, respects `prefers-reduced-motion`, and never autoplays media.
- [ ] If dark mode exists, both themes use separate tokens and pass contrast review.

## 6. Accessibility

- [ ] Text contrast meets WCAG AA (4.5:1 body, 3:1 large text/UI).
- [ ] Touch targets are at least 44×44px equivalent.
- [ ] Focus order is logical and focus-visible is obvious.
- [ ] The page works at 200% text zoom without clipping.
- [ ] Headings follow a logical hierarchy; landmarks and labels are present.
- [ ] Keyboard-only operation is verified for menus, sheets, composer and filters.
- [ ] Status messages are announced (`role="status"` or equivalent).

## 7. Responsive evidence

- [ ] Desktop 1440px screenshot.
- [ ] Tablet 1024px (or 768px) screenshot.
- [ ] Phone 390px screenshot.
- [ ] No horizontal scroll at any width.
- [ ] Navigation, composer and feed remain usable at each width.
- [ ] The page is reviewed in a real browser against the built binary or the Vite dev
      server, not only in unit tests.

## 8. Content and tone

- [ ] All user-facing copy is Indonesian, plain and short.
- [ ] The canonical validation terms are exact: `Belum divalidasi`, `Sudah divalidasi`.
- [ ] The copy does not imply official Maulid event coverage, verified religious authority
      or government endorsement.
- [ ] No Nextdoor copy, brand line, neighbourhood terminology or ad format is reused.
- [ ] Error messages say what happened and what to do next.

## 9. Boundaries

- [ ] No advertising, sponsored card, paywall, member-only gate or promo countdown.
- [ ] No streaks, points, badges, leaderboards or engagement scoring.
- [ ] No fintech-style identity verification (ID photo, selfie, document upload).
- [ ] No autoplay video, heavy carousel or decorative gradient that competes with content.
- [ ] No US hyperlocal framing (garage sales, neighbourhood classifieds, local deals).
- [ ] No `dangerouslySetInnerHTML` and no browser storage for credentials.

## 10. Evidence and completion

- [ ] Component/behavior tests added or updated and pass.
- [ ] Coverage remains strictly above 80% for the frontend metrics.
- [ ] `make lint` passes.
- [ ] Screenshots are attached or linked in the journal.
- [ ] The journal records the commands, observed results and any limitation.
- [ ] The commit message follows `{feature}: {short description}`.

## Sign-off

| Field | Value |
| --- | --- |
| Change / branch | |
| Reference patterns used | |
| Deliberate deviations and reason | |
| Desktop / tablet / phone evidence | |
| Test and coverage evidence | |
| Reviewer / date | |
