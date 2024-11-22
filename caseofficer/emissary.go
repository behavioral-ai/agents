package caseofficer

import (
	"github.com/advanced-go/common/core"
	"github.com/advanced-go/common/messaging"
	"github.com/advanced-go/resiliency/guidance"
)

type newServiceAgent func(origin core.Origin, handler messaging.OpsAgent) messaging.Agent

func emissaryAttend(agent *caseOfficer, guide *guidance.Guidance, newAgent newServiceAgent) {
	agent.dispatch(messaging.StartupEvent)
	createAssignments(agent, guide, newAgent)
	agent.startup()

	for {
		select {
		case <-agent.ticker.C():
			updateAssignments(agent, guide, newAgent)
			agent.dispatch(messaging.TickEvent)
		default:
		}
		select {
		case msg := <-agent.emissary.C:
			agent.setup(msg.Event())
			switch msg.Event() {
			case messaging.ShutdownEvent:
				agent.finalize()
				agent.dispatch(msg.Event())
				return
			case messaging.DataChangeEvent:
				if msg.IsContentType(guidance.ContentTypeCalendar) {
					agent.serviceAgents.Broadcast(msg)
					agent.dispatch(msg.Event())
				}
			default:
				agent.Notify(messaging.EventErrorStatus(agent.Uri(), msg))
			}
		default:
		}
	}
}
