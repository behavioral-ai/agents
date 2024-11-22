package service

import (
	"github.com/advanced-go/common/messaging"
)

// master attention
func masterAttend(agent *service) {
	agent.masterDispatch(messaging.StartupEvent)

	for {
		// message processing
		select {
		case msg := <-agent.master.C:
			agent.masterSetup(msg.Event())
			switch msg.Event() {
			case messaging.ShutdownEvent:
				agent.masterFinalize()
				agent.masterDispatch(msg.Event())
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
				agent.masterDispatch(msg.Event())
			default:
				agent.handler.Notify(messaging.EventErrorStatus(agent.Uri(), msg))
			}
		default:
		}
	}
}
