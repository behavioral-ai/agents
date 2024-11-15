package service

import (
	"github.com/advanced-go/common/messaging"
	"github.com/advanced-go/resiliency/common"
)

// emissary attention
func emissaryAttend(r *service, observe *common.Observation) {
	ticker := messaging.NewTicker(r.duration)
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
		default:
		}
		// message processing
		select {
		case msg := <-r.emissary.C:
			switch msg.Event() {
			case messaging.ShutdownEvent:
				ticker.Stop()
				r.emissary.Close()
				return
			case messaging.DataChangeEvent:
				if p := common.GetCalendar(r.handler, r.agentId, msg); p != nil {
				}
			default:
				r.handler.Handle(common.MessageEventErrorStatus(r.agentId, msg))
			}
		default:
		}
	}
}
