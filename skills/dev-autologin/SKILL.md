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

Rather than creating custom bypasses, test-only login buttons, or fake mock users that diverge from real production security, `webtyp.com/auth` replays real credentials from `.env` directly through the active authenticators, minting a genuine session.

---

## Configuration via `.env`

Add `DEV_AUTOLOGIN` to your application's `.env` file with the exact JSON payload expected by the configured authenticator:

### 1. Trusted IP Authenticator (`trusted_ip`)
Used when authentication is bound to an identity code (e.g. RUT, employee ID, account number):

```bash
DEV_AUTOLOGIN={"code":"12345678-5"}
```

### 2. Email & Password Authenticator (`email_password`)
Used when authentication requires username/email and password:

```bash
DEV_AUTOLOGIN={"email":"admin@example.com","password":"secretpassword"}
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
   If credentials fail or the client IP is not authorized, the server logs:
   ```text
   auth: DEV_AUTOLOGIN was rejected by every enabled authenticator — showing the normal login.
   ```
   and serves the normal login form without retrying until the next server restart.

---

## Production Safety

`DEV_AUTOLOGIN` is strictly guarded by Go build tags:
- In production builds (`//go:build prod`), `devlogin_prod.go` compiles a stub where `devAutologinBody` is unconditionally `nil`.
- The environment variable is **never read** in production, ensuring zero risk of accidental credential exposure or unauthorized session creation in deployed binaries.

---

## Rules & Constraints

- **Never create fake dev login buttons in UI:** Do not build custom bypass buttons or mock UI routes in modules. Let `DEV_AUTOLOGIN` issue the session natively.
- **Valid JSON only:** The variable value must be valid JSON matching the authenticator's DTO format.
- **Match database seed data:** The credentials in `DEV_AUTOLOGIN` must correspond to a user seeded in your development database.
