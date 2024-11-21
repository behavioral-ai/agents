package service

import (
	"github.com/advanced-go/common/messaging"
)

type masterT struct{}

func newTestMasterDispatcher() dispatcher {
	d := new(masterT)
	return d
}

func (d *masterT) setup(_ *service, _ string) {}

func (d *masterT) dispatch(agent *service, event string) {
	switch event {
	case messaging.DataChangeEvent:
		agent.handler.Trace(agent, messaging.MasterChannel, event, "Broadcast() -> calendar data change event")
	}
}
