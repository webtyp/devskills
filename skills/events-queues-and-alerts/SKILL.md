---
name: events-queues-and-alerts
description: The events domain of webtyp — typed publish/subscribe (events, events/inprocess), the security event vocabulary every layer reports through (securityevent), the durable queue of calls to other systems (externalcalls), retry delays (backoff) and alerts with their delivery channels (notify) — its repos, who imports whom, each one's single job and the open defects. Use before publishing or listening to an event, reporting something suspicious, queueing a call to another system, retrying, or sending an alert.
---

# Events, queues and alerts

Map of the whole framework: [framework-map](../framework-map/SKILL.md). Last link: the code (`go doc webtyp.com/events`).

```mermaid
flowchart TD
    appkit[appkit: wires the broker, listens to security events] --> events
    appkit --> externalcalls
    appkit -->|logs every event, counts failures, raises alerts| securityevent
    appkit --> notify
    gate[server/httpd: 401/403 & /_csp-report, router/security: 413, crud: rejected write] --> securityevent
    app[app: BUILD/SECURITY view in webtyp dev & MCP] -.reads.-> appkit
    auth[auth, crud, rpc, mcp, sse, money] --> events
    authmodes[auth and its modes] --> securityevent
    securityevent[securityevent: what a security event is called and carries] --> events
    securityevent --> msgtype[msgtype: Info, Warning, Error]
    notify[notify: alerts, recipients, delivery channels] --> events
    notify --> msgtype
    externalcalls[externalcalls: durable queue of calls to other systems] --> backoff
    externalcalls --> orm[orm, ddl: its own table]
    outbox[outbox, sse, goflare, app] --> backoff
    backoff[backoff: how long to wait before the next retry] --> time
    inprocess[events/inprocess: the broker inside one process] --> events
    events[events: Publisher, Subscriber, Event] --> model
```

| Repo | Its one job |
|---|---|
| `events` | the publish/subscribe contract (`Publisher`, `Subscriber`, `Event`); `events/inprocess` is the broker inside one process |
| `securityevent` | the one vocabulary of security events (16 `Kind`s with OWASP names and levels); `Reporter.Report` and `Listen`; never a secret |
| `externalcalls` | a durable queue of calls to other systems: record first, call later, retry with `backoff` |
| `backoff` | how long to wait before the next retry (doubles up to a longest delay); never sleeps |
| `notify` | alerts (`Alert` with a `msgtype` severity), recipients and the channels that deliver them |

## Rules

- A module publishes through `events.Publisher` and listens through `events.Subscriber`; it never imports a concrete broker.
- Anything suspicious is reported with `securityevent` (`Reporter.Report`), never as a printed line or a local event type.
- A call to another system goes through `externalcalls` (record, then call), never straight from a rule.
- A retry delay comes from `backoff`, never a hand-written loop.

## Open defects (each leaves when fixed)

- `notify`: `appkit` raises security alerts on it and persists `Warning`/`Error` alerts to `logs.log` (`!wasm`, 7d / 1 MiB via `notify.NewCriticalFileLog`); external delivery (email/WhatsApp) is wired by the app on `notify.TopicAlert`.
