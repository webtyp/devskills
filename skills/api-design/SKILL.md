---
name: api-design
description: The gate every new or changed public API must pass BEFORE it is written — five answers in the PLAN (prior art, novice-name test, complexity ledger, where it belongs, what it deletes), the harness rules that shape the code, and the zero-debt rules that close the change. Use when designing, adding, renaming, or reviewing any exported symbol, CLI surface, file convention, or declaration format — and when starting a library from scratch.
---

# API Design — the gate

The framework's promise is that an application built on it needs **fewer**
concepts than the alternatives. Every exported symbol pays for that promise or
erodes it, so this check runs before the first line is written.

**Applies to:** a new exported symbol, a changed signature, a new CLI subcommand
or flag, a new file/directory convention, a new declaration format, a new
library. **Not** to unexported code or to a bug fix that changes no signature.

**Three moments:** the five gate answers in `docs/PLAN.md` before coding · the
harness rules while writing · zero debt and the acid test before closing.

## SOLID — the spine, enforced in one place each

| | Principle | Checked in |
|---|---|---|
| **S** | One concern per library; the unit is the **repository**, not the file | [gate 4](#4-where-does-it-belong) |
| **O** | New capability = new implementation, never an edit to its consumer | [gate 4](#4-where-does-it-belong) |
| **L** | Substitutability is proved by a conformance suite, not promised | [publication](#publication) |
| **I** | An interface is as wide as its narrowest caller needs | [harness 8](#the-harness--the-api-is-the-documentation) |
| **D** | Only the composition root knows a concrete type | [harness 9](#the-harness--the-api-is-the-documentation) |

When a rule below and a principle seem to disagree, the principle decides.

## The gate — five answers in `docs/PLAN.md`, before coding

A plan that adds or changes public API without them is incomplete: do not
dispatch it.

### 1. Prior art

Name **at least three** established frameworks, what each does, and why this
ecosystem differs. "Nobody else does it this way" is valid only followed by the
reason it is right here. You need not copy them; you must know them — an API
invented without looking is how a framework accumulates novel vocabulary for
solved problems, and that is what makes a developer's first week expensive.

> *Worked example — route declaration.* Three families: **file-system routing**
> (Next.js, Nuxt, SvelteKit, Remix — the path *is* the file path); **central
> manifest** (Rails `routes.rb`, Django `include()`, Laravel, Phoenix, chi
> `Mount`); **decentralised discovery** (Spring `@GetMapping`, NestJS, ASP.NET —
> found by scanning at runtime). The first two make routes knowable **without
> running the app**; the third is why Spring needs an endpoint just to answer
> "what routes exist?". So `routes/routes.go` + `Mount(prefix, fn)` is family
> two — justified not by novelty but by the properties we need.

### 2. The novice-name test

Read every new name aloud as a sentence a **junior with no context** would say.
If it needs the documentation to be understood, it is the wrong name.

- The word the ecosystem already uses: `Mount` is chi's word for this;
  inventing `Attach` costs every Go developer a lookup for zero gain.
- Intent, not mechanism: `Public()`, not `SetAuthFlag(false)`.
- No abbreviation that is not already universal in Go (`ctx`, `id` yes;
  `authz`, `mdl`, `hdlr` no).
- No boolean parameters — `Test(args, false, 0, false, false)` says nothing.
  Typed options or separate methods.

**A wrong name is renamed in this change**, never scheduled: the cost is the
number of consumers, which only grows, and a name published in a tag is frozen
there permanently. If a file is touched anyway, the rename rides along. The only
reason to defer is external users who would break — count them; zero means now.

### 3. The complexity ledger

State the change as a balance sheet; a feature that only adds must justify
itself loudly.

```
Concepts the developer must learn   +N / −M
Files they must touch to do X       +N / −M
Lines at the call site              +N / −M
Ways to do the same thing           +N / −M   ← must never end positive
```

The last row is absolute — **one intent, one path**. If the new path is better,
the old one dies in the same change: not "deprecated for one version", not kept
behind a fallback. Be honest when a row gets worse; a design reported as
all-better is an unreviewed design.

### 4. Where does it belong

- **One concern per library (S).** A second concern is a new package or repo,
  not another file. The unit is the **repository**: one that exposes a contract
  *and* ships a concrete implementation has two responsibilities — every
  consumer pulls the contract, none of those who chose another backend pull the
  implementation. Read the `go.mod`, not the file list; a dependency most
  consumers never instantiate is the symptom.
- **New capability = new implementation (O).** Adding the second backend,
  driver or adapter must touch **zero lines** of the consumer. If it touches
  one, the contract was incomplete — fix the contract, do not add a branch.
- **Fix at the root, never at the leaf.** A consumer never re-creates a missing
  symbol locally, never wraps a library to patch its behaviour (a wrapper over a
  defect is a fork with a friendlier name), never works around a defect one
  layer down — each consumer that does encodes the workaround into its own
  contract and passes it on. A defect found in a base piece **blocks this
  plan**: declare the dependency on the base's fix and wait for its published
  tag. Never a "known issue", never a local patch with a comment.
- **A missing contract at a boundary is a defect upstream.** If two libraries
  meet and no type names what crosses between them, the type is missing in the
  library — do not declare a local intersection.
- **No `internal/` folders** — the signature of a duplicated dependency instead
  of a contribution upstream.
- **The glue is written once, in the library that owns it.** Wiring every
  application would repeat belongs to a piece.

### 5. What does this change delete?

Name it. A change that deletes nothing is either genuinely new capability — say
so — or a parallel path being accumulated.

## The harness — the API is the documentation

Long-form rationale: `CONSTRUCTION_HARNESS.md` in `webtyp/app-releases/docs/`.

1. **Typed over `any`.** No `func(...any)` or `interface{}` in a public API;
   `any` only at the I/O edge, never in the data. Reuse the types that exist
   (`fmt.KeyValue`, `model.Resource`) instead of duplicating them.
2. **Explicit over implicit.** Reading the call must be enough to know what it
   does; runtime discovery by type assertion is not — prefer a written line.
3. **Illegal states unrepresentable.** One intent = one path, typed to demand
   what it needs.
4. **Minimal surface.** Export exactly what the author uses.
5. **Fail at compile time.** Compile error → loud development diagnostic →
   never a silent failure.
6. **Self-describing signatures.** Autocomplete must be enough to build.
7. **Closed by default.** Where an API governs access, the zero value denies;
   opening is an explicit, greppable line.
8. **Narrow interfaces (I).** A wide interface forces the implementation that
   needs half of it to stub the rest, and a stub returning `nil` is exactly the
   silent failure rule 5 forbids. If any implementation would stub a method, the
   contract is two contracts: segregate, then compose — `type ReadWriter
   interface { Reader; Writer }` keeps the name consumers already use, so no
   call site changes. Prior art: `io.Reader`/`Writer`/`Closer`, `fs.FS` +
   `fs.ReadDirFS`. Counting test: name, per method, an implementation that would
   stub it — two in one domain means the split is overdue.
9. **The composition root wires, nothing else (D).** Logic depends on contracts;
   one place — the constructor — knows a concrete type. Read the `go.mod`: if
   the library defining the abstraction requires the concrete thing that
   abstraction hides, the inversion is in the prose, not in the build graph.

### Publication

An API is not published until a **consumer-shaped test inside the library
itself** proves it, through the real stack a consumer will use, with real
collaborators and a fake only at the edge. An awkward test means an awkward API,
found before shipping instead of after.

Once a contract has a **second implementation**, it is not published until a
**conformance suite in the repository owning the contract** proves them
substitutable **(L)** and every implementation passes it — as
`storage/conformance` and `ddl/conformance` do for `mem`, `sqlt`, `postgres` and
`indexdb`. Two implementations with no shared suite means substitutability is a
claim, and the divergence will be found by a consumer, in production, on the
second backend. The suite is also what makes "fix at the root" cheap: the fix
lands once and every implementation is re-proved.

## DRY — an enforceable check, not a slogan

- A repeated string (env key, path, prefix, flag name, error message) is a named
  constant in the owning package. String literals are forbidden in logic.
- A value the library computes is **read** by `cmd/`, never re-derived; a check
  it performs is called, never re-implemented at the call site.
- Two functions differing only in a constant are one function and a parameter.
- Documentation repeats a rule only when the reader cannot be assumed to have
  read the other document; otherwise it links. Two copies drift, and the stale
  one is the one someone follows.

## Zero technical debt — what must not survive the change

- A fallback that silently restores the old behaviour when the new one finds
  nothing. A missing input is an error or an empty result, never a guess.
- A `TODO`, a commented-out block, or a stub for an undecided command — an
  undecided thing fails loudly, pointing at the document that will decide it.
- A test that only exercises doubles, with no consumer-shaped case.
- An exported symbol nothing outside the package calls.
- A promise the code cannot keep because a contract is missing downstream (an
  ARIA role whose keyboard behaviour is unimplementable, a documented option
  with no effect). Narrow the claim or fix the contract; never ship the claim.

Before closing, run `grep -rn "TODO\|FIXME\|Deprecated" --include='*.go' .` on
the touched packages and confirm every hit predates the change.

## The acid test

An agent, or a junior, **with no context on the library** produces correct code
guided only by autocomplete and a few-line example. If they must read a manual
to avoid a mistake, something is still untyped. If they must **ask a question**
to proceed — "is it OK to declare this interface locally?" — the harness has
failed by its own definition: the signature did not guide, and the compiler
rejected the correct intent.
