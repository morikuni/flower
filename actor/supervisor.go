package actor

import (
	"log"
)

type paniced struct {
	actor Actor
}

type createRequest struct {
	name     string
	receiver Receiver
	c        chan Actor
}

type shutdown struct{}

type Supervisor interface {
	Actor
	ActorOf(Receiver, string) Actor
	Shutdown()
}

type supervisor struct {
	Actor
	children []Actor
}

func (sv *supervisor) ActorOf(receiver Receiver, name string) Actor {
	c := make(chan Actor)
	sv.Send(createRequest{
		name:     name,
		receiver: receiver,
		c:        c,
	})
	return <-c
}

func (sv *supervisor) Shutdown() {
	sv.Send(shutdown{})
}

func (sv *supervisor) Init() {
}

func (sv *supervisor) Receive(_ Actor, msg interface{}) {
	switch msg := msg.(type) {
	case paniced:
		log.Println(msg.actor.Name(), "paniced")
		msg.actor.init()
		msg.actor.start()
	case createRequest:
		a := newActor(msg.name, sv, msg.receiver)
		a.init()
		a.start()
		sv.children = append(sv.children, a)
		msg.c <- a
	case shutdown:
		for _, a := range sv.children {
			a.stop()
		}
		sv.stop()
	default:
		panic("Supervisor error: received unexpected message")
	}
}

func NewSupervisor() Supervisor {
	sv := &supervisor{}
	sv.Actor = newActor("supervisor", _guardian, sv)
	sv.start()
	return sv
}
