---
name: bug-fix-flow
description: The default, always-on response the instant a bug is confirmed to live in a shared library (webtyp/* or veltylabs/modules/*), not the consuming app — prove it with a failing test in that library, then fix inline (one file) or via docs/PLAN.md + immediate dispatch. Never stop to ask what to do next; this flow IS the answer.
---

# Bug fix flow

Applies once the root cause is **traced** to a shared library. A bug in the app's own `config/` or module wiring is fixed in the app. Do not report and wait: run the flow.

## 1. Prove it — a test that fails for the real reason

- Use the library's own conventions (`AGENTS.md`, `tests/`). Search the tests first (**testing** → «Search before you write»): an existing test of that scenario is **extended** with the failing assertion, not duplicated; a green one that misses the bug is too weak — tighten it.
- Real stack, fake only at the edge (`httptest.NewServer` + real client; `orm.New(mem.New())`). A synthetic shortcut that passes with the bug present proves nothing.
- UI bugs: try a stdlib test with a fake `Presenter` before a `_wasm_test.go`.
- `gotest` must fail **with the traced symptom**, and the full suite must be otherwise green.

## 2. Fix — by size

**One file:** fix now → `gotest` green → `gopush 'msg'` → bump affected apps (`go get <module>@<new>`), rebuild, run their tests → verify live if tooling exists. No plan.

**More than one file, a design decision, or public API:** do not touch code. Write `docs/PLAN.md` (**agents-workflow**, **plan-authoring**, **api-design** gate) with the failing test as acceptance criterion, then `codejob dispatch` immediately — a traced bug with a ready plan is the trigger. Move on to the next bug while it runs.

## 3. Close

PR ready → **plan-review**; correct → `codejob close 'msg'`, then bump the consuming apps.

Several bugs: route each through step 2 independently before checking dispatched plans.

## Stop and ask only when

- The fix needs a product decision the code and docs cannot answer (note it in the plan; don't block the rest).
- A new public API has two defensible designs (Q&A gate of **agents-workflow**).
- The action is destructive beyond what is understood.
