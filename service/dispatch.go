package service

import "github.com/advanced-go/common/messaging"

type dispatcher interface {
	setup(agent *service, event string)
	dispatch(agent *service, event string)
}

type master struct {
	test    bool
	filter  *messaging.TraceFilter
	channel string
}

func newMasterDispatcher(filter *messaging.TraceFilter, test bool) dispatcher {
	d := new(master)
	d.channel = messaging.MasterChannel
	d.filter = filter
	if d.filter == nil {
		d.filter = messaging.NewTraceFilter("", "")
	}
	d.test = test
	return d
}

func (d *master) setup(_ *service, _ string) {}

func (d *master) trace(agent *service, event, activity string) {
	if d.filter.Access(d.channel, event) {
		agent.handler.Trace(agent, d.channel, event, activity)
	}
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
	filter  *messaging.TraceFilter
}

func newEmissaryDispatcher(filter *messaging.TraceFilter, test bool) dispatcher {
	d := new(emissary)
	d.channel = messaging.EmissaryChannel
	d.test = test
	d.filter = filter
	if d.filter == nil {
		d.filter = messaging.NewTraceFilter("", "")
	}
	return d
}

func (d *emissary) setup(_ *service, _ string) {}

func (d *emissary) trace(agent *service, event, activity string) {
	if d.filter.Access(d.channel, event) {
		agent.handler.Trace(agent, d.channel, event, activity)
	}
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
