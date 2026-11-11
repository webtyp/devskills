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
5b. **The PR is judged against the plan, not against `gotest`:** count the tests of the plan's table against the `func Test…` delivered (a PR that skips the tests it could not make pass is green and wrong: `business_calendar` PR #20 shipped only helpers and the note «ignoring them»); a plan that ships a demo or a wasm file needs its wasm test, or `gotest` never runs the wasm destination and a broken `seed/demo.go` stays hidden.
6. Read the code against each stage. A storage/transport contract needs a consumer test on a real backend, not only the in-memory double.
7. Docs follow the mandatory templates of **documentation**: `README.md` is ≤ 150 lines with all 4 ordered sections (`Why This Exists` + Without/With and inline links, `Before You Start`, `Step-by-Step Guide` with compilable code and links to `.go` files instead of constant/error tables, `Documentation Map` table), `docs/ARCHITECTURE.md` has all 4 sections, and every `docs/diagrams/*.md` has context + `flowchart TD` + numbered walkthrough. Reject oversized READMEs, bare link dumps, or code snippets with undeclared variables.
8. Run `gotest` yourself, **once**: it already builds, vets, races and runs every destination (stdlib + wasm). Never add `go build`, `go vet` or a second `gotest` run — the plan's own delivery check asks for a repeat only when the repeat is the thing being proven. A PR that does not compile means the agent **guessed signatures**: fix it locally against the real source, and add the missing `## APIs used` rows to the plan's lesson (**plan-authoring**).
9. **A change the user sees** (screen, menu, icon, text, form) is opened in the module's demo (`web/browser/main.go`, run with `webtyp` at the module root) and **looked at in a screenshot** before approving. Tests prove behaviour and never the look (**testing**), so this look is the only check that catches a wrong glyph or an empty screen.
10. **Translations and text**: if the repo has `lang.json` or touches user-facing copy/errors, run `langc sync .` (**langc**) — must report 0 untranslated keys, 0 reuse warnings (`share ... — extract it as its own key and reuse it`), and 0 identical values (`ERROR ... ends with ... unchanged`). Reuse warnings are defects, never acceptable exceptions.
11. **Domain tests** (**testing** «Domain tests»): every new or renamed `func Test…` is `Test<Area>_<Rule>` in `tests/<area>_test.go` and has its row in `docs/TESTS.md`; a test that belongs to another repo is rejected.
12. **Build tags**: `git diff main | grep '^+//go:build'` — a non-test file with a tag other than `wasm`, `browser`, `dev` (or Go's OS names in OS code) is rejected (**project-layout** «No other tag»).

13. **The consuming app**: with the PR branch checked out, the app that uses the library points at it with a local `replace` (**codejob** «gopush and local `replace`») and its `gotest` passes, before `codejob close`.

A deviation from the plan is judged on its merits, not rejected by default.

## Outcome — at most 2 attempts by the agent

| Situation | Action |
|---|---|
| All correct | `codejob close 'msg'` yourself, same turn — «revisa» is the order to close; never stop to ask |
| Small fix (1 file, few lines), any round | Edit it on the PR branch, leave it **uncommitted**, then `codejob close 'msg'` in the same turn (it commits "review: corrections before merge" and pushes). Never `git commit`/`git push` |
| Several fixes (tests the agent skipped, a schema line, docs), and **less than half of the PR is unusable** | Same: fix ON THE PR BRANCH (uncommitted) and `codejob close 'msg'`. **Never re-create the work on `main`** (the PR would end Closed instead of Merged and the agent's commit would vanish from the history). Rewriting elsewhere is only for a PR where more than 50% of the work is unusable, and then the PR is closed with a comment saying why |
| The frontmatter names another session than the PR (a relaunch left two sessions of the same repo) | The owner closes the newer session; then align `SESSION`, `PR` and `STATUS: review` of `docs/PLAN.md` with the PR being reviewed (the one permitted repair of machine fields) so `codejob pull` and `codejob close` act on it. Never close the PR and republish from `main` for that reason |
| Attempt 1 or 2 fails (PR wrong, tests red) | Feedback to the agent: `codejob reply "<findings>"` while `running`, or one `gh pr comment <url> --body-file <f>` on its PR — each finding with file:line, what is wrong, what the plan required, how to verify; no code edits |
| Attempt 2 still wrong, PR published | Fix it locally on the PR branch, `gotest`, `codejob close 'msg'` |
| Attempt 2 still wrong, nothing published | Wait for the agent to publish; if it never does, execute the plan locally (**agents-workflow** → «Local execution») |
| Error found after merge | New `docs/PLAN.md` with the exact fix |
| The plan's design was wrong | Q&A with the user first (**agents-workflow**) |
