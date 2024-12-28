package service

import (
	"github.com/behavioral-ai/core/messaging"
)

// master attention
func masterAttend(agent *service) {
	paused := false
	comms := agent.master
	comms.dispatch(agent, messaging.StartupEvent)

	for {
		// message processing
		select {
		case msg := <-comms.channel().C:
			comms.setup(agent, msg.Event())
			switch msg.Event() {
			case messaging.PauseEvent:
				paused = true
				comms.dispatch(agent, msg.Event())
			case messaging.ResumeEvent:
				paused = false
				comms.dispatch(agent, msg.Event())
			case messaging.ShutdownEvent:
				comms.finalize()
				comms.dispatch(agent, msg.Event())
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
					comms.dispatch(agent, msg.Event())
				}
			default:
				agent.handler.Notify(messaging.EventErrorStatus(agent.Uri(), msg))
			}
		default:
		}
	}
}
