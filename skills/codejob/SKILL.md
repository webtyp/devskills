---
name: codejob
description: Running a docs/PLAN.md through codejob — the required frontmatter and its keys, the explicit commands (dispatch, pull, reply, close), gopush vs codejob, the 15-plan concurrency limit, and _temp/ for source the executor cannot see. Use when dispatching, pulling, closing or publishing.
---

# codejob — the plan loop

`gopush` publishes. `codejob` runs the loop (`dispatch → running → reviewing → review → close`), driven by `STATUS` in the plan's own frontmatter, and calls `gopush` when closing.

## Frontmatter — first line of `docs/PLAN.md` is `---`, or dispatch fails

```markdown
---
PLAN: "feat: what this plan implements"
TAG: v0.2.0
EXECUTOR: jules
REVIEWER: none
---

> Read AGENTS.md at the repo root before any change.
```

| Key | Writer | Meaning |
|---|---|---|
| `PLAN` | agent (required) | Commit message on close |
| `TAG` | agent | Explicit version; omitted → auto-bump |
| `EXECUTOR` / `CORRECTOR` | agent | Implementer (default `jules`) / who applies review feedback (default executor) |
| `REVIEWER` | agent | Agent posting a GitHub review; `none` → human only |
| `REVIEW_GUIDE` | agent | Extra review criteria file |
| `STATUS`, `SESSION`, `REVIEW_SESSION`, `ROUND`, `PR` | **machine only** | Never hand-edited |

Full spec: `devflow` repo, `docs/CODEJOB.md`.

## Commands — nothing acts without its verb

| Command | Does | Valid when STATUS is |
|---|---|---|
| `codejob` | read-only status + suggested next command | any |
| `codejob dispatch` | sends the plan to the EXECUTOR (run it when the user says "despacha") | absent / `dispatch` |
| `codejob pull` | checks the PR branch out **in the same local repo** (never `gh repo clone` / `gh pr checkout`); re-run to fast-forward after a correction | `running`, `reviewing`, `review` |
| `codejob reply "text"` / `codejob approve` | answers / approves the agent session; the text sent is also appended to that repo's `docs/PLAN.md` under `## Additions` (never to a master) | `running` |
| `codejob close "msg" [tag]` | merge + `gopush` + delete `docs/PLAN.md` (local) | `review` |

Cloud setup: the human merges the PR and the Action publishes. Never call `--ci` by hand.

## When the agent stops — "Jules was unable to complete the task", a failed session, a question

Never leave it waiting. Read its last messages (`codejob pull` prints them), then:

| What happened | Do |
|---|---|
| The work is done and it asks permission (to commit, push, open the PR) | `codejob reply` telling it to commit, push and open the PR now |
| It asks something the plan should have answered | answer with exact structures (`codejob reply`), then add that answer to the plan so the next dispatch does not ask |
| More than 3 hours elapsed or spinning in regression loops | **cancel/close the session and rewrite**: if an agent takes >3 hours or spins in logic regressions, it is a direct symptom that the plan lacked detail in the *how* (steps too broad or ambiguous). Stop patching with endless replies: abort the session, rewrite the plan with exact structures/conversions and divide it into smaller, granular stages before dispatching again |
| "Unable to complete" / failed with no usable work | **review the plan**: find what made it stop (an ambiguous step, a guessed API, a missing tool or version, a command that fails); fix the plan, reset its frontmatter to undispatched (remove `STATUS`/`SESSION`), and `codejob dispatch` again |
| Failed twice on the same plan, and it is **small/simple** (≤2 files, one clear concept) | finish it locally (**agents-workflow** → «Local execution») |
| Failed twice on the same plan, and it is **large or touches an unfamiliar pattern** (many files, a shared chassis, a reactive/architectural model the executor has to infer) | do not finish it locally and do not ask the user — **split it**: rewrite as the smallest safe slice (e.g. new types with zero wiring), dispatch that, and queue the rest as follow-up plans once it lands. Smaller, narrower plans succeed more often with a cheaper/weaker executor model than one large one |

## Publish with

| You did… | Publish |
|---|---|
| Docs or a 1-file fix, no plan | `gopush 'msg'` |
| Dispatched plan | `codejob close 'msg'` or merge the PR |
| Plan executed locally | `gopush 'msg'` |

Never put `gopush` or `codejob` inside a plan.

## Concurrency

- ≤ 15 plans `running` at once; one plan per repo (park the next as `docs/PLAN_<TOPIC>.md`).
- Plans in `review`/`reviewing` do not count.
- **After dispatching, set a recurring 15-minute timer** that runs `codejob` (read-only) and `codejob pull` in every repo with a plan in flight, and applies **plan-review** to each delivered PR. Remove the timer when no plan is in flight. Never poll in a tight loop.
- **Skip the timer when another agent already owns the review** (e.g. a master plan says "Antigravity coordina la revisión"). Dispatch and move on; that agent's own session does the polling.

## Source the executor cannot see — `_temp/`

Copy only the files the plan needs to `<repo>/_temp/<source-repo>/<same path>` (never `.env`, `docs/`, `data/`, `config/`; grep for secrets). `_` makes Go ignore it. The plan's last stage deletes `_temp/` (`test ! -e _temp`).

## gopush and local `replace`

Working across libraries uses a temporary `replace webtyp.com/<dep> => ../<dep>` in the consumer's `go.mod`. **Do not remove them by hand**: publish in dependency order and `gopush` does it.

- `gopush` on the dependency runs a cascade over every dependent in the workspace: it **removes the `replace` of the module just published**, `go get <dep>@<tag>`, `go mod tidy`, runs the `gotest` gate (reverting `go.mod`/`go.sum` on failure).
- Then: **clean tree and no other `replace` left** → publishes the dependent too (commit + tag + push). **Uncommitted work or `replace`s for other modules** → *deps-only*: commits just `go.mod`/`go.sum`, pushes **without a tag**, stops there; replaces of unrelated modules are left in place (so that commit can carry them until their module is published).
- So the order is leaf-most dependency first; the last `gopush` of the chain leaves no `replace`. `--no-cascade` publishes only the module itself.

