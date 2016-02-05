package actor

type Panic struct {
	Actor  Actor
	Reason interface{}
}

type Supervise struct {
	Actor Actor
}
