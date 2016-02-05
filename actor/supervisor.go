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
	sv.Actor = sys.ActorOf(sv, name)
	sv.start()
	return sv
}

type SupervisorStrategy struct {
	onPanic func(self Supervisor, target Actor)
}

func oneForOne(_ Supervisor, target Actor) {
	target.stop()
	target.init()
	target.start()
}

var OneForOneStrategy = SupervisorStrategy{
	onPanic: oneForOne,
}

func allForOne(sv Supervisor, target Actor) {
	for _, a := range sv.Children() {
		a.stop()
		a.init()
		a.start()
	}
}

var AllForOneStrategy = SupervisorStrategy{
	onPanic: allForOne,
}
