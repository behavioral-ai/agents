package service

import "github.com/advanced-go/common/messaging"

type emissary struct{ channel string }

func newEmissaryDispatcher() dispatcher {
	d := new(emissary)
	d.channel = messaging.EmissaryChannel
	return d
}

func (d *emissary) setup(_ *service, _ string) {}

func (d *emissary) dispatch(agent *service, event string) {
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
