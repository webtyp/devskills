# devskills
<img src="docs/img/badges.svg">

Agent Skills for the WebTyp ecosystem — source of truth for all `SKILL.md`, synced via `devskills` CLI to `~/.claude`, `~/.gemini`, `~/.codex`, `~/.qwen`, `~/.config/opencode`, and `~/.agents`.

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

Skills are embedded at compile time (`//go:embed skills`) and installed to `~/skills/`, then symlinked per-skill into each detected agent's skills dir (`~/.claude/skills/<name> -> ~/skills/<name>`, and likewise for `~/.gemini`, `~/.codex`, `~/.qwen`, `~/.config/opencode`, `~/.agents`). Linking per-skill — instead of the whole directory — lets an agent's own skills dir also hold vendor-bundled skills (e.g. Codex's `~/.codex/skills/.system/`) without devskills overwriting them.

After editing any `SKILL.md` in `devskills/skills/<name>/`:

```bash
cd devskills && go install ./cmd/devskills && devskills -f
```

## Skills

Browse them here — the directory **is** the index:
**[github.com/webtyp/devskills/skills](https://github.com/webtyp/devskills/tree/main/skills)**

Each folder is one skill; what it covers and when an agent should load it are
the `name:` and `description:` in the frontmatter of its `SKILL.md`. To read
them all at once from a clone:

```bash
grep -h -A1 '^name:' skills/*/SKILL.md
```

Full docs: [docs/DEVSKILLS.md](docs/DEVSKILLS.md) ·
Writing a new skill (Anthropic's guide): [docs/SKILL_INSTRUC.md](docs/SKILL_INSTRUC.md)

## License

MIT
