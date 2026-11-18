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
    appkit -.🔴 not yet: D2 of SECURITY_OBSERVABILITY.-> securityevent
    auth[auth, crud, rpc, mcp, sse, money] --> events
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

- `notify`: nobody imports it yet; security alerts will reach it through the `appkit` subscriber (SECURITY_OBSERVABILITY D2).
- `securityevent` (v0.0.2): nobody reports through it yet; `auth` still uses its own `SecurityEventType` (SECURITY_OBSERVABILITY O2).
