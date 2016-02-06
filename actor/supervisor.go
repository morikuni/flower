package actor

import (
	"log"
)

type supervisor struct {
	handlePanic func(*supervisor, Actor, interface{})
	children    []Actor
}

func (sv *supervisor) Init() {
}

func (sv *supervisor) Receive(self Actor, msg interface{}) {
	switch msg := msg.(type) {
	case Panic:
		log.Println(msg.Actor.Path(), "paniced")
		sv.handlePanic(sv, msg.Actor, msg.Reason)
	case Supervise:
		sv.children = append(sv.children, msg.Actor)
		self.Monitor(msg.Actor)
	default:
		panic("Supervisor error: received unexpected message")
	}
}

func NewSupervisor(strategy SupervisorStrategy) Behavior {
	return &supervisor{
		handlePanic: strategy.handlePanic,
	}
}
