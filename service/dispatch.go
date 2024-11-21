package service

import "github.com/advanced-go/common/messaging"

type dispatcher interface {
	setup(agent *service, event string)
	dispatch(agent *service, event string)
}

type dispatch struct{}

func newDispatcher() dispatcher {
	d := new(dispatch)
	return d
}

func (d *dispatch) setup(_ *service, _ string) {}

func (d *dispatch) dispatch(agent *service, event string) {
	switch event {
	case messaging.StartupEvent:
		agent.handler.Trace(agent, event, "startup")
	case messaging.ShutdownEvent:
		agent.handler.Trace(agent, event, "shutdown")
	}
}
