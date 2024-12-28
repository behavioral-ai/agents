package service

import (
	"github.com/behavioral-ai/core/messaging"
)

// master attention
func masterAttend(agent *service) {
	paused := false
	agent.master.dispatch(agent, messaging.StartupEvent)

	for {
		// message processing
		select {
		case msg := <-agent.master.ch.C:
			agent.master.setup(agent, msg.Event())
			switch msg.Event() {
			case messaging.PauseEvent:
				paused = true
				agent.master.dispatch(agent, msg.Event())
			case messaging.ResumeEvent:
				paused = false
				agent.master.dispatch(agent, msg.Event())
			case messaging.ShutdownEvent:
				agent.master.finalize()
				agent.master.dispatch(agent, msg.Event())
				return
			case messaging.ObservationEvent:
				if !paused {
					observe, status := getObservation(agent.handler, agent.Uri(), msg)
					if status.OK() {
						if observe.Gradient > 10 {
						}
						/*
							inf := runInference(r, observe)
							if inf == nil {
								continue
							}
							action := newAction(inf)
							rateLimiting.Limit = action.Limit
							rateLimiting.Burst = action.Burst
							common1.AddRateLimitingExperience(r.handler, r.origin, inf, action, exp)


						*/
					}
					agent.master.dispatch(agent, msg.Event())
				}
			default:
				agent.handler.Notify(messaging.EventErrorStatus(agent.Uri(), msg))
			}
		default:
		}
	}
}
