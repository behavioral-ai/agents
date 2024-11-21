package caseofficer

import (
	"github.com/advanced-go/common/core"
	"github.com/advanced-go/common/messaging"
	"github.com/advanced-go/resiliency/guidance"
)

type newServiceAgent func(origin core.Origin, c *caseOfficer)

func emissaryAttend(agent *caseOfficer, fn *caseOfficerFunc, guide *guidance.Guidance, newAgent newServiceAgent) {
	agent.dispatch(messaging.StartupEvent)
	//fn.startup(agent, guide, newAgent)

	for {
		// new assignment processing
		select {
		case <-agent.ticker.C():
			//fn.update(agent, guide, newAgent)

		default:
		}
		// control channel processing
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
				agent.Notify(MessageEventErrorStatus(agent.Uri(), msg))
			}
		default:
		}
	}
}
