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
func Screens(caller router.Caller) []view.Screen
```

`modules/server.go` and `modules/browser.go` are the typed registries (one append per module). A screen is **data** (`view.Screen`): a leaf has `Presenter` + `Permission`, a group has `Children` + `Help` and no permission. The app draws them with `layout/app.BuildMenu`, which runs `Screen.Validate` and returns its `error` — never a silently missing screen.

| File | Tag | Why |
|---|---|---|
| `server.go` | `//go:build !browser` | in the native server **and** the Worker; stripped from the browser binary |
| `browser.go` | `//go:build browser` | only in `browser.wasm`; its view test runs with `-tags browser` in headless Chrome (`gotest`) |
| `module.go` | none | `ID`, `Label`, shared by both sides |
| `svg.go`, `css.go` | `//go:build !wasm` | read by `sitec` at build time; never in the browser nor the Worker |

## Rules

- The whole view is built in `browser.go`; `module.go` holds only `ID`, `Label`.
- **Lazy:** never fetch in `Init()`; load in `Activate()`, called by the chassis when the screen becomes active (`layout/app` wires it).
- STATUS (remove this line when APP_COMPOSITION_MASTER_PLAN closes): a screen that is not a list (wizard, chat) gets `view.Screen.Component` (+ `Help`); until it ships such a module keeps its old screen — never flatten it into a list. Its component may be built with `layout/chatview`, `layout/crudview` or `layout/rightpanel`: that module then imports those pieces (the one accepted exception), never `layout/platformd`.
- Flat: no subdirectories except `docs/` (`find modules -mindepth 2 -type d | grep -v /docs` → empty).
- Forbidden names: `init.go`, `backend.go`, `view.go`, `frontend.go`, `client.go`, `*_wasm.go`.
- Static check: `webtyp.com/router/layoutscan.VerifyLayout`.

## Reusable modules (`github.com/veltylabs/<m>`)

The app writes **no** `modules/<m>/`; it imports `<m>/ui` from its registry. The module repo carries:

- `<m>/ui` — `ID`, `Label`, `Screens(caller) []view.Screen` with **no tag** (only the browser main imports `ui`), its asset producers (`css.go`, `svg.go`, `image.go`…) `!wasm`. Neither the root nor `ui` imports `layout` (exception: a `Component` screen built with `layout/chatview`/`crudview`/`rightpanel`; never `platformd`), so a `layout` release does not force a module release.
- `<m>/seed` — `Load(...) (Data, error)`, demo data through the module's own methods.
- `<m>/web/browser/main.go` — runnable demo (`storage/mem` + `router/loopback` + seed). **A screen is built and tested there, never in the app.**

Module dependencies form a DAG; a screen mixing modules lives in the most-downstream one. When modules wire each other, `modules/server.go` returns a typed struct with `Operations() []router.OperationModule`.

## Before a module enters an app — the review it must pass

Each point comes from a real defect (mjosefa-cms audit, 2026-10-28). A module that fails one stays out of the app's registry until fixed in its own repo.

1. Saves with the wire shapes of **module-layout** (`{records}`, `{ids, fields, record}`, `{ids}`) and answers errors through the shared piece, never its own `errors.go`. STATUS (remove this line when MVP_DEPLOYMENT_ROADMAP closes): the shared piece is `webtyp/rpc` (its M2i); until then the 400 carries the field and a 500 never carries `err.Error()`.
2. One resource per table, named by the framework, never typed as a string. STATUS (remove this line when MVP_DEPLOYMENT_ROADMAP closes): until its M2h ships, a `res*` constant equal to the table's `ModelName()`.
3. Foreign keys declared: its band of the app's database diagram has no errors and no warnings of its own.
4. Every screen registered, opened in the module's demo and seen in a screenshot.
5. Tests that travel through the real codec (real JSON), never a fake writer.
6. English keys, translated: `langc sync .` reports nothing pending.

## Checklist

- [ ] No `modules/<m>/` duplicating a module that ships `<m>/ui`.
- [ ] Flat files, canonical names, view in `browser.go`.
- [ ] No function returns `[]any`.
- [ ] Tags as in the table above.
