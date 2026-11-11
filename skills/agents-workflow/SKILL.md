---
name: agents-workflow
description: The planning agent's process around docs/PLAN.md — role, when to write a plan vs edit directly, the Q&A gate, the execution queue when a plan already exists, and local execution (LAST_PLAN_EXECUTED.md). Use before writing any PLAN.md.
---

# Agents workflow — the process around `docs/PLAN.md`

Plan content: **plan-authoring**. Dispatch, frontmatter, publish: **codejob**. Reviewing a PR: **plan-review**. Multi-repo waves: **master-plan**.

## Role of the planning agent

- Edits `.md` files, and code only when **Plan or direct edit** (cost rule) says it is cheaper than a plan.
- Never renames, moves or deletes `docs/PLAN.md` — `codejob` owns its lifecycle (sole exception: local execution, below).

## Wrong content in a plan or doc — fix it at once

Stale vocabulary, a wrong signature or a contradicted rule found in any `PLAN.md`, master or doc is corrected **in the same turn**, never parked waiting for the owner's approval or for another repo to publish: it gets forgotten. Approval is for design choices, not for correcting errors.

## Cheap executors — slice for them

The `codejob` executors are cheap and not very capable: they do mechanical work well and design badly. Use them to the maximum **while the capacity lasts**: cut a big plan into small stage files (`docs/PLAN_STAGE_<N>_<TOPIC>.md`, one concern each, exact code, verified signatures), dispatch many repos in parallel, and keep the judgment (design, review, fixes) for the planning agent. A plan the executor questions twice was cut too big or not exact enough: slice it more, never answer one by one.

**No delivery after the deadline the owner sets** (e.g. a plan of more than 4 stages with no pull request): cut it into 2–4 sequential plans of at most 3–4 stages, each one leaving `gotest` green on its own (the first deletes the tests of removed code, the next write the new ones), dispatch the first, `codejob close`, dispatch the next; the later parts wait as `docs/PLAN_<PART>.md`. Never let the executor carry a 7-stage plan whose only green point is the last stage.

## Plan or direct edit

- Documentation (`.md`) → edit directly.
- Go code that passes the cost rule below → `docs/PLAN.md` at the module root (next to `go.mod`), reviewed by the user before dispatch; the rest is edited locally.
- **Architecture plans are always reviewed by the owner before dispatch.** Dispatch only a plan whose master row carries 🔵 (**master-plan** → «Owner review mark»); show the owner the first plan without 🔵, in the master's review order, and wait.
- Public API change → the plan carries a `## Design gate` with the five answers of **api-design**; without it, do not dispatch.
- **Cost rule (the goal is maximum efficiency in tokens and productivity): never spend more writing a plan than doing the work.** Estimate before writing: plan (written code and examples) + dispatch + review versus the change itself. If the plan would take about as long to write as the change takes to do, **do it locally**: edit, `gotest`, `gopush` (**codejob** «gopush and local `replace`»). Minor changes get NO plan: a type or constant swap, a version bump, a rename across consumers, a few lines in a few repos. Measured 2026-11-10: the D61 consumers (`outbox`, `goflare`, `sse`, 1–3 changed lines each) came back in minutes, and each plan took twice as long to write as the fix would have taken.
- **A plan pays only when the work is big or independent:** a module migration or any change that needs more than about an hour of local work, several large pieces that can run in parallel, or a design the owner must review first (🔵). Everything else: local. A change that blocks an active dispatch loop also lands locally at once.

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
5. **Code is written in the target repo, never outside it.** No prototype in a scratch dir, no copy of a module "to try the design first": create the repo first (`gonew`) and write there. Code that already works is the implementation — copy it in, do not turn it into a plan.

Rationale that must outlive a plan goes to permanent docs before dispatch: **documentation** → «Ephemeral vs Permanent».

## Agent asking questions — the plan is ambiguous: fix the plan, never ping-pong

**The second question on one plan is the signal.** (2026-11-10: nine of twelve migration sessions asked, one `chat_room` session ten times; answering each by `codejob reply` took hours and lost work when a session died unpushed.) Steps, in order:

1. **Answer the first question** (one `codejob reply`, a permission question is always «yes, proceed, open the PR»). **From the second question on, stop replying one by one.**
2. **Collect** every question the session asked and classify each: missing signature/argument shape, missing decision, missing tool or version, missing test data.
3. **Fold each answer into the stage that needs it**, as exact Go or an exact table — never into a `## Additions` appendix (an executor reads stages, not notes). Module migration: take the blocks of **crud-module-migration**.
4. **Split** the plan into stage files (**plan-authoring** → «Modular stage files»; granularity in **crud-module-migration** §5): one concern per stage, every non-trivial body written out. A session over 3 hours or with ≥3 questions is aborted, not nudged.
5. **Re-dispatch:** remove `STATUS` and `SESSION` from the frontmatter, `codejob dispatch`. Jules publishes a branch only when it submits the PR (a `git push` from its sandbox does not reach GitHub: 2026-11-10 four sessions ended with hours of work and no branch, lost): the plan tells the executor to OPEN THE PR after the third stage and keep committing to it. A session «finished without a PR» is answered at once with «open the pull request now with everything you have, list what fails» (it resumes and submits); work not submitted is lost when the session ends.
6. **Put each new answer in the skill that owns it** (same turn), so the next plan is born with it.
7. **A module too big for one session, or a session that failed twice:** split it in two plans, dispatch one, then the other (**crud-module-migration** §6).

Never tell the executor to run local dev tools (`webtyp`) or capture screenshots in headless CI (verify `web/browser` with `go build ./web/browser`).
