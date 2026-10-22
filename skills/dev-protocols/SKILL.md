---
name: dev-protocols
description: Development protocols including language rules (English docs, mirrored chat), strategic justification, explicit execution gates, and harness plan-mode mirroring to docs/PLAN.md. Use as general workflow context.
---

# Development Protocols

> **Any change that adds or alters a public symbol, CLI surface, file convention
> or declaration format** must pass skill **api-design** BEFORE code, with its
> five answers written into the `docs/PLAN.md`. That gate also carries the
> zero-technical-debt rules and "fix at the root, never at the leaf" — a defect
> found one library down blocks the plan instead of being patched here.

- **Language Protocol:** Plans and documentation must be generated in **English**. Chat conversations with the user must match the language used by the user (mirroring). **Exception:** master plans go in the user's language (skill **master-plan**).
- **Strategic Justification:** Always provide the rationale for your chosen solution based on best practices, offering alternatives and explaining why the selection fits the context.
- **Verify before asserting a rule:** never state a rule about a set of files, packages or repos ("only `svg.go` is `!wasm`", "nobody imports X") from one example or from memory. Enumerate it first (`grep -l '^//go:build' */ui/*.go`, `go list -deps`) and quote what you counted. A generalisation the user acts on without the code behind it breaks every repo it is applied to; if you have not checked, say so instead of asserting.
- **Explicit Execution:** Never start writing/modifying actual codebase source code unless explicitly told to "execute the plan", "ok", or "ejecuta".

- **Harness Plan Mode (any agent with a built-in plan mode — e.g. Claude Code, Gemini CLI):**
    - **Plan Mirroring:** When the harness's plan mode creates a plan file at a system path (e.g. `~/.claude/plans/`), you MUST ALSO create or update `docs/PLAN.md` inside the active project's repository with the same final plan content. The system plan file is harness infrastructure; `docs/PLAN.md` is the canonical source of truth for execution agents.
    - **Plan Location Priority:** Always prefer writing plans inside the project repository. If the plan mode forces a system path, treat it as a draft and mirror the final content to the project before requesting plan approval.
