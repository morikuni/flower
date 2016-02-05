package actor

type createRequest struct {
	name     string
	behavior Behavior
	c        chan Actor
}

type shutdown struct{}

type ActorSystem interface {
	Path() Path
	ActorOf(Behavior, string) Actor
	Shutdown()
}

type actorSystem struct {
	Actor
	actors []Actor
}

func (sys *actorSystem) ActorOf(behavior Behavior, name string) Actor {
	c := make(chan Actor)
	sys.Send() <- createRequest{
		name:     name,
		behavior: behavior,
		c:        c,
	}
	a := <-c
	return a
}

func (sys *actorSystem) Shutdown() {
	sys.Send() <- shutdown{}
}

func (sys *actorSystem) Init() {
}

func (sys *actorSystem) Receive(_ Actor, msg interface{}) {
	switch msg := msg.(type) {
	case createRequest:
		a := newActor(msg.name, sys, msg.behavior)
		a.init()
		a.start()
		sys.actors = append(sys.actors, a)
		msg.c <- a
	case shutdown:
		for _, a := range sys.actors {
			a.stop()
		}
		sys.actors = sys.actors[:0]
		sys.stop()
	default:
		panic("ActorSystem error: received unexpected message")
	}
}

func NewActorSystem(name string) ActorSystem {
	sys := &actorSystem{}
	sys.Actor = newActor(name, _guardian, sys)
	sys.init()
	sys.start()
	_guardian.Monitor(sys)
	return sys
}
