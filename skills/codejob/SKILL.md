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
| `codejob reply "text"` / `codejob approve` | answers / approves the agent session | `running` |
| `codejob close "msg" [tag]` | merge + `gopush` + delete `docs/PLAN.md` (local) | `review` |

Cloud setup: the human merges the PR and the Action publishes. Never call `--ci` by hand.

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
- Wait with a background timer (10–15 min), then `codejob pull`; never poll in a tight loop.

## Source the executor cannot see — `_temp/`

Copy only the files the plan needs to `<repo>/_temp/<source-repo>/<same path>` (never `.env`, `docs/`, `data/`, `config/`; grep for secrets). `_` makes Go ignore it. The plan's last stage deletes `_temp/` (`test ! -e _temp`).
