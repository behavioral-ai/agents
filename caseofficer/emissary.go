package caseofficer

import (
	"github.com/advanced-go/common/core"
	"github.com/advanced-go/common/messaging"
	"github.com/advanced-go/resiliency/guidance"
)

type newServiceAgent func(origin core.Origin, c *caseOfficer)

func emissaryAttend(agent *caseOfficer, fn *caseOfficerFunc, guide *guidance.Guidance, newAgent newServiceAgent) {
	fn.startup(agent, guide, newAgent)

	for {
		// new assignment processing
		select {
		case <-agent.ticker.C():
			fn.update(agent, guide, newAgent)
			agent.OnTick(agent, agent.ticker)
		default:
		}
		// control channel processing
		select {
		case msg := <-agent.emissary.C:
			agent.OnMessage(agent, msg, agent.emissary)
			switch msg.Event() {
			case messaging.ShutdownEvent:
				agent.finalize()
				return
			case messaging.DataChangeEvent:
				if msg.IsContentType(guidance.ContentTypeCalendar) {
					agent.serviceAgents.Broadcast(msg)
				}
			default:
				agent.Notify(agent, MessageEventErrorStatus(agent.Uri(), msg))
			}
		default:
		}
	}
}
