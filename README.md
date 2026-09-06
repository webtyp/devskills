# devskills

Agent Skills for the WebTyp ecosystem — source of truth for all `SKILL.md`, synced via `devskills` CLI to `~/.claude`, `~/.gemini`, etc.

## Installation

```bash
go install webtyp.com/devskills/cmd/devskills@latest
```

## Usage

```bash
# Sync all installed LLMs and update skills
devskills

# Sync only Claude
devskills -l claude

# Force refresh
devskills -f
```

Skills are embedded at compile time (`//go:embed skills`) and installed to `~/skills/`, then symlinked from each LLM dir (`~/.claude/skills -> ~/skills`).

After editing any `SKILL.md` in `devskills/skills/<name>/`:

```bash
cd devskills && go install ./cmd/devskills && devskills -f
```

## Skills

- `core-principles` — SRP, DI, framework-less web
- `testing` — gotest, gopush, WASM dual testing
- `documentation` — ARCHITECTURE.md, PLAN.md, etc.
- `wasm` — WASM env rules
- `agents-workflow` — PLAN.md orchestrator
- `dev-protocols` — dev protocols
- `components` — component standards
- `form-codegen` — model/form generation
- `plan-authoring` — PLAN.md authoring
- `webtyp-app` — MCP daemon
- `devskills` — self (how to sync)

See `docs/DEVSKILLS.md` for full docs.
See `docs/SKILL_INSTRUCTIONS.md` for writing new skills by anthropics guide.

## License

MIT
