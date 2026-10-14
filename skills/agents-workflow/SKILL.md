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

## Q&A first

1. Read the code before asking; never ask what it answers.
2. Check prior art: git history, predecessor repos, decision records (verify a rejection's real scope).
3. An architectural choice gets ≥2 options with trade-offs and a recommendation; wait for the user.
4. Write the plan only when every decision is resolved; it holds resolutions only. Single library → Q&A in chat; multi-repo → in the master plan (**master-plan**).

## Never clobber a pending plan — the execution queue

If `docs/PLAN.md` already holds a pending plan:

1. Move it to `docs/PLAN_<TOPIC>.md`; write the new one to `docs/PLAN_<NEW_TOPIC>.md`.
2. Rewrite `docs/PLAN.md` as a queue that opens with:

   ```markdown
   # PLAN — execution queue for `<module>`

   > If you were told to "execute the plan described in docs/PLAN.md", execute
   > **ALL the plans below, in order**. Finish one (acceptance criteria green)
   > before the next. Never mix changes between plans.

   | Order | Plan | Subject |
   |-------|------|---------|
   | 1 | [PLAN_<TOPIC>.md](PLAN_<TOPIC>.md) | ... |

   After all plans, run `gotest ./...` once more: everything green.
   ```

## Local execution — `LAST_PLAN_EXECUTED.md`

When the user decides a plan runs locally instead of via `codejob`:

1. Rename `docs/PLAN.md` → `docs/LAST_PLAN_EXECUTED.md` at that moment.
2. Execute with it as the spec; commit it with the code via `gopush`.
3. Next local execution overwrites it (never `_2` variants; git keeps history).
4. Other docs may point to it only as "the most recent locally executed plan", never to its sections.

Rationale that must outlive a plan goes to permanent docs before dispatch: **documentation** → «Ephemeral vs Permanent».
