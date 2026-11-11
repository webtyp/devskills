---
name: plan-authoring
description: How to write the body of a docs/PLAN.md so a less-capable execution agent implements it without design judgment — links AGENTS.md instead of copying rules, exact files, verbatim errors, explicit deletions, delivery check, modular stage files. Use when drafting or reviewing the content of a PLAN.md.
---

# Plan authoring — the content of `docs/PLAN.md`

Process: **agents-workflow**. Frontmatter and dispatch: **codejob**. Multi-repo orchestration: **master-plan**.

## The executor has zero context and needs zero design judgment

An open choice in the plan is a Q&A gap, not a detail to delegate.

- Ecosystem rules are **linked, never copied**: the plan opens with `> Read AGENTS.md at the repo root before any change.` and, right below it, the **«Finish by yourself»** line (frontmatter example in **codejob**): Jules finished `site_manager` (2026-10-30), asked «should I proceed?» and the session died unpublished — a delivery rule at the end of the plan is not enough Missing or stale `AGENTS.md` → write it before dispatch (**agents-md**). The plan holds only this change's decisions.
- Cross-repo links are GitHub URLs (`https://github.com/webtyp/<repo>/blob/main/...`), never `../../`; content the executor needs from another repo is restated in the plan.
- Code the executor cannot reach goes to `_temp/` (**codejob**).
- **Every tool the plan names is installed in the plan.** The executor has none of the local CLIs: a `## Tools` section, first, with one `go install` line per tool the plan uses (catalog below), followed by `export PATH="$PATH:$(go env GOPATH)/bin"` (the executor otherwise reports «gotest: command not found»). A command the plan runs without its install line is a defect.
- Name every file per stage, including new ones.
- A plan that changes a convention (a build tag, a file or binary name, a constant shared by tools) lists **every site that knows it**, found by grep (`grep -rn '"browser"\|tags=' --include=*.go`), and a stage per site. Two copies of one function in the same repo count as two sites: changing one is how `style.css` shipped empty.
- Quote error messages verbatim; they become typed constants.
- List every deletion, each with a grep-verifiable criterion (`grep -rn "oldFn" .` → empty). A weak agent never removes what the plan did not order.
- State anti-footguns where an ecosystem rule could be misapplied (e.g. "this repo is backend tooling: stdlib is legitimate, do NOT 'fix' it").
- Risky mechanics (toolchain, environment) state the fallback and require an injectable seam for tests.
- Type-precise wording: a cache key is a typed struct; identify packages by import path, never by selector text.
- A plan that creates a library, module or exported name → `## Ecosystem search` first (what was searched, what exists, reuse / extend / new — **api-design** → «Step 0»), then, for a public API change, `## Design gate` before the stages, five answers of **api-design** (read that skill; do not restate its criteria). Name the existing packages the executor must import; the executor has no checkout of the ecosystem and cannot find them.
- Sequential stages, a stages table at the end.
- **Every plan ends with a documentation stage that brings `README.md`, `docs/ARCHITECTURE.md` and `docs/diagrams/*.md` to the mandatory templates of documentation**: `README.md` ≤ 150 lines with 4 ordered sections (`Why This Exists` + Without/With comparison and inline concept links, `Before You Start`, `Step-by-Step Guide` with complete compilable code and links to source files instead of constant/error tables, `Documentation Map` table), 4 ordered sections in `docs/ARCHITECTURE.md`, and 3 parts in each diagram file (context header, `flowchart TD`, numbered walkthrough). Write the exact target sections in the plan so the executor never writes a bare link dump or snippet with undeclared variables.
- **A plan for a `veltylabs/*` repo ends with a language stage**: translate to Spanish every README, `docs/*` file, diagram and code comment of the repo that is still in English (only those; never rewrite what is already Spanish), and add a usage example to exported-symbol comments that lack one (**documentation**).
- **A plan for a `webtyp/*` library that touches any `fmt.Err`/`lang.Translate`/`Label`/`Help` text adds `langc` to `## Tools` and a translation stage**: run `langc sync .`, then fill every empty value it reports — including columns opened by widening `languages` to `lang.Supported()` minus `en` (**langc**), not only the keys the plan's own change added. Delivery check requires `langc sync .` to report zero untranslated keys and **zero reuse warnings** (`share ... — extract it as its own key and reuse it`); repeated phrases must be decomposed at the call site so `fmt.Err` passes shared vocabulary, never dismissed as acceptable. webtyp is a global framework: a library's dictionary rides along with whatever change already touches it.
- Never mention `gopush` or `codejob` inside the plan.

## `## APIs used` — the executor must never guess a signature

The executor cannot see the code it calls and **invents what the plan only names** (a real plan said "call `goflare.NewTarget`" and got code against a signature that returned one value instead of two, an `Item` of the wrong type and methods that did not exist; the PR did not compile).

- Before writing the plan, **read the source at the published tag** (module cache or `go doc <pkg> <Symbol>`), never from memory.
- A `## APIs used` section lists every symbol of another repo the stages call: `symbol | exact signature | tag | file`. Include return values (`(*Target, error)`), the named types of fields and constants (`deploy.TargetID`, not `string`), interface method sets, and the zero-value or error behaviour the stages rely on.
- When a stage wires several of them together, add the exact call sequence as a code block that compiles.
- Close the section with: «If a signature in the published tag differs from this table, stop and name the symbol in the PR; never adapt it by guessing.»
- The Delivery check starts with `gotest` (it builds and vets; never add `go build`/`go vet`); the plan states that **a PR that does not compile, or whose tests could not run, is a failed delivery**, and gives the environment (`GOPRIVATE=webtyp.com/*,github.com/veltylabs/*`, `GOFLAGS=-mod=mod`, `go mod tidy` before generators).

## The plan carries the REAL API — the executor never guesses how to use it (rule, 2026-11-11)

Every plan that calls code (this repo's, a dependency's, a generator's) ships the **real, verified** usage:

- **Read the source first, at the tag the plan pins** (`go doc <pkg> <Symbol>`, or the file): exact signature, field names, return values, error values, what the framework does around the call. Never write an API from memory or from an older plan: constants that became functions, removed packages and renamed fields caused whole days of failed sessions (`return crud.ApplyAndPublish, nil`, `router/loopback`, `view.Ops{List:…}`, `orm.NewMigrator`).
- **Code blocks are code that compiles.** Verify before dispatch: put the block in the target repo (scratch branch or working tree), run `gotest`/`go vet`, then revert or keep it. A plan block that was never compiled is a guess.
- **Never write «adapt it if the compiler disagrees», «check with go doc», «if X does not exist…».** That hands the guess to the weakest agent. Resolve it yourself and write the answer.
- **A test is code, not a table row**: give at least the helper (module builder, fakes, the call helper) and the first two tests as Go that compiles; the table lists the rest with exact values.
- **A generator's behaviour is a fact to state** (what `definitionc` writes and where, what makes it write a skeleton, what it rejects), found by running it in a scratch copy.

## A plan is structures and examples, not a description

A plan that says *what* without *the exact shape* makes the executor ask, stall or invent (the catalog modules took hours: lost files, a hand-written FEFO, a wrong `lang.json`, an unknown ORM call, every one a question the plan could have answered). For a small module the plan must carry:

- **Every type written out**: the full Go struct or `model.Definition` literal with field types and constraints, the constants with their values, the error constants verbatim. Never "an entity with the usual fields".
- **A worked example per operation**: a table `input → effect/output`, with at least one error row, and for the non-obvious ones the body as Go that compiles (the transaction, the call to the library, the revision check).
- **Each new file listed with its exported symbols**, or a **reference file to copy** (path in a published repo) and exactly what to rename.
- **Literal samples of data files** (`lang.json`, JSON, SQL): paste three lines of the real format, and the command that generates it (`langc sync .`).
- **Each test named `Test<Area>_<Rule>`, in `tests/<area>_test.go`, with its row in `docs/TESTS.md`** (**testing** «Domain tests»; the row is `| [Test<Area>_<Rule>](../tests/<area>_test.go) | rule |` under header `| Test | Rule |` — the link IS the file; a separate File column fails the `gotest` tests-index check and stalls the executor); a test of another repo's behaviour is not written here.
- **Each test as setup → call → expected**, not "test the rules". An expected value read from a list names the index and the list's order (`ListAudit()[0]`, newest first) and the exact field — never "the last entry" (`access`, 2026-10-30: Jules read the oldest row).
- **Diagram first, then its tests (TDD)**: a plan that adds or changes a flow writes its `docs/diagrams/<FLOW>.md` (mermaid, **documentation**) in the first stage and, before the code, the mocked tests that walk **every branch** of it — each success path and each failure path (dependency down, bad input, rejected signature…) with the exact error expected. The code stage comes after the red tests (**testing** → Diagram-Driven Testing). A flow plan with only happy-path tests is incomplete: do not dispatch it.
- **Code pasted from the repo inherits its defects.** Before copying existing code into a plan («unchanged parts are copied verbatim»), check its imports and names against the package's rules (**isomorphic-libraries**, **wasm**, **api-design** naming): a stdlib `strings` in a package the browser imports, or a name that fails the novice test, becomes a stage of the plan that fixes it, never a copy of it.
- **Toolchain limits already known** (e.g. `definitionc` rejects module constants and local kinds inside `input.Select(...)`; no `map`, no stdlib `time`) go in the stage they affect.
- **The delivery rules, in the plan itself**: never pause to ask permission or report progress — when the delivery check passes, commit, push and open the PR; stuck for 30 minutes, open the PR anyway listing what fails. Never edit a generated file (`model_orm.go`): a build error there means the hand-written code is wrong. Without these lines the agent stops ("unable to complete") or "fixes" the generator's output. Also: **commit after every stage** and never run `git reset`, `git restore`, `git checkout -- <file>` or `git stash` on its own work — a debugging reset once wiped a finished stage.
- **The 3-hour rule (granularity and detail in the *how*)**: if an agent takes >3 hours on a plan, it is a symptom that the plan was not detailed enough in the *how*. Ambiguous directions, broad stages, or missing exact signatures force the executor to invent or get trapped in test regressions. When that happens: do not keep sending band-aids; split the plan into smaller, tighter stages with explicit examples and exact contracts, and re-dispatch.
- **The plan forbids hand-written `go.sum`/`go.mod` lines**: only `go get` and `go mod tidy` write them. An executor that could not download a module once wrote `h1:XXXX…` hashes into `go.sum` (`app`, 2026-11-10) and `gotest` died with «SECURITY ERROR»; the plan says: if a module cannot be fetched, stop and say so in the PR, never invent a hash.
- **Ambiguity audit before dispatch.** For each stage, write down the questions an executor with zero context would ask — which argument type, which signature, which test data, which file, what to do when the old behaviour and the new rule disagree — and answer each **inside that stage**, as code or a table. A stage whose «how» is a verb («port the checks of X») instead of a body is not dispatchable. The plan also tells the executor to open the PR after the third stage (Jules publishes only by submitting it) and never to ask a yes/no about progress. Module migrations: paste the blocks of **crud-module-migration**.
- **Pin every dependency version the plan relies on** (`go get webtyp.com/storage@v0.1.10`) in `## Tools`: a fix published while the plan runs is invisible to an executor that resolved the old version.
- **The body of every non-trivial function is written in the plan** as Go that compiles, not as an ordered list of steps: an executor given steps re-derives the code, gets it subtly wrong and debugs itself into destructive commands.
- When the executor asks twice the same kind of question, the plan lacked a section: answer it, then add the section to the plan.

## Tool catalog — copy the lines the plan needs

| Tool | Install |
|---|---|
| `gotest` (always: every plan runs the suite) | `go install webtyp.com/devflow/cmd/gotest@latest` |
| `langc` (translations, `lang.json`) | `go install webtyp.com/lang/cmd/langc@latest` |
| `definitionc` (regenerates `definition_gen.go`) | `go install webtyp.com/definitionc/cmd/definitionc@latest` |
| `sitec` (SSR / asset extraction) | `go install webtyp.com/sitec/cmd/sitec@latest` |
| `weightsc` | `go install webtyp.com/weightsc/cmd/weightsc@latest` |

Never in a plan: `gopush`, `codejob`, `devskills`, `webtyp` (local developer tools).

## Modular stage files

Complex plans: `docs/PLAN.md` is the index; stages go to `docs/PLAN_STAGE_<N>_<TOPIC>.md`, each opening with `← [Stage N-1](...) | Next → [Stage N+1](...)`.

When porting code, append the original snippets at the end of the stage file and say what to recycle and what to replace.
