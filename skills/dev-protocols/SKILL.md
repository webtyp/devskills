---
name: dev-protocols
description: Development protocols including language rules (English docs in webtyp, Spanish docs and comments in veltylabs, mirrored chat), strategic justification, explicit execution gates, and harness plan-mode mirroring to docs/PLAN.md. Use as general workflow context.
---

# Development Protocols

> **Any change that adds or alters a public symbol, CLI surface, file convention
> or declaration format** must pass skill **api-design** BEFORE code, with its
> five answers written into the `docs/PLAN.md`. That gate also carries the
> zero-technical-debt rules and "fix at the root, never at the leaf" — a defect
> found one library down blocks the plan instead of being patched here.

- **Language Protocol:** by tree. `webtyp/*`: documentation in **English**. `veltylabs/*` (modules, apps, adapters): **everything in Spanish** — `README.md`, `docs/*`, diagrams, `AGENTS.md` `## This repo`, and **code comments** — except code itself (identifiers, `lang.json` keys, log/error keys stay English, skill **translations**). `docs/PLAN.md` for an executor stays English everywhere. Chat mirrors the user's language; master plans go in the user's language (skill **master-plan**).
- **Plain words to the developer:** every explanation, question or report to the user is written for a junior with no context — you never know who reads it. No insider jargon or invented labels («comodín», «costura», «gate B»): name the real thing (the file, the screen, what the user sees) and say what happens in one concrete sentence. A design question shows a **short example**: the code the developer writes today and the code they would write with each option, so the choice is visible, not argued. The goal of every design choice is an API that is easy to use, with one way to do each thing and no duplicated code — say how each option moves toward or away from that.
- **Strategic Justification:** Always provide the rationale for your chosen solution based on best practices, offering alternatives and explaining why the selection fits the context.
- **Verify before asserting a rule:** never state a rule about a set of files, packages or repos ("only `svg.go` is `!wasm`", "nobody imports X") from one example or from memory. Enumerate it first (`grep -l '^//go:build' */ui/*.go`, `go list -deps`) and quote what you counted. A generalisation the user acts on without the code behind it breaks every repo it is applied to; if you have not checked, say so instead of asserting.
- **Minimal chat (never duplicate state):** never duplicate or transcribe repo/module statuses, progress tables, or review diaries into the chat. The master plan (`webtyp/docs/<TOPIC>_MASTER_PLAN.md`) or repo `docs/PLAN.md` is the **single source of truth** for all states. Chat messages must be minimal (actions executed, files touched, or blocking questions only) so the chat can be cleared/compacted at any moment and work resumed solely by reading the master plan.
- **Explicit Execution:** Never start writing/modifying actual codebase source code unless explicitly told to "execute the plan", "ok", or "ejecuta".

- **Harness Plan Mode (any agent with a built-in plan mode — e.g. Claude Code, Gemini CLI):**
    - **Plan Mirroring:** When the harness's plan mode creates a plan file at a system path (e.g. `~/.claude/plans/`), you MUST ALSO create or update `docs/PLAN.md` inside the active project's repository with the same final plan content. The system plan file is harness infrastructure; `docs/PLAN.md` is the canonical source of truth for execution agents.
    - **Plan Location Priority:** Always prefer writing plans inside the project repository. If the plan mode forces a system path, treat it as a draft and mirror the final content to the project before requesting plan approval.
- **Write directly to disk in IDE harnesses (Antigravity / VS Code):** IDE file-edit tools (`edit_file`, `write_file`) open an unsaved buffer in VS Code; if the user has the file open or closes it without saving, the on-disk file reverts or stays stale. Write and edit files directly on disk via shell commands (`python3`, `cat <<'EOF' > ...`, `mv`) and verify on disk before running `codejob`, `gotest`, or `gopush`.
