---
name: agents-md
description: The AGENTS.md every repo carries at its root for execution agents — the ecosystem base block (copied verbatim between markers) plus repo-specific constraints. Use before dispatching any plan (create or refresh it), or when creating a repo.
---

# AGENTS.md — the executor's rules, once per repo

Execution agents do not read skills. Every repo has `AGENTS.md` at its root (next to `go.mod`); every `docs/PLAN.md` links it and never restates it (**plan-authoring**).

- Missing, or base block differs from the one below → write/refresh it **before** `codejob dispatch` (dispatch publishes the working tree).
- Structure: the base block verbatim, then `## This repo` with only what is specific (architecture limits, house exceptions, anti-footguns). Never restate the base.
- The base block is edited **only here**; repos copy it between the markers.

## Base block

```markdown
<!-- webtyp:base — copied from devskills skill agents-md; do not edit here -->
## Ecosystem rules — read before any change

`docs/PLAN.md` is the spec: no design judgment; anything unclear is a question in the PR, not a guess.

- Public API: typed (no `any`/`interface{}`), explicit, minimal; never export a symbol only a test uses.
- One intent, one path: the old path is deleted in the same change. No silent fallback, `TODO`, commented code or stub.
- Fail at compile time, else a loud diagnostic — never silently. Zero values deny.
- A defect in a dependency is reported, never wrapped or patched here.
- No string literals in logic: named constants. The same logic never lives in two repos.
- `cmd/*/main.go`: parse flags, inject, print/exit only. stdout = data, stderr = logs, exit codes are the contract.
- WASM code: `webtyp/fmt`, `webtyp/time`, `webtyp/json` instead of stdlib; no `map`; `dom.Element` embedded by value.
- `css.go`, `svg.go`, `js.go`, `html.go` carry `//go:build !wasm`; no `ssr.go`, no `front.go`.
- Code and keys in English; UI text is `lang.Text`, translations live in `lang.json`.
- Tests cover behaviour and domain rules, never the look (CSS, colours, markup order, class strings). In `tests/` (external package), real stack with fakes only at the edge; run `gotest`.
- Docs (`README.md`, `docs/*`) updated with the API; permanent docs never link `docs/PLAN.md`.
- Before the PR: `git diff --stat main` lists every file of the plan's stages, and nothing the plan did not order is deleted.
<!-- /webtyp:base -->
```

Each line's full rule lives in **api-harness**, **zero-debt**, **core-principles**, **wasm**, **components**, **translations**, **testing**, **documentation**.
