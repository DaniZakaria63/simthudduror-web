# AGENTS.md

## Project

maulid.id is a community news-sharing platform for Maulid Rasulullah events in Indonesia.

It is intentionally a **monolith** with three surfaces sharing one codebase and deployment:
- Public landing page
- Mobile API
- Backend + backoffice for staff

Do not split these into separate apps/services without an explicit decision from the user.

## Domain rules

- Any user can interact with posts.
- Only **validated users** may create posts. Enforce this in the backend/API, not only in the UI.
- Validation status is a core account concept; keep its naming and behavior consistent across the API, backoffice, and landing page.

## Architecture docs

- `docs/` is the planned home for system design and architecture notes; treat it as the source of truth once populated.
- Read the relevant doc before structural changes. If docs and code disagree, surface the conflict instead of guessing.

## Status

- The repo currently has no scaffold, git history, or build/test/lint commands. Do not assume a framework or toolchain; derive exact commands from the manifests once they exist and replace this section with them.
