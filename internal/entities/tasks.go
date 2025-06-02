package entities

type Func = func()

type Task struct {
	ID       string
	Func     Func
	Priority int
}
