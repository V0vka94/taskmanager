package task

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
