package service

import (
	"github.com/advanced-go/agents/common"
	"github.com/advanced-go/common/core"
	"github.com/advanced-go/common/messaging"
	"time"
)

const (
	Class           = "service"
	defaultDuration = time.Second * 10
)

type service struct {
	running bool
	agentId string
	origin  core.Origin

	duration       time.Duration
	emissary       *messaging.Channel
	master         *messaging.Channel
	handler        messaging.OpsAgent
	masterSender   dispatcher
	emissarySender dispatcher
}

func serviceAgentUri(origin core.Origin) string {
	return origin.Uri(Class)
}

// NewAgent - create a new service agent
func NewAgent(origin core.Origin, handler messaging.OpsAgent) messaging.Agent {
	return newAgent(origin, handler, newMasterDispatcher(), newEmissaryDispatcher())
}

func newAgent(origin core.Origin, handler messaging.OpsAgent, master, emissary dispatcher) *service {
	r := new(service)
	r.origin = origin
	r.agentId = serviceAgentUri(origin)
	r.duration = defaultDuration

	r.emissary = messaging.NewEmissaryChannel(true)
	r.master = messaging.NewMasterChannel(true)
	r.handler = handler
	r.masterSender = master
	r.emissarySender = emissary
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

func (s *service) emissarySetup(event string) {
	s.emissarySender.setup(s, event)
}

func (s *service) emissaryDispatch(event string) {
	s.emissarySender.dispatch(s, event)
}

func (s *service) masterSetup(event string) {
	s.masterSender.setup(s, event)
}

func (s *service) masterDispatch(event string) {
	s.masterSender.dispatch(s, event)
}
