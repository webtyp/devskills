---
name: plan-review
description: Reviewing the PR of a dispatched docs/PLAN.md — evidence checks commit by commit, docs verification, running the suite, at most 2 feedback attempts to the agent then local fix or local execution, and what to do with errors after merge. Use when a plan's STATUS is review/reviewing or the user asks to review a plan's PR.
---

# Reviewing a plan's PR

Checkout: `codejob pull` (**codejob**). The spec is `docs/PLAN.md` on the PR branch.

## Checks — the PR is its commits, not the agent's claim

1. Commit by commit: `gh pr view <url> --json commits --jq '.commits[]|.oid[0:7]+" "+.messageHeadline'`, then `git show --stat <sha>`. A later commit can revert an earlier one.
2. Every file of the plan's stages table appears in `git diff --stat main`.
3. Nothing deleted that the plan did not order: `git diff --stat main --diff-filter=D`.
4. No stray files (`*.orig`, `*.rej`, scratch).
5. Local branch equals the PR head.
6. Read the code against each stage. A storage/transport contract needs a consumer test on a real backend, not only the in-memory double.
7. Docs follow the mandatory templates of **documentation**: `README.md` is ≤ 150 lines with all 4 ordered sections (`Why This Exists` + Without/With and inline links, `Before You Start`, `Step-by-Step Guide` with compilable code and links to `.go` files instead of constant/error tables, `Documentation Map` table), `docs/ARCHITECTURE.md` has all 4 sections, and every `docs/diagrams/*.md` has context + `flowchart TD` + numbered walkthrough. Reject oversized READMEs, bare link dumps, or code snippets with undeclared variables.
8. Run `gotest` yourself, **once**: it already builds, vets, races and runs every destination (stdlib + wasm). Never add `go build`, `go vet` or a second `gotest` run — the plan's own delivery check asks for a repeat only when the repeat is the thing being proven. A PR that does not compile means the agent **guessed signatures**: fix it locally against the real source, and add the missing `## APIs used` rows to the plan's lesson (**plan-authoring**).
9. **A change the user sees** (screen, menu, icon, text, form) is opened in the module's demo (`web/browser/main.go`, run with `webtyp` at the module root) and **looked at in a screenshot** before approving. Tests prove behaviour and never the look (**testing**), so this look is the only check that catches a wrong glyph or an empty screen.
10. **Translations and text**: if the repo has `lang.json` or touches user-facing copy/errors, run `langc sync .` (**langc**) — must report 0 untranslated keys, 0 reuse warnings (`share ... — extract it as its own key and reuse it`), and 0 identical values (`ERROR ... ends with ... unchanged`). Reuse warnings are defects, never acceptable exceptions.
11. **Build tags**: `git diff main | grep '^+//go:build'` — a non-test file with a tag other than `wasm`, `browser`, `dev` (or Go's OS names in OS code) is rejected (**project-layout** «No other tag»).

A deviation from the plan is judged on its merits, not rejected by default.

## Outcome — at most 2 attempts by the agent

| Situation | Action |
|---|---|
| All correct | Tell the user to merge, or `codejob close 'msg'` |
| Small fix (1 file, few lines), any round | Edit it on the PR branch and leave it **uncommitted**: `codejob close` commits it ("review: corrections before merge") and pushes. Never `git commit`/`git push` |
| Attempt 1 or 2 fails (PR wrong, tests red) | Feedback to the agent: `codejob reply "<findings>"` while `running`, or one `gh pr comment <url> --body-file <f>` on its PR — each finding with file:line, what is wrong, what the plan required, how to verify; no code edits |
| Attempt 2 still wrong, PR published | Fix it locally on the PR branch, `gotest`, `codejob close 'msg'` |
| Attempt 2 still wrong, nothing published | Wait for the agent to publish; if it never does, execute the plan locally (**agents-workflow** → «Local execution») |
| Error found after merge | New `docs/PLAN.md` with the exact fix |
| The plan's design was wrong | Q&A with the user first (**agents-workflow**) |
