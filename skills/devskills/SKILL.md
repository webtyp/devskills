---
name: devskills
description: Sync skills from tinywasm/devskills/skills to all installed LLM agents (Claude, Gemini). Run after creating or modifying any SKILL.md file in devskills/skills/.
---

# devskills

`devskills` is a CLI tool in `tinywasm/devskills` that synchronizes skill files from the devskills repository to all installed LLM agent configurations.

## When to Run

After creating or modifying any `SKILL.md` in `tinywasm/devskills/skills/<name>/`:

```bash
cd tinywasm/devskills && go install ./cmd/devskills && devskills -f
```

**The rebuild step is mandatory**: skills are EMBEDDED in the binary at
compile time (`//go:embed skills` in `devskills.go`) — `devskills` does NOT
read the working directory. Running the old binary reinstalls the OLD
embedded skills and still prints "Skills updated". Always verify the change
landed:

```bash
grep -c "<some new phrase>" ~/.claude/skills/<name>/SKILL.md
```

## How It Works

1. Skills live in `tinywasm/devskills/skills/` (source of truth) and are
   embedded into the `devskills` binary when it is built.
2. `devskills` installs the embedded skills to `~/skills/`.
3. It then symlinks `~/skills/` from each detected LLM config dir
   (`~/.claude/skills/`, `~/.gemini/skills/`); if the target already exists
   as a real directory, it falls back to copying into it.

## Installation

`devskills` is part of `github.com/tinywasm/devskills`.

**Step 1 — Install the binary:**

```bash
go install github.com/tinywasm/devskills/cmd/devskills@latest
```

Or install all devskills binaries at once (requires the repo to be cloned first):

```bash
go install github.com/tinywasm/devskills/cmd/devskills@latest && devskills
```

**Step 2 — Clone the devskills repo** (required to EDIT skills — the binary
carries an embedded copy from build time):

```bash
git clone https://github.com/tinywasm/devskills
cd devskills
devskills
```

To pick up local skill edits: `git pull` (or edit) + `go install ./cmd/devskills` + `devskills -f`.

## Usage

```bash
# Sync all installed LLMs
devskills

# Sync only Claude
devskills -l claude

# Force overwrite with backup
devskills -f
```

## Skill File Format

Each skill is a folder with a single `SKILL.md`:

```
tinywasm/devskills/skills/
└── myskill/
    └── SKILL.md
```

`SKILL.md` frontmatter:
```markdown
---
name: myskill
description: One-line description used by the agent to decide when to apply this skill.
---

# Skill content here
```

## After Creating a New Skill

1. Write `SKILL.md` in `tinywasm/devskills/skills/<skillname>/`
2. Run `devskills` from the shell
3. The skill is immediately active in all installed agents

## Important

- `devskills` is a **local developer tool** — never include it in `PLAN.md` files sent to external agents.
- Skills are read-only context for agents — agents do not run `devskills` themselves.
