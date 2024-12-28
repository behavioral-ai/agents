package service

import "github.com/behavioral-ai/core/messaging"

type communications struct {
	name   string
	ch     *messaging.Channel
	global messaging.Dispatcher
	local  dispatcher
}

func newComms(master bool, global messaging.Dispatcher, local dispatcher) *communications {
	c := new(communications)
	if master {
		c.name = messaging.MasterChannel
		c.ch = messaging.NewEmissaryChannel(true)
	} else {
		c.name = messaging.EmissaryChannel
		c.ch = messaging.NewMasterChannel(true)
	}
	c.global = global
	c.local = local
	return c
}

func newMasterComms(global messaging.Dispatcher, local dispatcher) *communications {
	return newComms(true, global, local)
}

func newEmmissaryComms(global messaging.Dispatcher, local dispatcher) *communications {
	return newComms(false, global, local)
}

func (c *communications) channel() *messaging.Channel { return c.ch }

func (c *communications) isFinalized() bool { return c.ch.IsFinalized() }

func (c *communications) finalize() { c.ch.Close() }

func (c *communications) enable() { c.ch.Enable() }

func (c *communications) send(m *messaging.Message) { c.ch.C <- m }

func (c *communications) setup(agent *service, event string) { c.local.setup(agent, event) }

func (c *communications) dispatch(agent *service, event string) {
	if c.global != nil {
		c.global.Dispatch(agent, c.name, event, "")
	}
	if c.local != nil {
		c.local.dispatch(agent, event)
	}
}
