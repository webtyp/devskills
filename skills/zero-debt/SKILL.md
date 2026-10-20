---
name: zero-debt
description: What must not survive a change — DRY taken to the extreme because the WASM binary ships to the browser (same logic in two repos is a third repo, no string literals in logic) and the zero-technical-debt list (no silent fallbacks, TODOs, stubs, doubles-only tests, unkeepable promises). Use before closing any change.
---

# Zero debt — before closing

## DRY — taken to the extreme, because the binary ships to the browser

The client is a TinyGo WASM binary the user downloads before seeing anything. Two libraries that
each carry their own copy of the same logic both land in it — the linker cannot merge them — so every
duplicated string, struct, method or function is bytes on the wire. Nobody loads megabytes for a
hello world; that is the main reason WASM adoption stalls in the browser, and it is why Go in the
browser is only competitive if **nothing is written twice anywhere in the ecosystem**.

- **The same functionality in two repos is a third repo**, imported by both — never a copy, never a
  local helper "just for here". This is the lego rule measured in bytes. **Search before you write**:
  **api-design** → «Step 0» is the procedure, and the plan records its result.
- A repeated string (env key, path, prefix, flag name, error message) is a named constant in the
  owning package; string literals are forbidden in logic.
- A value the library computes is **read** by `cmd/`, never re-derived; a check it performs is
  called, never re-implemented. Two functions differing only in a constant are one function and a
  parameter.
- A library's docs are an "I want X → use Y" table, one example per use case, and the signatures —
  never "remember to call…", which is a hole in the harness. A rule written twice drifts: link.

## Zero technical debt — what must not survive the change

- A fallback that silently restores the old behaviour when the new one finds nothing; a missing input is
  an error or an empty result, never a guess.
- A `TODO`, a commented-out block, or a stub for an undecided command — an undecided thing fails loudly,
  pointing at the document that decides it.
- A test that only exercises doubles, with no consumer-shaped case; an exported symbol nothing
  outside the package calls — including one whose only callers are tests.
- A promise the code cannot keep because a contract is missing downstream (an ARIA role whose keyboard
  behaviour is unimplementable, a documented option with no effect): narrow it or fix the contract.

Before closing: `grep -rn "TODO\|FIXME\|Deprecated" --include='*.go' .` on the touched packages, every
hit predating the change.
