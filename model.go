package app

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"

	"taskman/internal/task"
)

type Model struct {
	Tasks    []task.Task
	Cursor   int
	Filter   task.FilterMode
	Input    textinput.Model
	Adding   bool
	Editing  bool
	EditIdx  int
	Width    int
	Height   int
	FilePath string
}

func New(filePath string, tasks []task.Task) Model {
	ti := textinput.New()
	ti.Placeholder = "Купить спаржу..."
	ti.Focus()
	ti.CharLimit = 200
	ti.Width = 50

	return Model{
		Tasks:    tasks,
		Input:    ti,
		FilePath: filePath,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) filteredTasks() []task.Task {
	var result []task.Task
	for _, t := range m.Tasks {
		switch m.Filter {
		case task.FilterAll:
			result = append(result, t)
		case task.FilterActive:
			if !t.Done {
				result = append(result, t)
			}
		case task.FilterDone:
			if t.Done {
				result = append(result, t)
			}
		}
	}
	return result
}

func (m Model) realIndex(filteredIdx int) int {
	count := 0
	for i, t := range m.Tasks {
		matches := false
		switch m.Filter {
		case task.FilterAll:
			matches = true
		case task.FilterActive:
			matches = !t.Done
		case task.FilterDone:
			matches = t.Done
		}
		if matches {
			if count == filteredIdx {
				return i
			}
			count++
		}
	}
	return -1
}
