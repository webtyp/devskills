---
name: translations
description: Authoring user-facing text in webtyp — English-only keys in Go, noun-before-adjective word order, gender-agnostic phrasing, lang.Text for fixed UI copy vs plain string for data. Use when writing any UI text, label, button or error message. The lang.json dictionary and langc CLI: skill langc.
---

# Translations — writing text

Dictionary, extraction and sync: **langc**.

## Code is English only

The Go string is the lookup key; an untranslated key passes through unchanged.

```go
return fmt.Err("Rol no encontrado") // ❌
return fmt.Err("Role not found")    // ✅ Spanish lives in lang.json
```

## Where a translation is written — once, at the lowest owner

- A text is translated in the `lang.json` of the library that writes it; never in an app (an app's `config/lang.json` is configuration only — **langc**).
- A shared column or word (`tenant id`, `created at`, `name`, `status`) is translated once in `webtyp/model/lang.json`; a module never repeats it, `langc sync` drops a copy that agrees and reports one that disagrees.
- A reference to an entity (`room id`, `device id`) belongs to the module that owns the entity.
- Leave a language **empty** when no person who reads it reviewed it. Never fill a slot with a placeholder (`"name (zh)"`, `"fake …"`, `"Traducción"`, `"kind español"`) or copy the Spanish into other languages: an empty slot shows the English key and is reported; a placeholder is shown to users and hides the gap.
- STATUS (remove this line when TUI_SIGNAL_MASTER_PLAN closes): local, unpublished (`lang/langc`, `model/lang.json`).

## Think in Spanish first — noun before adjective

Multi-word keys and variadic `lang.Translate(...)` must read naturally in Spanish:

| Wrong | Right |
|---|---|
| `"invalid", "format"` → "Inválido Formato" | `"format", "invalid"` → "Formato Inválido" |
| `"not", "initialized", "directory"` | `"directory", "not", "initialized"` |

## Gender-agnostic phrasing

Never masculine inflection toward the user; use verbs or neutral phrasing:

| English | ❌ | ✅ |
|---|---|---|
| Are you sure? | ¿Estás seguro? | ¿Confirmas? |
| Welcome | Bienvenido | Te damos la bienvenida |
| User registered | Usuario registrado | Registro completado |

## `lang.Text` vs `string`

- Fixed UI copy (labels, buttons, titles, placeholders) is typed `lang.Text`; translate in `Render()` with `lang.Translate(x)`. Callers still write literals: `Modal{Title: "Delete project"}`.
- Data stays `string` and never passes through `lang.Translate`: user data, ids, timestamps, format samples (`"12345678-9"`), prefixes and type names. Pass non-prose as variables so `langc` does not scoop them into `lang.json`.
- Reuse vocabulary: decompose shared error patterns (`fmt.Err(opPrefix, "requires at least one", "field")`) so multi-word phrases are shared across operations instead of repeating full sentences.
- Inputs: `SetExample(...)` takes the format sample; `SetPlaceholder(...)` only translatable words.
