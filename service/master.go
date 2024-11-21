package service

import (
	"github.com/advanced-go/agents/common"
	"github.com/advanced-go/common/messaging"
)

// master attention
func masterAttend(agent *service) {
	agent.dispatch(messaging.StartupEvent)

	for {
		// message processing
		select {
		case msg := <-agent.master.C:
			agent.setup(msg.Event())
			switch msg.Event() {
			case messaging.ShutdownEvent:
				agent.masterFinalize()
				agent.dispatch(msg.Event())
				return
			case messaging.ObservationEvent:
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
				agent.dispatch(msg.Event())
			default:
				agent.handler.Notify(common.MessageEventErrorStatus(agent.Uri(), msg))
			}
		default:
		}
	}
}
