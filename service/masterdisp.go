package service

import "github.com/advanced-go/common/messaging"

type dispatcher interface {
	setup(agent *service, event string)
	dispatch(agent *service, event string)
}

type master struct{}

func newMasterDispatcher() dispatcher {
	d := new(master)
	return d
}

func (d *master) setup(_ *service, _ string) {}

func (d *master) dispatch(agent *service, event string) {
	switch event {
	case messaging.StartupEvent:
		agent.handler.Trace(agent, messaging.MasterChannel, event, "startup")
	case messaging.ShutdownEvent:
		agent.handler.Trace(agent, messaging.MasterChannel, event, "shutdown")
	}
}
