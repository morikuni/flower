package actor

type Message interface{}

type Panic struct {
	Actor  Actor
	Reason interface{}
}

type Supervise struct {
	Actor Actor
}
