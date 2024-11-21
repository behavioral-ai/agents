package service

import (
	"github.com/advanced-go/agents/common"
	"github.com/advanced-go/common/core"
	"github.com/advanced-go/common/messaging"
	"time"
)

const (
	Class           = "service"
	defaultDuration = time.Minute * 5
)

type service struct {
	running bool
	agentId string
	origin  core.Origin

	duration time.Duration
	emissary *messaging.Channel
	master   *messaging.Channel
	handler  messaging.OpsAgent
	sender   dispatcher
}

func serviceAgentUri(origin core.Origin) string {
	return origin.Uri(Class)
}

// NewAgent - create a new service agent
func NewAgent(origin core.Origin, handler messaging.OpsAgent) messaging.Agent {
	return newAgent(origin, handler, newDispatcher())
}

func newAgent(origin core.Origin, handler messaging.OpsAgent, sender dispatcher) *service {
	r := new(service)
	r.origin = origin
	r.agentId = serviceAgentUri(origin)
	r.duration = defaultDuration

	r.emissary = messaging.NewEmissaryChannel(true)
	r.master = messaging.NewMasterChannel(true)
	r.handler = handler
	r.sender = sender
	return r
}

// String - identity
func (s *service) String() string { return s.Uri() }

// Uri - agent identifier
func (s *service) Uri() string { return s.agentId }

// Message - message the agent
func (s *service) Message(m *messaging.Message) {
	if m == nil {
		return
	}
	switch m.To() {
	case messaging.EmissaryChannel:
		s.emissary.C <- m
	case messaging.MasterChannel:
		s.master.C <- m
	default:
		s.emissary.C <- m
	}
}

// Run - run the agent
func (s *service) Run() {
	if s.running {
		return
	}
	go masterAttend(s)
	go emissaryAttend(s, common.Observe)
	s.running = true
}

// Shutdown - shutdown the agent
func (s *service) Shutdown() {
	if !s.running {
		return
	}
	s.running = false
	msg := messaging.NewControlMessage(s.Uri(), s.Uri(), messaging.ShutdownEvent)
	s.emissary.Enable()
	s.emissary.C <- msg
	s.master.Enable()
	s.master.C <- msg
}

func (s *service) IsFinalized() bool {
	return s.emissary.IsFinalized() && s.master.IsFinalized()
}

func (s *service) isFinalizedEmissary() bool {
	return s.emissary.IsFinalized()
}

func (s *service) emissaryFinalize() {
	s.emissary.Close()
}

func (s *service) isFinalizedMaster() bool {
	return s.master.IsFinalized()
}

func (s *service) masterFinalize() {
	s.master.Close()
}

func (s *service) setup(event string) {
	s.sender.setup(s, event)
}

func (s *service) dispatch(event string) {
	s.sender.dispatch(s, event)
}
