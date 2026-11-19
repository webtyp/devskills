---
name: dev-autologin
description: Protocol and configuration for DEV_AUTOLOGIN in webtyp.com/auth. Explains how development environments automatically open sessions via environment variables without UI login forms or boilerplate.
---

# Development Autologin (`DEV_AUTOLOGIN`)

## What this is about

`DEV_AUTOLOGIN` is a built-in zero-boilerplate authentication protocol provided by `webtyp.com/auth`. It allows developers to automatically authenticate into local development environments without manually filling out login forms or clicking login buttons each time the server restarts.

## When the reader will run into it

You will use this whenever you are developing or testing an application that requires authenticated access (e.g., an administration panel, protected API endpoints, or user settings) and you want your browser and API calls to authenticate immediately upon startup.

## Why this exists

Manual logins during local development waste time: every server restart or session expiration forces the developer to re-enter credentials or navigate through auth redirect flows. 

Rather than creating custom bypasses, test-only login buttons, or fake mock users that diverge from real production security, `webtyp.com/auth` replays real credentials from the project's keyring directly through the active authenticators, minting a genuine session.

---

## Configuration

Declare `DEV_AUTOLOGIN` in `config/deploy.go`:

```go
deploy.Secret{
	Key:    "DEV_AUTOLOGIN",
	What:   "Credentials replayed to open a dev session automatically",
	Source: deploy.Typed,
	Guide:  []lang.Text{"Type the JSON payload of your authenticator, e.g. {\"code\":\"12345678-5\"}"},
	To:     []deploy.Place{deploy.ToLocal},
}
```

Each developer types the value in the **SECRETS** tab of `webtyp`; `webtyp.com/env` reads it from the project's keyring (`keyring/auto.OpenForModule`). Never store it in `.env`.

Example values expected by the configured authenticator:

### 1. Trusted IP Authenticator (`trusted_ip`)
Used when authentication is bound to an identity code (e.g. RUT, employee ID, account number):

```json
{"code":"12345678-5"}
```

### 2. Email & Password Authenticator (`email_password`)
Used when authentication requires username/email and password:

```json
{"email":"admin@example.com","password":"secretpassword"}
```

---

## How It Works (Protocol Internals)

1. **Startup Inspection:**
   When `authority.New(...)` initializes in a non-production build, it reads `os.Getenv("DEV_AUTOLOGIN")`.
2. **First Unauthenticated Request:**
   When an incoming HTTP request lacks a valid session cookie, the `Authn` middleware intercepts the request before it reaches protected handlers or redirects to the login page.
3. **Authenticator Replay:**
   The middleware passes the `DEV_AUTOLOGIN` JSON payload to the enabled authenticator (verifying credentials and checking that the client IP is allowed for that user, typically `127.0.0.1` or `::1`).
4. **Session Issuance:**
   If authentication succeeds, a real cryptographic session cookie is issued in the response headers, emitting `EventDevAutologin`:
   ```text
   auth: DEV_AUTOLOGIN opened a session for user <user_id>
   ```
5. **Rejection Handling:**
   If credentials fail or the client IP is not authorized, the authenticator reports the reason as a SECURITY line and the server adds one short line:
   ```text
   auth: DEV_AUTOLOGIN rejected, login form shown on localhost
   ```
   That host serves the normal login form without retrying until the next server restart.

---

## Several users — one tab each

To test how the app behaves for different users, list their credentials; each one gets its own browser tab, already signed in:

```json
[{"code":"12345678-5"},{"code":"9876543-2"}]
```

- Credential 1 signs in on `localhost:PORT`, credential N on `userN.localhost:PORT` (Chrome resolves `*.localhost` to 127.0.0.1, so `trusted_ip` sees the same IP, and each host keeps its own session cookie). A single object still works as before, on any host.
- When it opens the browser (and when you switch project or show it again), webtyp asks the app's server `auth.PathDevAutologinTabs` (only in `-tags dev`, only from loopback) and opens one tab per host (`devbrowser.OpenTab`). MCP browser tools follow the tab you have in front (each tab keeps its own console, network and errors; request interception, the selected element and device emulation stay on the first tab); every tab reloads on rebuild.
- **Both factors:** autologin issues the main session and the second-factor one (so screens with `Definition.MFA`, e.g. `access`, open in dev), and upgrades a one-factor session on a host whose credential is live. Production never has this path (`devlogin_off.go`).
- **Sign-out** on a tab pauses autologin for that host only: it shows the normal login (sign in by hand with any code; that session keeps one factor) until the server restarts.
- A credential rejected once stops only its own host: one short line names the host (`rejected, login form shown on user2.localhost`); the reason is in the SECURITY line the authenticator reports. Invalid JSON disables autologin entirely, with one warning on the first request.
- With `trusted_ip`, every listed user must be active staff assigned to the device of the IP you test from (`staff_manager` `IsTrustedIP`).

---

## Production Safety

`DEV_AUTOLOGIN` is strictly guarded by Go build tags:
- Only the webtyp dev environment compiles with `-tags dev` (`devlogin.go`). Every other build — production is simply the normal build — compiles `devlogin_off.go` (`//go:build !dev`), where `devAutologinBody` is unconditionally `nil`. An app follows the same convention: production behaviour by default, development relaxations behind `dev` (never a `prod` tag someone must remember to pass; `gotest` runs the default build and a `-tags dev` pass).
- The environment variable is **never read** in production, ensuring zero risk of accidental credential exposure or unauthorized session creation in deployed binaries.

---

## Rules & Constraints

- **Never create fake dev login buttons in UI:** Do not build custom bypass buttons or mock UI routes in modules. Let `DEV_AUTOLOGIN` issue the session natively.
- **Valid JSON only:** The variable value must be valid JSON matching the authenticator's DTO format.
- **Match database seed data:** The credentials in `DEV_AUTOLOGIN` must correspond to a user seeded in your development database.
