package caseofficer

import (
	"github.com/behavioral-ai/agents/service"
	"github.com/behavioral-ai/core/core"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/resiliency/guidance"
	"time"
)

const (
	Class              = "case-officer"
	assignmentDuration = time.Second * 15
)

type caseOfficer struct {
	running bool
	agentId string
	origin  core.Origin

	ticker        *messaging.Ticker
	emissary      *messaging.Channel
	serviceAgents *messaging.Exchange
	handler       messaging.OpsAgent
	dispatcher    messaging.Dispatcher
	sender        dispatcher
}

func AgentUri(origin core.Origin) string {
	return origin.Uri(Class)
}

// NewAgent - create a new case officer agent
func NewAgent(origin core.Origin, handler messaging.OpsAgent, dispatcher messaging.Dispatcher) messaging.OpsAgent {
	return newAgent(origin, handler, dispatcher, newDispatcher(false))
}

// newAgent - create a new case officer agent
func newAgent(origin core.Origin, handler messaging.OpsAgent, dispatcher messaging.Dispatcher, sender dispatcher) *caseOfficer {
	c := new(caseOfficer)
	c.agentId = AgentUri(origin)
	c.origin = origin
	c.ticker = messaging.NewPrimaryTicker(assignmentDuration)
	c.emissary = messaging.NewEmissaryChannel(true)
	c.handler = handler
	c.serviceAgents = messaging.NewExchange()
	c.sender = sender
	c.dispatcher = dispatcher
	return c
}

// String - identity
func (c *caseOfficer) String() string { return c.Uri() }

// Uri - agent identifier
func (c *caseOfficer) Uri() string { return c.agentId }

// Message - message the agent
func (c *caseOfficer) Message(m *messaging.Message) {
	if m == nil {
		return
	}
	c.emissary.C <- m
}

// Notify - notifier
func (c *caseOfficer) Notify(status *core.Status) *core.Status { return c.handler.Notify(status) }

// Trace - activity tracing
func (c *caseOfficer) Trace(agent messaging.Agent, channel, event, activity string) {
	c.handler.Trace(agent, channel, event, activity)
}

// Run - run the agent
func (c *caseOfficer) Run() {
	if c.running {
		return
	}
	c.running = true
	go emissaryAttend(c, guidance.Assign, service.NewAgent, newFeedbackAgent)
}

// Shutdown - shutdown the agent
func (c *caseOfficer) Shutdown() {
	if !c.running {
		return
	}
	c.running = false
	msg := messaging.NewControlMessage(c.Uri(), c.Uri(), messaging.ShutdownEvent)
	c.serviceAgents.Shutdown()
	c.emissary.C <- msg
}

func (c *caseOfficer) IsFinalized() bool {
	return c.emissary.IsFinalized() && c.ticker.IsFinalized() && c.serviceAgents.IsFinalized()
}

func (c *caseOfficer) startup() {
	c.ticker.Start(-1)
}

func (c *caseOfficer) finalize() {
	c.emissary.Close()
	c.ticker.Stop()
	c.serviceAgents.Shutdown()
}

func (c *caseOfficer) reviseTicker(newDuration time.Duration) {
	c.ticker.Start(newDuration)
}

func (c *caseOfficer) setup(event string) {
	c.sender.setup(c, event)
}

func (c *caseOfficer) dispatch(event string) {
	if c.dispatcher != nil {
		c.dispatcher.Dispatch(c, messaging.EmissaryChannel, event, "")
		return
	}
	c.sender.dispatch(c, event)
}
