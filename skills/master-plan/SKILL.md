---
name: master-plan
description: Multi-repo master plan as the only conversation state — a short fixed-format file in the user's language, questions first, decisions and pending work only, never history — and the «consolida el master» command. Use when a task spans more than one library, when the user answers inside a master plan, or says "consolida (el master)".
---

# Master plan

Per-library plans: **plan-authoring**; their loop: **codejob**.

## What a master is

**Working state, not an archive.** It holds only what someone needs to keep the wave moving: the questions the user must answer, the decisions not yet built, and the work not yet done. It **shrinks** as the wave advances, because each finished thing leaves for the place that owns it. When the wave ends it is deleted (and its row in `MASTER_PLANS.md`).

**Hard limit: 150 lines.** Check with `wc -l` after every edit. Over the limit means something that belongs elsewhere is still in it: move it out before adding anything.

## Fixed format (sections in this order, no others)

```markdown
# <Tema> — master plan
> **Status:** <PENDING|IN PROGRESS|CLOSED> · <date> · <one line: what is blocked on whom>

## Preguntas abiertas
**P<n> — <title>.** <one or two lines of context>. A) … B) … **Recomiendo X**: <why, with file:line>. Respuesta:

## Decisiones por implementar
| # | Decisión | Dónde se implementa |

## Fases
| Fase | Repo | Qué | Estado |

## Bloqueos
- <date> · <what is stuck now, and on whom>   ← only current ones; a solved blocker is deleted

## Relación con otras olas
- <one line each>
```

- **Questions first**, because they are what the user reads. No questions → write «Ninguna.»
- Each question: ≥2 options, one recommendation with evidence (`file:line`, prior art), at most ~6 lines. Never ask what the code answers.
- A decision row is **one line** with where it gets built; it leaves when that repo ships it.
- A phase row is one line; a closed phase stays as one line («cerrada vX.Y.Z») until the wave ends.
- Written in **the language the user writes in**; dispatched `docs/PLAN.md` stay in English.
- Status line first-class: it must say what the wave is waiting for.

## What never goes in a master

- **Conversation history**: answers to the agent (`codejob reply`), review diaries, "what was tried", timestamps of notices, the user's raw answers. A reply sent to a running agent is recorded **in that repo's `docs/PLAN.md`** (section `## Additions`), which dies with the plan.
- **Evidence, rationale, prior art, rejected options of closed decisions** → the library's `docs/ARCHITECTURE.md` or a permanent `docs/<TOPIC>_DESIGN.md` (**documentation**).
- **Rules and lessons** («the agent guesses signatures», «time is nanoseconds») → the skill that owns them (**devskills**), same turn.
- **Implemented decisions and shipped phases' detail** → the module's own docs; the master keeps one line.
- Long tables of facts, research summaries, verbatim plans or code. Link the file instead.

## Rules

- File: `webtyp/docs/<TOPIC>_MASTER_PLAN.md` — descriptive name, never a bare `MASTER_PLAN.md`.
- Indexed as a row in `webtyp/docs/MASTER_PLANS.md`; read the index before starting a wave (it may exist or be blocked).
- Every phase that **creates** a library or module says, in one line, what the **ecosystem search** found (**api-design** → «Step 0»): reuse, extend, or new.
- Chat reply = list of files touched. No summaries.
- One `docs/PLAN.md` per library, only once its decisions are closed.
- **Nothing is left in the conversation.** The moment decisions close, write every plan of the wave on disk (or exact instructions as a phase row). A session can end at any time.
- A new working rule from the user goes into the skill that owns it, same turn.
- Numbering never restarts: new questions continue after the highest `P` (and `D`) ever used in the file.
- Permanent docs and skills never link a master; while it lives, a reference carries `STATUS (remove this line when <TOPIC>_MASTER_PLAN closes): …`.

## «Consolida el master»

1. Read the whole file, including images the user added.
2. One option left or "respuesta: X" → chosen. Free text → the user's intent, made precise. "No entiendo" / "explícame" → the question stays open, **rewritten simpler** (an example instead of jargon). "Ya contesté" → derive the answer from the other answers and say so in the decision.
3. Closed `P` → a row in «Decisiones por implementar» (`D<n> (P<n>)`); remove the question and the raw answer.
4. Re-read the code the decisions touch: an obvious implication goes inside the decision; a real trade-off becomes a new `P`.
5. Propagate to «Fases», the status line and the `MASTER_PLANS.md` row.
6. Move out everything finished, verified or historical (see «What never goes in a master»), then `wc -l` ≤ 150.
