---
name: bug-fix-flow
description: The default, always-on response the instant a bug is confirmed to live in a shared library (webtyp/* or veltylabs/modules/*), not the consuming app. Find or write a failing test in the affected library that proves the bug, then fix it — inline for a single-file fix, via docs/PLAN.md + immediate codejob dispatch otherwise. Never stop to ask what to do next; this flow IS the answer. Use whenever a bug, defect, or unexpected behavior traces back to a shared library.
---

# Bug Fix Flow

> This is the execution flow. Skill **plan-authoring** covers what a
> `docs/PLAN.md` must contain; skill **agents-workflow** covers dispatch,
> review, and closing the loop with `codejob`. This skill is the glue that
> decides, the moment a bug is found, which of those paths to take —
> automatically, without waiting to be told.

## Trigger

You have just confirmed — by tracing code, not by guessing — that a bug's
root cause lives in a shared library the app only consumes: any
`webtyp/*` repo or any `veltylabs/modules/*` repo. This is distinct from a
bug in the consuming app's own composition code (its `config/`, its
`modules/<m>/browser.go`/`server.go` wiring) — those are fixed in the app
directly, this flow does not apply to them.

The instant you reach this point, **do not report the finding and wait**.
Run the flow below. Reporting "I found X, want me to fix it?" is the wrong
response to a confirmed, traced bug — it is the right response only to a
genuinely open design question (see "When to actually stop" below).

## Step 1 — Find the affected library's own tests

`cd` into the library and locate its test suite (`tests/` directory, or
`*_test.go` beside the code, per that repo's own convention — read its
`AGENTS.md` if one exists, it states the layout). Search for a test that
already exercises the broken behavior.

- **A relevant test exists and is green despite the bug** → the test's
  assertion is too weak to catch it (a "consumer-shaped" test would have
  failed — see skill **api-design**'s zero-debt rule against doubles-only
  tests). Note this, and still do Step 2 — tighten or replace it so it fails
  for the right reason.
- **No relevant test exists** → go to Step 2.

## Step 2 — Write a test that fails, for the real reason

Write it using that library's own real conventions and real collaborators —
never a synthetic shortcut that would pass even with the bug still present:

- A wire/transport bug (client↔server encoding, RBAC gate, op routing):
  drive the REAL stack — a real server (`httptest.NewServer`), the real
  client/caller, real collaborators — a fake only at the very edge (e.g.
  `storage/mem` for the DB, `router/mock` for pure unit-level registry
  tests). A hand-rolled raw HTTP request that skips the actual client code
  proves nothing about the actual bug.
- A model/validation bug: build the module the way its own tests already
  do (`orm.New(mem.New())`, its own `Deps`), and assert the exact behavior
  that is wrong today.
- A UI/component bug in a WASM package: check whether a stdlib (`!wasm`)
  test can reach it (many `crudview`/`layout` bugs can, via a fake
  `Presenter`/`Lister`) before reaching for a `_wasm_test.go`.

Run the library's test runner (`gotest`, never `go test` directly, once
`devflow/cmd/gotest` is installed) and **confirm it fails, and fails with
the exact symptom you traced** — not a different error, not a compile
error for the wrong reason. A test that doesn't compile because it targets
a not-yet-existing API (a genuine new-field addition) is an acceptable red
state — see skill **api-design**'s Design Gate for that case — but a test
failing for an unrelated reason means you haven't isolated the bug yet.

Run the FULL suite (`gotest`, not just the new test) to confirm nothing
else is currently broken and nothing else regresses once you add it.

## Step 3 — The fork: how many files does the fix touch?

**Single file** (the correction is confined to exactly one file in this
library — a model definition, one function's body, one widget assignment):

1. Apply the fix directly, right now, in this session.
2. Re-run `gotest` — must be fully green, including the new test.
3. Publish: `gopush 'message'` (never `codejob` — there is no plan, nothing
   for it to close).
4. Bump every consuming app you know is affected to the new version
   (`go get <module>@<new-version>`), rebuild (server + wasm client if
   applicable), and re-run that app's own tests too.
5. If you have live-app verification available (a running dev instance,
   browser tooling), confirm the fix end-to-end — don't stop at "it
   compiles."

No `docs/PLAN.md` is written for a single-file fix. Writing one anyway is
its own kind of waste — see skill **plan-authoring**'s litmus test: a plan
exists so a less-capable executor can implement without design judgment; a
one-file, already-understood fix needs no such document.

**More than one file, a real design decision, or any new/changed public
API** (an exported symbol, a signature, a CLI flag, a declaration format):

1. Do **not** touch the code directly.
2. Write (or replace, per skill **agents-workflow**'s "never clobber an
   existing plan" rule) `docs/PLAN.md` at that library's module root,
   following skill **plan-authoring** in full: self-contained, exact file
   names, quoted error messages, a `## Design gate` section if any public
   API changes (skill **api-design** — prior art, novice-name test,
   complexity ledger, where it belongs, what it deletes), the failing test
   already committed as the acceptance criterion, and a stages table.
3. Dispatch it **immediately**: `cd` into that repo, run bare `codejob`
   (no arguments, no message). Do not wait for the user to say "despacha" —
   a confirmed, traced, multi-file bug with a plan ready IS the dispatch
   trigger.
4. Move on to the next found bug (see "Working multiple bugs" below) while
   it runs in the background — do not idle waiting for one plan to finish
   before starting the next repo's investigation.

## Step 4 — Closing a dispatched plan

When you learn a dispatched plan's PR is ready (the user tells you, or you
check with `codejob` per skill **agents-workflow**):

1. Pull it in: `cd` into that repo, bare `codejob` (advances
   `running → review`, checks out the PR branch).
2. Read `docs/PLAN.md` on that branch, inspect the actual diff against every
   stage, and run `gotest`.
3. A deviation from the plan is not automatically wrong — judge it on its
   own merits (e.g. an executor discovering and fixing a necessary
   side-effect the plan didn't anticipate, with a sound justification, is
   correct work, not a defect to reject).
4. Correct and green → close it yourself: bare `codejob` again (this is the
   merge + `gopush` + delete-`docs/PLAN.md` step — **not** a manual
   `gopush`, the plan loop owns its own publish). Then bump every
   consuming app, same as the single-file path's step 4.
5. Broken or incomplete → write a new `docs/PLAN.md` with the specific
   fix (skill **agents-workflow**'s "Error Handling After Agent Execution"
   table) and dispatch that instead of editing code by hand.

## Working multiple bugs in the same session

Bugs found together are not necessarily in the same library. Apply this
flow to each independently: single-file ones get fixed and pushed right
away; multi-file ones get a plan and a dispatch each, fired off in
sequence without waiting on one another. Only after every found bug has
been routed through Step 3 do you circle back to check on dispatched
plans or report a summary.

## When to actually stop and ask

This flow is the default and should not be narrated or asked-permission-for
mid-execution. Stop and ask only when:

- The fix requires a genuine, un-made product/design decision this session
  cannot infer from the code or existing docs (e.g. "should creating a
  funcionario also provision a login account?") — write that question down
  as an explicit out-of-scope note in the plan, and ask the user
  separately; do not block the rest of the flow on it.
- Two or more defensible designs exist for a new public API and skill
  **agents-workflow**'s Q&A gate requires presenting the trade-off before
  writing the plan.
- The action is destructive or affects shared state beyond what's already
  understood (force-push, deleting another team's branch) — ordinary
  behavior per the harness's own safety rules, not specific to this flow.

A request for "should I fix this?" about a bug you have already traced to
its root cause is not one of these — that question has already been
answered by finding the bug. Fix it.
