---
name: documentation
description: Documentation standards for ARCHITECTURE.md, PLAN.md, DESIGN.md, SPECS.md, SKILL.md, diagrams, and README indexing. Use when creating or updating project documentation.
---

# Documentation

## Who every document is written for

**This applies to every file in `docs/` and to `README.md` — not only to the
README.**

The reader is someone with **little to no context**: a junior developer, a new
hire, or the original author six months later. They did not attend the
conversation that produced the document. They do not know the vocabulary of this
project yet. Assume nothing beyond general programming knowledge.

The goal for WebTyp is that **any person can understand the framework**, so
documentation is written closer to a tutorial than to a set of meeting minutes.
A document that is only legible to whoever was in the room has failed, however
correct its content is.

### Every document and every section opens with context, before its content

Before the first claim, table, decision, or API listing, the reader must be told,
in plain sentences:

1. **What this is about** — name the thing being discussed and say what it *is*,
   even when the name looks self-explanatory to you. "`truststore` is a Go
   library that installs a certificate into the operating system's list of
   trusted authorities."
2. **When the reader will run into it** — the concrete situation that makes this
   page relevant. "You ran `webtyp dev`, the browser opened on `https://` and did
   not warn you about the certificate. This is the piece that made that happen."
3. **Why the document exists** — what question it settles, or what would go wrong
   without it.

Only then the detail. A section that opens with **"Question."** or with a
comparison table, and leaves the reader to infer the subject from the title, is
exactly the failure this rule exists to close.

### Introduce the vocabulary you use

The first time a document uses a term the reader may not own — trust store, CA,
leaf certificate, SAN, SPKI, RBAC, adapter, strategy — define it in one short
clause on the spot, or link to the place that does. Never let the explanation of
a decision depend on jargon the reader was never given.

### Prefer the concrete over the abstract

Show the command the reader would type, the error they would see, the file that
would appear. An abstract statement of a rule plus one concrete example beats two
paragraphs of abstraction.

## README — a newcomer's entry point, not just an index

A README that only indexes `docs/` files (see "Readme Indexing" below) answers
"where is everything," never "what do I actually do." Both are required — the
index is not optional, but it is not the first thing a reader needs either.

### Getting Started, staged by what the reader wants to do

Every `README.md` MUST open (above the doc index) with a "Getting Started"
section broken into the concrete goals a reader arrives with, most common
first — typically:

1. **See it running locally.** The shortest real path to a working UI/output
   on the reader's own machine. Exact commands, in order, copy-pasteable.
   State plainly what does NOT work yet at this stage (e.g., "login is
   disabled here," "uses fake data") instead of leaving the reader to
   discover a silent limitation.
2. **Deploy it.** Link to `docs/DEPLOY.md` (or equivalent) — do not inline
   deploy steps here if they already live there.
3. **Develop a feature / contribute code.** Link to `docs/ARCHITECTURE.md`/
   `AGENTS.md` — do not re-explain architecture here.

A reader must be able to tell which section matches what they want to do
**right now** from the headers alone, without reading the other two.

### Every required configuration value gets a plain-language line, not just a name

A `.env` key, CLI flag, or secret listed as required MUST document three
things, not just its name:

- **What it actually is**, in one plain sentence — not a restatement of the
  name. `IAM_CLIENT_SECRET` is not self-explanatory; "the password this app
  uses to prove its own identity to the shared login service" is.
- **Where it comes from** — generated how, by whom, or fetched from where.
- **What breaks if it is missing** — the exact error the reader will see, so
  a stuck reader can match their symptom back to the cause.

Never list a bare variable name and assume the reader already knows the
system well enough to infer its purpose. That assumption is exactly the
failure mode this rule exists to close.

### A complicated setup is a design defect, not a documentation gap to paper over

Before writing an elaborate "configure these 5 prerequisites first" section,
ask whether the software itself could remove the need for some of them. A
local dev/preview path that requires production credentials just to render a
UI is usually the thing to fix, not the thing to document more thoroughly.
Documentation describes the simplest real path — it does not make a needlessly
complicated one sound simpler than it is.

## Everything else

- **Documentation First:** You MUST update the documentation *before* coding or running `gopush`.
- **Context-Aware Rule Compilation:** Every `PLAN.md` MUST start with a "Development Rules" section that copies/pastes the relevant constraints from this skill file (e.g., WASM restrictions, DI rules).

- **Standard Documents:**
    - **`docs/ARCHITECTURE.md`:** Defines WHAT & WHY (abstract design, constraints). NO implementation code.
    - **`docs/PLAN.md`:** Defines HOW (steps, reference code, test strategy). It is the master orchestrator for execution. **Ephemeral**: `codejob` renames it to `CHECK_PLAN.md` and deletes it when the loop closes (see skill agents-workflow).
    - **`docs/DESIGN.md`:** (On demand) Justifies technical decisions and explores alternatives. Must NOT duplicate `ARCHITECTURE.md`. Heavily linked by `ARCHITECTURE.md` to keep the main document clean and focused on abstract structure rather than debate. Every decision entry MUST open with the background a reader needs to care about it — what problem the code was solving, what the candidates even are — before naming the choice. A decision entry that starts at the comparison has skipped the only part a newcomer needed.
    - **`docs/SPECS.md`:** (On demand) Strict functional requirements, exact inputs/outputs, and data logic. Must NOT duplicate `ARCHITECTURE.md`. `PLAN.md` consumes it to derive exact test cases and assertions (link direction: plan → specs, never the reverse).
    - **`docs/SKILL.md`:** (On demand) Provides an LLM-friendly, highly condensed summary of the library's context and constraints.
    - **Modular Docs:** If `ARCHITECTURE.md` or `PLAN.md` become too large, they must be divided into domain-specific, uppercase, underscore-separated files (e.g., `docs/BUS_ARCHITECTURE.md`, `docs/CHART_BAR_PLAN.md`).

- **API changes carry a `## Design gate` section:** when a `PLAN.md` adds or
  changes any public API, that section holds the five answers of skill
  **api-design** (prior art, novice-name test, complexity ledger, where it
  belongs, what it deletes). Do not restate those rules in the document — link
  to the skill; two copies of a rule drift, and the stale one is the one someone
  follows.

- **Diagram Standards:**
    - **Format & Location:** Markdown files (`*.md`) containing Mermaid code, stored in `docs/diagrams/` and linked from the architecture documents.
    - **Simplicity:** Use simple, vertical, linear flowcharts (`flowchart TD`). **NEVER** use the `subgraph` directive (ruins TUI rendering). Use `<br/>` for line breaks inside standard nodes instead of quoting text strings.

- **Ephemeral vs Permanent:** Permanent docs (`README.md`, `ARCHITECTURE.md`, `DESIGN.md`, `SPECS.md`, diagrams) must **NEVER link to or cite `PLAN.md`/`CHECK_PLAN.md`** — not even section references like "PLAN §8" — those files are deleted at loop close, so every such reference is a guaranteed dead link. Rationale and rejected alternatives belong in `DESIGN.md`; contracts in `ARCHITECTURE.md`. A permanent doc written documentation-first (before the implementation lands) may carry a self-deleting marker — `STATUS (remove this note when X lands): …` — whose removal is an explicit task of the plan.

- **Readme Indexing:** The `README.md` must act as an index. Every file in `docs/` must be linked from `README.md` — **except the ephemeral lifecycle files** (`PLAN.md`, `PLAN_*.md`, `CHECK_PLAN.md`), which are never indexed. Cross-link logically related permanent documents (e.g., `ARCHITECTURE.md` linking to `SPECS.md` or `DESIGN.md`) to avoid duplicating information across files.
