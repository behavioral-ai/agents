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
			agent.emissaryDispatch(messaging.TickEvent)
			//agent.onTick(agent, ticker)
			//		actual, status := observe.PercentThresholdQuery(r.handler, r.origin, time.Now().UTC(), time.Now().UTC())
			//		if status.OK() {
			//			m := messaging.NewRightChannelMessage("", r.agentId, messaging.ObservationEvent, common1.NewObservation(actual, limit))
			//			r.Message(m)
			//			}
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
