package actor

type Behavior interface {
	Init()
	Receive(self Actor, msg interface{})
}

type Actor interface {
	Path() Path
	Send() chan<- interface{}

	init()
	start()
	stop()
	receive(msg interface{})
}

type actor struct {
	path     Path
	behavior Behavior
	parent   Supervisor
	msgChan  chan interface{}
}

func (actor *actor) Send() chan<- interface{} {
	return actor.msgChan
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
				actor.parent.Send() <- paniced{actor}
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

func (actor *actor) Path() Path {
	return actor.path
}

func newActor(name string, parent Supervisor, behavior Behavior) *actor {
	c := make(chan interface{})
	return &actor{
		path:     parent.Path().join(name),
		behavior: behavior,
		parent:   parent,
		msgChan:  c,
	}
}
