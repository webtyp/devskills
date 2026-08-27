---
name: documentation
description: Documentation standards for ARCHITECTURE.md, PLAN.md, DESIGN.md, SPECS.md, SKILL.md, diagrams, and README indexing. Use when creating or updating project documentation.
---

# Documentation

## README — a newcomer's entry point, not just an index

The audience is someone with **little to no context** on this specific
project — a junior developer, a new hire, or the original author six months
later. A README that only indexes `docs/` files (see "Readme Indexing" below)
answers "where is everything," never "what do I actually do." Both are
required — the index is not optional, but it is not the first thing a reader
needs either.

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
    - **`docs/DESIGN.md`:** (On demand) Justifies technical decisions and explores alternatives. Must NOT duplicate `ARCHITECTURE.md`. Heavily linked by `ARCHITECTURE.md` to keep the main document clean and focused on abstract structure rather than debate.
    - **`docs/SPECS.md`:** (On demand) Strict functional requirements, exact inputs/outputs, and data logic. Must NOT duplicate `ARCHITECTURE.md`. `PLAN.md` consumes it to derive exact test cases and assertions (link direction: plan → specs, never the reverse).
    - **`docs/SKILL.md`:** (On demand) Provides an LLM-friendly, highly condensed summary of the library's context and constraints.
    - **Modular Docs:** If `ARCHITECTURE.md` or `PLAN.md` become too large, they must be divided into domain-specific, uppercase, underscore-separated files (e.g., `docs/BUS_ARCHITECTURE.md`, `docs/CHART_BAR_PLAN.md`).

- **Diagram Standards:**
    - **Format & Location:** Markdown files (`*.md`) containing Mermaid code, stored in `docs/diagrams/` and linked from the architecture documents.
    - **Simplicity:** Use simple, vertical, linear flowcharts (`flowchart TD`). **NEVER** use the `subgraph` directive (ruins TUI rendering). Use `<br/>` for line breaks inside standard nodes instead of quoting text strings.

- **Ephemeral vs Permanent:** Permanent docs (`README.md`, `ARCHITECTURE.md`, `DESIGN.md`, `SPECS.md`, diagrams) must **NEVER link to or cite `PLAN.md`/`CHECK_PLAN.md`** — not even section references like "PLAN §8" — those files are deleted at loop close, so every such reference is a guaranteed dead link. Rationale and rejected alternatives belong in `DESIGN.md`; contracts in `ARCHITECTURE.md`. A permanent doc written documentation-first (before the implementation lands) may carry a self-deleting marker — `STATUS (remove this note when X lands): …` — whose removal is an explicit task of the plan.

- **Readme Indexing:** The `README.md` must act as an index. Every file in `docs/` must be linked from `README.md` — **except the ephemeral lifecycle files** (`PLAN.md`, `PLAN_*.md`, `CHECK_PLAN.md`), which are never indexed. Cross-link logically related permanent documents (e.g., `ARCHITECTURE.md` linking to `SPECS.md` or `DESIGN.md`) to avoid duplicating information across files.
