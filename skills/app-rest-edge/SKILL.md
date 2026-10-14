---
name: app-rest-edge
description: Variant of the webtyp app layout for plain REST and Cloudflare Workers — config/routes.go as the explicit mount table with security policy, modules/panel/ kept out of edge.wasm, edge/main.go and cmd/migrate/main.go. Use when a project serves REST instead of MCP or deploys with goflare.
---

# REST and edge variant

Canonical transport is MCP (`mcp.HarvestOps` server-side, `mcp.NewCaller(mcp.NewClient(origin, ""))` client-side); HTTP stays for uploads, static SSR and auth. Everything in **project-layout** still applies.

## Plain REST

- `config/routes.go` is the explicit mount table: every `r.Get/r.Post/…` with `.Public()`/`.Authenticated()`. Security headers: `r.Use(security.Policy{}.AllowImages("https://…").Middleware())`.
- Renderer-touching code moves to `modules/panel/` (`//go:build wasm`), which `config/routes.go` never imports — the missing import keeps `dom`/`layout` out of `edge.wasm` (`goflare`/`sitec` pass no `-tags`).
- Server handlers stay flat in `modules/<domain>/` (`handler.go`, `server.go`).

## Optional dirs

- `edge/main.go` (`//go:build wasm`): D1/R2 bindings → `config.Register` (or MCP mount) → `edge.Serve(r)`. Never schema DDL.
- `cmd/migrate/main.go` (`//go:build !wasm`): reads `CLOUDFLARE_*`/`D1_*`, runs each package's `Migrate`/`Sync`; run from the deploy workflow or by hand. Thin main.
