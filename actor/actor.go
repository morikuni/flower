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
	sys      ActorSystem
	monitors []Actor
	msgChan  chan interface{}
	stopChan chan struct{}
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
	select {
	case actor.stopChan <- struct{}{}:
	default: // default means the Actor has already stopped.
	}
}

func (actor *actor) init() {
	actor.behavior.Init()
}

func (actor *actor) start() {
	go func() {
		defer func() {
			err := recover()
			if err != nil {
				p := paniced{
					actor:  actor,
					reason: err,
				}

				for _, m := range actor.monitors {
					m.Send() <- p
				}
			}
		}()

		for {
			select {
			case msg := <-actor.msgChan:
				actor.receive(msg)
			case <-actor.stopChan:
				break
			}
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

func newActor(name string, sys ActorSystem, behavior Behavior) *actor {
	return &actor{
		path:     sys.Path().join(name),
		behavior: behavior,
		sys:      sys,
		monitors: []Actor{},
		msgChan:  make(chan interface{}),
		stopChan: make(chan struct{}),
	}
}
