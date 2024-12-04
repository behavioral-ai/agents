package service

import "github.com/behavioral-ai/core/messaging"

type dispatcher interface {
	setup(agent *service, event string)
	dispatch(agent *service, event string)
}

type master struct {
	test    bool
	channel string
}

func newMasterDispatcher(test bool) dispatcher {
	d := new(master)
	d.channel = messaging.MasterChannel
	d.test = test
	return d
}

func (d *master) setup(_ *service, _ string) {}

func (d *master) trace(agent *service, event, activity string) {
	agent.handler.Trace(agent, d.channel, event, activity)
}

func (d *master) dispatch(agent *service, event string) {
	switch event {
	case messaging.StartupEvent:
		d.trace(agent, event, "")
	case messaging.ShutdownEvent:
		d.trace(agent, event, "")
	case messaging.ObservationEvent:
		d.trace(agent, event, "")
	case messaging.DataChangeEvent:
		if d.test {
			d.trace(agent, event, "Broadcast() -> calendar data change event")
		}
	}
}

type emissary struct {
	test    bool
	channel string
}

func newEmissaryDispatcher(test bool) dispatcher {
	d := new(emissary)
	d.channel = messaging.EmissaryChannel
	d.test = test
	return d
}

func (d *emissary) setup(_ *service, _ string) {}

func (d *emissary) trace(agent *service, event, activity string) {
	agent.handler.Trace(agent, d.channel, event, activity)
}

func (d *emissary) dispatch(agent *service, event string) {
	switch event {
	case messaging.StartupEvent:
		d.trace(agent, event, "")
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
	case messaging.ObservationEvent:
		d.trace(agent, event, "")
	}
}
