package app

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"

	"taskman/internal/storage"
	"taskman/internal/task"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		if m.Width > 4 {
			m.Input.Width = m.Width - 4
		} else {
			m.Input.Width = 1
		}
		return m, nil

	case tea.KeyMsg:
		if m.Adding {
			return m.handleInput(msg)
		}
		if m.Editing {
			return m.handleEdit(msg)
		}
		return m.handleNormal(msg)
	}

	return m, nil
}

func (m Model) handleInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEnter:
		text := strings.TrimSpace(m.Input.Value())
		if text != "" {
			m.Tasks = append(m.Tasks, task.Task{Text: text})
			_ = storage.Save(m.FilePath, m.Tasks)
		}
		m.Input.Reset()
		m.Adding = false
		return m, nil

	case tea.KeyEsc:
		m.Input.Reset()
		m.Adding = false
		return m, nil
	}

	var cmd tea.Cmd
	m.Input, cmd = m.Input.Update(msg)
	return m, cmd
}

func (m Model) handleEdit(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEnter:
		text := strings.TrimSpace(m.Input.Value())
		if text != "" && m.EditIdx >= 0 && m.EditIdx < len(m.Tasks) {
			m.Tasks[m.EditIdx].Text = text
			_ = storage.Save(m.FilePath, m.Tasks)
		}
		m.Input.Reset()
		m.Editing = false
		m.EditIdx = -1
		return m, nil

	case tea.KeyEsc:
		m.Input.Reset()
		m.Editing = false
		m.EditIdx = -1
		return m, nil
	}

	var cmd tea.Cmd
	m.Input, cmd = m.Input.Update(msg)
	return m, cmd
}

func (m Model) handleNormal(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	visible := m.filteredTasks()

	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "up", "k":
		if m.Cursor > 0 {
			m.Cursor--
		}

	case "down", "j":
		if m.Cursor < len(visible)-1 {
			m.Cursor++
		}

	case "enter", " ":
		if m.Cursor >= 0 && m.Cursor < len(visible) {
			realIdx := m.realIndex(m.Cursor)
			if realIdx >= 0 {
				m.Tasks[realIdx].Done = !m.Tasks[realIdx].Done
				_ = storage.Save(m.FilePath, m.Tasks)
			}
		}

	case "n", "a":
		m.Adding = true
		m.Input.Focus()
		m.Input.Reset()
		return m, textinput.Blink

	case "e":
		if m.Cursor >= 0 && m.Cursor < len(visible) {
			realIdx := m.realIndex(m.Cursor)
			if realIdx >= 0 {
				m.Editing = true
				m.EditIdx = realIdx
				m.Input.Focus()
				m.Input.SetValue(m.Tasks[realIdx].Text)
				return m, textinput.Blink
			}
		}

	case "d", "delete":
		if m.Cursor >= 0 && m.Cursor < len(visible) {
			realIdx := m.realIndex(m.Cursor)
			if realIdx >= 0 {
				m.Tasks = append(m.Tasks[:realIdx], m.Tasks[realIdx+1:]...)
				_ = storage.Save(m.FilePath, m.Tasks)
				if m.Cursor >= len(m.filteredTasks()) && m.Cursor > 0 {
					m.Cursor--
				}
			}
		}

	case "tab":
		m.Filter = (m.Filter + 1) % 3
		m.Cursor = 0
	}

	return m, nil
}
