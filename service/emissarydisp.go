package service

import "github.com/advanced-go/common/messaging"

type emissary struct{}

func newEmissaryDispatcher() dispatcher {
	d := new(emissary)
	return d
}

func (d *emissary) setup(_ *service, _ string) {}

func (d *emissary) dispatch(agent *service, event string) {
	switch event {
	case messaging.StartupEvent:
		agent.handler.Trace(agent, messaging.EmissaryChannel, event, "startup")
	case messaging.ShutdownEvent:
		agent.handler.Trace(agent, messaging.EmissaryChannel, event, "shutdown")
	}
}
