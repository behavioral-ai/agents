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
	ctrlC         chan *messaging.Message
	handler       messaging.OpsAgent
	serviceAgents *messaging.Exchange
	shutdownFunc  func()
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
	c.ticker = messaging.NewTicker(assignmentDuration)
	c.ctrlC = make(chan *messaging.Message, messaging.ChannelSize)
	c.handler = handler
	c.serviceAgents = messaging.NewExchange()
	return c
}

// String - identity
func (c *caseOfficer) String() string { return c.Uri() }

// Uri - agent identifier
func (c *caseOfficer) Uri() string { return c.agentId }

// Message - message the agent
func (c *caseOfficer) Message(m *messaging.Message) { c.ctrlC <- m }

// Handle - error handler
func (c *caseOfficer) Handle(status *core.Status) *core.Status {
	// TODO : do we need any processing specific to a case officer? If not then forward to handler
	return c.handler.Handle(status)
}

// AddActivity - add activity
func (c *caseOfficer) AddActivity(agentId string, content any) {
	// TODO : Any operations specific processing ??  If not then forward to handler
	c.handler.AddActivity(agentId, content)
}

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
	if c.ctrlC != nil {
		c.ctrlC <- msg
	}

}

func (c *caseOfficer) startup() {
	c.ticker.Start(-1)
}

func (c *caseOfficer) shutdown() {
	close(c.ctrlC)
	c.ticker.Stop()
}

func (c *caseOfficer) reviseTicker(newDuration time.Duration) {
	c.ticker.Start(newDuration)
}
