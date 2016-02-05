package actor

import (
	"log"
)

type paniced struct {
	actor  Actor
	reason interface{}
}

type superviseMe struct {
	actor Actor
}

type Supervisor interface {
	Supervise(Actor)
	Children() []Actor
}

type supervisor struct {
	Actor
	onPanic  func(Supervisor, Actor)
	children []Actor
}

func (sv *supervisor) Supervise(target Actor) {
	sv.Send() <- superviseMe{target}
}

func (sv *supervisor) Children() []Actor {
	return sv.children
}

func (sv *supervisor) Init() {
}

func (sv *supervisor) Receive(_ Actor, msg interface{}) {
	switch msg := msg.(type) {
	case paniced:
		log.Println(msg.actor.Path(), "paniced")
		sv.onPanic(sv, msg.actor)
	case superviseMe:
		sv.children = append(sv.children, msg.actor)
		sv.Monitor(msg.actor)
	default:
		panic("Supervisor error: received unexpected message")
	}
}

func NewSupervisor(name string, strategy SupervisorStrategy, sys ActorSystem) Supervisor {
	sv := &supervisor{
		onPanic: strategy.onPanic,
	}
	sv.Actor = sys.ActorOf(name, sv)
	return sv
}
