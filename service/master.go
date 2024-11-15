package service

import (
	"github.com/advanced-go/agents/common"
	"github.com/advanced-go/common/messaging"
)

// master attention
func masterAttend(r *service) {
	//rateLimiting := action1.NewRateLimiting()
	//common1.SetRateLimitingAction(r.handler, r.origin, rateLimiting, exp)

	for {
		// message processing
		select {
		case msg := <-r.master.C:
			switch msg.Event() {
			case messaging.ShutdownEvent:
				r.master.Close()
				r.handler.AddActivity(r.agentId, messaging.ShutdownEvent)
				return
			case messaging.ObservationEvent:
				r.handler.AddActivity(r.agentId, messaging.ObservationEvent)
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
			default:
				r.handler.Handle(common.MessageEventErrorStatus(r.agentId, msg))
			}
		default:
		}
	}
}
