---
name: testing
description: What deserves a test (behaviour and domain rules, never the look) and how tests are written and run — gotest, tests/ location, never export for a test, mocks, WASM/stdlib dual pattern, diagram-driven tests. Use before writing, planning or reviewing any test.
---

# Testing

> **Which test an API owes before it may ship** is not here: skill
> **api-design** («Publication») requires a **consumer-shaped test inside the
> library itself**, through the real stack a consumer will use, with a fake only
> at the edge — plus a **conformance suite** in the repo owning the contract once
> a second implementation exists. This skill covers HOW tests are written and
> run; that one covers WHICH test an API owes.

## What deserves a test — behaviour that must keep working, never the look

Ask before writing one: **"does it fail when the product is broken, and only then?"** No → do not write it.

- **Test:** domain rules (validation, state machines, calculations, permissions, tenant isolation); contracts between pieces (encoding, routing, a public API's observable behaviour); behaviour a user relies on (form validation timing, submit, keyboard and focus, the ARIA contract a role promises); every confirmed bug (**bug-fix-flow**).
- **Do not test:** the look — CSS declarations, colours, sizes, spacing, markup order, class strings, copy; trivial code (a recipe composing two options, a getter); implementation details instead of outcomes. These are *change-detector* tests: they break on every intended change and never catch a real defect.
- The look is verified **by eye in the running app** (webtyp MCP screenshot), not by the suite.
- **Integration with real models runs in `webtyp/agenteval`'s laboratory** (agent + avatar + UI in the browser, never a lab page per app); every run's findings go to `agenteval/docs/REAL_TESTS.md`.
- Structural guards are written **once per repo** (a conformance suite: stylesheets validate, classes match, Kind allows states), never copied per component.
- One test per behaviour, through the public API, the way a consumer uses it.
- **Search before you write.** Before adding a test, grep the repo's tests for the behaviour, the symbol and the message it would assert. A test that already pins the scenario is **extended** (same fixture, one more assertion); a test that pins a behaviour you are deliberately changing is **updated or deleted** — never left failing, never shadowed by a second test. A new test is justified only by a scenario no existing one reaches; say which tests were extended, updated and new.

## Running and placing tests

- **Testing Runner (`gotest`):** For Go tests, ALWAYS use the globally installed `gotest` CLI command. **DO NOT** use `go test` directly, and **DO NOT** invoke it via `go run webtyp.com/devflow/cmd/gotest`. Simply type `gotest` (no arguments) for the full suite, or `gotest -run TestName`. It automatically handles `-vet`, `-race`, `-cover`, WASM tests, and README badges.
- **`gotest` in a plan:** its install line goes in the plan's `## Tools` (**plan-authoring** → «Tool catalog»).
- **Publishing (`gopush`):** If tests pass and docs are updated, ALWAYS use the globally installed `gopush 'your commit message'` CLI command to deploy. **DO NOT** use standard `git commit` / `git push`, and **DO NOT** invoke it via `go run`. It handles testing, tagging, pushing, and updating dependencies automatically.

- **Test location — `tests/`, the root only when justified:**
    - Tests live in `tests/` (never `test/`), as an external package (`package <lib>_test`) that uses **only the exported API** — the same view a consumer has.
    - A test stays at the module root (internal package) **only** when it needs an unexported identifier and the behaviour cannot be observed through the public API. Its first line justifies it: `// Root-level test (justified): exercises <unexported identifiers> — <why the behaviour is not observable through the exported API>.` A root test without that line is a defect.
    - Mechanical criterion for plans and reviews: only exported identifiers → `tests/`; needs an unexported one → first try the public API; still impossible → root + justification. List each root file with its reason in the PR.
- **Never export API for a test:** a symbol whose only user is a test (a `NewWithRunner`, a `SetHome`, a `XFrom(…, home)` seam, a `Seed`) is unnecessary public API and a defect (see skill **api-design**). Control the environment with test tools instead — `t.Setenv("HOME", …)`, `t.Chdir`, `t.TempDir`, build tags (`//go:build !wasm`), or testing the same algorithm through another public entry point (e.g. an explicit-base variant). If none works, the test is a justified root-level test, never a new export.

- **Standard Library Only:** **NEVER** use external assertion libraries (e.g., `testify`, `gomega`). Use only the standard `testing`, `net/http/httptest`, and `reflect` APIs.
- **Mocking (No I/O):** Tests MUST use Mocks for all external interfaces to remain fast, deterministic, and side-effect free.

- **Automated tests own correctness; manual review is ONLY for look & feel.** Every functional behaviour — auth gates, validation, error codes, state transitions, data isolation, secret/token verification — MUST be covered by an automated test, exercised with the same severity locally and in the deployed target. The **only** thing a developer verifies by hand is the visual result and the user experience (layout, spacing, copy, flow), because taste cannot be automated. This forbids the tempting shortcuts: no `if dev` branch that skips a check, no seeded fixtures that production lacks, no behaviour that differs between local and deployed "to make local testing easier". If something is awkward to test locally, the answer is an injectable fake or a test seam — never a divergent code path. The seam is either a real dependency a consumer also injects, or it stays unexported and is used from a justified root-level test; it is never a public symbol added for the test.

- **WASM/Stlib Dual Testing Pattern:**
    - **Separation:** Use build tags for isomorphic code (`frontWasm_test.go` -> `//go:build wasm`, `backStlib_test.go` -> `//go:build !wasm`).
    - **Shared Logic:** Both files MUST call a shared test runner (e.g., `RunAPITests(t)`) to avoid duplication.
    - **Unified Setup:** Use a single `setup_test.go` to initialize the library/test server once.

- **Diagram-Driven Testing (DDT) & Black-Box Validation:**
    - **Flow Coverage:** Logic flows defined in `docs/diagrams/*.md` MUST have corresponding Integration Tests covering all branches/diamonds.
    - **Visual/Binary Outputs:** Never write tests that perform exact string/byte matching on complex binary formats (e.g., PDFs, SVGs). Floating-point variations make them brittle. Use black-box integration testing that outputs a file to disk (`os.WriteFile`) to be verified visually by the developer.
