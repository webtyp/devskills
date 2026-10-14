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
- Data stays `string` and never passes through `lang.Translate`: user data, ids, timestamps, format samples (`"12345678-9"`).
- Inputs: `SetExample(...)` takes the format sample; `SetPlaceholder(...)` only translatable words.
