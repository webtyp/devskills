---
name: translations
description: Multilingual translation architecture in webtyp/lang and langc — English source key rule, positional lang.json dictionary format, 8 AST extraction rules, langc sync CLI, and runtime dictionary merging. Use when authoring UI text, error messages, handling missing or untranslated keys, or managing lang.json.
---

# Translations (`webtyp/lang` & `langc`)

The `webtyp/lang` package and its companion CLI `langc` provide the multilingual translation system across the WebTyp ecosystem.

Translations in WebTyp follow a strict separation of concerns:
- **Code writes English**: Go source code only ever contains English keys and literals.
- **Translations are data**: Translated phrases live in JSON dictionary files (`lang.json`), never in Go code.
- **Extraction is automatic**: `langc sync` scans the AST of Go files and extracts English keys into `lang.json`.
- **Translations are manual**: Developers/translators fill in the target language values in `lang.json`.

---

## The Golden Rule — Code is English Only

**Never write Spanish, French, or another non-English language as string literals or keys in Go code.**

```go
// ❌ WRONG — Spanish hardcoded as source key
return fmt.Err("Rol no encontrado")
btn := button.New().SetText("Guardar")

// ✅ CORRECT — English key in code; translated via lang.json
return fmt.Err("Role not found")
btn := button.New().SetText("Save")
```

The Go source string is the canonical **lookup key**. When the application or browser renders the text, `webtyp/lang` looks up this English key in the active dictionary. If no translation exists, the English key passes through unchanged.

### Think in Spanish first (Noun + Adjective Word Order)

English word order is flexible ("Invalid Format" and "Format Invalid" are both understandable). Spanish is not: adjectives follow the noun, and word order dictates natural phrasing.

When composing multi-word keys or variadic arguments in `lang.Translate(...)`, **verify the sequence produces natural Spanish output**:

| Goal | Wrong EN order | Why it breaks in ES | Correct order |
|------|---------------|---------------------|---------------|
| Format error | `"invalid", "format"` | "Inválido Formato" ❌ | `"format", "invalid"` → "Formato Inválido" ✓ |
| Directory not initialized | `"not", "initialized", "directory"` | "No Inicializado Directorio" ❌ | `"directory", "not", "initialized"` → "Directorio No Inicializado" ✓ |
| Project root missing | `"missing", "project", "root"` | "Falta Proyecto Raíz" ❌ | `"project", "root", "missing"` → "Proyecto Raíz Falta" ✓ |

### Gender-Agnostic Phrasing (Lenguaje Agnóstico al Género)

**UI text and user-facing questions must NEVER assume gender or use grammatical masculine.**

Never address the user with masculine inflection in confirmation questions or prompts:

| English Prompt | ❌ Wrong (Masculine) | ✅ Correct (Gender-Agnostic) | Reason |
|----------------|----------------------|------------------------------|--------|
| `"Are you sure?"` | `"¿Estás seguro?"` | `"¿Confirmas?"` / `"¿Deseas continuar?"` | Direct verb; completely gender-neutral. |
| `"Welcome"` | `"Bienvenido"` | `"Te damos la bienvenida"` | Avoids masculine inflection. |
| `"User registered"` | `"Usuario registrado"` | `"Registro completado"` | Focuses on the action, not the person's gender. |

Use verbs, infinitives, or neutral phrasing instead of gendered participles and adjectives.

### Separation of Translatable Text vs Technical Data (api-design)

- **Translatable UI text is typed `lang.Text`**: Fixed labels, buttons, dialog titles, placeholders.
- **Data & format samples stay plain `string`**: Format examples (`"48:F1:7F:D9:D7:B7"`, `"192.168.1.1"`, `"12345678-9"`), user data, identifiers, and configuration values are plain `string` and must never be passed as translatable words. In input widgets, format sample data is set via `SetExample(...)`, while `SetPlaceholder(...)` receives only translatable prompt words (e.g. `"example:"`).
- **Test fixtures and mock models**: Directories `tests/` and `conformance/` in dependencies are skipped by `langc sync` (`skipDirs`) — test fixture labels (e.g., `"Action A"`, `"Action B"`) never leak into an application's production dictionary.
- **Pass-through guarantee**: Unknown words and untranslated technical terms pass through `lang.Translate` unchanged at runtime. They do not belong in `lang.json`.

---

## Where Translations Live

Translations live in two distinct layers:

### 1. Applications (`<root>/config/lang.json`)
Every WebTyp application stores its dictionary in `config/lang.json`. It specifies the `default` fallback language and the list of active `languages`.

### 2. Libraries (`<root>/lang.json`)
Shared modules and libraries (`webtyp/components`, `webtyp/input`, `veltylabs/iam`, etc.) ship their own `lang.json` at their repository module root (without a `"default"` field).

### Runtime Merging ("Project Wins")
When `sitec` builds the web client bundle, `langc.BundleTranslations` merges all dependency libraries' `lang.json` files with the application's `config/lang.json`. The application dictionary overrides library definitions (`project wins`). The merged dictionary is injected into the HTML as:

```html
<script type="application/json" id="webtyp-lang">{ ... }</script>
```

In the browser (Wasm), `webtyp/lang` reads this embedded tag once at startup. A backend server never translates HTML strings directly; translation is client-side or tool-side.

---

## Dictionary Format (`lang.json`)

The dictionary file is strictly positional:

```json
{
  "default": "es",
  "languages": ["es", "fr"],
  "keys": {
    "Delete": ["Eliminar", "Supprimer"],
    "Invalid name": ["Nombre inválido", "Nom invalide"],
    "Send": ["", ""]
  }
}
```

### Positional Correspondence Rule
- `languages`: Array of supported language codes (`"en"`, `"es"`, `"zh"`, `"hi"`, `"ar"`, `"pt"`, `"fr"`, `"de"`, `"ru"`).
- `keys`: Map where the key is the exact **English source string** extracted from Go code.
- Each value is a string array: `value[i]` corresponds strictly to `languages[i]`.
- An empty string `""` denotes an **untranslated** key.
- **Never reorder or remove an item in `languages`** without manually shifting all corresponding elements in `keys`.

---

## How `langc sync` Works

`langc` is the backend CLI and code generator in `webtyp.com/lang/cmd/langc`.

### Commands

```bash
# Sync project or library translations (scans WASM/client code)
langc sync [dir]

# Sync backend CLI tools (scans native backend code and backend fmt.Err)
langc sync -tool [dir]
```

When run:
1. `langc sync` walks the AST of all `.go` files in the module (skipping `*_test.go`, vendor, hidden dirs).
2. It extracts all English keys according to the **8 AST Extraction Rules**.
3. In `lang.json`, it appends newly discovered keys with empty strings `["", ...]` for each configured language.
4. For projects (`config/lang.json`), it deletes unused keys that no longer appear anywhere in the codebase. (For libraries, unused keys are kept because libraries may construct dynamic runtime keys like month names).
5. It formats `lang.json` alphabetically sorted by byte order.
6. It reports untranslated keys to the terminal:
   ```
   lang: 3 untranslated (es) — see config/lang.json
   ```

---

## The 8 AST Extraction Rules

`langc/scan.go` extracts keys matching these 8 AST patterns:

1. **`lang.Translate("...")` and `lang.Text("...")`**: Direct translation calls or type conversions.
2. **`fmt.Err("...")`**: Error formatting strings within the module (in projects, scoped to the main module).
3. **`model.Field` literals (Label & Name)**:
   - If `Label: "..."` is defined, `Label` is extracted.
   - If `Label` is absent, `Name: "..."` is extracted with underscores replaced by spaces (e.g., `client_secret` → `"client secret"`).
4. **`fmt.KeyValue{Value: "..."}`**: Values passed directly to `webtyp.com/input` functions.
5. **`input.SetPlaceholder("...")` & `input.SetTitle("...")`**: Direct placeholder and title setters inside `webtyp.com/input`.
6. **`model.Field` Help text**: `Help: "..."` in field definitions.
7. **Struct literals with `lang.Text` fields**: Any composite literal assigning to a field whose type is `lang.Text` (e.g. `Dialog{Title: "Confirm"}`).
8. **Functions returning `lang.Text`**: String literals returned by any function whose sole return type is `lang.Text`.

---

## `lang.Text` vs `string`

Use `lang.Text` to distinguish static UI copy from dynamic data:

```go
// In a component or layout struct:
type Modal struct {
    Title       lang.Text  // ✅ Fixed UI label: translated in Render()
    Description lang.Text  // ✅ Fixed UI label: translated in Render()
    UserID      string     // Data: never translated
    UserName    string     // Data: never translated
}
```

- Callers still write plain string literals: `Modal{Title: "Delete Project"}`.
- `langc sync` automatically discovers `Title` as a key because of the `lang.Text` type (Rule 7).
- In the component's `Render()`, translate using `lang.Translate(m.Title)`.
- Dynamic runtime data (emails, user names, IDs, timestamps) must remain plain `string` and must **never** pass through `lang.Translate`.

---

## Step-by-Step Translation Workflow

### When adding or editing code:
1. **Write code with English keys**:
   ```go
   if name == "" {
       return fmt.Err("Name cannot be empty")
   }
   ```
2. **Run `langc sync`**:
   ```bash
   langc sync .
   ```
3. **Inspect and edit `lang.json`**:
   Open `config/lang.json` (or `lang.json` in a library), locate the new key:
   ```json
   "Name cannot be empty": ["", ""]
   ```
   Add the translations:
   ```json
   "Name cannot be empty": ["El nombre no puede estar vacío", "Le nom ne peut pas être vide"]
   ```
4. **Verify**:
   Run `gotest`. If any translations were missing, `langc sync` output will indicate the count.

---

## Summary Checklist

- [ ] All string literals in `fmt.Err(...)`, `lang.Translate(...)`, and `lang.Text(...)` are written in **English**.
- [ ] Multi-word phrases in `Translate(...)` follow noun + adjective ordering ("think in Spanish first").
- [ ] `langc sync` was run to synchronize AST keys into `lang.json`.
- [ ] All empty strings `[""]` in `lang.json` have been filled with their corresponding translations.
- [ ] Languages array was not reordered without updating keys.
- [ ] The root package `webtyp.com/lang` was not polluted with `langc` imports (it must compile to WASM).
