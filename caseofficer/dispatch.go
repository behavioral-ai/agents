package caseofficer

import "github.com/advanced-go/common/messaging"

type dispatcher interface {
	setup(agent *caseOfficer, channel, event string)
	dispatch(agent *caseOfficer, channel, event string)
}

type dispatch struct{}

func newDispatcher() dispatcher {
	d := new(dispatch)
	return d
}

func (d *dispatch) setup(_ *caseOfficer, _, _ string) {}

func (d *dispatch) dispatch(agent *caseOfficer, channel, event string) {
	switch event {
	case messaging.StartupEvent:
		agent.handler.Trace(agent, channel, event, "startup")
	case messaging.ShutdownEvent:
		agent.handler.Trace(agent, channel, event, "shutdown")
	}
}
