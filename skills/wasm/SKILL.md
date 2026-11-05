---
name: wasm
description: Go rules for code compiled to WASM with TinyGo — stdlib replacements (webtyp/fmt, webtyp/time, webtyp/json), no maps, goroutines and channels under asyncify. Use when writing frontend Go code.
---

# WASM (TinyGo)

Dev daemon and MCP: **webtyp-app**.

- **`wasm` is a platform, not "the browser":** the Cloudflare Worker is wasm too. Browser-only code is `//go:build browser`; tags table in **project-layout**.
- **stdlib replacements:** `webtyp/fmt` for `fmt`/`strings`/`strconv`/`errors`; `webtyp/time` for `time`; `webtyp/json` for `encoding/json`. Import them without alias so call sites read `fmt.*`, `time.*`, `json.*` — never `tinytime`, `tinyfmt` or `tinyjson` (there is no `tinytime`).
  STATUS (remove this line when CODE_GUARD_MASTER_PLAN closes): measured with TinyGo in `gocheck/docs/STDLIB_WEIGHT.md` and applied to ALL code (browser, Worker, server — one way, no duplicated code): replace only `fmt`, `encoding/json`, `reflect`, `strconv`; use the stdlib `strings` and `errors` (the `webtyp/fmt` string helpers are being removed; `fmt.Err` stays for translatable messages); never compare two interface values (`err == X`, `errors.Is`, +27 KB): check the type with `errors.AsType`; no maps anywhere.
- **Local time only through a `time.Zone`:** `zone, err := time.LoadZone(deps.Timezone)` then `zone.Date(u)`, `zone.Clock(u)`, `zone.At(date, minutes)`, `zone.Day(date)`. There is no implicit local time; what is UTC says so in its name (`DateUTC`, `WeekdayUTC`, `MidnightUTC`). `time.MachineZone()` only for tools. Store an instant (`time.Unix`), never a local date or hour as text beside it.
- **No `map` in WASM code** (binary bloat): use structs or slices.
- **Goroutines and channels are safe:** TinyGo's wasm target uses the `asyncify` scheduler; `js.FuncOf` callbacks sending into a channel is the proven bridge (`webtyp/indexdb`).
