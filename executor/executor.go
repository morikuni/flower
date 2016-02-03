package executor

type Executor interface {
	Execute(func())
	Stop()
}
