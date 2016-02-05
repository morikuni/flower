package actor

type SupervisorStrategy struct {
	onPanic func(self Supervisor, target Actor)
}

func oneForOne(_ Supervisor, target Actor) {
	target.stop()
	target.init()
	target.start()
}

var OneForOneStrategy = SupervisorStrategy{
	onPanic: oneForOne,
}

func allForOne(sv Supervisor, target Actor) {
	for _, a := range sv.Children() {
		a.stop()
		a.init()
		a.start()
	}
}

var AllForOneStrategy = SupervisorStrategy{
	onPanic: allForOne,
}
