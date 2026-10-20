---
name: api-design
description: The design gate for any public API — an ecosystem search first, then five answers (prior art, novice-name test, complexity ledger, where it belongs, what it deletes) written in docs/PLAN.md before coding. Use when adding, changing, renaming or reviewing an exported symbol, CLI surface, file convention or declaration format, or starting a library.
---

# API design — the gate

**The typed, explicit code is the harness**: an agent that does not know the library must be correct
from the signatures alone, and the compiler must reject what is wrong. Writing rules: **api-harness**.
What must not survive the change: **zero-debt**.

Applies to a new/changed exported symbol, CLI subcommand or flag, file/directory convention,
declaration format or library — not to unexported code or a fix that changes no signature.
SOLID is enforced where it is checked: **S**, **O** in gate 4; **L**, **I**, **D** in **api-harness**
(publication, rules 8–9). When a rule and a principle disagree, the principle wins.

## Step 0 — search the ecosystem first, before the five answers

The same logic in two repos is a third repo (**zero-debt**), but nobody can obey that without looking. Before designing
anything that adds a library, a module, an exported type or an exported function, **search what already exists** and write the
result in the plan under `## Ecosystem search` (3–6 lines):

1. **List what exists:** `ls ~/Dev/Project/webtyp ~/Dev/Project/veltylabs ~/Dev/Project/veltylabs/modules`, and read the README of every repo whose name is near the idea.
2. **Grep every noun and verb of the idea** over both trees: `grep -rli '<word>' --include=*.go --include=*.md ~/Dev/Project/webtyp ~/Dev/Project/veltylabs`. Search the names you are about to invent (`Fanout`, `Notify`, `Publisher`, `Urgency`, `Currency`, `Key`) and their synonyms (severity/level, queue/outbox, bus/broker/events).
3. **Decide each concept:** **reuse** (import it; the plan names the package), **extend** (what is missing is a plan against *that* repo, never a local copy or a wrapper), or **new** (name the closest existing thing and the concrete reason it does not fit).
4. A plan that creates a library, module or exported name with no `## Ecosystem search` is incomplete: do not dispatch it. Re-run the search when the plan changes scope.

Misses this step exists for: a notification library that re-invented the severity of `msgtype` and the pub/sub of `events`; a money library that planned its own `Fanout` while `events.Fanout` existed.

## The gate — five answers in `docs/PLAN.md`, before coding

A plan that adds or changes public API without them is incomplete: do not dispatch it.
### 1. Prior art

Name **at least three** established frameworks, what each does, and why this ecosystem differs.
"Nobody else does it this way" is valid only followed by the reason it is right here — an API invented
without looking accumulates novel vocabulary for solved problems.

> *Routes:* file-system (Next.js, Remix) · central manifest (Rails, Django, chi `Mount`) · runtime
> discovery (Spring, NestJS). Only the first two make routes knowable **without running the app** — so
> `config/routes.go` + `Mount(prefix, fn)` is family two, chosen for that property.

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
  tag — never a "known issue", never a local patch with a comment (why: **api-harness** → «Lego pieces»).
- **A missing contract at a boundary is a defect upstream.** If two libraries meet and no type names
  what crosses between them, the type is missing there — never a local intersection.
- **No `internal/` folders** — the signature of a duplicated dependency instead of a contribution
  upstream.
- **The glue is written once, in the library that owns it.** Wiring every application would repeat belongs
  to a piece, not to the applications.

### 5. What does this change delete?

Name it. A change that deletes nothing is either genuinely new capability — say so — or a parallel path
being accumulated.
