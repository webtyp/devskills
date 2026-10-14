package devskills_test

import (
	"strings"
	"testing"
	"testing/fstest"

	"webtyp.com/devskills"
)

func skillOf(lines int) *fstest.MapFile {
	return &fstest.MapFile{Data: []byte(strings.Repeat("x\n", lines))}
}

func TestCheckSkills(t *testing.T) {
	fsys := fstest.MapFS{
		"skills/good/SKILL.md":  skillOf(devskills.SkillLinesGood),
		"skills/long/SKILL.md":  skillOf(devskills.SkillLinesGood + 1),
		"skills/limit/SKILL.md": skillOf(devskills.SkillLinesMax),
		"skills/good/notes.md":  skillOf(devskills.SkillLinesMax + 1), // not a SKILL.md
	}
	warnings, err := devskills.CheckSkills(fsys)
	if err != nil {
		t.Fatalf("no skill is over the max, got error: %v", err)
	}
	if len(warnings) != 2 {
		t.Fatalf("want 2 warnings (limit, long), got %v", warnings)
	}
	if !strings.Contains(warnings[0], "limit") || !strings.Contains(warnings[1], "long") {
		t.Errorf("warnings must name the skill, longest first: %v", warnings)
	}
}

func TestCheckSkills_OverMaxFails(t *testing.T) {
	fsys := fstest.MapFS{
		"skills/huge/SKILL.md": skillOf(devskills.SkillLinesMax + 1),
		"skills/ok/SKILL.md":   skillOf(10),
	}
	_, err := devskills.CheckSkills(fsys)
	if err == nil || !strings.Contains(err.Error(), "huge (151)") {
		t.Fatalf("want error naming huge (151), got %v", err)
	}
}

func TestCheckSkills_NoTrailingNewline(t *testing.T) {
	fsys := fstest.MapFS{"skills/a/SKILL.md": {Data: []byte(strings.Repeat("x\n", devskills.SkillLinesGood) + "x")}}
	warnings, _ := devskills.CheckSkills(fsys)
	if len(warnings) != 1 {
		t.Fatalf("a last line without newline still counts: %v", warnings)
	}
}
