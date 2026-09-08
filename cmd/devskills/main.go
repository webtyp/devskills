package main

import (
	"flag"
	"fmt"
	"os"

	"webtyp.com/devskills"
)

func main() {
	fs := flag.NewFlagSet("devskills", flag.ExitOnError)

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, `devskills - Sync WebTyp skills to LLMs

Usage:
    devskills              Sync all installed LLMs
    devskills -l claude    Sync only Claude
    devskills -f           Force overwrite (with backup)
    devskills -h           Show this help

Detects LLMs by directory:
    ~/.claude, ~/.gemini, ~/.codex, ~/.qwen, ~/.config/opencode, ~/.agents

Skills source: webtyp.com/devskills/skills
Installed to: ~/skills, then symlinked per-skill into each detected
dir's skills/ subfolder (~/.claude/skills/<name>, ~/.codex/skills/<name>, etc.)
`)
	}

	llmFlag := fs.String("l", "", "Sync specific LLM (claude, gemini, codex, qwen, opencode, agents)")
	fs.StringVar(llmFlag, "llm", "", "Sync specific LLM (alias)")
	forceFlag := fs.Bool("f", false, "Force overwrite with backup")
	fs.BoolVar(forceFlag, "force", false, "Force overwrite with backup (alias)")
	helpFlag := fs.Bool("h", false, "Show help")
	fs.BoolVar(helpFlag, "help", false, "Show help")

	fs.Parse(os.Args[1:])

	if *helpFlag {
		fs.Usage()
		os.Exit(0)
	}

	llm := devskills.NewLLM()

	summary, err := llm.Sync(*llmFlag, *forceFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if summary != "" {
		fmt.Println(summary)
	}
}
