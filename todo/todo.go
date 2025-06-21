package todo

type Task struct {
	ID   int
	Text string
	Done bool
}

// List holds tasks in memory.
type List struct {
	tasks []Task
}

// NewList creates an empty List.
func NewList() *List {
	return &List{}
}

// Add a task to the list and return its ID.
func (l *List) Add(text string) int {
	id := len(l.tasks) + 1
	l.tasks = append(l.tasks, Task{ID: id, Text: text})
	return id
}

// Complete marks a task as done.
func (l *List) Complete(id int) bool {
	for i, t := range l.tasks {
		if t.ID == id {
			l.tasks[i].Done = true
			return true
		}
	}
	return false
}

// Remove deletes a task from the list.
func (l *List) Remove(id int) bool {
	for i, t := range l.tasks {
		if t.ID == id {
			l.tasks = append(l.tasks[:i], l.tasks[i+1:]...)
			return true
		}
	}
	return false
}

// Tasks returns a copy of tasks.
func (l *List) Tasks() []Task {
	cp := make([]Task, len(l.tasks))
	copy(cp, l.tasks)
	return cp
}
