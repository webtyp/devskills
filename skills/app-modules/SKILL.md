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
- A screen that is not a list (wizard, chat) is a `view.Screen.Component` leaf (+ `Help`) — never flatten it into a list. Its component may be built with `layout/chatview`, `layout/crudview` or `layout/rightpanel`: that module then imports those pieces (the one accepted exception), never `layout/platformd`.
  STATUS (remove this line when ARCHITECTURE_MASTER_PLAN closes): D16 deletes this exception — a module imports only `webtyp/view` contracts (+ `components` in a bespoke screen) and the chassis picks `crudview`/`chatview`/`rightpanel`; only the app (and the daemon, for the demo) imports `layout`.
- Flat: no subdirectories except `docs/` (`find modules -mindepth 2 -type d | grep -v /docs` → empty).
- Forbidden names: `init.go`, `backend.go`, `view.go`, `frontend.go`, `client.go`, `*_wasm.go`.
- Static check: `webtyp.com/router/layoutscan.VerifyLayout`.

## Reusable modules (`github.com/veltylabs/<m>`)

The app writes **no** `modules/<m>/`; it imports `<m>/ui` from its registry. The module repo carries:

- `<m>/ui` — `ID`, `Label`, `Screens(caller) []view.Screen` with **no tag** (only the browser main imports `ui`), its asset producers (`css.go`, `svg.go`, `image.go`…) `!wasm`. Neither the root nor `ui` imports `layout` (exception: a `Component` screen built with `layout/chatview`/`crudview`/`rightpanel`, or implementing `layout.UsesNotifier`; never `platformd`), so a `layout` release does not force a module release. A `Component` screen tells the user «saved» or an error through the platform's notices (`UseNotifier(layout.Notifier)` + `Notify(msgtype.Success|Error, msg, layout.Auto())`), never a status line of its own.
- `<m>/seed` — `Load(...) (Data, error)`, demo data through the module's own methods.
- `<m>/web/browser/main.go` — runnable demo (`storage/mem` + `router/loopback` + seed). **A screen is built and tested there, never in the app.**
  STATUS (remove this line when ARCHITECTURE_MASTER_PLAN closes): D15 replaces it with `seed.Demo() (view.Demo, error)`; the `webtyp/app` daemon generates the main in memory with `layout`, so the module's `go.mod` never lists `layout`.
  STATUS (remove this line when ARCHITECTURE_MASTER_PLAN closes): D54 — the generated main lives in `os.UserCacheDir()/webtyp/demos/<hash>/`, never in the repo, and the CSS/icon scan (`sitec`) reads that same folder, so the demo carries the styles of `layout` and `components`.

Modules never import each other (skill **module-layout**): a screen mixing modules lives in the app, which wires each module's narrow interfaces. When modules wire each other, `modules/server.go` returns a typed struct with `Operations() []router.OperationModule`.

## Webhooks and provider returns — never trust the request

A provider calling back (payment return, signed webhook) is a route of **that provider's adapter package**, mounted by the app
(`AppServer.Routes` or an `Infra`), never a crud table: its format is the provider's and it carries no session.
- The adapter passes the module **identification only** (payment id, provider token). The module's method never accepts an
  outcome from the request: it obtains the truth itself with its secret — **fetch-back** to the provider API (`Gateway.Confirm`), or a
  **signature** over the raw body (`webtyp.com/crypto/hmac`: `HMACSHA256` + constant-time `HMACEqual`, timestamp within 5 minutes) — and
  settles only what the provider reports.
- Idempotent: providers retry; a repeated event changes nothing (state machine or provider event id).
- Answer 2xx fast; slow work runs after commit (outbox, D25 of ARCHITECTURE_MASTER_PLAN while it lives). Secrets: skill **secrets**.
- **One verification core for every signed webhook** (WhatsApp, Facebook/Meta, Stripe-like): raw body, signature schemes (`X-Hub-Signature-256`
  `sha256=` HMAC…), timestamp tolerance, provider handshakes (Meta's `hub.challenge`), body size limit and event-id dedup live in **one**
  `webtyp` library; an adapter only declares its scheme and secret. It does not exist yet: the first signed-webhook consumer creates it
  (**api-design**, reusing `webtyp/crypto/hmac` and `webtyp/idempotency`) — never a verifier written inside an adapter.

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
