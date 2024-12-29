package feedback

import (
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/resiliency/guidance"
)

func emissaryAttend(agent *feedback, assignments *guidance.Assignments) {
	paused := false
	agent.startup()
	agent.dispatch(messaging.StartupEvent)

	for {
		select {
		case <-agent.ticker.C():
			agent.dispatch(messaging.TickEvent)
		default:
		}
		select {
		case msg := <-agent.emissary.C:
			agent.setup(msg.Event())
			switch msg.Event() {
			case messaging.PauseEvent:
				paused = true
			case messaging.ResumeEvent:
				paused = false
			case messaging.ShutdownEvent:
				agent.finalize()
				agent.dispatch(msg.Event())
				return
			case messaging.DataChangeEvent:
				if paused {
				}
				if !msg.IsContentType(guidance.ContentTypeCalendar) {
					//agent.dispatch(msg.Event())
				}
			default:
				agent.Notify(messaging.EventErrorStatus(agent.Uri(), msg))
			}
			agent.dispatch(msg.Event())
		default:
		}
	}
}
