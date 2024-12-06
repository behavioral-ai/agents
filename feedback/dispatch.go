package feedback

import (
	"fmt"
	"github.com/behavioral-ai/core/messaging"
)

type dispatcher interface {
	setup(agent *feedback, event string)
	dispatch(agent *feedback, event string)
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

func (d *dispatch) setup(_ *feedback, _ string) {}

func (d *dispatch) trace(agent *feedback, event, activity string) {
	agent.handler.Trace(agent, d.channel, event, activity)
}

func (d *dispatch) dispatch(agent *feedback, event string) {
	switch event {
	case messaging.StartupEvent:
		if d.test {
			d.trace(agent, event, fmt.Sprintf("count:%v", 1))
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
