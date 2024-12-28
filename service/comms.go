package service

import "github.com/behavioral-ai/core/messaging"

type comms struct {
	channel string
	agent   *service
	ch      *messaging.Channel
	global  messaging.Dispatcher
	local   dispatcher
}

func newComms(master bool, agent *service, global messaging.Dispatcher, local dispatcher) *comms {
	c := new(comms)
	c.agent = agent
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

func newMasterComms(agent *service, global messaging.Dispatcher, local dispatcher) *comms {
	return newComms(true, agent, global, local)
}

func newEmmissaryComms(agent *service, global messaging.Dispatcher, local dispatcher) *comms {
	return newComms(false, agent, global, local)
}

func (c *comms) isFinalized() bool { return c.ch.IsFinalized() }

func (c *comms) finalize() { c.ch.Close() }

func (c *comms) enable() { c.ch.Enable() }

func (c *comms) send(m *messaging.Message) { c.ch.C <- m }

func (c *comms) setup(event string) { c.local.setup(c.agent, event) }

func (c *comms) dispatch(event string) {
	if c.global != nil {
		c.global.Dispatch(c.agent, c.channel, event, "")
		//return
	}
	if c.local != nil {
		c.local.dispatch(c.agent, event)
	}
}
