package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

// ─── Модель данных ───────────────────────────────────────────

type Task struct {
	Text string `json:"text"`
	Done bool   `json:"done"`
}

type FilterMode int

const (
	FilterAll FilterMode = iota
	FilterActive
	FilterDone
)

func (f FilterMode) String() string {
	return [...]string{"Все", "Активные", "Выполненные"}[f]
}

// ─── Стили (Lipgloss) ─────────────────────────────────────────

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FF79C6")).
			Padding(0, 1)

	doneStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#50FA7B")).
			Strikethrough(true)

	activeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F8F8F2"))

	cursorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFB86C")).
			Bold(true)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6272A4"))

	filterStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8BE9FD"))

	inputPromptStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#BD93F9"))
)

// ─── Модель приложения ───────────────────────────────────────

type model struct {
	tasks    []Task
	cursor   int
	filter   FilterMode
	input    textinput.Model
	adding   bool
	editing  bool
	editIdx  int
	width    int
	height   int
	filePath string
}

func initialModel(path string) model {
	ti := textinput.New()
	ti.Placeholder = "Оформить документ..."
	ti.Focus()
	ti.CharLimit = 200
	ti.Width = 50

	m := model{
		tasks:    []Task{},
		input:    ti,
		filePath: path,
	}
	m.load()
	return m
}

func (m *model) filteredTasks() []Task {
	var result []Task
	for _, t := range m.tasks {
		switch m.filter {
		case FilterAll:
			result = append(result, t)
		case FilterActive:
			if !t.Done {
				result = append(result, t)
			}
		case FilterDone:
			if t.Done {
				result = append(result, t)
			}
		}
	}
	return result
}

// ─── Сохранение/загрузка ──────────────────────────────────────

func (m *model) save() {
	data, err := json.MarshalIndent(m.tasks, "", "  ")
	if err != nil {
		return
	}
	_ = os.WriteFile(m.filePath, data, 0644)
}

func (m *model) load() {
	data, err := os.ReadFile(m.filePath)
	if err != nil {
		return
	}
	_ = json.Unmarshal(data, &m.tasks)
}

// ─── Init ────────────────────────────────────────────────────

func (m model) Init() tea.Cmd {
	return nil
}

// ─── Update ──────────────────────────────────────────────────

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// обновляем ширину input: оставляем по 2 символа на отступы
		if m.width > 4 {
			m.input.Width = m.width - 4
		} else {
			m.input.Width = 1 // минимальная ширина
		}
	return m, nil

	case tea.KeyMsg:
		// Режим ввода новой задачи
		if m.adding {
			return m.handleInput(msg)
		}
		// Режим редактирования
		if m.editing {
			return m.handleEdit(msg)
		}
		// Обычный режим
		return m.handleNormal(msg)
	}

	return m, nil
}

func (m model) handleInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEnter:
		text := strings.TrimSpace(m.input.Value())
		if text != "" {
			m.tasks = append(m.tasks, Task{Text: text})
			m.save()
		}
		m.input.Reset()
		m.adding = false
		return m, nil

	case tea.KeyEsc:
		m.input.Reset()
		m.adding = false
		return m, nil
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m model) handleEdit(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEnter:
		text := strings.TrimSpace(m.input.Value())
		if text != "" && m.editIdx >= 0 && m.editIdx < len(m.tasks) {
			m.tasks[m.editIdx].Text = text
			m.save()
		}
		m.input.Reset()
		m.editing = false
		m.editIdx = -1
		return m, nil

	case tea.KeyEsc:
		m.input.Reset()
		m.editing = false
		m.editIdx = -1
		return m, nil
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m model) handleNormal(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	visible := m.filteredTasks()

	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit

	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}

	case "down", "j":
		if m.cursor < len(visible)-1 {
			m.cursor++
		}

	case "enter", " ":
		if m.cursor >= 0 && m.cursor < len(visible) {
			// Находим реальный индекс в m.tasks
			realIdx := m.realIndex(m.cursor)
			if realIdx >= 0 {
				m.tasks[realIdx].Done = !m.tasks[realIdx].Done
				m.save()
			}
		}

	case "n", "a":
		m.adding = true
		m.input.Focus()
		m.input.Reset()
		return m, textinput.Blink

	case "e":
		if m.cursor >= 0 && m.cursor < len(visible) {
			realIdx := m.realIndex(m.cursor)
			if realIdx >= 0 {
				m.editing = true
				m.editIdx = realIdx
				m.input.Focus()
				m.input.SetValue(m.tasks[realIdx].Text)
				return m, textinput.Blink
			}
		}

	case "d", "delete":
		if m.cursor >= 0 && m.cursor < len(visible) {
			realIdx := m.realIndex(m.cursor)
			if realIdx >= 0 {
				m.tasks = append(m.tasks[:realIdx], m.tasks[realIdx+1:]...)
				m.save()
				if m.cursor >= len(m.filteredTasks()) && m.cursor > 0 {
					m.cursor--
				}
			}
		}

	case "tab":
		m.filter = (m.filter + 1) % 3
		m.cursor = 0
	}

	return m, nil
}

// realIndex переводит индекс в отфильтрованном списке в индекс в m.tasks
func (m model) realIndex(filteredIdx int) int {
	count := 0
	for i, t := range m.tasks {
		matches := false
		switch m.filter {
		case FilterAll:
			matches = true
		case FilterActive:
			matches = !t.Done
		case FilterDone:
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

// ─── View ────────────────────────────────────────────────────

func (m model) View() string {
	if m.adding || m.editing {
		return m.renderInput()
	}
	return m.renderList()
}
var borderStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("#8BE9FD")). // цвет рамки (как фильтр)
	Padding(1, 2)                               // отступ внутри рамки

func (m model) renderInput() string {
	label := "Новая задача"
	if m.editing {
		label = "Редактирование"
	}

	header := titleStyle.Render("📋 Менеджер задач")
	prompt := inputPromptStyle.Render(fmt.Sprintf("%s: ", label))

	// textinput сам рисует себя, но без рамки. Мы оборачиваем его в borderStyle
	inputBox := borderStyle.Render(m.input.View())

	hint := helpStyle.Render("Enter — сохранить · Esc — отмена")

	return fmt.Sprintf("%s\n\n%s\n%s\n\n%s", header, prompt, inputBox, hint)
}

func (m model) renderList() string {
	var b strings.Builder

	// Заголовок
	b.WriteString(titleStyle.Render("📋 Менеджер задач"))
	b.WriteString("\n\n")

	// Фильтр
	filterLine := filterStyle.Render(fmt.Sprintf("Фильтр: %s (Tab для смены)", m.filter.String()))

	// Подсчёт
	total := len(m.tasks)
	done := 0
	for _, t := range m.tasks {
		if t.Done {
			done++
		}
	}
	stats := helpStyle.Render(fmt.Sprintf("Всего: %d · Выполнено: %d · Активных: %d", total, done, total-done))

	b.WriteString(fmt.Sprintf("%s  %s\n\n", filterLine, stats))

	// Список задач
	visible := m.filteredTasks()

	if len(visible) == 0 {
		b.WriteString(helpStyle.Render("  Список пуст. Нажми 'n' чтобы добавить задачу."))
		b.WriteString("\n")
	} else {
		for i, t := range visible {
			cursor := "  "
			if i == m.cursor {
				cursor = cursorStyle.Render("▶ ")
			}

			var text string
			if t.Done {
				text = doneStyle.Render(fmt.Sprintf("✔ %s", t.Text))
			} else {
				text = activeStyle.Render(fmt.Sprintf("○ %s", t.Text))
			}

			b.WriteString(fmt.Sprintf("%s%s\n", cursor, text))
		}
	}

	// Подсказки
	b.WriteString("\n")
	b.WriteString(helpStyle.Render(
		"n — добавить · e — редактировать · d — удалить · Space — отметить · Tab — фильтр · q — выход",
	))

	return b.String()
}

// ─── main ────────────────────────────────────────────────────

func main() {
	home, _ := os.UserHomeDir()
	path := home + "/.tasks.json"

	// Если передан аргумент — используем его как путь к файлу
	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	p := tea.NewProgram(initialModel(path), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Ошибка: %v\n", err)
		os.Exit(1)
	}
}

