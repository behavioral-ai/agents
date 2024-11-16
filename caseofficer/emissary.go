package caseofficer

import (
	"github.com/advanced-go/common/core"
	"github.com/advanced-go/common/messaging"
	"github.com/advanced-go/resiliency/guidance"
)

type newServiceAgent func(origin core.Origin, c *caseOfficer)

func emissaryAttend(c *caseOfficer, fn *caseOfficerFunc, guide *guidance.Guidance, newAgent newServiceAgent) {
	fn.startup(c, guide, newAgent)

	for {
		// new assignment processing
		select {
		case <-c.ticker.C():
			c.handler.AddActivity(c.agentId, "onTick()")
			fn.update(c, guide, newAgent)
		default:
		}
		// control channel processing
		select {
		case msg := <-c.ctrlC:
			switch msg.Event() {
			case messaging.ShutdownEvent:
				c.shutdown()
				c.handler.AddActivity(c.agentId, messaging.ShutdownEvent)
				return
			case messaging.DataChangeEvent:
				if msg.IsContentType(guidance.ContentTypeCalendar) {
					c.serviceAgents.Broadcast(msg)
				}
			default:
				c.handler.Handle(MessageEventErrorStatus(c.agentId, msg))
			}
		default:
		}
	}
}
