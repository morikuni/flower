package actor

import (
	"log"
)

var _guardian *guardian

func init() {
	_guardian = &guardian{
		msgChan: make(chan interface{}),
	}
	_guardian.start()
}

type guardian struct {
	msgChan chan interface{}
}

func (g *guardian) ActorOf(_ Behavior, _ string) Actor {
	return nil
}

func (g *guardian) Shutdown() {
}

func (g *guardian) Path() Path {
	return rootPath
}

func (g *guardian) Send() chan<- interface{} {
	return g.msgChan
}

func (g *guardian) Monitor(_ Actor) {
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

func (g *guardian) receive(msg interface{}) {
	log.Println("Guardian received:", msg)
}
