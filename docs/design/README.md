# Design reference — index

This directory holds maulid.id's UI design reference. The reference product is
**Nextdoor**, selected by the project owner on 2026-11-15. It was researched from Mobbin
captures of Nextdoor's web and iOS products.

## Source of truth

| Document | Owns |
| --- | --- |
| `nextdoor-reference.md` | product context, information architecture, layout and design principles |
| `nextdoor-components.md` | component anatomy, states and composition rules |
| `nextdoor-principles.md` | the 21 design principles and their evidence |
| `nextdoor-screen-inventory.md` | the captured screen index, grouped by flow |
| `nextdoor-content.md` | voice registers, microcopy, dates, empty states and validation copy |
| `nextdoor-visual-language.md` | colour roles, typography, spacing, radii, imagery and motion |
| `ui-review-checklist.md` | the review gate every UI change must pass |

The agent-facing operational version is the project skill
[`nextdoor-design`](../../.agents/skills/nextdoor-design/SKILL.md). The skill references
are shorter, task-oriented extracts. If the skill and this directory disagree, this
directory wins and the skill is corrected in the same change.

## What this reference is

- A structural and interaction reference: how Nextdoor organises a local community feed,
  a composer, posts, events, trust signals and moderation.
- An adaptation guide for maulid.id's domain: Indonesian Maulid news, public reading,
  staff-validated publishing, one shared API.
- Evidence-based: every screen claim points to a Mobbin URL.

## What it is not

- Not a pixel-for-pixel Nextdoor clone.
- Not permission to reuse Nextdoor's logo, brand assets, illustrations, photography,
  copy, advertising formats or proprietary layouts.
- Not a claim that Nextdoor's patterns improve retention or trust. Mobbin shows what
  shipped, not what worked.
- Not a replacement for accessibility, performance or security verification.

## Mandatory use

Before any UI design, layout, component, visual or design-review work:

1. Load the `nextdoor-design` skill.
2. Read the relevant section of this directory.
3. Apply the checklist in `ui-review-checklist.md`.
4. Record deviations and their reasons in the plan or journal.
5. Attach desktop, tablet and phone screenshots to the completion evidence.

A UI change that cannot show this evidence is not complete.
