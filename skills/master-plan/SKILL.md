---
name: master-plan
description: Multi-repo master plan as the only conversation state — user's language, questions and decisions in the file, never in chat — and the «consolida el master» command. Use when a task spans more than one library, when the user answers inside a master plan, or says "consolida (el master)".
---

# Master plan

Per-library plans: **plan-authoring**; their loop: **codejob**.

## Rules

- File: `webtyp/docs/<TOPIC>_MASTER_PLAN.md` — descriptive name, never a bare `MASTER_PLAN.md`.
- Second line is the status: `> **Status:** <PENDING|IN PROGRESS|SHIPPED|CLOSED> · <date> · <one line>`.
- Indexed as a row in `webtyp/docs/MASTER_PLANS.md`; read the index before starting a wave (it may exist or be blocked).
- States its relation to prior waves it touches (extends, supersedes, independent).
- The libraries table marks gates (block the next) vs parallel phases.
- Written in **the language the user writes in**; dispatched `docs/PLAN.md` stay in English.
- Questions live in the file: `P<n>`, ≥2 options, one recommendation with evidence (`file:line`, prior art). Never ask what the code answers.
- Chat reply = list of files touched. No summaries.
- One `docs/PLAN.md` per library, only once its decisions are closed; the master plan tracks them in a table.
- **Nothing is left in the conversation.** The moment decisions close, write **every** plan of the wave on disk — including follow-ups and small "in-session" tasks (each gets its own plan or exact instructions in the master). A plan whose repo already has a pending `docs/PLAN.md` goes to `docs/PLAN_<TOPIC>.md` and the master says so. Why: a session can end at any time; a plan that only exists in context makes later plans inconsistent.
- Every addition sent to a running agent (`codejob reply`) is also recorded in the master, verbatim.
- A new working rule from the user goes into the skill that owns it, same turn.
- Numbering never restarts: new questions continue after the highest `P` used.

Order: status · objective/scope · **Decisiones** · **Preguntas abiertas** · evidence · api-design gate · libraries table · other waves.

## «Consolida el master»

1. Read the whole file, including images the user added.
2. One option left or "respuesta: X" → chosen. Free text → the decision is the user's intent, made precise. "Explícame" / "propón" → stays open, rewritten.
3. Closed `P` → `D<n> (cierra P<n>)` with its consequences (what is deleted, which library); remove the question and the raw answer.
4. Re-read the code the decisions touch: an obvious implication goes inside the decision; a real trade-off becomes a new `P`.
5. Propagate to gate, libraries table, status line and the `MASTER_PLANS.md` row.
