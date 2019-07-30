package watcher

type OpType int

const (
	_ OpType = iota
	Create
	Update
	Delete
)

type Event struct {
	Op   OpType
	Path string
}

type Watcher interface {
	AddWatch(path string, fn EventHandler) (err error)
	RemoveWatch(path string)
}

type EventHandler func(event Event) error
