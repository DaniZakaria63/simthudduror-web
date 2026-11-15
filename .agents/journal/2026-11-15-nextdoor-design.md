# Nextdoor design extraction — 2026-11-15

## Session

- Agent: OpenCode (DeepSeek V4.1 Flash 1017).
- Request: after selecting Nextdoor as the reference, extract all screen components, design
  patterns, UI principles, screen composition and components into the project as design
  documentation and as an agent skill; make future UI work follow it; push to origin.
- Worktree: `/home/dani/WebstormProjects/maulid.id-feat-designing`
- Branch: `feat/designing`
- Base for the previous slice: `b4617dd` (Mobbin reference research).

## Research workers

| Session | Brief | Result |
| --- | --- | --- |
| `ses_e1fef93c9ffewh4jtP9X6CDjeI` | Nextdoor web app and website extraction | ~100 screens + 20 flows; measured 3-zone geometry, modal detail, events/groups/search/settings, web anti-patterns |
| `ses_e1fef93c8ffe60zeKEZ1SdS6mJ` | Nextdoor iOS app extraction | 84 screens + 20 flows; IA, composition, component metrics, motion, anti-patterns |
| `ses_e1fef93c7ffeKX2OYK7xHwU6ub` | Cross-cutting trust, moderation, content, states, seasonal, accessibility, monetisation | ~70 screens; validation/trust registers, moderation menus, copy, empty states, settings, ads, accessibility |

### Correction: no Nextdoor seasonal module

The earlier `2026-11-15-designing` journal and the broader Mobbin report cited a Nextdoor
“Ramadan Mubarak” card at `feb2cdda-699f-4be5-a11a-244612039997`. The Mobbin page title for
that capture is **Nextdoor iOS Local Deals — local deals explanation page**. It is not
seasonal evidence. Searches for `Nextdoor Ramadan`, seasonal greetings and holiday modules
returned no Nextdoor capture. All current design documents now state that the Maulid
seasonal card is a maulid.id editorial decision, not a Nextdoor precedent.

### Key findings accepted

- Nextdoor has no visible “verified neighbour” badge. Trust is an address-gated signup
  (implicit, not displayed) plus role labels (`Local publisher`, `Admin`, agency, business,
  `Sponsored`). maulid.id maps this to implicit account state + explicit
  `Sudah divalidasi` text, with no badge.
- Validation is a named, resumable checklist: invitation → explicit in-flight
  acknowledgement → ticked item, with no ambiguous pending state.
- The iOS root shell is four labelled destinations (`Home`, `Search`, `For Sale`, `Faves`)
  plus a raised green Post FAB; notifications, messages and profile are header utilities,
  not tabs.
- The feed is one column of stacked cards; text-only posts are complete; media optional.
- Two primaries, two jobs: green for social commit, navy for operational commit, red for
  destructive.
- Moderation is an overflow menu with consequence-first items; report only when a review
  path exists.
- Empty states are quiet; an empty feed may have no CTA; illustrations are reserved for
  actionable surfaces.
- Ads occupy the same author line and card geometry as editorial content; this is the
  clearest anti-pattern for maulid.id.
- Web geometry (measured at 1920): left nav ~17%, feed ~38–39%, right rail ~19%, rest
  margin; feed:rail ≈2:1. Top bar ~82px; active nav is light grey + filled icon, not green.
- Web post detail is a centred modal overlay; publisher content is a split-pane overlay.
- Web events use a 3-up grid; For Sale uses a 6-up grid; the feed stays a single column.
- Mobbin has no Nextdoor marketing-site sections; do not infer a website palette or layout.
- Web ads reuse post anatomy (author row, body, media, CTA) and are distinguished only by a
  small “Sponsored” line; another reason to reject ads.
- Accessibility: icon+label navigation, error icon+text+border, visible focus, dark mode;
  icon-only card actions and carousels are risky and must be fixed in maulid.id.

### Deliverables

```text
docs/design/
├── README.md
├── nextdoor-reference.md
├── nextdoor-components.md
├── nextdoor-principles.md
├── nextdoor-screen-inventory.md
├── nextdoor-content.md
├── nextdoor-visual-language.md
└── ui-review-checklist.md

.agents/skills/nextdoor-design/
├── SKILL.md
└── references/
    ├── sources.md
    ├── web-app.md
    ├── mobile-app.md
    ├── components.md
    ├── patterns.md
    ├── content.md
    ├── visual-language.md
    └── adaptation.md
```

## Enforcement added

- `AGENTS.md`: `nextdoor-design` is mandatory before any UI/layout/component/visual/design
  review; added a Design reference section and UI verification bullet.
- `docs/skills.md`: skill indexed.
- `.agents/skills/react-feature-architecture/SKILL.md`: cross-reference to the design skill
  and review checklist.
- `.agents/skills/monolith-delivery/SKILL.md`: UI changes require the checklist and
  screenshots; definition of done updated.
- `README.md`: docs map includes `docs/design/`.

## Checks

- Skill frontmatter present (`name`, `description`); `nextdoor-design` advertised by
  OpenCode after the skill file was added.
- Local link check across `docs/design/` and the skill: 17 files, 0 missing relative links.
- Mobbin source check: 161 unique URLs extracted from the docs/skill; all returned HTTP 200
  on 2026-11-16.
- One invalid URL found during validation (`72f65c0b-…` 404) and replaced with a verified
  capture (`5380096e-…`); the full URL set was re-checked.
- `git diff --cached --check` -> exit 0.
- Staged diff is docs/agents only: 27 files, 3,685+ insertions before the final edits; no
  Go/React source changed, so the `make` gates are not run for this docs/agents slice.
- Push to `origin/feat/designing` and remote-head verification: see the final session
  report; the branch is pushed after this journal entry.

## Limitations

- Mobbin is a capture library, not a design-token export.
- Not every Nextdoor screen is represented; the skill records the captured set only.
- Web captures are desktop-only (~1920); tablet/phone breakpoints are maulid.id decisions.
- Mobbin has no Nextdoor marketing-site sections; the reference is the web and iOS app.
- Nextdoor's brand/assets are proprietary; only patterns and measurements are adapted.
- The current maulid.id `Post` model cannot support title/image/category/event components;
  those are documented as future schema work.
