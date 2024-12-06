package feedback

import (
	"github.com/behavioral-ai/agents/service"
	"github.com/behavioral-ai/core/core"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/resiliency/guidance"
	"time"
)

const (
	Class    = "feedback"
	duration = time.Second * 15
)

type feedback struct {
	running bool
	agentId string
	origin  core.Origin

	ticker     *messaging.Ticker
	emissary   *messaging.Channel
	handler    messaging.OpsAgent
	dispatcher messaging.Dispatcher
	sender     dispatcher
}

func AgentUri(origin core.Origin) string {
	return origin.Uri(Class)
}

// NewAgent - create a new feedback agent
func NewAgent(origin core.Origin, handler messaging.OpsAgent, dispatcher messaging.Dispatcher) messaging.Agent {
	return newAgent(origin, handler, dispatcher, newDispatcher(false))
}

// newAgent - create a new feedback agent
func newAgent(origin core.Origin, handler messaging.OpsAgent, dispatcher messaging.Dispatcher, sender dispatcher) *feedback {
	c := new(feedback)
	c.agentId = AgentUri(origin)
	c.origin = origin
	c.ticker = messaging.NewPrimaryTicker(duration)
	c.emissary = messaging.NewEmissaryChannel(true)
	c.handler = handler
	c.sender = sender
	c.dispatcher = dispatcher
	return c
}

// String - identity
func (c *feedback) String() string { return c.Uri() }

// Uri - agent identifier
func (c *feedback) Uri() string { return c.agentId }

// Message - message the agent
func (c *feedback) Message(m *messaging.Message) {
	if m == nil {
		return
	}
	c.emissary.C <- m
}

// Notify - notifier
func (c *feedback) Notify(status *core.Status) *core.Status { return c.handler.Notify(status) }

// Trace - activity tracing
func (c *feedback) Trace(agent messaging.Agent, channel, event, activity string) {
	c.handler.Trace(agent, channel, event, activity)
}

// Run - run the agent
func (c *feedback) Run() {
	if c.running {
		return
	}
	c.running = true
	go emissaryAttend(c, guidance.Assign, service.NewAgent)

}

// Shutdown - shutdown the agent
func (c *feedback) Shutdown() {
	if !c.running {
		return
	}
	c.running = false
	msg := messaging.NewControlMessage(c.Uri(), c.Uri(), messaging.ShutdownEvent)
	c.emissary.C <- msg
}

func (c *feedback) IsFinalized() bool {
	return c.emissary.IsFinalized() && c.ticker.IsFinalized()
}

func (c *feedback) startup() {
	c.ticker.Start(-1)
}

func (c *feedback) finalize() {
	c.emissary.Close()
	c.ticker.Stop()
}

func (c *feedback) setup(event string) {
	c.sender.setup(c, event)
}

func (c *feedback) dispatch(event string) {
	if c.dispatcher != nil {
		c.dispatcher.Dispatch(c, messaging.EmissaryChannel, event, "")
		return
	}
	c.sender.dispatch(c, event)
}
