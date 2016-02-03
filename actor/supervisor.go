package actor

import (
	"log"
)

type Supervisor interface {
	ActorOf(Receiver, string) Actor

	onCrashChild(*actor, interface{})
}

type supervisor struct {
	children []Actor
}

func (sv *supervisor) onCrashChild(a *actor, err interface{}) {
	log.Println("[", a.name, "]", err)
	a.init()
	a.start()
}

func (sv *supervisor) ActorOf(receiver Receiver, name string) Actor {
	receiver.Init()
	a := newActor(name, sv, receiver)
	a.start()
	sv.add(a)
	return a
}

func (sv *supervisor) add(a *actor) {
	sv.children = append(sv.children, a)
}

func NewSupervisor() Supervisor {
	return &supervisor{}
}
