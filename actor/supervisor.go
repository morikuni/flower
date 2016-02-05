package actor

import (
	"log"
)

type paniced struct {
	actor Actor
}

type createRequest struct {
	name     string
	behavior Behavior
	c        chan Actor
}

type shutdown struct{}

type Supervisor interface {
	ActorOf(Behavior, string) Actor
	Shutdown()
}

type supervisor struct {
	Actor
	children []Actor
}

func (sv *supervisor) ActorOf(behavior Behavior, name string) Actor {
	c := make(chan Actor)
	sv.Send() <- createRequest{
		name:     name,
		behavior: behavior,
		c:        c,
	}
	a := <-c
	sv.Monitor(a)
	return a
}

func (sv *supervisor) Shutdown() {
	sv.Send() <- shutdown{}
}

func (sv *supervisor) Init() {
}

func (sv *supervisor) Receive(_ Actor, msg interface{}) {
	switch msg := msg.(type) {
	case paniced:
		log.Println(msg.actor.Path(), "paniced")
		msg.actor.init()
		msg.actor.start()
	case createRequest:
		a := newActor(msg.name, sv, msg.behavior)
		a.init()
		a.start()
		sv.children = append(sv.children, a)
		msg.c <- a
	case shutdown:
		for _, a := range sv.children {
			a.stop()
		}
		sv.children = sv.children[:0]
		sv.stop()
	default:
		panic("Supervisor error: received unexpected message")
	}
}

func NewSupervisor(name string) Supervisor {
	sv := &supervisor{}
	sv.Actor = newActor(name, _guardian, sv)
	sv.start()
	return sv
}
