# Project agent skills

`.agents/skills/` holds the repository's own agent skills, discovered automatically by
OpenCode (project-compatibility source) and readable by any agent working in this
checkout. They extract the guidance from the prompt's reference lists and the owner-selected
Nextdoor design reference into loadable, project-specific rules.

| Skill ID | Load it when | Extracted from |
| --- | --- | --- |
| [`react-feature-architecture`](../.agents/skills/react-feature-architecture/SKILL.md) | building or reviewing any React screen, component, hook, state or test | mdskills.ai React tag, callstack agent-skills, Ashwathvia's React architecture article, patterns.dev skills catalog |
| [`go-feature-architecture`](../.agents/skills/go-feature-architecture/SKILL.md) | building or reviewing Go API, persistence, auth or staff code | Leapcell clean architecture, Frederich Blessy's Go architecture, Reintech Go structure 2026, CyberAgent ABEMA architecture, idiomatic Go patterns |
| [`security-audit`](../.agents/skills/security-audit/SKILL.md) | security review, pre-release audit, auth/session/validation changes | Cloudflare `security-audit-skill` (MIT, vendored at a pinned revision) |
| [`monolith-delivery`](../.agents/skills/monolith-delivery/SKILL.md) | starting a prompt-sized task, committing, pushing, closing a worktree | `docs/plans/initial.prompt.txt` branching/orchestration/journal rules |
| [`nextdoor-design`](../.agents/skills/nextdoor-design/SKILL.md) | **mandatory** before any UI, layout, component, visual or design-review work | Mobbin captures of Nextdoor web and iOS, recorded in the skill's `references/sources.md` |

## How they are structured

```text
.agents/skills/<skill-id>/
  SKILL.md              # frontmatter (name, description) + the rules
  references/           # sources, worked examples, testing patterns
```

OpenCode advertises a skill by its `description`, then loads the body on demand. Paths
inside a skill are relative to the directory containing `SKILL.md`.

## Rules that apply to the skills themselves

- Keep them project-specific: cite the real files and commands, not generic advice.
- Record what was adapted and what was rejected from each reference
  (`references/sources.md` in each skill).
- When a convention changes, update the skill in the same delivery; a stale skill is
  worse than none.
- Do not vendor third-party skills without their license and pinned revision. The
  Cloudflare skill is MIT and vendored under
  `security-audit/references/upstream/` with its `LICENSE`.
- These skills never override the project prompt or `AGENTS.md`. If they disagree, the
  prompt/plan wins and the skill is corrected.
