package actor

type SupervisorStrategy struct {
	handlePanic func(self *supervisor, target Actor, reason interface{})
}

func oneForOne(_ *supervisor, target Actor, reason interface{}) {
	target.stop()
	target.restart(reason)
}

var OneForOneStrategy = SupervisorStrategy{
	handlePanic: oneForOne,
}

func allForOne(sv *supervisor, target Actor, reason interface{}) {
	for _, a := range sv.children {
		a.stop()
		a.restart(reason)
	}
}

var AllForOneStrategy = SupervisorStrategy{
	handlePanic: allForOne,
}
