---
name: agents-workflow
description: The planning agent's process around docs/PLAN.md — role, when to write a plan vs edit directly, the Q&A gate, the execution queue when a plan already exists, and local execution (LAST_PLAN_EXECUTED.md). Use before writing any PLAN.md.
---

# Agents workflow — the process around `docs/PLAN.md`

Plan content: **plan-authoring**. Dispatch, frontmatter, publish: **codejob**. Reviewing a PR: **plan-review**. Multi-repo waves: **master-plan**.

## Role of the planning agent

- Edits only `.md` files unless the user explicitly asks for code.
- A code change touching more than one file → `docs/PLAN.md`, never a direct edit (exceptions: **plan-review** → «Outcome»).
- Never renames, moves or deletes `docs/PLAN.md` — `codejob` owns its lifecycle (sole exception: local execution, below).

## Plan or direct edit

- Documentation (`.md`) → edit directly.
- Go code → `docs/PLAN.md` at the module root (next to `go.mod`), reviewed by the user before dispatch.
- Public API change → the plan carries a `## Design gate` with the five answers of **api-design**; without it, do not dispatch.
- **Default to dispatching, even a single-file fix.** The planning agent runs on a constrained token budget while the user's `codejob`/Jules capacity is comparatively cheap: orchestrate (write the plan, dispatch, pull, review, close) rather than editing code directly. Reserve a direct edit for something that would take longer to write as a plan than to fix (a one-line label or constant), or for work that blocks an active dispatch loop and must land immediately.

## Q&A first

1. Read the code before asking; never ask what it answers.
2. Search the ecosystem for what already exists (**api-design** → «Step 0»), then check prior art: git history, predecessor repos, decision records (verify a rejection's real scope). Reuse or extend before creating.
3. An architectural choice gets ≥2 options with trade-offs and a recommendation; wait for the user.
4. Write the plan only when every decision is resolved; it holds resolutions only. Single library → Q&A in chat; multi-repo → in the master plan (**master-plan**).

## A repo that already has a plan — merge, never queue

New work for a repo whose `docs/PLAN.md` exists:

1. **Not dispatched yet** (no `STATUS`): add the new work as **more stages of that plan**, so one
   dispatch does everything. Then reconcile both halves: dependencies/prerequisites, files touched
   twice, deletions, docs that one half rewrites and the other contradicts, header, stages table,
   acceptance. Fix every inconsistency in the same edit.
2. **Running or in review** (`STATUS` set): it cannot change. Write the new work to
   `docs/PLAN_<TOPIC>.md` and record it in the master plan; when the running plan closes, that file
   becomes `docs/PLAN.md` (merging, per step 1, with any other pending plan of the repo).
3. **Obsolete plan files** (`docs/PLAN_*.md`, `docs/PLAN_STAGE_*.md`, drafts, whose work is already in
   the code or superseded): verify against the code, then **delete them now** — they are docs, the
   planning agent deletes them directly and publishes. Never leave one for a later plan: a leftover
   pointing at `PLAN.md` misleads the executor.

## Local execution — `LAST_PLAN_EXECUTED.md`

When the user decides a plan runs locally instead of via `codejob`:

1. Rename `docs/PLAN.md` → `docs/LAST_PLAN_EXECUTED.md` at that moment.
2. Execute with it as the spec; commit it with the code via `gopush`.
3. Next local execution overwrites it (never `_2` variants; git keeps history).
4. Other docs may point to it only as "the most recent locally executed plan", never to its sections.

Rationale that must outlive a plan goes to permanent docs before dispatch: **documentation** → «Ephemeral vs Permanent».
