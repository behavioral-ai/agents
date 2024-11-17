package caseofficer

import (
	"github.com/advanced-go/common/core"
	"github.com/advanced-go/common/messaging"
	"github.com/advanced-go/resiliency/guidance"
)

type newServiceAgent func(origin core.Origin, c *caseOfficer)

func emissaryAttend[T messaging.Notifier](agent *caseOfficer, fn *caseOfficerFunc, guide *guidance.Guidance, newAgent newServiceAgent) {
	var notify T

	fn.startup(agent, guide, newAgent)

	for {
		// new assignment processing
		select {
		case <-agent.ticker.C():
			agent.handler.AddActivity(agent.Uri(), "onTick()")
			fn.update(agent, guide, newAgent)
			notify.OnTick(agent, agent.ticker)
		default:
		}
		// control channel processing
		select {
		case msg := <-agent.emissary.C:
			switch msg.Event() {
			case messaging.ShutdownEvent:
				agent.shutdown()
				agent.handler.AddActivity(agent.Uri(), messaging.ShutdownEvent)
				notify.OnMessage(agent, msg, agent.emissary)
				return
			case messaging.DataChangeEvent:
				if msg.IsContentType(guidance.ContentTypeCalendar) {
					agent.serviceAgents.Broadcast(msg)
				}
				notify.OnMessage(agent, msg, agent.emissary)
			default:
				notify.OnError(agent, agent.handler.Handle(MessageEventErrorStatus(agent.Uri(), msg)))
			}
		default:
		}
	}
}
