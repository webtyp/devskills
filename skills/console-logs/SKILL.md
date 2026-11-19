---
name: console-logs
description: Formatting and reducing noise in webtyp terminal/MCP logs — ~ home shortening, CompactList path grouping, no redundant handler prefixes, TUI/MCP handler parity, and OperationID accumulation. Use when writing or reviewing log emissions in devflow, app, server, devtui, or any TUI/MCP handler.
---

# Console & MCP Logs

Every log line consumes terminal space and LLM context. Keep output compact, non-redundant, and free of raw home paths.

## TUI / MCP Parity

- **Every handler is equal for TUI and MCP:** `SSEPublisher` keeps one ring per `TAB/HANDLER` (no handler-specific rings or filter special cases), merges all rings chronologically for the combined view, and reroutes derived lines (e.g. `SECURITY`) without duplicating them in the source handler.
- **Discrete runtime output:** `SERVER` keeps a single live line (`OperationID`) only while compiling/starting/restarting; runtime app stdout/stderr lines are discrete so every line is visible in both TUI and MCP.

## Paths

- **Never log raw `/home/...` paths:** `devtui` applies `filepath.Short` (`.`-relative inside project) and `filepath.Tilde` (`~/...` outside), and `app` / `devflow` callers should pass external paths through `filepath.Tilde` or `filepath.CompactList`.
- **Group path lists:** pass multi-path slices to `filepath.CompactList(paths)` so the common parent prints once with `~` and differing tails appear in braces (`~/Dev/Project/webtyp/{form, model}`). Never `strings.Join(paths, ", ")` on raw paths.

## Batching & Accumulation

- **One line per batch:** collect added/removed/processed items in a slice and emit one summary log after the loop, never one log per item.
- **Collapse build bursts with `OperationID`:** build-progress handlers keep one live TUI line via `GetLastOperationID() string` / `SetLastOperationID(id string)` + `tui.LogOpen`/`tui.LogClose` so intermediate hot-reload steps animate in place and resolved errors clear when the next cycle closes OK; clear `OperationID` before streaming discrete runtime events.
- **No duplicate `WATCH` echo:** handlers that already logged an error under their own badge return an error implementing `AlreadyLogged() bool` (`devwatch.AlreadyLoggedError`) so `devwatch` skips browser reload without reprinting the message.
- **Compact compiler diagnostics:** format raw `go build` / SSR errors with `gobuild.CompactErrors` (groups repeated `file:line: msg` into one line with `~`).
- **Identical consecutive messages** from the same handler are deduplicated automatically; never add manual repeat-suppression wrappers at call sites.

## Message Text

- **Never repeat the badge name:** `devtui` already renders the handler badge (`GOMOD`, `WATCH`, `SERVER`); omit prefixes like `GoModHandler:` in log strings (`devtui` also strips `<Word>Handler:` when it matches the badge).
- **Fixed UI copy uses `lang.Translate`;** dynamic paths and counts are passed as separate arguments (**translations**).
