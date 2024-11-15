package service

import (
	"github.com/advanced-go/common/core"
	"github.com/advanced-go/common/messaging"
	"github.com/advanced-go/resiliency/common"
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

	// Emissary channel
	r.duration = defaultDuration
	r.emissary = messaging.NewEnabledChannel()

	// Master channel
	r.master = messaging.NewEnabledChannel()
	r.handler = handler
	return r
}

// String - identity
func (r *service) String() string { return r.Uri() }

// Uri - agent identifier
func (r *service) Uri() string { return r.agentId }

// Message - message the agent
func (r *service) Message(m *messaging.Message) {
	if m == nil {
		return
	}
	// Specifically for the lhc or profile content
	if m.Channel() == messaging.ChannelLeft {
		r.emissary.C <- m
	} else {
		r.master.C <- m
	}
}

// Add - add a shutdown function
func (r *service) Add(f func()) { r.shutdownFunc = messaging.AddShutdown(r.shutdownFunc, f) }

// Run - run the agent
func (r *service) Run() {
	if r.running {
		return
	}
	go masterAttend(r)
	go emissaryAttend(r, common.Observe)
	r.running = true
}

// Shutdown - shutdown the agent
func (r *service) Shutdown() {
	if !r.running {
		return
	}
	r.running = false
	if r.shutdownFunc != nil {
		r.shutdownFunc()
	}
	msg := messaging.NewControlMessage(r.agentId, r.agentId, messaging.ShutdownEvent)
	r.emissary.Enable()
	r.emissary.C <- msg
	r.master.Enable()
	r.master.C <- msg
}
