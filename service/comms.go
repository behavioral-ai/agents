package service

import "github.com/behavioral-ai/core/messaging"

type comms struct {
	channel string
	ch      *messaging.Channel
	global  messaging.Dispatcher
	local   dispatcher
}

func newComms(master bool, global messaging.Dispatcher, local dispatcher) *comms {
	c := new(comms)
	if master {
		c.channel = messaging.MasterChannel
		c.ch = messaging.NewEmissaryChannel(true)
	} else {
		c.channel = messaging.EmissaryChannel
		c.ch = messaging.NewMasterChannel(true)
	}
	c.global = global
	c.local = local
	return c
}

func newMasterComms(global messaging.Dispatcher, local dispatcher) *comms {
	return newComms(true, global, local)
}

func newEmmissaryComms(global messaging.Dispatcher, local dispatcher) *comms {
	return newComms(false, global, local)
}

func (c *comms) isFinalized() bool { return c.ch.IsFinalized() }

func (c *comms) finalize() { c.ch.Close() }

func (c *comms) enable() { c.ch.Enable() }

func (c *comms) send(m *messaging.Message) { c.ch.C <- m }

func (c *comms) setup(agent *service, event string) { c.local.setup(agent, event) }

func (c *comms) dispatch(agent *service, event string) {
	if c.global != nil {
		c.global.Dispatch(agent, c.channel, event, "")
	}
	if c.local != nil {
		c.local.dispatch(agent, event)
	}
}
