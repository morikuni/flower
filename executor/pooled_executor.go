package executor

import (
	"log"
)

type pooledExecutor struct {
	taskChan chan<- func()
}

func (executor *pooledExecutor) Execute(task func()) {
	executor.taskChan <- task
}

func (executor *pooledExecutor) Stop() {
	close(executor.taskChan)
}

func undeadWorker(tc <-chan func()) *worker {
	return newWorker(tc, func(err interface{}) {
		log.Println(err)
		w := undeadWorker(tc)
		go w.run()
	})
}

func NewPooledExecutor(n uint) Executor {
	tc := make(chan func())

	for i := uint(0); i < n; i++ {
		w := undeadWorker(tc)
		go w.run()
	}

	return &pooledExecutor{tc}
}
