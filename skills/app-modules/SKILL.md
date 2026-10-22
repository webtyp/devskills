---
name: app-modules
description: How a module is built inside a webtyp app or a veltylabs module repo — the typed server.go/browser.go pair and its build tags (!browser / browser), the two registries, lazy Activate(), flat files, forbidden names, and reusable modules that ship their own ui/, seed/ and web/browser/main.go demo. Use when adding or reviewing a module.
---

# App modules

Directory around them and the two tag axes: **project-layout**.

## A module is `server.go` + `browser.go`, typed — never `[]any`

```go
// modules/<m>/server.go — //go:build !browser
func Server(db *orm.DB, ids model.IDGenerator, pub events.Publisher, tenantID string) (router.OperationModule, error)

// modules/<m>/browser.go — //go:build browser
func Browser(caller router.Caller, ids model.IDGenerator, tenantID string) (platformd.UIModule, error)
```

`modules/server.go` and `modules/browser.go` are the typed registries (one append per module). A construction failure is a returned `error`, never a silently missing module.

| File | Tag | Why |
|---|---|---|
| `server.go` | `//go:build !browser` | in the native server **and** the Worker; stripped from the browser binary |
| `browser.go` | `//go:build browser` | only in `browser.wasm`; its view test runs with `-tags browser` in headless Chrome (`gotest`) |
| `module.go` | none | `ID`, `Label`, shared by both sides |
| `svg.go`, `css.go` | `//go:build !wasm` | read by `sitec` at build time; never in the browser nor the Worker |

## Rules

- The whole view is built in `browser.go`; `module.go` holds only `ID`, `Label`.
- **Lazy:** never fetch in `Init()`; load in `Activate()`, called by `platformd` when the module becomes active.
- Flat: no subdirectories except `docs/` (`find modules -mindepth 2 -type d | grep -v /docs` → empty).
- Forbidden names: `init.go`, `backend.go`, `view.go`, `frontend.go`, `client.go`, `*_wasm.go`.
- Static check: `webtyp.com/router/layoutscan.VerifyLayout`.

## Reusable modules (`github.com/veltylabs/<m>`)

The app writes **no** `modules/<m>/`; it imports `<m>/ui` from its registry. The module repo carries:

- `<m>/ui` — `ID`, `Label`, `Browser(...)` with **no tag** (only the browser main imports `ui`), its asset producers (`css.go`, `svg.go`, `image.go`…) `!wasm`. The domain root never imports `layout`.
- `<m>/seed` — `Load(...) (Data, error)`, demo data through the module's own methods.
- `<m>/web/browser/main.go` — runnable demo (`storage/mem` + `router/loopback` + seed). **A screen is built and tested there, never in the app.**

Module dependencies form a DAG; a screen mixing modules lives in the most-downstream one. When modules wire each other, `modules/server.go` returns a typed struct with `Operations() []router.OperationModule`.

## Checklist

- [ ] No `modules/<m>/` duplicating a module that ships `<m>/ui`.
- [ ] Flat files, canonical names, view in `browser.go`.
- [ ] No function returns `[]any`.
- [ ] Tags as in the table above.
