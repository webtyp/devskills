package devskills

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

//go:embed skills
var embeddedSkills embed.FS

// LLMConfig representa la configuración de un LLM específico
type LLMConfig struct {
	Name string // "claude", "gemini", "codex", "opencode", "qwen", "agents"
	Dir  string // "~/.claude", "~/.gemini" — skills are linked at Dir/skills
}

// LLM handles synchronization of LLM configuration files and Agent Skills
type LLM struct {
	log func(...any)
}

// NewLLM creates a new LLM handler
func NewLLM() *LLM {
	return &LLM{
		log: func(...any) {},
	}
}

// SetLog sets the logger function
func (l *LLM) SetLog(fn func(...any)) {
	if fn != nil {
		l.log = fn
	}
}

// GetSupportedLLMs retorna la lista de LLMs soportados
func (l *LLM) GetSupportedLLMs() []LLMConfig {
	home, _ := os.UserHomeDir()
	return []LLMConfig{
		{Name: "claude", Dir: filepath.Join(home, ".claude")},
		{Name: "gemini", Dir: filepath.Join(home, ".gemini")},
		{Name: "codex", Dir: filepath.Join(home, ".codex")},
		{Name: "qwen", Dir: filepath.Join(home, ".qwen")},
		{Name: "opencode", Dir: filepath.Join(home, ".config", "opencode")},
		// ~/.agents/skills is the vendor-neutral convention opencode (and
		// others) also auto-load, independent of any single LLM's own dir.
		{Name: "agents", Dir: filepath.Join(home, ".agents")},
	}
}

// DetectInstalledLLMs detecta qué LLMs están instalados
func (l *LLM) DetectInstalledLLMs() []LLMConfig {
	var installed []LLMConfig
	for _, llm := range l.GetSupportedLLMs() {
		if _, err := os.Stat(llm.Dir); err == nil {
			installed = append(installed, llm)
			l.log("Detected LLM:", llm.Name, "at", llm.Dir)
		}
	}
	return installed
}

// InstallSkills instala los Agent Skills embebidos en ~/skills/
func (l *LLM) InstallSkills() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	destRoot := filepath.Join(home, "skills")
	l.log("Installing skills to:", destRoot)

	err = fs.WalkDir(embeddedSkills, "skills", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Determinar ruta de destino
		relPath, _ := filepath.Rel("skills", path)
		destPath := filepath.Join(destRoot, relPath)

		if d.IsDir() {
			return os.MkdirAll(destPath, 0755)
		}

		// Leer contenido embebido
		data, err := embeddedSkills.ReadFile(path)
		if err != nil {
			return err
		}

		// Escribir archivo (sobrescribir siempre para actualizar)
		return os.WriteFile(destPath, data, 0644)
	})
	return destRoot, err
}

// Sync sincroniza todos los LLMs instalados e instala los skills
func (l *LLM) Sync(specificLLM string, force bool) (string, error) {
	// 1. Instalar skills
	destRoot, err := l.InstallSkills()
	if err != nil {
		return "", fmt.Errorf("failed to install skills: %w", err)
	}

	installed := l.DetectInstalledLLMs()
	if len(installed) == 0 {
		return "⚠️  No LLMs detected (skills installed in ~/skills/)", nil
	}

	// Filtrar por LLM específico si se proporcionó
	if specificLLM != "" {
		var filtered []LLMConfig
		for _, llm := range installed {
			if llm.Name == specificLLM {
				filtered = append(filtered, llm)
				break
			}
		}
		if len(filtered) == 0 {
			return "", fmt.Errorf("LLM '%s' not found or not installed", specificLLM)
		}
		installed = filtered
	}

	var updated []string
	var skipped []string

	for _, llm := range installed {
		changed, err := l.linkSkills(llm.Dir, destRoot)
		if err != nil {
			return "", fmt.Errorf("failed to sync %s: %w", llm.Name, err)
		}

		if changed {
			updated = append(updated, llm.Name)
		} else {
			skipped = append(skipped, llm.Name)
		}
	}

	// Construir resumen
	summary := "Skills updated. "
	if len(updated) > 0 {
		summary += fmt.Sprintf("✅ Config updated: %v", updated)
	}
	if len(skipped) > 0 {
		if len(updated) > 0 {
			summary += ", "
		}
		summary += fmt.Sprintf("⏭️  Config already had reference: %v", skipped)
	}

	return summary, nil
}

// CopyFile copia un archivo (helper para backup)
func CopyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}

// linkSkills links each skill under skillsSource into llmDir/skills individually
// (skillsSource/<name> -> llmDir/skills/<name>), rather than symlinking the whole
// directory. Per-skill linking lets an LLM's skills dir also hold vendor-owned
// skills (e.g. ~/.codex/skills/.system, opencode's bundled skills) without
// devskills clobbering them. Falls back to copying a skill's files if
// symlinking fails (Windows without Developer Mode).
func (l *LLM) linkSkills(llmDir, skillsSource string) (bool, error) {
	target := filepath.Join(llmDir, "skills")

	// Migrate away from the old whole-directory symlink approach.
	if info, err := os.Lstat(target); err == nil && info.Mode()&os.ModeSymlink != 0 {
		os.Remove(target)
	}

	if err := os.MkdirAll(target, 0755); err != nil {
		return false, err
	}

	entries, err := os.ReadDir(skillsSource)
	if err != nil {
		return false, err
	}

	changed := false
	present := make(map[string]bool, len(entries))

	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		name := e.Name()
		present[name] = true
		src := filepath.Join(skillsSource, name)
		link := filepath.Join(target, name)

		if dest, err := os.Readlink(link); err == nil {
			if dest == src {
				continue // already linked
			}
			os.Remove(link) // stale symlink pointing elsewhere; relink below
		} else if _, err := os.Lstat(link); err == nil {
			// A real file/dir already sits at this name. Since the name
			// belongs to one of our own skills, this can only be stale
			// output from an older devskills version (e.g. its pre-symlink
			// copy fallback, or the old whole-directory copy approach) —
			// never a third-party vendor skill, since our skill names
			// (dev-protocols, tinywasm-app, form-codegen, ...) don't
			// collide with any vendor's own. Replace it so the agent picks
			// up the current skill. Names that aren't one of ours are
			// never visited by this loop, so genuine vendor skills living
			// under different names are untouched.
			if err := os.RemoveAll(link); err != nil {
				return changed, err
			}
		}

		if err := os.Symlink(src, link); err != nil {
			if err := copyDir(src, link); err != nil {
				return changed, err
			}
		}
		changed = true
	}

	// Prune symlinks devskills previously created for skills that no longer
	// exist in the source (renamed/removed skill). Never touch anything that
	// isn't our own symlink.
	stale, err := os.ReadDir(target)
	if err != nil {
		return changed, err
	}
	for _, e := range stale {
		name := e.Name()
		if present[name] {
			continue
		}
		link := filepath.Join(target, name)
		dest, err := os.Readlink(link)
		if err != nil || filepath.Dir(dest) != skillsSource {
			continue // not a symlink, or not one of ours
		}
		if err := os.Remove(link); err != nil {
			return changed, err
		}
		changed = true
	}

	return changed, nil
}

func copyDir(src, dst string) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, _ := filepath.Rel(src, path)
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0755)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, 0644)
	})
}
