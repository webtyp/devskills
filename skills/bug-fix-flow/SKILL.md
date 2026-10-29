---
name: bug-fix-flow
description: The default, always-on response the instant a bug is confirmed to live in a shared library (webtyp/* or veltylabs/modules/*), not the consuming app — prove it with a failing test, name the hole in the framework that let it be written, close that hole so the same bug cannot compile or pass the gate anywhere, sweep every other instance, then fix inline (one file) or via docs/PLAN.md + immediate dispatch. A fix at the call site alone is a patch, never the final state. Never stop to ask what to do next; this flow IS the answer.
---

# Bug fix flow

Applies once the root cause is **traced** to a shared library. A bug in the app's own `config/` or module wiring is fixed in the app — and still goes through step 2, because the app could write it only because the framework let it. Do not report and wait: run the flow.

## 1. Prove it — a test that fails for the real reason

- Use the library's own conventions (`AGENTS.md`, `tests/`). Search the tests first (**testing** → «Search before you write»): an existing test of that scenario is **extended** with the failing assertion, not duplicated; a green one that misses the bug is too weak — tighten it.
- Real stack, fake only at the edge (`httptest.NewServer` + real client; `orm.New(mem.New())`; a real codec, never a fake writer). A synthetic shortcut that passes with the bug present proves nothing.
- UI bugs: try a stdlib test with a fake `Presenter` before a `_wasm_test.go`.
- `gotest` must fail **with the traced symptom**, and the full suite must be otherwise green.

## 2. Name the hole — why could this be written at all?

Fixing the one place it happened is half the work. The fix must **remove** code or rules, never pile them up: the simplest fix that closes the hole wins. Answer in writing (plan or closing message) **which of these let it through**, and close it in the library that owns it (**api-design** gate 4: at the root, never at the leaf):

| The hole | Closed by |
|---|---|
| The signature accepts the wrong state: a free `string` where a type exists, a `Close()`/`Init()` to remember, a `bool`, an optional field that must be set | change the type so the wrong code **does not compile** (**api-harness** rules 3, 5) |
| The rule lives only in a comment, a doc or a skill | move it into a signature, or into a gate check (`gotest`, `langc`, `sitec`, `ormc`) that fails |
| A test double hides the real transport or codec | a conformance test on the real one, run by every implementation |
| A convention nobody checks (a file name, a wire shape) | the tool that reads it rejects the violation — only if the check needs no list of exceptions; otherwise the rule goes in its owning skill and in **plan-review** (build tags: **project-layout**) |

Then **sweep**: grep both trees (`~/Dev/Project/webtyp`, `~/Dev/Project/veltylabs`) for every other instance. The fix closes them all in the same wave — or the class fix makes them stop compiling, and each one becomes a row of the wave.

**Done means:** re-introducing the bug anywhere in the ecosystem fails at compile time or in the gate. Say it in the close: «instance fixed in X; class closed by Y». If Y is not shipped yet, its plan is dispatched and it is a row in the master — never a «known issue».

## 3. Fix — by size

**One file:** fix now → `gotest` green → `replace` in the affected app's `go.mod` (**codejob** «gopush and local `replace`») → rebuild, run its tests, verify live if tooling exists → only then `gopush 'msg'` (its cascade drops the `replace` and bumps the app). Never publish to test a fix. No plan.

**More than one file, a design decision, or public API** (step 2 usually is): do not touch code. Write `docs/PLAN.md` (**agents-workflow**, **plan-authoring**, **api-design** gate) with the failing test and the step-2 guard as acceptance criteria, then `codejob dispatch` immediately — a traced bug with a ready plan is the trigger. Move on to the next bug while it runs.

## 4. Close

PR ready → **plan-review**; correct → `codejob close 'msg'`, then bump the consuming apps.

Several bugs: route each through steps 2–3 independently before checking dispatched plans.

## Stop and ask only when

- The fix needs a product decision the code and docs cannot answer (note it in the plan; don't block the rest).
- A new public API has two defensible designs (Q&A gate of **agents-workflow**).
- The action is destructive beyond what is understood.
