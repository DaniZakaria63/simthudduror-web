---
name: Monolith Delivery Workflow
description: Follow this repository's delivery rules for prompts, plans, branches and worktrees, small feature-prefixed commits, agent journals, verification gates and the final push. Use when starting a new prompt-sized piece of work, preparing a commit, or deciding whether work is ready to deliver.
---

# Delivery workflow (maulid.id)

This skill turns the project prompt's workflow rules into concrete steps. It applies
to every prompt-sized unit of work, not to a one-line fix inside an existing plan.

## When to use it

- starting work from a `*.prompt.txt`;
- creating or finishing a branch/worktree;
- preparing commits or pushing;
- deciding whether a change is “done”;
- recording agent work in the journal.

## 1. Prompt → plan → work

The naming contract is fixed by the prompt:

```text
docs/plans/<topic>.prompt.txt   # original request, preserved verbatim
docs/plans/<topic>.plan.md      # interpretation, requirements, acceptance evidence
```

1. Read the prompt and the existing codebase first. Do not guess a stack or a product
   backlog.
2. Write the neighbouring `*.plan.md` **before** structural implementation. It must
   state the interpretation, required behaviour, boundaries, verification evidence and
   known limitations.
3. Keep the prompt byte-for-byte unchanged. Update the plan as evidence arrives; never
   mark acceptance for a check that did not run.

## 2. Branch and worktree

- One plan → one dedicated branch and, when the harness uses worktrees, one worktree.
- Use a descriptive branch name (`feat/initialize`, `feat/<topic>`).
- Keep the worktree until the human confirms the branch is merged. Only then remove it.
- Never delete a branch or worktree that contains unaudited or unpushed work.

## 3. Commits

Format is exact:

```text
{feature name}: {short changes description}
```

Examples from this repository:

```text
backend: add feature-modular Go monolith with SQLite accounts, posts and staff validation
web: add editorial landing page, account screens and staff backoffice
tooling: add build, coverage gates, boundary checks and developer docs
docs: record plan evidence, architecture, research sources and security audit
```

Rules:

- Commit small, coherent slices; avoid one giant final commit.
- Do not mix unrelated formatting or generated files into a feature commit.
- Never commit `node_modules`, `dist`, coverage output, databases, `.env` files,
  build binaries or audit scratch artifacts (see `.gitignore`).
- Commit after the slice passes the checks that apply to it; do not commit known
  red gates.

## 4. Journal

Every agentic work session writes under `.agents/journal/`, named by date and topic
(`2026-11-14-implementation.md`). Record:

- what was attempted and by which agent/model (including the verified model variant
  when the task requires it);
- files/areas changed and the reason;
- commands actually run and the observed results (numbers, not adjectives);
- failures, how they were diagnosed, and the remediation;
- decisions that deviate from a reference or plan, and why;
- remaining limitations and open items.

The journal is evidence, not a diary of intentions. “Tests pass” without the command
and result is not evidence.

## 5. Verification gate

A change is deliverable only when the applicable checks pass **and** their output is
recorded:

```sh
make lint        # gofmt check, go vet, frontend boundary + source-order
make test        # go test ./... + vitest
make coverage    # backend >80% statements; frontend statements/branches/functions/lines >80%
make race        # go test -race ./...
make security    # govulncheck + npm audit
make build       # one binary with the embedded React build
```

Plus runtime/browser evidence when a surface changed: the three surfaces served by the
built binary, validation invariant through the API, and a real-browser check of the
page at desktop/tablet/phone widths.

A non-zero exit means the work is not done. A coverage value of exactly 80 fails: the
requirement is strictly greater than 80%.

## 6. Push and hand-off

- Push only when all work in the plan is done and the gates pass:
  `git push -u origin <branch>`.
- Verify with `git ls-remote origin refs/heads/<branch>` that the remote head equals the
  reviewed local head.
- Report the branch, commit range, verification results and remaining limitations.
- Keep the worktree until the human confirms the merge; then remove it and delete the
  merged branch.

## 7. Agent roles (when the Orca harness is in use)

The prompt names the roles; the harness is opt-in:

| Role | Agent | Owns | Never |
| --- | --- | --- | --- |
| Principal | Codex `gpt-6.1-sol`, high | analysis, architecture, delegation, review, acceptance | edits app code |
| Implementation | OpenCode DeepSeek V4.1 Flash, high | all app code and tests | launches agents or treats advice as permission |
| Research | Antigravity | source/reference verification, provenance | edits app code |
| Advisor | Jev (TypeSafe) | typed advisory judgments | decides, authorizes or replaces evidence |

Jev confidence is never permission. Record its advisory output in the journal, then
decide from evidence. Do not recursively launch the harness or create a second
principal.

## 8. Documentation duties

When architecture, commands or security posture change, update the matching source of
truth in the same delivery:

- `README.md` — commands, configuration, API, repository map;
- `docs/architecture.md` — module boundaries, request lifecycle, data model;
- `docs/security-audit.md` — audit run, findings, limits;
- `AGENTS.md` — rules and verification expectations for future agents;
- `.agents/skills/*` — when a reference or convention changes.

## Definition of done

- [ ] Plan updated with observed evidence and explicit limitations.
- [ ] Journal records the session, commands and results.
- [ ] Lint, tests, coverage, race, security and build pass.
- [ ] Runtime/browser evidence captured for changed surfaces.
- [ ] Small feature-prefixed commits, no generated artifacts.
- [ ] Pushed; remote head matches local head.
- [ ] Worktree kept until the human confirms merge.
