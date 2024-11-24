package caseofficer

import (
	"fmt"
	"github.com/advanced-go/common/messaging"
)

type dispatcher interface {
	setup(agent *caseOfficer, event string)
	dispatch(agent *caseOfficer, event string)
}

type dispatch struct {
	test    bool
	channel string
}

func newDispatcher(test bool) dispatcher {
	d := new(dispatch)
	d.channel = messaging.EmissaryChannel
	d.test = test
	return d
}

func (d *dispatch) setup(_ *caseOfficer, _ string) {}

func (d *dispatch) trace(agent *caseOfficer, event, activity string) {
	agent.handler.Trace(agent, d.channel, event, activity)
}

func (d *dispatch) dispatch(agent *caseOfficer, event string) {
	switch event {
	case messaging.StartupEvent:
		if d.test {
			d.trace(agent, event, fmt.Sprintf("count:%v", agent.serviceAgents.Count()))
		} else {
			d.trace(agent, event, "")
		}
	case messaging.ShutdownEvent:
		d.trace(agent, event, "")
	case messaging.TickEvent:
		d.trace(agent, event, "")
	case messaging.DataChangeEvent:
		if d.test {
			d.trace(agent, event, "Broadcast() -> calendar data change event")
		} else {
			d.trace(agent, event, "")
		}
	}
}
