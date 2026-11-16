---
name: console-logs
description: Formatting and reducing noise in webtyp terminal/MCP logs — ~ home shortening, CompactList path grouping, no redundant handler prefixes, and OperationID accumulation. Use when writing or reviewing log emissions in devflow, app, devtui, or any TUI/MCP handler.
---

# Console & MCP Logs

Every log line consumes terminal space and LLM context. Keep output compact, non-redundant, and free of raw home paths.

## Paths

- **Never log raw `/home/...` paths:** `devtui` applies `filepath.Short` (`.`-relative inside project) and `filepath.Tilde` (`~/...` outside), and `app` / `devflow` callers should pass external paths through `filepath.Tilde` or `filepath.CompactList`.
- **Group path lists:** pass multi-path slices to `filepath.CompactList(paths)` so the common parent prints once with `~` and differing tails appear in braces (`~/Dev/Project/webtyp/{form, model}`). Never `strings.Join(paths, ", ")` on raw paths.

## Batching & Accumulation

- **One line per batch:** collect added/removed/processed items in a slice and emit one summary log after the loop, never one log per item.
- **Collapse bursts with `OperationID`:** handlers that may fire multiple times for one logical operation (e.g. file-watcher bursts within the same second) must implement `GetLastOperationID() string` / `SetLastOperationID(id string)` so `devtui` (local and remote SSE via `HeadlessTUI`) updates the existing line in place instead of appending duplicates.
- **Start a new line explicitly** by changing the ID via `SetLastOperationID(newID)`.
- **Identical consecutive messages** from the same handler are deduplicated automatically; never add manual repeat-suppression wrappers at call sites.

## Message Text

- **Never repeat the badge name:** `devtui` already renders the handler badge (`GOMOD`, `WATCH`, `SERVER`); omit prefixes like `GoModHandler:` in log strings (`devtui` also strips `<Word>Handler:` when it matches the badge).
- **Fixed UI copy uses `lang.Translate`;** dynamic paths and counts are passed as separate arguments (**translations**).
