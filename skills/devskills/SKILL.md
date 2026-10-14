---
name: devskills
description: How skills are written (one domain, link never repeat, ≤50 lines good, >150 refused) and synced from webtyp/devskills/skills to every installed agent (Claude, Gemini, Antigravity, Codex, opencode, Qwen). Use when creating, editing, splitting or auditing any SKILL.md, and run the sync after every change.
---

# devskills — writing and syncing skills

## A skill is a small memory of how ONE thing is done, not a manual

- **One domain per skill (SRP).** Content of another domain moves to that skill.
- **Link, never repeat (DRY):** name the owning skill in bold (`**codejob**`); duplication found → split or delete.
- **Size:** good ≤ 50 lines; > 50 warns; > 150 makes `devskills` refuse to sync (`CheckSkills`, `SkillLinesGood`/`SkillLinesMax`).
- Imperative rules, one line each; rationale only when it prevents a misuse. The reader is an agent.
- Format: `skills/<name>/SKILL.md` with frontmatter `name` + `description` (what it covers + "Use when …" triggers).
- English (master plans are the exception: **master-plan**).

## Sync — after every SKILL.md change

```bash
cd webtyp/devskills && go install ./cmd/devskills && devskills -f
grep -c "<a new phrase>" ~/.claude/skills/<name>/SKILL.md   # verify it landed
```

- Rebuilding is mandatory: skills are embedded (`//go:embed skills`); the old binary reinstalls old skills and still prints "Skills updated".
- Size check runs first: an oversized skill prints the error and nothing is installed; long ones print `⚠️ long: <name> (N lines > 50)`.
- Installs to `~/skills/`, then symlinks each skill into every detected agent dir (`~/.claude/skills/<name>`, `~/.gemini/…`, `~/.codex/…`, `~/.qwen/…`, `~/.config/opencode/…`, `~/.agents/…`); copies when a symlink fails. Vendor skills in those dirs are preserved.
- `devskills -l claude` syncs one agent.
- Never put `devskills` in a `PLAN.md` for an external agent.
