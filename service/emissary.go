package service

import (
	"github.com/advanced-go/agents/common"
	"github.com/advanced-go/common/messaging"
	"github.com/advanced-go/resiliency/guidance"
)

// emissary attention
func emissaryAttend(agent *service, observe *common.Observation) {
	agent.emissaryDispatch(messaging.StartupEvent)
	ticker := messaging.NewPrimaryTicker(agent.duration)

	ticker.Start(-1)
	for {
		select {
		case <-ticker.C():
			e, status := observe.Timeseries(agent.handler, agent.origin)
			if status.OK() {
				m := messaging.NewControlMessage(messaging.MasterChannel, agent.Uri(), messaging.ObservationEvent)
				m.SetContent(contentTypeObservation, observation{
					Latency:  e.Latency,
					Gradient: e.Gradient})
				agent.Message(m)
				agent.emissaryDispatch(messaging.ObservationEvent)
			}
		default:
		}
		select {
		case msg := <-agent.emissary.C:
			agent.emissarySetup(msg.Event())
			switch msg.Event() {
			case messaging.ShutdownEvent:
				ticker.Stop()
				agent.emissaryFinalize()
				agent.emissaryDispatch(msg.Event())
				return
			case messaging.DataChangeEvent:
				if p := guidance.GetCalendar(agent.handler, agent.Uri(), msg); p != nil {
					agent.emissaryDispatch(msg.Event())
				}
			default:
				agent.handler.Notify(messaging.EventErrorStatus(agent.Uri(), msg))
			}
		default:
		}
	}
}
