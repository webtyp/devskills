---
name: secrets
description: Where secrets live in the webtyp ecosystem — in the destination's store in production and in webtyp/keyring locally, entered through the webtyp/app wizard, never in .env, the repo, a plan, a shell command or the chat; how a deploy destination (CI, Worker) receives them. Use when a project needs a token, key, password or client credential, when writing deploy docs or a plan that mentions one, or when tempted to tell the user to paste or `export` a value.
---

# Secrets

A secret is any value that grants access: API tokens, private keys, client secrets, passwords, webhook/HMAC keys, `DEV_AUTOLOGIN` credentials.

## Rules

- **Two places, no others.** In production it lives in the **destination's** store (GitHub Actions secrets, Worker secrets). On the developer's machine it lives in **`webtyp/keyring`** (OS credential store; browser: WebCrypto), the source the wizard feeds destinations from. Never in `.env`, a committed file, `deploy.yaml`, a `docs/PLAN.md`, a log line, a test fixture or a chat message.
- **Destinations are write-only**: the keyring never fetches a secret back from the cloud. Missing locally → the wizard re-creates it (mint or generate, then replace at the destination) or asks for it. Flow: `webtyp/docs/diagrams/DEPLOY_SECRETS_FLOW.md`.
- **The developer types it once, into the `webtyp/app` wizard** (a `Sensitive` step). Docs and agents never say "add `X=...` to `.env`", "`export X=...`" or "`gh secret set X`" with a value.
- **Destinations are fed by tooling, from the keyring**: CI secrets, Worker secrets, etc. Copying a value by hand into a dashboard is the documented exception only when the provider has no API (say so, with the screen and field).
- **The project declares which secrets it needs; the framework stores and delivers them.** A library never hardcodes a list of secret names.
- **An agent never asks for, reads back or prints a secret value.** It names the key, says what it is, where it comes from and what breaks without it (skill **documentation**), guides the developer step by step through any manual provider screen, and points them to the wizard to type the value.
- A secret found in `.env`, git history or a doc is a finding: report it, move it to the keyring, rotate it.
- Tests generate their own throwaway keys (`asym.GenerateKeyPair()`, random bytes); never a real secret.

## Current mechanism

`.env` may hold the reference marker `KEY=keyring://` (`keyring.Scheme`) — never the value; the app's **SECRETS** tab asks for every reference with no keyring entry (`app/missing_secrets.go`). How deploy targets declare and receive secrets is being redesigned: `webtyp/docs/DEPLOY_WIZARD_MASTER_PLAN.md`.
