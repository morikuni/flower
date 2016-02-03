package actor

type Receiver interface {
	Init()
	Receive(self Actor, msg interface{})
}

type Actor interface {
	Supervisor
	Send(msg interface{})
	Stop()

	init()
	start()
}

type actor struct {
	name       string
	receiver   Receiver
	parent     Supervisor
	supervisor Supervisor
	msgChan    chan interface{}
}

func (actor *actor) Send(msg interface{}) {
	actor.msgChan <- msg
}

func (actor *actor) Stop() {
	close(actor.msgChan)
}

func (actor *actor) ActorOf(receiver Receiver, name string) Actor {
	return actor.supervisor.ActorOf(receiver, name)
}

func (actor *actor) onCrashChild(child *actor, err interface{}) {
	actor.supervisor.onCrashChild(child, err)
}

func (actor *actor) init() {
	actor.receiver.Init()
}

func (actor *actor) start() {
	go func() {
		defer func() {
			err := recover()
			if err != nil {
				actor.parent.onCrashChild(actor, err)
			}
		}()

		for msg := range actor.msgChan {
			actor.receiver.Receive(actor, msg)
		}
	}()
}

func newActor(name string, parent Supervisor, receiver Receiver) *actor {
	c := make(chan interface{})
	return &actor{
		name:       name,
		receiver:   receiver,
		parent:     parent,
		supervisor: &supervisor{},
		msgChan:    c,
	}
}
