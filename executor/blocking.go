package executor

type blockingExecuter struct{}

func (executor *blockingExecuter) Execute(task func()) {
	task()
}

func (executor *blockingExecuter) Stop() {
}

func NewBlockingExecutor() Executor {
	return &blockingExecuter{}
}
