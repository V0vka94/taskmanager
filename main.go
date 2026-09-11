package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbletea"

	"taskman/internal/app"
	"taskman/internal/storage"
)

func main() {
	home, _ := os.UserHomeDir()
	path := home + "/.tasks.json"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	tasks, err := storage.Load(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Не удалось загрузить задачи: %v\n", err)
	}

	m := app.New(path, tasks)

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка запуска: %v\n", err)
		os.Exit(1)
	}
}
