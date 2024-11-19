package service

import (
	"github.com/advanced-go/agents/common"
	"github.com/advanced-go/common/messaging"
)

// master attention
func masterAttend(agent *service) {
	//rateLimiting := action1.NewRateLimiting()
	//common1.SetRateLimitingAction(r.handler, r.origin, rateLimiting, exp)

	for {
		// message processing
		select {
		case msg := <-agent.master.C:
			switch msg.Event() {
			case messaging.ShutdownEvent:
				agent.masterFinalize()
				agent.handler.Trace(agent, messaging.ShutdownEvent)
				agent.handler.OnMessage(agent, msg, agent.master)
				return
			case messaging.ObservationEvent:
				agent.handler.Trace(agent, messaging.ObservationEvent)
				/*
					observe, ok := msg.Body.(*common1.Observation)
					if !ok {
						continue
					}
					inf := runInference(r, observe)
					if inf == nil {
						continue
					}
					action := newAction(inf)
					rateLimiting.Limit = action.Limit
					rateLimiting.Burst = action.Burst
					common1.AddRateLimitingExperience(r.handler, r.origin, inf, action, exp)

				*/
				agent.handler.OnMessage(agent, msg, agent.master)
			default:
				agent.handler.OnError(agent, agent.handler.Notify(common.MessageEventErrorStatus(agent.Uri(), msg)))
			}
		default:
		}
	}
}
