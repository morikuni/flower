package actor

type SupervisorStrategy struct {
	handlePanic func(self *supervisor, target Actor)
}

func oneForOne(_ *supervisor, target Actor) {
	target.stop()
	target.init()
	target.start()
}

var OneForOneStrategy = SupervisorStrategy{
	handlePanic: oneForOne,
}

func allForOne(sv *supervisor, target Actor) {
	for _, a := range sv.children {
		a.stop()
		a.init()
		a.start()
	}
}

var AllForOneStrategy = SupervisorStrategy{
	handlePanic: allForOne,
}
