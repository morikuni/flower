package actor

import (
	"log"
)

var _guardian Supervisor = &guardian{
	msgChan: make(chan interface{}),
}

type guardian struct {
	msgChan chan interface{}
}

func (g *guardian) ActorOf(_ Behavior, _ string) Actor {
	return nil
}

func (g *guardian) Shutdown() {
}

func (g *guardian) Send(msg interface{}) {
	g.msgChan <- msg
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

func (g *guardian) Name() string {
	return ""
}
