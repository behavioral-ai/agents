package caseofficer

import (
	"github.com/behavioral-ai/core/core"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/resiliency/guidance"
)

type newServiceAgent func(origin core.Origin, handler messaging.OpsAgent, dispatcher messaging.TraceDispatcher) messaging.Agent

func emissaryAttend(agent *caseOfficer, assignments *guidance.Assignments, newAgent newServiceAgent) {
	createAssignments(agent, assignments, newAgent)
	agent.startup()
	agent.dispatch(messaging.StartupEvent)

	for {
		select {
		case <-agent.ticker.C():
			updateAssignments(agent, assignments, newAgent)
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
