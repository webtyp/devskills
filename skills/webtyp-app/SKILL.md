---
name: webtyp-app
description: The global webtyp MCP daemon — starting it, starting a project headless, the human TUI, IDE auto-config, auth token, HTTP endpoints and diagnostics when the MCP does not answer. Use when working with the webtyp development environment.
---

# webtyp/app — the dev daemon

- **The agent talks to projects only through the global MCP server on `:6060`.** Not running → `webtyp mcp` (a subcommand, as `webtyp help` shows in v0.7.33; run it from inside a Go module, e.g. the project root — elsewhere it exits with "not a Go module").
- Up to 5 projects run at once on ports `:8080`–`:8084`, each on its own port. `start_development` adds a project; `list_development` lists running projects; `stop_development` stops one (or all without argument). `project` argument selects the target for `app_info` and `app_get_logs` (omitted = project used last). Closing a client terminal (`Ctrl+C`) stops only that terminal's project; the daemon stays up until `webtyp stop`. When the binary changes on disk, the daemon snapshots running projects (`projects.json`), relaunches, and restores them.
- STATUS (remove this line when TUI_STAGES_MASTER_PLAN closes): today each case calls `SetActiveTab` and migration runs inside BUILD. Decided: tabs in stage order WIZARD → SECRETS → DATABASE → BUILD → DEPLOY → MCP, `webtyp dev` opens the first with something pending (else BUILD), DATABASE lists leftover columns with data and asks before dropping; MCP `db_migrate_plan` / `db_drop_column`.
- Start a project with the MCP tool `start_development` (headless). **Never run `webtyp` in a shell to start a project.**
- Tools are registered once at startup and self-describe via `tools/list`; `browser_*` tools answer "not ready" until a project is active.
- The human runs `webtyp dev` (TUI, SSE log viewer).

## IDE config — written automatically at startup (`app.ConfigureIDEs`)

| IDE | File |
|---|---|
| Claude Code | `~/.claude.json` → `"webtyp": {"url": "http://localhost:6060/mcp", "type": "http"}` |
| VS Code | `~/.config/Code/User/mcp.json` (or profiles) |
| Antigravity | `~/.gemini/antigravity/mcp_config.json` |

Auth: the daemon persists an API key (`WEBTYP_API_KEY`); requests carry `Authorization: Bearer <key>`. No key → open.

## Endpoints

| Route | Auth | Purpose |
|---|---|---|
| `POST /mcp` | Bearer | JSON-RPC 2.0 MCP + `webtyp/*` methods |
| `GET /logs` | — | SSE log stream |
| `GET /webtyp/state` · `POST /webtyp/action` | Bearer | TUI state · key actions (`q`, `r`, start, stop, restart) |
| `GET /version` | — | daemon version |

## MCP does not answer

```bash
curl http://localhost:6060/mcp -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"initialize","id":"1","params":{"protocolVersion":"2024-11-05","clientInfo":{"name":"t","version":"0"}}}'
webtyp mcp                                    # not running → start it (from a Go module dir)
grep -A5 mcpServers ~/.claude.json            # check the IDE entry
```

## Rebuild on save — the order is the contract

A saved `.go` file runs the watcher handlers **one after another, in this order**: `go.mod` → `ormc` (regenerates
`*_orm.go`) → the WASM client (front) → the server → browser reload (`app/section-build.go`, `devwatch`). Code shared
by front and server rebuilds **both, front first**; code of one side rebuilds only that side. Guarded by
`devwatch/tests/build_gate_test.go` (`TestBuildGate_SequentialOrdering`): a change that breaks this order is a bug,
never a new behaviour.
