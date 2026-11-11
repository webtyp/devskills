---
name: isomorphic-libraries
description: What a webtyp library is — isomorphic by default (the same code runs in the browser over indexdb and on the server over SQL), what must stay out of it (server transports, generator/validator tooling), how kit packages like appkit separate front and back with build tags, and `tests/go.mod` for heavy test dependencies. Use before judging a package "server-only", before adding an import to a library, when a library test needs httpd, or when placing code in model/router/crud/appkit.
---

# Isomorphic libraries

Tags and layout: **project-layout**. Tests: **testing**. Contracts and lego pieces: **api-harness**.

## What a library is

- **Isomorphic by default.** `model`, `storage`, `orm`, `crud`, `router`, `events`, `rbac`, `rpc`, `mcp`… run in the browser (over `indexdb`) **and** on the server (over `sqlite`/`postgres`), so an app is full stack with one declaration. Never conclude "this is server code" because a package imports `storage`, `orm`, `crud` or `router`: the browser uses them too. When unsure, read the library's `AGENTS.md`; do not guess from the import list.
- **Platform split only when the API does not exist there** (`wasm`/`!wasm`) or the file belongs to another binary (`browser`/`!browser`). Everything else is untagged and imports only isomorphic packages.

## What stays out of a library

- **Server transports are injected by the consumer.** `webtyp/server/httpd` (and any `net/http` server) is imported only by `web/server/main.go` and `config/`; a library accepts `router.Router` or a narrow interface. A library whose `go.mod` lists `httpd`, `sqlite` or `postgres` is polluted: remove it.
- **A test that really needs `httpd`** (critical, end to end) lives in `tests/` with **its own `go.mod`** (`replace` to the library under test), so the library's `go.mod` never lists it. Everything else uses fakes at the edge.
- **`model` and `router` hold contracts and isomorphic structures only.** Code that only a generator or a validator tool runs (checkers, parsers) goes in its own library, imported by whoever needs it; its error texts must not enter the dictionary of every WASM consumer.

## Contracts never depend on their implementations — not even in tests

A library that is a contract (`storage`, `orm`, `crud`, `router`, `rbac`, `changelog`…) has no database driver, no server and no other implementation in its `go.mod`, **tests included**. What every implementation must show — atomic transactions, one winner among concurrent writers, a duplicate key reported, a revision that rejects a stale update, a migration that applies — is a **conformance suite**: an exported `conformance` package in the contract's own repo that imports no implementation and takes a factory (`conformance.Run(t, Factory{New: …})`).
Each implementation (`sqlite`, `postgres`, `indexdb`, …) runs that suite in its own `tests/` module and **must pass it**; the contract's own tests run it against the in-memory implementation that ships with it (`storage/mem`). A missing guarantee is added to the suite once, never to a test inside a contract that imports a driver. A library that needs a driver in its `go.mod` to be tested is a defect: move the test into a suite.

## Kit packages (`appkit`)

- A kit unites the contracts of **both sides in one place** for a full-stack app, with the sides separated by build tags: untagged = isomorphic (`App`, `Screens`, `Brand`, module contracts); server composition `!browser`; browser composition `browser`.
- A file named `*_server.go` that compiles into the browser binary is a defect: tag it `!browser`.
- A kit never builds an HTTP server: the server `main` creates it (`httpd.New`) and hands it in.
