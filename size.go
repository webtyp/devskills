package devskills

import (
	"bytes"
	"fmt"
	"io/fs"
	"path"
	"sort"
	"strings"
)

// A skill is a small memory of how one thing is done, not a manual: past
// SkillLinesGood lines it warns, past SkillLinesMax lines Sync refuses to run.
const (
	SkillLinesGood = 50
	SkillLinesMax  = 150
)

const skillFile = "SKILL.md"

// CheckSkills counts the lines of every SKILL.md in fsys. It returns one
// warning per skill longer than SkillLinesGood, and an error naming every
// skill longer than SkillLinesMax.
func CheckSkills(fsys fs.FS) (warnings []string, err error) {
	type size struct {
		name  string
		lines int
	}
	var sizes []size
	walkErr := fs.WalkDir(fsys, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || d.Name() != skillFile {
			return nil
		}
		data, err := fs.ReadFile(fsys, p)
		if err != nil {
			return err
		}
		sizes = append(sizes, size{name: path.Base(path.Dir(p)), lines: countLines(data)})
		return nil
	})
	if walkErr != nil {
		return nil, walkErr
	}
	sort.Slice(sizes, func(i, j int) bool { return sizes[i].lines > sizes[j].lines })

	var over []string
	for _, s := range sizes {
		switch {
		case s.lines > SkillLinesMax:
			over = append(over, fmt.Sprintf("%s (%d)", s.name, s.lines))
		case s.lines > SkillLinesGood:
			warnings = append(warnings, fmt.Sprintf("⚠️ long: %s (%d lines > %d)", s.name, s.lines, SkillLinesGood))
		}
	}
	if len(over) > 0 {
		return warnings, fmt.Errorf("skills over %d lines — split them (one domain per skill, link instead of repeating): %s",
			SkillLinesMax, strings.Join(over, ", "))
	}
	return warnings, nil
}

// countLines counts lines the way wc -l does, plus a last line with no newline.
func countLines(data []byte) int {
	n := bytes.Count(data, []byte("\n"))
	if len(data) > 0 && data[len(data)-1] != '\n' {
		n++
	}
	return n
}
