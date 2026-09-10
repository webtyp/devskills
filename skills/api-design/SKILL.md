---
name: api-design
description: The gate every new or changed public API must pass BEFORE it is written — prior-art comparison, the novice-name test, the complexity ledger, harness/DRY/SRP, and the zero-technical-debt rules. Use when designing, adding, renaming, or reviewing any exported symbol, CLI surface, file convention, or declaration format — and when starting a library from scratch.
---

# API Design — the gate

The framework's promise is that an application built on it needs **fewer**
concepts than the alternatives, not more. Every exported symbol either pays for
that promise or erodes it. This skill is the check that must pass before the
first line is written.

**It applies to:** a new exported symbol, a changed signature, a new CLI
subcommand or flag, a new file/directory convention, a new declaration format,
a new library. It does **not** apply to unexported code or to a bug fix that
changes no signature.

---

## The gate — answer all five in writing, in the `docs/PLAN.md`, before coding

A plan that adds or changes public API without these answers is incomplete and
must not be dispatched.

### 1. Prior art — how do established frameworks solve this?

Name **at least three**, say what each does, and state explicitly why this
ecosystem differs. "Nobody else does it this way" is a valid answer only when
followed by the reason it is right here.

You are not obliged to copy them. You are obliged to know them. An API invented
without looking is how a framework accumulates novel vocabulary for solved
problems — the exact thing that makes a developer's first week expensive.

> **Worked example — route declaration.** The landscape has three families:
> **file-system routing** (Next.js, Nuxt, SvelteKit, Remix, Astro — the path *is*
> the file path); **central manifest** (Rails `config/routes.rb`, Django
> `urls.py` + `include()`, Laravel `routes/`, Phoenix `router.ex`, chi
> `r.Mount`); **decentralised discovery** (Spring `@GetMapping`, NestJS
> decorators, ASP.NET attributes — declared beside the handler, found by
> scanning at runtime).
>
> The first two make routes knowable **without running the app**; the third does
> not, and is why Spring needs a dedicated endpoint just to answer "what routes
> exist?". WebTyp's `MountAPI` + type assertion is family three. Moving to
> `routes/routes.go` + `Mount(prefix, fn)` is family two — the same shape as
> Django's `include()`, Laravel's route groups and chi's `Mount`. That is the
> justification: not novelty, but the family whose properties we need.

### 2. The novice-name test

Read every new name aloud as a sentence a **junior with no context** would say.
If the name needs the documentation to be understood, it is the wrong name.

- Prefer the word the surrounding ecosystem already uses. `Mount` is chi's word
  for exactly this; inventing `Attach` would cost every Go developer a lookup
  for zero gain.
- The name states **intent**, not mechanism: `Public()` not `SetAuthFlag(false)`.
- No abbreviation that is not already universal in Go (`ctx`, `fmt`, `id` yes;
  `authz`, `mdl`, `hdlr` no).
- A boolean parameter at a call site is unreadable: `Test(args, false, 0, false,
  false)` says nothing. Prefer typed options or separate methods.

### 3. The complexity ledger

State the change as a balance sheet. This is the complexity tax — a feature that
only adds must justify itself loudly.

```
Concepts the developer must learn   +N  / −M
Files they must touch to do X       +N  / −M
Lines at the call site              +N  / −M
Ways to do the same thing           +N  / −M    ← must never end positive
```

The last row is absolute: **the change must not leave two ways to do one thing.**
If the new path is better, the old one is deleted in the same change, not
"deprecated for now".

Be honest when a row gets worse. A design with one row worse and four better is
a good design; a design reported as all-better is an unreviewed design.

### 4. Where does it belong — SRP and lego pieces

- **One concern per library.** If the new thing is a second concern, it is a new
  package or a new repo, not another file in the existing one.
- **A consumer never re-creates a missing symbol locally.** If a library does not
  expose what you need, stop and fix it upstream. Recreating it downstream forks
  that library's responsibility and the copy can never be reused.
- **A missing contract at a boundary is a defect in the library, not in the
  consumer.** If two libraries meet and no type names what crosses between them,
  the type is missing upstream — do not declare a local intersection.
- **Never wrap a library to fix its behaviour.** A wrapper that patches a defect
  is a fork with a friendlier name.
- **No `internal/` folders.** They are the signature of a duplicated dependency
  instead of a contribution upstream.
- **The glue is written once, in the library that owns it.** If every application
  would write the same wiring, that wiring belongs to a piece.

### 4b. If the name is wrong, it is renamed **now**

A rename never gets cheaper. Its cost is the number of consumers, which only
grows, and a name published in a tag is frozen there permanently. "We will
rename it when things settle" means "we will pay more later for the same work,
and live with the bad name until then".

So: a bad name found while touching a package is fixed **in the same change**,
not scheduled. If a consumer must be touched anyway for another reason, the
rename rides along — one migration over those files instead of two.

The only valid reason to defer is external users who would break. Count them.
If the answer is zero, there is no argument for waiting.

### 5. What does this change delete?

Name it. A change that deletes nothing is either genuinely new capability — say
so — or it is accumulating a parallel path.

---

## DRY — as an enforceable check, not a slogan

- A repeated string (env key, path, prefix, flag name, error message) is a named
  constant in the owning package. String literals are forbidden in logic.
- A value the library already computes is **read** by `cmd/`, never re-derived.
- A check the library already performs is called, never re-implemented at the
  call site.
- Two functions differing only in a constant are one function and a parameter.
- Documentation repeats a rule only when the reader cannot be assumed to have
  read the other document. Otherwise it links — two copies of a rule drift, and
  the stale one is the one someone follows.

---

## The harness — the API is the documentation

Long-form rationale: `CONSTRUCTION_HARNESS.md` in `webtyp/app-releases/docs/`.
The normative rules:

1. **Typed over `any`.** No `func(...any)` or `interface{}` in a public API —
   methods typed by intent. `any` only at the I/O edge, never in the data. Reuse
   types that already exist (`fmt.KeyValue`, `model.Resource`) instead of
   duplicating them.
2. **Explicit over implicit.** Reading the call must be enough to know what it
   does. Runtime discovery by type assertion is implicit — prefer a written line.
3. **Illegal states unrepresentable.** One intent = one path, typed to demand
   what it needs.
4. **One way to do each thing.**
5. **Minimal surface.** Export exactly what the author uses.
6. **Fail at compile time, not at runtime.** Compile error → loud development
   diagnostic → never a silent failure.
7. **Self-describing signatures.** Autocomplete must be enough to build. If using
   the API requires reading a long document, the API is incomplete.
8. **Closed by default.** Where an API governs access, the zero value denies.
   Opening is an explicit, greppable line.
9. **Lego pieces, never forks.** One concern per library, exposed as a typed
   contract. Consumers assemble; they do not re-implement, wrap, or copy.

**The publication rule:** an API is not published until a **consumer-shaped
test, inside the library itself**, proves it — through the real stack a consumer
will use, with real collaborators and a fake only at the edge. If that test is
awkward to write, the API is awkward to use, and the defect was found before
shipping instead of after.

---

## Zero technical debt — what must not survive the change

The framework carries no debt forward. When the change lands, none of these may
exist:

- A deprecated path kept "for one version". Either the change is safe to make
  now, or it is not made.
- A fallback that silently restores the old behaviour when the new one finds
  nothing. A missing input is an error or an empty result — never a guess.
- A `TODO`, a commented-out block, or a stub for an undecided command. An
  undecided thing fails loudly pointing at the document that will decide it.
- Two documents stating the same rule.
- A test that only exercises doubles, with no consumer-shaped case.
- An exported symbol nothing outside the package calls.
- A workaround for a defect that belongs to a library one layer down. **Fix at the root, never
  at the leaf.** A defect in a base piece compounds through every piece built on it: each new
  consumer inherits it, works around it, or encodes the workaround into its own contract and
  passes it on. So a defect found in a base piece **blocks this plan** — the plan declares the
  dependency on the base's fix and waits for its published tag, exactly like every gate/phase
  master plan in this repo. It never becomes a "known issue" or a local patch with a comment.
- A promise the code cannot keep because a contract is missing downstream (an ARIA role whose
  keyboard behaviour is unimplementable, a documented option with no effect). Narrow the claim
  or fix the contract — never ship the claim.

Run before closing: `grep -rn "TODO\|FIXME\|Deprecated" --include='*.go' .` on
the touched packages, and confirm every hit predates the change.

---

## The acid test

An agent, or a junior, **with no context on the library** produces correct code
guided only by autocomplete and a few-line example.

If they must read a manual to avoid a mistake, something is still untyped. And
if they must **ask a question** to proceed — "is it OK to declare this interface
locally?" — the harness has already failed by its own definition: the signature
did not guide, and the compiler rejected the correct intent.
