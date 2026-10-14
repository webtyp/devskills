---
name: api-harness
description: Rules while writing a public API — lego pieces (one concern, typed contract, capability bag + type assertion, never patch at the leaf), the nine harness rules (typed over any, explicit, illegal states unrepresentable, minimal surface, compile-time failure, closed by default, narrow interfaces, composition root), publication with consumer-shaped and conformance tests, the refactor checklist and the acid test. Use while writing or reviewing exported code.
---

# API harness — while writing

The gate before coding: **api-design**. Before closing: **zero-debt**.

## Lego pieces — one concern, one typed contract

Each concern is a single-responsibility library exposing a **typed contract**; applications assemble
pieces and never re-implement, wrap or copy one.

```go
type APIModule interface {    // the seam: a module declares what it can do
	model.ModuleNaming        // identity: ModelName()
	MountAPI(r Router)        // the module registers its own routes
}
if api, ok := m.(router.APIModule); ok { api.MountAPI(r) } // each boundary asserts what it needs
```

**Capability bag + type assertion at the seam** is the assembly pattern: swapping an implementation
(another `Router`, a fake `Caller` in tests) never touches the modules. **When a library does not
expose what you need, stop and report it** — a copy downstream forks that library's responsibility,
and a wrapper over a defect is a fork with a friendlier name.

**Why the leaf must not patch.** An API gap surfaces at the application, where the agent has no
authority to publish upstream — so it patches locally, and the debt stops being an accident. And a
piece is the base of the next: each consumer inherits the defect, works around it, or **encodes the
workaround into its own contract** and passes it on, so cost and blast radius grow with every piece
built since. The leaf waits for the base's tag; it never ships a promise a missing contract prevents
it from keeping, it narrows the claim.

> *Measured here:* `dom.Element` exposed one event method, `On(t string, h func(Event))` — the type an
> untyped string, so `dom.Event` never grew a keyboard contract. A calendar could not implement WAI-ARIA
> arrow-key navigation yet shipped `role="grid"`: an accessibility defect whose root cause is one
> untyped parameter two libraries down.

## The harness — the API is the documentation

1. **Typed over `any`.** No `func(...any)` or `interface{}` in a public API; `any` only at the I/O
   edge, never in the data; reuse the types that exist (`fmt.KeyValue`, `model.Resource`). House
   pattern: `webtyp/json`, one method per primitive (`String`, `Int`, `Object`, …).
   ```go
   func (b *Builder) Add(items ...any) *Builder        // ❌ hole: intent lost, fails at runtime
   func (b *Builder) Text(s string) *Builder           // ✅ typed by intent: misuse does not compile
   func (b *Builder) Set(kv ...fmt.KeyValue) *Builder  // ✅ reuses a type declared in fmt
   ```
2. **Explicit over implicit.** Reading the call must be enough to know what it does; runtime discovery
   by type assertion is not — prefer a written line.
3. **Illegal states unrepresentable.** One intent = one path, typed to demand what it needs.
4. **Minimal surface.** Export exactly what the author uses; plumbing stays unexported. **Tests are
   not authors**: a symbol exported so a test can reach it is not surface, it is a defect (skill
   **testing** → «Never export API for a test»).
5. **Fail at compile time.** Compile error → loud development diagnostic → never a silent failure.
6. **Self-describing signatures.** Autocomplete must be enough to build; if the API needs a long document,
   the API is incomplete.
7. **Closed by default.** Where an API governs access, the zero value **denies**; opening is an explicit,
   greppable, typed line. Reachable because *nobody said otherwise* breaks rule 5.
8. **Narrow interfaces (I).** If any implementation would stub a method the contract is two — a stub
   returning `nil` is the silent failure rule 5 forbids. Segregate, then compose: `type ReadWriter
   interface { Reader; Writer }` keeps the name consumers use (`io.Reader`, `fs.FS` + `fs.ReadDirFS`).
9. **The composition root wires, nothing else (D).** Logic depends on contracts; one place — the
   constructor — knows a concrete type. If the library defining the abstraction requires in its `go.mod`
   the concrete thing it hides, the inversion is prose, not build graph.

### Publication

An API is not published until a **consumer-shaped test inside the library itself** proves it, through
the real stack a consumer will use with a fake only at the edge: a CRUD layout that takes a model,
generates a form and ships it through a caller is tested with a real model, the real form package and a
fake caller. Doubles-only tests hide the gaps until a consumer hits them, and a consumer can only patch;
an awkward test means an awkward API.

A contract with a **second implementation** is not published until a **conformance suite in the
repository owning the contract** proves them substitutable **(L)** and all of them pass it — as
`storage/conformance` and `ddl/conformance` do for `mem`, `sqlt`, `postgres` and `indexdb`. Without it
substitutability is a claim, found false by a consumer in production on the second backend; it also
makes "fix at the root" cheap — the fix lands once, every implementation is re-proved.

## Aligning an existing library — the refactor checklist

Hunt for and fix: every `any`/`interface{}`/`...any` in the public API · invariants checked with an
`if` plus an error or panic that a type could have made unwritable · **things you "have to remember"**
(call order, "don't forget X") · more than one way to do the same thing · exported plumbing only the
library uses · misuse that produces neither an error nor a visible effect · zero values that permit
instead of denying · boundaries where a consumer would declare a local interface · glue every app
writes identically. Afterwards the only failure modes are a **compile error** or a **loud diagnostic**.

## The acid test

An agent, or a junior, **with no context on the library** produces correct code guided only by
autocomplete and a few-line example. If they must read a manual to avoid a mistake, something is still
untyped. If they must **ask a question** to proceed — "is it OK to declare this interface locally?" —
the harness failed by its own definition.
