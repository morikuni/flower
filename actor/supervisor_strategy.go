package actor

type SupervisorStrategy struct {
	handlePanic func(self *supervisor, target Actor)
}

func oneForOne(_ *supervisor, target Actor) {
	target.stop()
	target.restart()
}

var OneForOneStrategy = SupervisorStrategy{
	handlePanic: oneForOne,
}

func allForOne(sv *supervisor, target Actor) {
	for _, a := range sv.children {
		a.stop()
		a.restart()
	}
}

var AllForOneStrategy = SupervisorStrategy{
	handlePanic: allForOne,
}
