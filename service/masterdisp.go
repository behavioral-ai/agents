package service

import "github.com/advanced-go/common/messaging"

type dispatcher interface {
	setup(agent *service, event string)
	dispatch(agent *service, event string)
}

type master struct{ channel string }

func newMasterDispatcher() dispatcher {
	d := new(master)
	d.channel = messaging.MasterChannel
	return d
}

func (d *master) setup(_ *service, _ string) {}

func (d *master) dispatch(agent *service, event string) {
	switch event {
	case messaging.StartupEvent:
		agent.handler.Trace(agent, d.channel, event, "")
	case messaging.ShutdownEvent:
		agent.handler.Trace(agent, d.channel, event, "")
	case messaging.ObservationEvent:
		agent.handler.Trace(agent, d.channel, event, "")
	}
}
