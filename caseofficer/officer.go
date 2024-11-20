package caseofficer

import (
	"github.com/advanced-go/common/core"
	"github.com/advanced-go/common/messaging"
	"github.com/advanced-go/resiliency/guidance"
	"time"
)

const (
	CaseOfficerClass   = "case-officer"
	assignmentDuration = time.Minute * 30
)

type caseOfficer struct {
	running bool
	agentId string
	origin  core.Origin

	ticker        *messaging.Ticker
	emissary      *messaging.Channel
	handler       messaging.OpsAgent
	serviceAgents *messaging.Exchange
	shutdownFunc  func()

	// testing/debugging
	onMessage       func(agent any, msg *messaging.Message, ch *messaging.Channel)
	onTick          func(agent any, ticker *messaging.Ticker)
	preStateChange  func()
	postStateChange func()
}

func AgentUri(origin core.Origin) string {
	return origin.Uri(CaseOfficerClass)
}

// NewAgent - create a new case officer agent
func NewAgent(origin core.Origin, handler messaging.OpsAgent) messaging.OpsAgent {
	return newAgent(origin, handler)
}

// newAgent - create a new case officer agent
func newAgent(origin core.Origin, handler messaging.OpsAgent) *caseOfficer {
	c := new(caseOfficer)
	c.agentId = AgentUri(origin)
	c.origin = origin
	c.ticker = messaging.NewPrimaryTicker(assignmentDuration)
	c.emissary = messaging.NewEmissaryChannel(true)
	c.handler = handler
	c.serviceAgents = messaging.NewExchange()

	c.onMessage = func(agent any, msg *messaging.Message, src *messaging.Channel) {}
	c.onTick = func(agent any, ticker *messaging.Ticker) {}
	c.preStateChange = func() {}
	c.postStateChange = func() {}
	return c
}

// String - identity
func (c *caseOfficer) String() string { return c.Uri() }

// Uri - agent identifier
func (c *caseOfficer) Uri() string { return c.agentId }

// Message - message the agent
func (c *caseOfficer) Message(m *messaging.Message) { c.emissary.C <- m }

// Notify - notifier
func (c *caseOfficer) Notify(status *core.Status) *core.Status {
	// TODO : do we need any processing specific to a case officer? If not then forward to handler
	return c.handler.Notify(status)
}

// Trace - activity tracing
func (c *caseOfficer) Trace(agent, activity any) { c.handler.Trace(agent, activity) }

// Add - add a shutdown function
func (c *caseOfficer) Add(f func()) { c.shutdownFunc = messaging.AddShutdown(c.shutdownFunc, f) }

// Run - run the agent
func (c *caseOfficer) Run() {
	if c.running {
		return
	}
	c.running = true
	go emissaryAttend(c, officer, guidance.Guide, initAgent)
}

// Shutdown - shutdown the agent
func (c *caseOfficer) Shutdown() {
	if !c.running {
		return
	}
	c.running = false
	// Removes agent from its exchange if registered
	if c.shutdownFunc != nil {
		c.shutdownFunc()
	}
	msg := messaging.NewControlMessage(c.agentId, c.agentId, messaging.ShutdownEvent)
	c.serviceAgents.Shutdown()
	c.emissary.C <- msg
}

func (c *caseOfficer) IsFinalized() bool {
	return c.emissary.IsFinalized() && c.ticker.IsFinalized()
}

func (c *caseOfficer) startup() {
	c.ticker.Start(-1)
}

func (c *caseOfficer) finalize() {
	c.emissary.Close()
	c.ticker.Stop()
}

func (c *caseOfficer) reviseTicker(newDuration time.Duration) {
	c.ticker.Start(newDuration)
}
