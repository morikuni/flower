package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/morikuni/worker/actor"
)

// import (
// 	"fmt"
// 	"runtime"
// 	"time"
//
// 	"github.com/morikuni/worker/executor"
// )
//
// func main() {
// 	const n = 100000
// 	const x = 40
// 	runtime.GOMAXPROCS(1)
// 	executor.NewFixedExecutor(n)
//
// 	fmt.Println("wait")
// 	time.Sleep(10 * time.Second)
// }

type Greet struct {
	msg string
}

type Stop struct{}

type Crash struct{}

type Check struct {
	c chan int
}

type CountActor struct {
	c int
}

func (ca *CountActor) Init() {
	ca.c = 0
}

func (ca *CountActor) Receive(self actor.Actor, msg interface{}) {
	switch msg := msg.(type) {
	case Greet:
		ca.c++
	case Check:
		msg.c <- ca.c
	case Stop:
		self.Stop()
	case Crash:
		panic("Crashshs")
	}
}

func main() {
	fmt.Println(runtime.NumGoroutine())
	supervisor := actor.NewSupervisor()
	actor := supervisor.ActorOf(&CountActor{}, "counter")
	c := make(chan int)
	actor.Send(Check{c})
	fmt.Println("count", <-c)
	actor.Send(Greet{"Hello"})
	actor.Send(Greet{"Hello"})
	actor.Send(Check{c})
	fmt.Println("count", <-c)
	fmt.Println(runtime.NumGoroutine())
	actor.Send(Crash{})
	actor.Send(Greet{"Hello"})
	actor.Send(Check{c})
	fmt.Println("count", <-c)
	fmt.Println(runtime.NumGoroutine())
	actor.Send(Stop{})
	time.Sleep(time.Second * 1)
	fmt.Println(runtime.NumGoroutine())
}
