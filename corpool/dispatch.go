package corpool

type Dispatcher struct {
	WorkerPool chan chan Job
	MaxWorkers int
	cache      chan Job
}

func NewDispatcher(maxWorkers int, cacheSize int) *Dispatcher {
	p := make(chan chan Job, maxWorkers)
	return &Dispatcher{
		WorkerPool: p,
		MaxWorkers: maxWorkers,
		cache:      make(chan Job, cacheSize),
	}
}

func (d *Dispatcher) Run() {
	for i := 0; i < d.MaxWorkers; i++ {
		worker := NewWorker(d.WorkerPool)
		worker.Start()
	}
	go d.dispatch()
}

func (d *Dispatcher) dispatch() {
	for {
		select {
		case job := <-d.cache:
			go func(job Job) {
				jobChannel := <-d.WorkerPool
				jobChannel <- job
			}(job)
		}
	}
}

func (d *Dispatcher) RegisteJob(job Job) {
	d.cache <- job
}

func (d *Dispatcher) RegisteTask(f FUNC) {
	d.cache <- GenerateJob(f)
}
