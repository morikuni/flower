package actor

type Behavior interface {
	OnRestart(reason interface{})
	OnReceive(self Actor, msg Message)
}
