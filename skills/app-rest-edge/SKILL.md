---
name: app-rest-edge
description: Variant of the webtyp app layout for plain REST and Cloudflare Workers — config/routes.go as the explicit mount table with security policy, modules/panel/ (browser) kept out of edge.wasm, web/edge.go (wasm && !browser) and cmd/migrate/main.go. Use when a project serves REST instead of rpc operations or deploys with goflare.
---

# REST and edge variant

Canonical transport is rpc (`rpc.Mount(r, prefix, modules...)` server-side, `rpc.NewCaller(origin)` client-side); MCP is the agents' surface only, closed by default (`config/mcp.go`). HTTP stays for uploads, static SSR and auth. Everything in **project-layout** still applies, tags included.

## Plain REST

- `config/routes.go` (`!browser`) is the explicit mount table: every `r.Get/r.Post/…` with `.Public()`/`.Authenticated()`. Security headers: `r.Use(security.Policy{}.AllowImages("https://…").Middleware())`.
- Renderer-touching code lives in `modules/panel/` (`//go:build browser`): it never enters `edge.wasm`, because `goflare` passes no destination tag and refuses `-tags browser`.
- Server handlers stay flat in `modules/<domain>/` (`handler.go`, `server.go` with `!browser`).

## Optional files

- `web/edge.go` (`//go:build wasm && !browser`): D1/R2 bindings → `config.Register` (or rpc mount) → `edge.Serve(r)`. `goflare` builds it into `web/edge.wasm` + `web/edge.js` (never into `web/public/`, which is served to anyone). Never schema DDL.
- `cmd/migrate/main.go` (`//go:build !wasm`): reads `CLOUDFLARE_*`/`D1_*`, runs each package's `Migrate`/`Sync`; run as the deploy workflow's `pre-deploy` or by hand. Thin main.
