---
name: wasm
description: Go rules for code compiled to WASM with TinyGo — stdlib replacements (webtyp/fmt, webtyp/time, webtyp/json), no maps, goroutines and channels under asyncify. Use when writing frontend Go code.
---

# WASM (TinyGo)

Dev daemon and MCP: **webtyp-app**.

- **`wasm` is a platform, not "the browser":** the Cloudflare Worker is wasm too. Browser-only code is `//go:build browser`; tags table in **project-layout**.
- **stdlib replacements:** `webtyp/fmt` for `fmt`/`strings`/`strconv`/`errors`; `webtyp/time` for `time`; `webtyp/json` for `encoding/json`.
- **No `map` in WASM code** (binary bloat): use structs or slices.
- **Goroutines and channels are safe:** TinyGo's wasm target uses the `asyncify` scheduler; `js.FuncOf` callbacks sending into a channel is the proven bridge (`webtyp/indexdb`).
