package service

import (
	"github.com/advanced-go/agents/common"
	"github.com/advanced-go/common/messaging"
)

// master attention
func masterAttend[T messaging.Notifier](agent *service) {
	var notify T

	//rateLimiting := action1.NewRateLimiting()
	//common1.SetRateLimitingAction(r.handler, r.origin, rateLimiting, exp)

	for {
		// message processing
		select {
		case msg := <-agent.master.C:
			switch msg.Event() {
			case messaging.ShutdownEvent:
				agent.master.Close()
				agent.handler.AddActivity(agent.Uri(), messaging.ShutdownEvent)
				notify.OnMessage(agent, msg, agent.master)
				return
			case messaging.ObservationEvent:
				agent.handler.AddActivity(agent.Uri(), messaging.ObservationEvent)
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
				notify.OnMessage(agent, msg, agent.master)
			default:
				notify.OnError(agent, agent.handler.Handle(common.MessageEventErrorStatus(agent.Uri(), msg)))
			}
		default:
		}
	}
}
