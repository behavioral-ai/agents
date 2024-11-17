package service

import (
	"github.com/advanced-go/agents/common"
	"github.com/advanced-go/common/messaging"
	"github.com/advanced-go/resiliency/guidance"
)

// emissary attention
func emissaryAttend[T messaging.Notifier](agent *service, observe *common.Observation) {
	var notify T

	ticker := messaging.NewPrimaryTicker(agent.duration)
	//limit := timeseries1.Threshold{}
	//common1.SetPercentileThreshold(r.handler, r.origin, &limit, observe)

	ticker.Start(-1)
	for {
		// observation processing
		select {
		case <-ticker.C():
			//		actual, status := observe.PercentThresholdQuery(r.handler, r.origin, time.Now().UTC(), time.Now().UTC())
			//		if status.OK() {
			//			m := messaging.NewRightChannelMessage("", r.agentId, messaging.ObservationEvent, common1.NewObservation(actual, limit))
			//			r.Message(m)
			//			}
			notify.OnTick(agent, ticker)
		default:
		}
		// message processing
		select {
		case msg := <-agent.emissary.C:
			switch msg.Event() {
			case messaging.ShutdownEvent:
				ticker.Stop()
				agent.emissary.Close()
				notify.OnMessage(agent, msg, agent.emissary)
				return
			case messaging.DataChangeEvent:
				if p := guidance.GetCalendar(agent.handler, agent.agentId, msg); p != nil {
				}
				notify.OnMessage(agent, msg, agent.emissary)

			default:
				notify.OnError(agent, agent.handler.Handle(common.MessageEventErrorStatus(agent.Uri(), msg)))
			}
		default:
		}
	}
}
