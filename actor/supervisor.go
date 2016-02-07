package actor

import (
	"sync"
)

type supervisor struct {
	handlePanic func(*supervisor, Actor, interface{})
	children    []Actor
}

func (sv *supervisor) OnRestart(_ interface{}) {
}

func (sv *supervisor) Receive(self Actor, msg Message) {
	switch msg := msg.(type) {
	case Panic:
		sv.handlePanic(sv, msg.Actor, msg.Reason)
	case Supervise:
		sv.children = append(sv.children, msg.Actors...)
		for _, a := range msg.Actors {
			self.Monitor(a)
		}
		msg.Done()
	default:
		panic("Supervisor error: received unexpected message")
	}
}

func NewSupervisor(strategy SupervisorStrategy) Behavior {
	return &supervisor{
		handlePanic: strategy.handlePanic,
	}
}

func SuperviseSync(supervisor Actor, children []Actor) {
	var wg sync.WaitGroup
	wg.Add(1)
	supervisor.Send() <- Supervise{&wg, children}
	wg.Wait()
}
