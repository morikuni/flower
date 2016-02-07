package actor

import (
	"github.com/morikuni/flower/log"
	"sync"
)

type Behavior interface {
	Init()
	Receive(self Actor, msg interface{})
}

type Actor interface {
	Path() Path
	Send() chan<- Message
	Monitor(target Actor)

	init()
	start()
	restart(reason interface{})
	stop()
	receive(msg Message)
}

type notifyMe struct {
	actor Actor
}

type actor struct {
	path        Path
	behavior    Behavior
	monitors    []Actor
	msgChan     chan Message
	stopChan    chan struct{}
	stoppedChan chan struct{}

	mu      sync.Mutex
	running bool
}

func (actor *actor) Path() Path {
	return actor.path
}

func (actor *actor) Send() chan<- Message {
	return actor.msgChan
}

func (actor *actor) Monitor(target Actor) {
	target.Send() <- notifyMe{actor}
}

func (actor *actor) stop() {
	actor.mu.Lock()
	running := actor.running
	actor.mu.Unlock()
	if running {
		actor.stopChan <- struct{}{}
		<-actor.stoppedChan
	}
}

func (actor *actor) init() {
	actor.behavior.Init()
}

func (actor *actor) start() {
	actor.mu.Lock()
	defer actor.mu.Unlock()
	if actor.running {
		return
	}
	actor.running = true
	log.Debug(actor.path, "start")
	go func() {
		defer func() {
			actor.mu.Lock()
			actor.running = false
			log.Debug(actor.path, "stop")
			actor.mu.Unlock()

			err := recover()
			if err != nil {
				p := Panic{
					Actor:  actor,
					Reason: err,
				}

				log.Error(actor.path, "panic")
				for _, m := range actor.monitors {
					m.Send() <- p
				}
			} else {
				actor.stoppedChan <- struct{}{}
			}
		}()

	LOOP:
		for {
			select {
			case <-actor.stopChan:
				break LOOP
			case msg := <-actor.msgChan:
				actor.receive(msg)
			}
		}
	}()
}

func (actor *actor) restart(_ interface{}) {
	actor.init()
	actor.start()
}

func (actor *actor) receive(msg Message) {
	if req, ok := msg.(notifyMe); ok {
		actor.monitors = append(actor.monitors, req.actor)
		log.Info(actor.path, "notify events to", req.actor.Path())
		return
	}
	actor.behavior.Receive(actor, msg)
}

func newActor(name string, behavior Behavior, path Path) *actor {
	a := &actor{
		path:        path.join(name),
		behavior:    behavior,
		monitors:    []Actor{},
		msgChan:     make(chan Message),
		stopChan:    make(chan struct{}),
		stoppedChan: make(chan struct{}),
	}
	a.init()
	a.start()
	return a
}
