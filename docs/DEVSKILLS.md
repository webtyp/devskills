# devskills - LLM Skills Sync

Synchronizes WebTyp Agent Skills to all installed LLM configuration directories.

## Installation

```bash
go install webtyp.com/devskills/cmd/devskills@latest
```

## Usage

### Sync all installed LLMs and update skills
```bash
devskills
```

### Sync specific LLM
```bash
devskills -l claude
devskills --llm gemini
```

### Force refresh (reinstalls all skills)
```bash
devskills -f
devskills --force
```

## How it works

1. **Installs Modular Skills**: Copies embedded Agent Skills to `~/skills/`.
2. **Detects installed LLMs**: Checks for `~/.claude/`, `~/.gemini/`, `~/.codex/`,
   `~/.qwen/`, `~/.config/opencode/`, and `~/.agents/` directories.
3. **Creates Symlinks**: Symlinks each skill individually from `~/skills/<name>` into
   `~/.claude/skills/<name>` (and the same for the other detected dirs). Linking
   per-skill, rather than the whole directory, lets an agent's own skills dir also
   hold vendor-bundled skills (e.g. Codex's `~/.codex/skills/.system/`) without
   devskills overwriting them. This allows LLMs to natively discover and use all
   domain-specific skills without any text-based configuration.

### Why `~/.agents` and opencode

- `~/.agents/skills/` is a vendor-neutral convention that opencode (and other
  agents) auto-load from, independent of any single tool's own config dir — it
  covers future agents without new code.
- `~/.config/opencode/skills/` is opencode's own dedicated skills dir. opencode
  already auto-discovers `~/.claude/skills/` and `~/.agents/skills/`, so this is
  mostly a hedge in case that auto-discovery is disabled (`disableClaudeCodeSkills`
  / `disableExternalSkills` in opencode's config).
- Cursor, GitHub Copilot, and Windsurf are intentionally **not** targeted: they
  don't discover `SKILL.md` files at all — they read their own formats
  (`.cursor/rules/*.mdc`, `.github/copilot-instructions.md`, `CONVENTIONS.md`).
  Supporting them would mean writing a format converter, not adding a sync path.

## Supported Skills

The skills are divided by domain to minimize context usage:

- `core-principles`: SRP, DI, Framework-less development.
- `testing`: gotest, gopush, WASM dual testing.
- `documentation`: doc standards, diagrams, readme indexing.
- `wasm`: webtyp MCP, frontend Go compatibility.
- `agents-workflow`: PLAN.md orchestrator, stage-driven execution.
- `dev-protocols`: Language rules, justification, Claude plan mirroring.
- `devskills`: Self - how to sync skills after editing SKILL.md
- `components`: WebTyp component standards.
- `form-codegen`: Model/form generation.
- `plan-authoring`: PLAN.md authoring.
- `webtyp-app`: MCP daemon.

## Library Usage

```go
import "webtyp.com/devskills"

llm := devskills.NewLLM()

// Sync all and install skills
summary, err := llm.Sync("", false)
if err != nil {
    log.Fatal(err)
}
fmt.Println(summary)
```

## Troubleshooting

### "No LLMs detected"
Skills are installed in `~/skills/` even if no specific LLM config is found.

### Skills not appearing in LLM
Ensure the LLM directory exists (`~/.claude`, `~/.gemini`, `~/.codex`, `~/.qwen`,
`~/.config/opencode`, or `~/.agents`). If symlinks are broken, run `devskills -f`
to force a refresh.

## See Also

- [gotest](GOTEST.md) - Run tests with coverage and badges
- [gopush](GOPUSH.md) - Complete workflow: test + push + update

## Migration from devflow/llmskill

Previously this tool lived as `llmskill` in `webtyp.com/devflow`:

```bash
# Old
go install webtyp.com/devflow/cmd/llmskill@latest && llmskill

# New
go install webtyp.com/devskills/cmd/devskills@latest && devskills
```
