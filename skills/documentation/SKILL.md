---
name: documentation
description: Documentation standards for README.md, ARCHITECTURE.md, PLAN.md, DESIGN.md, SPECS.md, SKILL.md, and diagrams. Grounded in Minimalist Instruction, Cognitive Load Theory, and Diátaxis. Use when creating, updating, or reviewing project documentation.
---

# Documentation

## Who every document is written for (Pedagogical Basis)

Every file in `docs/` and `README.md` is written for a **junior developer with zero prior context** who did not attend the conversation that produced the repo. We combine three proven frameworks:

1. **Minimalist Instruction (John M. Carroll):** Open with a concrete real-world problem before any abstraction. State prerequisites before Step 1. Pair every action with its visible output.
2. **Cognitive Load Theory & Progressive Disclosure (Sweller / Nielsen):** Keep `README.md` short (**≤ 150 lines**, never > 200). An oversized README overwhelms the reader. Never duplicate error lists, env var catalogs, or glossary tables in the README — **link concepts inline on first mention** to `docs/*.md`, the Go source file (`layout.go`, `install.go`), or the external standard.
3. **Diátaxis (Daniele Procida):** Follow cognitive order — *Why it exists* → *Before you start* → *Step-by-step tutorial* → *Documentation map*.

Language follows **dev-protocols**: English in `webtyp/*`, Spanish in `veltylabs/*`. In `webtyp/*`, never use client names or domains; use RFC 2606 (`example.com`, `acme`). Update documentation *before* coding or running `gopush`.

## Mandatory `README.md` Template (4 Sections in Exact Order, ≤ 150 Lines)

Every repository `README.md` MUST stay concise (≤ 150 lines) and contain these 4 sections in order (translated to Spanish in `veltylabs/*`):

### 1. `## Why This Exists` (`## ¿Por qué existe?`)
- Plain-language paragraph naming what the package/app is and when the reader runs into it.
- A concise **Without this vs. With this** (*Sin esto vs. Con esto*) comparison showing the concrete failure mode this repository solves.
- **Inline links on first mention:** The first time any domain concept, cryptographic format, sibling repo (`https://github.com/webtyp/...`), or internal flow appears in prose, link it right there to its section in `docs/ARCHITECTURE.md`, `docs/diagrams/*.md`, or external spec. Do NOT create a separate glossary table in the README.

### 2. `## Before You Start` (`## Antes de comenzar`)
- Short checklist of what must exist *before* Step 1, split by environment (e.g., publisher machine vs. target server), with inline links to the tools or source files that define those requirements.

### 3. `## Step-by-Step Guide` (`## Guía paso a paso`)
- Numbered steps (`### Step 1: ...`, `### Step 2: ...`) walking from zero to a working result.
- Every step MUST include three parts:
  1. **Goal:** What you are doing and why, in one plain sentence.
  2. **Command or runnable code:** Copy-pasteable shell commands or a **complete, compilable Go file** (`package`, `import`, declared variables — fragments with undeclared variables are forbidden). Always import `"webtyp.com/fmt"` in Go examples (never stdlib `"fmt"` or `"log"`), and exit `main()` on error with `fmt.Println(err); return` (or `os.Exit(1)` in host CLIs).
  3. **Expected result:** What appears on screen, on disk, or over HTTP.
- **Link to code instead of listing constants:** Never copy tables of environment variables, secret files, or exported `Err*` constants into `README.md` — link directly to the Go source files (`[layout.go](layout.go)`, `[install.go](install.go)`, `[apply.go](apply.go)`) and `[docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)` where they live.

### 4. `## Documentation Map` (`## Mapa de documentación`)
- An ordered reading table (`Order | Document | What question it answers & when to read it`) indexing every permanent file in `docs/` and `docs/diagrams/`. Bare comma-separated link dumps (`See [A], [B], [C]`) are forbidden.

## Mandatory `docs/ARCHITECTURE.md` Template (4 Sections in Exact Order)

Defines WHAT & WHY (no implementation code blocks):

1. **`## 1. Overview in One Minute` (`## 1. Visión general en un minuto`)**: Plain-language mental model of the pieces and who decides what, before any table.
2. **`## 2. Repository Pieces & On-Disk Layout` (`## 2. Piezas y estructura`)**: Tables mapping each source file, service, or on-disk file to its single responsibility and owner.
3. **`## 3. Step-by-Step Flows` (`## 3. Flujos paso a paso`)**: One subsection per core flow, narrating the happy path and failure/rollback path in plain words and linking to its diagram in `docs/diagrams/`.
4. **`## 4. Security Boundaries & Contracts` (`## 4. Contratos y límites de seguridad`)**: Why each invariant exists (least-privilege UIDs, digest pinning, network isolation) and what is forbidden.

## Mandatory `docs/diagrams/*.md` Template (3 Parts in Exact Order)

1. **Context header (above the diagram):** 2–3 plain sentences stating what scenario the diagram shows, what triggers it, and which test file verifies its branches.
2. **Mermaid diagram:** Simple vertical `flowchart TD`. **NEVER** use `subgraph` (breaks TUI rendering). Use `<br/>` for line breaks inside nodes.
   - STATUS (remove this line when DB_DIAGRAM_MASTER_PLAN closes): the database diagram becomes `docs/diagrams/database.svg`, generated from the models on save (never edited by hand); until `webtyp/erd` ships, keep the existing `database.md`.
3. **Step-by-step walkthrough (below the diagram):** Numbered list explaining every step and error branch in plain words.

## Other Standard Documents & Rules

- **`docs/PLAN.md`:** Ephemeral execution plan (`codejob close` deletes it). Permanent docs (`README.md`, `ARCHITECTURE.md`, `DESIGN.md`, `SPECS.md`, diagrams) must **NEVER** link to or cite `PLAN.md`/`CHECK_PLAN.md`. A doc written before implementation lands may carry `STATUS (remove this note when X lands): …`.
- **`docs/DESIGN.md` / `docs/SPECS.md` / `docs/SKILL.md`:** Created on demand. `DESIGN.md` opens every decision with background before comparing candidates; `SPECS.md` holds strict I/O contracts; `SKILL.md` summarizes library usage for LLMs. Split large docs into `docs/<DOMAIN>_ARCHITECTURE.md`.
- **Code comments:** Every exported symbol comment explains its purpose in plain words and includes a short usage example (`// Example: …` or `// Ejemplo: …`).

## Self-documenting code — nothing taken as obvious

webtyp is self-documenting: a developer or agent who opens any file must reach the why and the architecture by
following links, never by guessing or searching. Every generated file opens with a header that says **what made
it, from what, why it exists, and where to read more** — a bare `Code generated by https://webtyp.com` is not enough:

```go
// Code generated by webtyp.com/server (maingen.go) from config/app.go. DO NOT EDIT.
// Why: the server binary of this app; dev tools come from webtyp.com/devtools only with -tags dev.
// Architecture: https://github.com/webtyp/appkit/blob/main/docs/ARCHITECTURE.md
// To change what is generated: https://github.com/webtyp/server/blob/main/maingen.go
```

The first line keeps Go's `Code generated … DO NOT EDIT.` form (tools detect it). The same rule holds for every
exported symbol's doc comment and every `docs/` page: link to the file or repo that owns the next answer.
