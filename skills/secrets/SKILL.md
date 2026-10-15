---
name: secrets
description: Where secrets live in the webtyp ecosystem — production secrets only in the destination's store (GitHub Actions, Worker), personal and development secrets in each person's webtyp/keyring, entered through the webtyp/app wizard; never in .env, the repo, a plan, a shell command or the chat. Built for teams. Use when a project needs a token, key, password or client credential, when writing deploy docs or a plan that mentions one, or when tempted to tell the user to paste or `export` a value.
---

# Secrets

A secret is any value that grants access: API tokens, private keys, client secrets, passwords, webhook/HMAC keys, `DEV_AUTOLOGIN` credentials.

## Where it lives — built for a team

- **Production secret → only the destination's store** (GitHub Actions secrets, Worker secrets). **Nobody keeps a copy**, so two people never overwrite each other with a stale value.
- **Personal or development secret → that person's `webtyp/keyring`** (OS credential store; browser: WebCrypto): their Cloudflare bootstrap token, their own dev credentials. Never shared between people.
- **Never** in `.env`, a committed file, a `docs/PLAN.md`, a log line, a test fixture or a chat message.

## Rules

- **Destinations are write-only**: a production secret is never recovered, it is **rotated** (create a new one, replace). Whether it exists is read from destination metadata (names, `updated_at`), never from the value.
- **A person types a value once, into the `webtyp/app` wizard** (a `Sensitive` step). Docs and agents never say "add `X=...` to `.env`", "`export X=...`" or "`gh secret set X`" with a value.
- **Tooling delivers it**, from the wizard straight to the destination's API. A hand copy into a dashboard is acceptable only when the provider has no API (say so, with the screen and field).
- **The project declares which secrets it needs** (`config/deploy.go`); the framework creates, asks for and delivers them. A library never hardcodes a list of secret names.
- **An agent never asks for, reads back or prints a secret value.** It names the key, says what it is, where it comes from and what breaks without it (skill **documentation**), guides the person step by step through any manual provider screen, and points them to the wizard to type the value.
- **Offboarding rotates nothing**: remove the person from the Cloudflare account and the GitHub organization; per-project tokens belong to the account.
- A secret found in `.env`, git history or a doc is a finding: report it, rotate it.
- Tests generate their own throwaway keys (`asym.GenerateKeyPair()`, random bytes); never a real secret.

Flow: `webtyp/docs/diagrams/DEPLOY_SECRETS_FLOW.md`. Tooling status: `webtyp/docs/DEPLOY_WIZARD_MASTER_PLAN.md` — until it ships, `.env` may still hold the `KEY=keyring://` marker read by the app's SECRETS tab.
