package actor

import (
	"github.com/morikuni/flower/log"
	"sync"
)

type Behavior interface {
	OnRestart(reason interface{})
	Receive(self Actor, msg Message)
}

type Actor interface {
	Path() Path
	Send() chan<- Message
	Monitor(target Actor)

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

func (actor *actor) restart(reason interface{}) {
	actor.behavior.OnRestart(reason)
	actor.start()
}

func (actor *actor) receive(msg Message) {
	switch msg := msg.(type) {
	case notifyMe:
		actor.monitors = append(actor.monitors, msg.actor)
		log.Info(actor.path, "notify events to", msg.actor.Path())
	default:
		actor.behavior.Receive(actor, msg)
	}
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
	a.start()
	return a
}
