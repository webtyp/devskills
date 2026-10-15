---
name: plan-authoring
description: How to write the body of a docs/PLAN.md so a less-capable execution agent implements it without design judgment — links AGENTS.md instead of copying rules, exact files, verbatim errors, explicit deletions, delivery check, modular stage files. Use when drafting or reviewing the content of a PLAN.md.
---

# Plan authoring — the content of `docs/PLAN.md`

Process: **agents-workflow**. Frontmatter and dispatch: **codejob**. Multi-repo orchestration: **master-plan**.

## The executor has zero context and needs zero design judgment

An open choice in the plan is a Q&A gap, not a detail to delegate.

- Ecosystem rules are **linked, never copied**: the plan opens with `> Read AGENTS.md at the repo root before any change.` Missing or stale `AGENTS.md` → write it before dispatch (**agents-md**). The plan holds only this change's decisions.
- Cross-repo links are GitHub URLs (`https://github.com/webtyp/<repo>/blob/main/...`), never `../../`; content the executor needs from another repo is restated in the plan.
- Code the executor cannot reach goes to `_temp/` (**codejob**).
- **Every tool the plan names is installed in the plan.** The executor has none of the local CLIs: a `## Tools` section, first, with one `go install` line per tool the plan uses (catalog below). A command the plan runs without its install line is a defect.
- Name every file per stage, including new ones.
- Quote error messages verbatim; they become typed constants.
- List every deletion, each with a grep-verifiable criterion (`grep -rn "oldFn" .` → empty). A weak agent never removes what the plan did not order.
- State anti-footguns where an ecosystem rule could be misapplied (e.g. "this repo is backend tooling: stdlib is legitimate, do NOT 'fix' it").
- Risky mechanics (toolchain, environment) state the fallback and require an injectable seam for tests.
- Type-precise wording: a cache key is a typed struct; identify packages by import path, never by selector text.
- Public API change → `## Design gate` before the stages, five answers of **api-design** (read that skill; do not restate its criteria).
- Sequential stages, a stages table at the end.
- Never mention `gopush` or `codejob` inside the plan.

## Tool catalog — copy the lines the plan needs

| Tool | Install |
|---|---|
| `gotest` (always: every plan runs the suite) | `go install webtyp.com/devflow/cmd/gotest@latest` |
| `langc` (translations, `lang.json`) | `go install webtyp.com/lang/cmd/langc@latest` |
| `ormc` (regenerates `model_orm.go`) | `go install webtyp.com/ormc/cmd/ormc@latest` |
| `sitec` (SSR / asset extraction) | `go install webtyp.com/sitec/cmd/sitec@latest` |
| `weightsc` | `go install webtyp.com/weightsc/cmd/weightsc@latest` |

Never in a plan: `gopush`, `codejob`, `devskills`, `webtyp` (local developer tools).

## Modular stage files

Complex plans: `docs/PLAN.md` is the index; stages go to `docs/PLAN_STAGE_<N>_<TOPIC>.md`, each opening with `← [Stage N-1](...) | Next → [Stage N+1](...)`.

When porting code, append the original snippets at the end of the stage file and say what to recycle and what to replace.
