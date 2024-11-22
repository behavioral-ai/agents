package service

import (
	"github.com/advanced-go/common/messaging"
)

type masterT struct{ channel string }

func newTestMasterDispatcher() dispatcher {
	d := new(masterT)
	d.channel = messaging.MasterChannel
	return d
}

func (d *masterT) setup(_ *service, _ string) {}

func (d *masterT) dispatch(agent *service, event string) {
	switch event {
	case messaging.DataChangeEvent:
		agent.handler.Trace(agent, d.channel, event, "Broadcast() -> calendar data change event")
	}
}
