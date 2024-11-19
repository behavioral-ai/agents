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
			agent.Trace(agent.Uri(), "onTick()")
			fn.update(agent, guide, newAgent)
			agent.OnTick(agent, agent.ticker)
		default:
		}
		// control channel processing
		select {
		case msg := <-agent.emissary.C:
			switch msg.Event() {
			case messaging.ShutdownEvent:
				agent.finalize()
				agent.Trace(agent.Uri(), msg)
				agent.OnMessage(agent, msg, agent.emissary)
				return
			case messaging.DataChangeEvent:
				if msg.IsContentType(guidance.ContentTypeCalendar) {
					agent.serviceAgents.Broadcast(msg)
				}
				agent.Trace(agent.Uri(), msg)
				agent.OnMessage(agent, msg, agent.emissary)
			default:
				agent.OnError(agent, agent.Notify(MessageEventErrorStatus(agent.Uri(), msg)))
			}
		default:
		}
	}
}
