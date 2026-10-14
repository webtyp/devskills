---
name: webtyp-app
description: The global webtyp MCP daemon — starting it, starting a project headless, the human TUI, IDE auto-config, auth token, HTTP endpoints and diagnostics when the MCP does not answer. Use when working with the webtyp development environment.
---

# webtyp/app — the dev daemon

- **The agent talks to projects only through the global MCP server on `:6060`.** Not running → `webtyp -mcp`.
- Start a project with the MCP tool `start_development` (headless). **Never run `webtyp` in a shell to start a project.**
- Tools are registered once at startup and self-describe via `tools/list`; `browser_*` tools answer "not ready" until a project is active.
- The human runs `webtyp` (TUI, SSE log viewer): `Ctrl+C` detaches, the project keeps running; `q` stops the project.

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
webtyp -mcp                                   # not running → start it
grep -A5 mcpServers ~/.claude.json            # check the IDE entry
```
