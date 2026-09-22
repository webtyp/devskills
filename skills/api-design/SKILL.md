---
name: api-design
description: The construction harness — how every public API is designed: five gate answers in docs/PLAN.md before coding, the lego-piece and typed-harness rules while writing, zero debt before closing. Use when adding, changing, renaming or reviewing any exported symbol, CLI surface, file convention or declaration format; when starting a library; or when a defect traces to a missing contract upstream.
---

# API Design — the construction harness

**The typed, explicit code is itself the harness.** Most code here is written by an agent that does
not know the library: it must be correct from the signatures alone, and the compiler must **reject**
what is wrong. A manual pushes correctness onto the reader; a harness pushes it onto the compiler —
the right path is the only one that exists. That is what lets an application need **fewer** concepts
than the alternatives; every exported symbol pays for that or erodes it.

**Applies to** a new exported symbol, a changed signature, a new CLI subcommand or flag, a new
file/directory convention, a new declaration format, a new library — **not** to unexported code or a bug
fix that changes no signature. **Three moments:** the gate in `docs/PLAN.md` before coding · lego and
harness rules while writing · zero debt and the acid test before closing. SOLID is enforced where it is
checked, never restated: **S** and **O** in [gate 4](#4-where-does-it-belong), **L** in
[publication](#publication), **I** and **D** in harness rules 8–9; when a rule below and a principle
seem to disagree, the principle wins.

## The gate — five answers in `docs/PLAN.md`, before coding

A plan that adds or changes public API without them is incomplete: do not dispatch it.
### 1. Prior art

Name **at least three** established frameworks, what each does, and why this ecosystem differs.
"Nobody else does it this way" is valid only followed by the reason it is right here — an API invented
without looking accumulates novel vocabulary for solved problems.

> *Routes:* file-system (Next.js, Remix) · central manifest (Rails, Django, chi `Mount`) · runtime
> discovery (Spring, NestJS). Only the first two make routes knowable **without running the app** — so
> `routes/routes.go` + `Mount(prefix, fn)` is family two, chosen for that property.

### 2. The novice-name test

Read every new name aloud as a sentence a **junior with no context** would say; if it needs the docs
to be understood, it is the wrong name.

- The word the ecosystem already uses: `Mount` is chi's word; `Attach` costs every Go developer a
  lookup for zero gain.
- Intent, not mechanism: `Public()`, not `SetAuthFlag(false)`.
- No abbreviation that is not universal in Go (`ctx`, `id` yes; `authz`, `hdlr` no).
- No boolean parameters — `Test(args, false, 0, false, false)` says nothing. Typed options or
  separate methods.

**A wrong name is renamed in this change**, never scheduled: the cost is the number of consumers,
which only grows, and a name in a published tag is frozen there. Defer only for external users who
would break — zero means now.

### 3. The complexity ledger

```
Concepts the developer must learn   +N / −M
Files they must touch to do X       +N / −M
Lines at the call site              +N / −M
Ways to do the same thing           +N / −M   ← must never end positive
```

The last row is absolute — **one intent, one path**. If the new path is better, the old one dies in the
same change: not "deprecated for one version", not behind a fallback. A ledger reported as all-better
is unreviewed.

### 4. Where does it belong

- **One concern per library (S).** A second concern is a new package or repo, not another file. The
  unit is the **repository**: one that exposes a contract *and* ships an implementation has two
  responsibilities — a dependency most consumers never instantiate is the symptom, in the `go.mod`.
- **New capability = new implementation (O).** The second backend or adapter must touch **zero lines**
  of the consumer; if it touches one, fix the contract, never add a branch.
- **Fix at the root, never at the leaf.** A consumer never re-creates a missing symbol, never
  wraps a library to patch its behaviour, never works around a defect one layer down. A defect in
  a base piece **blocks this plan**: declare the dependency on its fix and wait for the published
  tag — never a "known issue", never a local patch with a comment ([why](#lego-pieces--one-concern-one-typed-contract)).
- **A missing contract at a boundary is a defect upstream.** If two libraries meet and no type names
  what crosses between them, the type is missing there — never a local intersection.
- **No `internal/` folders** — the signature of a duplicated dependency instead of a contribution
  upstream.
- **The glue is written once, in the library that owns it.** Wiring every application would repeat belongs
  to a piece, not to the applications.

### 5. What does this change delete?

Name it. A change that deletes nothing is either genuinely new capability — say so — or a parallel path
being accumulated.

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
4. **Minimal surface.** Export exactly what the author uses; plumbing stays unexported.
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

## DRY — taken to the extreme, because the binary ships to the browser

The client is a TinyGo WASM binary the user downloads before seeing anything. Two libraries that
each carry their own copy of the same logic both land in it — the linker cannot merge them — so every
duplicated string, struct, method or function is bytes on the wire. Nobody loads megabytes for a
hello world; that is the main reason WASM adoption stalls in the browser, and it is why Go in the
browser is only competitive if **nothing is written twice anywhere in the ecosystem**.

- **The same functionality in two repos is a third repo**, imported by both — never a copy, never a
  local helper "just for here". This is the lego rule measured in bytes.
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
  outside the package calls.
- A promise the code cannot keep because a contract is missing downstream (an ARIA role whose keyboard
  behaviour is unimplementable, a documented option with no effect): narrow it or fix the contract.

Before closing: `grep -rn "TODO\|FIXME\|Deprecated" --include='*.go' .` on the touched packages, every
hit predating the change.

## The acid test

An agent, or a junior, **with no context on the library** produces correct code guided only by
autocomplete and a few-line example. If they must read a manual to avoid a mistake, something is still
untyped. If they must **ask a question** to proceed — "is it OK to declare this interface locally?" —
the harness failed by its own definition.
