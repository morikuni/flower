package actor

type Behavior interface {
	Init()
	Receive(self Actor, msg interface{})
}

type Actor interface {
	Name() string
	Send(msg interface{})

	init()
	start()
	stop()
	receive(msg interface{})
}

type actor struct {
	name     string
	behavior Behavior
	parent   Supervisor
	msgChan  chan interface{}
}

func (actor *actor) Send(msg interface{}) {
	actor.msgChan <- msg
}

func (actor *actor) stop() {
	close(actor.msgChan)
}

func (actor *actor) init() {
	actor.behavior.Init()
}

func (actor *actor) start() {
	go func() {
		defer func() {
			err := recover()
			if err != nil {
				actor.parent.Send(paniced{actor})
			}
		}()

		for msg := range actor.msgChan {
			actor.receive(msg)
		}
	}()
}

func (actor *actor) receive(msg interface{}) {
	actor.behavior.Receive(actor, msg)
}

func (actor *actor) Name() string {
	return actor.name
}

func newActor(name string, parent Supervisor, behavior Behavior) *actor {
	c := make(chan interface{})
	return &actor{
		name:     parent.Name() + "/" + name,
		behavior: behavior,
		parent:   parent,
		msgChan:  c,
	}
}
