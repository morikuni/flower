package actor

import (
	"sync"
)

type Message interface{}

type Panic struct {
	Actor  Actor
	Reason interface{}
}

type Supervise struct {
	*sync.WaitGroup
	Actors []Actor
}
