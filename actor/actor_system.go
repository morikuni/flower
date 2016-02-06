package actor

type createRequest struct {
	name     string
	behavior Behavior
	c        chan Actor
}

type shutdown struct{}

type ActorSystem interface {
	Path() Path
	ActorOf(name string, behavior Behavior) Actor
	Shutdown()
}

type actorSystem struct {
	Actor
	actors []Actor
}

func (sys *actorSystem) ActorOf(name string, behavior Behavior) Actor {
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

func (sys *actorSystem) Receive(self Actor, msg interface{}) {
	switch msg := msg.(type) {
	case createRequest:
		a := newActor(msg.name, msg.behavior, sys.Path())
		sys.actors = append(sys.actors, a)
		msg.c <- a
	case shutdown:
		for _, a := range sys.actors {
			a.stop()
		}
		sys.actors = sys.actors[:0]
		go self.stop() // need goroutine to handle stopChan
	default:
		panic("ActorSystem error: received unexpected message")
	}
}

func NewActorSystem(name string) ActorSystem {
	sys := &actorSystem{}
	sys.Actor = _guardian.ActorOf(name, sys)
	return sys
}
