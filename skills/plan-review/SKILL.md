---
name: plan-review
description: Reviewing the PR of a dispatched docs/PLAN.md — evidence checks commit by commit, docs verification, running the suite, review rounds (comment first, fix second) and what to do with errors after execution. Use when a plan's STATUS is review/reviewing or the user asks to review a plan's PR.
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
7. Docs updated where the change requires (`docs/API.md`, `ARCHITECTURE.md`, `README.md`, `docs/SKILL.md`, any doc the plan lists).
8. Run `gotest ./...` yourself (stdlib + wasm).

A deviation from the plan is judged on its merits, not rejected by default.

## Outcome

| Situation | Action |
|---|---|
| All correct | Tell the user to merge, or `codejob close 'msg'` |
| Small fix (1 file, few lines), any round | Fix it on the PR branch yourself |
| Round 1 errors | One `gh pr comment <url> --body-file <f>` listing each finding (file:line, wrong, required, how to verify); no code edits |
| Round 2, correction still wrong | Fix on the PR branch, `gotest`, `codejob close 'msg'` |
| Error found after merge | New `docs/PLAN.md` with the exact fix |
| The plan's design was wrong | Q&A with the user first (**agents-workflow**) |
