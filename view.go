package app

import (
	"fmt"
	"strings"

	"taskman/internal/ui"
)

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	runes := []rune(s)
	if max < 4 {
		max = 4
	}
	return string(runes[:max-3]) + "…"
}

func (m Model) View() string {
	if m.Adding || m.Editing {
		return m.renderInput()
	}
	return m.renderList()
}

func (m Model) renderInput() string {
	label := "Новая задача"
	if m.Editing {
		label = "Редактирование"
	}

	header := ui.TitleStyle.Render("📋 Менеджер задач")
	prompt := ui.InputPromptStyle.Render(fmt.Sprintf("%s: ", label))
	inputBox := ui.BorderStyle.Render(m.Input.View())
	hint := ui.HelpStyle.Render("Enter — сохранить · Esc — отмена")

	return fmt.Sprintf("%s\n\n%s\n%s\n\n%s", header, prompt, inputBox, hint)
}

func (m Model) renderList() string {
	var b strings.Builder

	b.WriteString(ui.TitleStyle.Render("📋 Менеджер задач"))
	b.WriteString("\n\n")

	filterLine := ui.FilterStyle.Render(fmt.Sprintf("Фильтр: %s (Tab для смены)", m.Filter.String()))

	total := len(m.Tasks)
	done := 0
	for _, t := range m.Tasks {
		if t.Done {
			done++
		}
	}
	stats := ui.HelpStyle.Render(fmt.Sprintf("Всего: %d · Выполнено: %d · Активных: %d", total, done, total-done))

	b.WriteString(fmt.Sprintf("%s  %s\n\n", filterLine, stats))

	visible := m.filteredTasks()

	maxLines := m.Height - 12
	if maxLines < 5 {
		maxLines = 5
	}

	if len(visible) == 0 {
		b.WriteString(ui.HelpStyle.Render("Список пуст. Нажми 'n' чтобы добавить задачу."))
		b.WriteString("\n")
	} else {
		showCount := len(visible)
		truncated := false
		if len(visible) > maxLines {
			showCount = maxLines
			truncated = true
		}

		maxTaskWidth := m.Width - 8
		if maxTaskWidth < 10 {
			maxTaskWidth = 10
		}

		for i := 0; i < showCount; i++ {
			t := visible[i]
			cursor := "  "
			if i == m.Cursor {
				cursor = ui.CursorStyle.Render("▶ ")
			}

			text := truncate(t.Text, maxTaskWidth)
			var line string
			if t.Done {
				line = ui.DoneStyle.Render(fmt.Sprintf("✔ %s", text))
			} else {
				line = ui.ActiveStyle.Render(fmt.Sprintf("○ %s", text))
			}

			b.WriteString(fmt.Sprintf("%s%s\n", cursor, line))
		}

		if truncated {
			hint := ui.HelpStyle.Render(fmt.Sprintf("... и ещё %d задач. Расширьте окно.", len(visible)-maxLines))
			b.WriteString(fmt.Sprintf("\n%s", hint))
		}
	}

	b.WriteString("\n")
	b.WriteString(ui.HelpStyle.Render(
		"n — добавить · e — редактировать · d — удалить · Space — отметить · Tab — фильтр · q — выход",
	))

	return b.String()
}
