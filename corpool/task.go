package corpool

type Tasker interface {
	Process()
}

type FUNC func()

func (f FUNC) Process() {
	f()
}
