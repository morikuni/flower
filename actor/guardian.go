package actor

import (
	"log"
)

var _guardian *guardian

func init() {
	_guardian = &guardian{
		msgChan: make(chan Message),
	}
	_guardian.start()
}

type guardian struct {
	msgChan chan Message
}

func (g *guardian) ActorOf(name string, behavior Behavior) Actor {
	a := newActor(name, behavior, g.Path())
	g.Monitor(a)
	return a
}

func (g *guardian) Shutdown() {
}

func (g *guardian) Path() Path {
	return rootPath
}

func (g *guardian) Send() chan<- Message {
	return g.msgChan
}

func (g *guardian) Monitor(actor Actor) {
	actor.Send() <- notifyMe{g}
}

func (g *guardian) stop() {
}

func (g *guardian) init() {
}

func (g *guardian) start() {
	go func() {
		defer func() {
			panic("Guardian error: dead")
		}()

		for msg := range g.msgChan {
			g.receive(msg)
		}
	}()
}

func (g *guardian) restart(_ interface{}) {
}

func (g *guardian) receive(msg Message) {
	log.Println("Guardian received:", msg)
}
