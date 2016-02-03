package executor

type worker struct {
	taskChan <-chan func()
	OnPanic  func(interface{})
}

func (w *worker) run() {
	defer func() {
		err := recover()
		if err != nil {
			w.OnPanic(err)
		}
	}()

	for task := range w.taskChan {
		task()
	}
}

func newWorker(tc <-chan func(), onPanic func(interface{})) *worker {
	return &worker{
		taskChan: tc,
		OnPanic:  onPanic,
	}
}
