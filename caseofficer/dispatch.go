package caseofficer

import "github.com/advanced-go/common/messaging"

type dispatcher interface {
	setup(agent *caseOfficer, event string)
	dispatch(agent *caseOfficer, event string)
}

type dispatch struct{ channel string }

func newDispatcher() dispatcher {
	d := new(dispatch)
	d.channel = messaging.EmissaryChannel
	return d
}

func (d *dispatch) setup(_ *caseOfficer, _ string) {}

func (d *dispatch) dispatch(agent *caseOfficer, event string) {
	switch event {
	case messaging.StartupEvent:
		agent.handler.Trace(agent, d.channel, event, "")
	case messaging.ShutdownEvent:
		agent.handler.Trace(agent, d.channel, event, "")
	case messaging.TickEvent:
		agent.handler.Trace(agent, d.channel, event, "")
	case messaging.DataChangeEvent:
		agent.handler.Trace(agent, d.channel, event, "")
	}
}
