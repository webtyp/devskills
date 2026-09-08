package devskills_test

import "webtyp.com/devskills"

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLLM_InstallSkills(t *testing.T) {
	tmpDir := t.TempDir()

	// Simular home directory temporal
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", oldHome)

	llm := devskills.NewLLM()
	if _, err := llm.InstallSkills(); err != nil {
		t.Fatalf("failed to install skills: %v", err)
	}

	// Verificar directorios creados
	skillsRoot := filepath.Join(tmpDir, "skills")
	requiredSkills := []string{
		"core-principles",
		"testing",
		"documentation",
		"wasm",
		"agents-workflow",
		"dev-protocols",
		"devskills",
	}

	for _, skill := range requiredSkills {
		skillPath := filepath.Join(skillsRoot, skill, "SKILL.md")
		if _, err := os.Stat(skillPath); os.IsNotExist(err) {
			t.Errorf("skill file not found: %s", skillPath)
		}
	}
}

func TestLLM_DetectInstalledLLMs(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", oldHome)

	llm := devskills.NewLLM()

	// Caso 1: No hay LLMs instalados
	installed := llm.DetectInstalledLLMs()
	if len(installed) != 0 {
		t.Errorf("expected 0 LLMs, got %d", len(installed))
	}

	// Caso 2: Solo Claude instalado
	claudeDir := filepath.Join(tmpDir, ".claude")
	os.Mkdir(claudeDir, 0755)

	installed = llm.DetectInstalledLLMs()
	if len(installed) != 1 {
		t.Fatalf("expected 1 LLM, got %d", len(installed))
	}
	if installed[0].Name != "claude" {
		t.Errorf("expected 'claude', got '%s'", installed[0].Name)
	}
}

// TestLLM_GetSupportedLLMs_IncludesNewTargets locks in the additional agent
// dirs devskills now syncs to: codex, opencode, qwen, and the vendor-neutral
// ~/.agents convention (opencode and others auto-load skills from there too).
func TestLLM_GetSupportedLLMs_IncludesNewTargets(t *testing.T) {
	llm := devskills.NewLLM()
	names := map[string]bool{}
	for _, cfg := range llm.GetSupportedLLMs() {
		names[cfg.Name] = true
	}

	for _, want := range []string{"claude", "gemini", "codex", "qwen", "opencode", "agents"} {
		if !names[want] {
			t.Errorf("expected %q in GetSupportedLLMs(), got %v", want, names)
		}
	}
}

func TestLLM_Sync(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", oldHome)

	// Simular instalación de LLM
	claudeDir := filepath.Join(tmpDir, ".claude")
	os.Mkdir(claudeDir, 0755)

	llm := devskills.NewLLM()
	summary, err := llm.Sync("", false)
	if err != nil {
		t.Fatalf("Sync failed: %v", err)
	}

	if !strings.Contains(summary, "Config updated") {
		t.Errorf("summary should mention updated config: %s", summary)
	}

	// Verificar que se instalaron skills
	skillsRoot := filepath.Join(tmpDir, "skills")
	if _, err := os.Stat(skillsRoot); os.IsNotExist(err) {
		t.Error("skills dir not created during Sync")
	}

	// El dir de skills del LLM debe ser un directorio real (no un symlink al
	// completo ~/skills), conteniendo un symlink por cada skill individual.
	claudeSkills := filepath.Join(claudeDir, "skills")
	info, err := os.Lstat(claudeSkills)
	if err != nil {
		t.Fatalf("failed to stat claude skills dir: %v", err)
	}
	if info.Mode()&os.ModeSymlink != 0 {
		t.Fatal("expected claude skills dir to be a real directory, not a whole-dir symlink")
	}

	link := filepath.Join(claudeSkills, "core-principles")
	dest, err := os.Readlink(link)
	if err != nil {
		t.Fatalf("failed to read per-skill symlink: %v", err)
	}
	if dest != filepath.Join(skillsRoot, "core-principles") {
		t.Errorf("expected per-skill symlink to %s, got %s", filepath.Join(skillsRoot, "core-principles"), dest)
	}

	// Segunda ejecución: debe skipear
	summary2, _ := llm.Sync("", false)
	if !strings.Contains(summary2, "already had reference") {
		t.Errorf("summary should mention skipped config: %s", summary2)
	}
}

func TestLLM_Sync_SpecificLLM(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", oldHome)

	os.Mkdir(filepath.Join(tmpDir, ".claude"), 0755)
	os.Mkdir(filepath.Join(tmpDir, ".gemini"), 0755)

	llm := devskills.NewLLM()
	summary, _ := llm.Sync("claude", false)

	if !strings.Contains(summary, "claude") {
		t.Error("summary should mention claude")
	}
	if strings.Contains(summary, "gemini") {
		t.Error("summary should NOT mention gemini")
	}

	// Verificar Claude creado, Gemini no
	if _, err := os.Stat(filepath.Join(tmpDir, ".claude", "skills")); os.IsNotExist(err) {
		t.Error("claude skills symlink not created")
	}
	if _, err := os.Stat(filepath.Join(tmpDir, ".gemini", "skills")); err == nil {
		t.Error("gemini skills symlink should not be created")
	}
}

// TestLLM_LinkSkills_PreservesVendorContent verifies that per-skill linking
// leaves vendor-owned entries in an LLM's skills dir untouched, e.g. codex's
// ~/.codex/skills/.system/ bundled skills.
func TestLLM_LinkSkills_PreservesVendorContent(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", oldHome)

	codexDir := filepath.Join(tmpDir, ".codex")
	vendorSkill := filepath.Join(codexDir, "skills", ".system", "skill-creator")
	os.MkdirAll(vendorSkill, 0755)
	os.WriteFile(filepath.Join(vendorSkill, "SKILL.md"), []byte("vendor content"), 0644)

	llm := devskills.NewLLM()
	if _, err := llm.Sync("codex", false); err != nil {
		t.Fatalf("Sync failed: %v", err)
	}

	// Vendor content must survive untouched.
	data, err := os.ReadFile(filepath.Join(vendorSkill, "SKILL.md"))
	if err != nil {
		t.Fatalf("vendor skill file missing after sync: %v", err)
	}
	if string(data) != "vendor content" {
		t.Errorf("vendor skill content was modified: %q", data)
	}

	// Our own skills must also be present, alongside it.
	link := filepath.Join(codexDir, "skills", "core-principles")
	if _, err := os.Lstat(link); err != nil {
		t.Errorf("expected core-principles symlink in codex skills dir: %v", err)
	}
}

// TestLLM_LinkSkills_ReplacesStaleSkillCopy reproduces the real bug that
// motivated per-skill linking: an LLM's skills dir held a real (non-symlink)
// directory for one of our own skills — leftover from an older devskills
// version's copy fallback — with outdated content. Since the name belongs to
// one of our skills, sync must replace it with a fresh symlink rather than
// silently deferring to the stale copy forever.
func TestLLM_LinkSkills_ReplacesStaleSkillCopy(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", oldHome)

	agentsDir := filepath.Join(tmpDir, ".agents")
	staleSkill := filepath.Join(agentsDir, "skills", "core-principles")
	os.MkdirAll(staleSkill, 0755)
	os.WriteFile(filepath.Join(staleSkill, "SKILL.md"), []byte("outdated content"), 0644)

	llm := devskills.NewLLM()
	if _, err := llm.Sync("agents", false); err != nil {
		t.Fatalf("Sync failed: %v", err)
	}

	link := filepath.Join(agentsDir, "skills", "core-principles")
	info, err := os.Lstat(link)
	if err != nil {
		t.Fatalf("failed to stat core-principles after sync: %v", err)
	}
	if info.Mode()&os.ModeSymlink == 0 {
		t.Error("expected stale core-principles copy to be replaced with a symlink")
	}
	dest, err := os.Readlink(link)
	if err != nil {
		t.Fatalf("failed to read symlink: %v", err)
	}
	if dest != filepath.Join(tmpDir, "skills", "core-principles") {
		t.Errorf("expected symlink to current skills source, got %s", dest)
	}
}

// TestLLM_LinkSkills_PrunesStaleSkillLinks verifies that a symlink devskills
// previously created for a skill that no longer exists in the source gets
// removed on the next sync, without touching unrelated entries.
func TestLLM_LinkSkills_PrunesStaleSkillLinks(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", oldHome)

	claudeDir := filepath.Join(tmpDir, ".claude")
	claudeSkills := filepath.Join(claudeDir, "skills")
	os.MkdirAll(claudeSkills, 0755)

	skillsRoot := filepath.Join(tmpDir, "skills")
	os.MkdirAll(skillsRoot, 0755)

	// Simulate a symlink left over from a removed/renamed skill.
	staleTarget := filepath.Join(skillsRoot, "removed-skill")
	os.Symlink(staleTarget, filepath.Join(claudeSkills, "removed-skill"))

	llm := devskills.NewLLM()
	if _, err := llm.Sync("claude", false); err != nil {
		t.Fatalf("Sync failed: %v", err)
	}

	if _, err := os.Lstat(filepath.Join(claudeSkills, "removed-skill")); !os.IsNotExist(err) {
		t.Error("expected stale skill symlink to be pruned")
	}
}

func TestLLM_LinkSkills_Fallback(t *testing.T) {
	tmpDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", oldHome)

	skillsSource := filepath.Join(tmpDir, "skills")
	os.MkdirAll(skillsSource, 0755)

	llmDir := filepath.Join(tmpDir, ".claude")
	os.MkdirAll(llmDir, 0755)

	target := filepath.Join(llmDir, "skills")
	os.MkdirAll(target, 0755) // Directorio preexistente (debe conservarse)
	os.WriteFile(filepath.Join(target, "unrelated.txt"), []byte("keep me"), 0644)

	llm := devskills.NewLLM()

	summary, err := llm.Sync("", false)
	if err != nil {
		t.Fatalf("Sync failed: %v", err)
	}
	if !strings.Contains(summary, "Config updated") {
		t.Errorf("expected config updated, got %s", summary)
	}

	info, err := os.Lstat(target)
	if err != nil {
		t.Fatalf("failed to stat target: %v", err)
	}
	if info.Mode()&os.ModeSymlink != 0 {
		t.Error("expected target to remain a real directory")
	}

	// Unrelated preexisting content must not be removed.
	if _, err := os.Stat(filepath.Join(target, "unrelated.txt")); err != nil {
		t.Errorf("expected unrelated preexisting file to survive: %v", err)
	}
}

func TestCopyFile(t *testing.T) {
	tmpDir := t.TempDir()
	src := filepath.Join(tmpDir, "source.txt")
	dst := filepath.Join(tmpDir, "dest.txt")

	content := "Test content"
	os.WriteFile(src, []byte(content), 0644)

	if err := devskills.CopyFile(src, dst); err != nil {
		t.Fatalf("devskills.CopyFile failed: %v", err)
	}

	dstContent, _ := os.ReadFile(dst)
	if string(dstContent) != content {
		t.Errorf("content mismatch")
	}
}
