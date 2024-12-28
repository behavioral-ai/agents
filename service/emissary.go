package service

import (
	"github.com/behavioral-ai/agents/common"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/resiliency/guidance"
)

// emissary attention
func emissaryAttend(agent *service, observe *common.Observation) {
	paused := false
	agent.emissary.dispatch(agent, messaging.StartupEvent)
	ticker := messaging.NewPrimaryTicker(agent.duration)

	ticker.Start(-1)
	for {
		select {
		case <-ticker.C():
			if !paused {
				e, status := observe.Timeseries(agent.handler, agent.origin)
				if status.OK() {
					m := messaging.NewControlMessage(messaging.MasterChannel, agent.Uri(), messaging.ObservationEvent)
					m.SetContent(contentTypeObservation, observation{
						Latency:  e.Latency,
						Gradient: e.Gradient})
					agent.master.send(m)
					agent.emissary.dispatch(agent, messaging.ObservationEvent)
				}
			}
		default:
		}
		select {
		case msg := <-agent.emissary.ch.C:
			agent.emissary.setup(agent, msg.Event())
			switch msg.Event() {
			case messaging.PauseEvent:
				paused = true
				agent.emissary.dispatch(agent, msg.Event())
			case messaging.ResumeEvent:
				paused = false
				agent.emissary.dispatch(agent, msg.Event())
			case messaging.ShutdownEvent:
				ticker.Stop()
				agent.emissary.finalize()
				agent.emissary.dispatch(agent, msg.Event())
				return
			case messaging.DataChangeEvent:
				if !paused {
					if p := guidance.GetCalendar(agent.handler, agent.Uri(), msg); p != nil {
						agent.emissary.dispatch(agent, msg.Event())
					}
				}
			default:
				agent.handler.Notify(messaging.EventErrorStatus(agent.Uri(), msg))
			}
		default:
		}
	}
}
