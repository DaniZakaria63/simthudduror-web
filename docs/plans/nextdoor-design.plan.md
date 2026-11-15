# Nextdoor design extraction plan

## Interpretation

The user selected Nextdoor as the reference product. This slice must:

1. Extract Nextdoor's screen inventory, components, screen composition, interaction
   patterns, design principles and visual language from Mobbin, in as much detail as the
   available captures allow.
2. Store that material in the repository twice:
   - as a human-readable design reference under `docs/design/`;
   - as an agent-loadable project skill under `.agents/skills/nextdoor-design/`.
3. Make the reference mandatory for future UI work, so later prompts do not silently
   ignore it.
4. Commit and push the result to `origin/feat/designing`.

This is **documentation and agent guidance**, not a UI rewrite. No production React or Go
code changes in this slice.

## What “extract” means

- Reconstruct Nextdoor's information architecture from actual Mobbin screens, flows and
  site sections.
- Describe component anatomy, states, responsive behaviour and interaction rules.
- Record the design principles the product demonstrates, with visible evidence.
- Record the visual language: hierarchy, spacing, radii, colour roles, typography, icons,
  imagery and motion, as observed. Approximate values are labelled as observed/estimated;
  no invented "official" tokens.
- State what maulid.id should adapt, and what it must not copy.
- Cite the Mobbin URL for every screen-level claim.

## What “extract” does not mean

- It is not a pixel-for-pixel clone of Nextdoor.
- It does not copy Nextdoor's logo, brand assets, proprietary illustrations, copy, ads or
  protected layouts.
- It does not transcribe Nextdoor's source code or imply access to it.
- Mobbin captures show what a product shipped; they do not prove the design improved a
  metric.

## Deliverables

```text
docs/design/
├── README.md                      # index and source-of-truth rules
├── nextdoor-reference.md          # human-readable, detailed reference
├── nextdoor-components.md         # component anatomy, states, measurements
├── nextdoor-principles.md         # design principles and composition rules
├── nextdoor-screen-inventory.md   # captured screen index
├── nextdoor-content.md            # voice, microcopy, dates, empty states
├── nextdoor-visual-language.md    # colour, type, spacing, imagery, motion
└── ui-review-checklist.md         # review gate for future UI changes

.agents/skills/nextdoor-design/
├── SKILL.md                       # mandatory operating rules for agents
└── references/
    ├── sources.md                 # Mobbin URLs, corrections and capture notes
    ├── web-app.md                 # web information architecture and screens
    ├── mobile-app.md              # iOS information architecture and screens
    ├── components.md              # component contracts
    ├── patterns.md                # interaction patterns
    ├── content.md                 # copy and tone rules
    ├── visual-language.md         # typography, colour, spacing, imagery, motion
    └── adaptation.md              # maulid.id adoption and anti-patterns
```

A single source of truth is maintained by:

- `docs/design/nextdoor-reference.md` as the readable index;
- `.agents/skills/nextdoor-design/` as the operational agent guidance;
- `AGENTS.md` and `docs/skills.md` naming the skill as mandatory before UI work.

## Status (2026-11-16)

| Requirement | Status | Evidence |
| --- | --- | --- |
| Screen inventory | Done | 84 iOS + ~40 web screens in `docs/design/nextdoor-screen-inventory.md` |
| Components | Done | `nextdoor-components.md`, skill `components.md` |
| Screen composition | Done | `nextdoor-reference.md` §4, `web-app.md`, `mobile-app.md` |
| Design principles | Done | `nextdoor-principles.md` (21 principles) |
| Visual language | Done | `nextdoor-visual-language.md`, skill `visual-language.md` |
| Content/tone | Done | `nextdoor-content.md`, skill `content.md` |
| Agent skill | Done | `.agents/skills/nextdoor-design/SKILL.md` + seven references |
| Enforcement | Done | `AGENTS.md`, `docs/skills.md`, `react-feature-architecture`, `monolith-delivery`, checklist |
| Source integrity | Done | 161 unique Mobbin URLs returned HTTP 200 (2026-11-16) |
| Production UI change | None | Docs/agents only; no Go/React source modified |

## Enforcement

- Add a **Design reference** section to `AGENTS.md`: before any UI, layout, visual or
  component work, load `nextdoor-design`, read the relevant references, and apply the
  documented principles.
- Add the skill to the project skill index in `docs/skills.md`.
- Add a review checklist that every UI change must satisfy, including desktop/tablet/mobile
  evidence and an explicit note on which Nextdoor-derived pattern was used or deliberately
  rejected.
- Update the delivery expectations: UI changes without the checklist evidence are not
  complete.
- Do not add a brittle source-code grep gate for visual rules; the checklist plus real
  browser evidence is the enforceable gate.

## Research method

1. Query Mobbin for Nextdoor web screens, iOS screens, web flows, iOS flows and any
   available site sections.
2. Inspect the returned captures; record product, platform, URL and visible behaviour.
3. Group findings by information architecture, component, pattern and visual language.
4. Separate observed fact from adaptation advice.
5. Reconcile overlapping captures; keep one canonical URL per component.
6. Record capture limitations (desktop-only web, ads, account-specific screens, no
   seasonal module, no Nextdoor marketing-site sections).

## Verification

- Every referenced screen has a Mobbin URL and platform.
- `docs/design/` and the skill agree on component names and rules.
- The skill's YAML frontmatter is valid and the skill is discoverable under
  `.agents/skills/`.
- `AGENTS.md`, `docs/skills.md` and the review checklist link the skill.
- `git diff --check` is clean; diff is docs/agents only.
- `make` gates are not required for a docs-only slice; if any source code is added, the full
  gate applies.
- Push to `origin/feat/designing`; verify the remote head equals local head.

## Limitations

- Mobbin is a capture library, not a design-system export; exact spacing, font sizes and
  colour tokens are estimated from images unless a token is visibly stated.
- Some screens are promoted (ads), account-state specific, or captured from a business
  account. No Nextdoor seasonal module was found.
- Web captures are desktop-only (~1920px); tablet/phone breakpoints are maulid.id decisions.
- Mobbin has no Nextdoor marketing-site sections; the reference is the web and iOS app only.
- Nextdoor's brand and assets are proprietary; only patterns, measurements and
  interaction logic may be adapted.
- The current maulid.id `Post` model supports text body, author, time and likes only.
  Title/image/category/event components are documented as future work, not implemented.
