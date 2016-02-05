package actor

type Behavior interface {
	Init()
	Receive(self Actor, msg interface{})
}

type Actor interface {
	Path() Path
	Send() chan<- interface{}
	Monitor(Actor)

	init()
	start()
	stop()
	receive(msg interface{})
}

type notifyMe struct {
	actor Actor
}

type actor struct {
	path     Path
	behavior Behavior
	monitors []Actor
	msgChan  chan interface{}
}

func (actor *actor) Path() Path {
	return actor.path
}

func (actor *actor) Send() chan<- interface{} {
	return actor.msgChan
}

func (actor *actor) Monitor(target Actor) {
	target.Send() <- notifyMe{actor}
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
				for _, m := range actor.monitors {
					m.Send() <- paniced{actor}
				}
			}
		}()

		for msg := range actor.msgChan {
			actor.receive(msg)
		}
	}()
}

func (actor *actor) receive(msg interface{}) {
	if req, ok := msg.(notifyMe); ok {
		actor.monitors = append(actor.monitors, req.actor)
		return
	}
	actor.behavior.Receive(actor, msg)
}

func newActor(name string, parent Actor, behavior Behavior) *actor {
	c := make(chan interface{})
	return &actor{
		path:     parent.Path().join(name),
		behavior: behavior,
		monitors: []Actor{},
		msgChan:  c,
	}
}
