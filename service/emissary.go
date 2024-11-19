package service

import (
	"github.com/advanced-go/agents/common"
	"github.com/advanced-go/common/messaging"
	"github.com/advanced-go/resiliency/guidance"
)

// emissary attention
func emissaryAttend(agent *service, observe *common.Observation) {
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
			agent.handler.OnTick(agent, ticker)
		default:
		}
		// message processing
		select {
		case msg := <-agent.emissary.C:
			switch msg.Event() {
			case messaging.ShutdownEvent:
				ticker.Stop()
				agent.emissaryFinalize()
				agent.handler.OnMessage(agent, msg, agent.emissary)
				return
			case messaging.DataChangeEvent:
				if p := guidance.GetCalendar(agent.handler, agent.Uri(), msg); p != nil {
				}
				agent.handler.OnMessage(agent, msg, agent.emissary)
			default:
				agent.handler.OnError(agent, agent.handler.Notify(common.MessageEventErrorStatus(agent.Uri(), msg)))
			}
		default:
		}
	}
}
