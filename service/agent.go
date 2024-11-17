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

	// Channels
	duration time.Duration
	emissary *messaging.Channel
	master   *messaging.Channel

	handler      messaging.OpsAgent
	shutdownFunc func()
}

func serviceAgentUri(origin core.Origin) string {
	return origin.Uri(Class)
}

// NewAgent - create a new service agent
func NewAgent(origin core.Origin, handler messaging.OpsAgent) messaging.Agent {
	return newAgent(origin, handler)
}

func newAgent(origin core.Origin, handler messaging.OpsAgent) *service {
	r := new(service)
	r.origin = origin
	r.agentId = serviceAgentUri(origin)
	r.duration = defaultDuration

	// Channels
	r.emissary = messaging.NewEmissaryChannel(true)
	r.master = messaging.NewMasterChannel(true)

	r.handler = handler
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
	// Specifically for the lhc or profile content
	if m.Channel() == messaging.ChannelLeft {
		s.emissary.C <- m
	} else {
		s.master.C <- m
	}
}

// Add - add a shutdown function
func (s *service) Add(f func()) { s.shutdownFunc = messaging.AddShutdown(s.shutdownFunc, f) }

// Run - run the agent
func (s *service) Run() {
	if s.running {
		return
	}
	go masterAttend[messaging.MutedNotifier](s)
	go emissaryAttend[messaging.MutedNotifier](s, common.Observe)
	s.running = true
}

// Shutdown - shutdown the agent
func (s *service) Shutdown() {
	if !s.running {
		return
	}
	s.running = false
	if s.shutdownFunc != nil {
		s.shutdownFunc()
	}
	msg := messaging.NewControlMessage(s.agentId, s.agentId, messaging.ShutdownEvent)
	s.emissary.Enable()
	s.emissary.C <- msg
	s.master.Enable()
	s.master.C <- msg
}

func (s *service) IsFinalized() bool {
	return s.emissary.IsFinalized() && s.master.IsFinalized()
}
