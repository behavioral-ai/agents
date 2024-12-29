package caseofficer

import (
	"github.com/behavioral-ai/agents/feedback"
	"github.com/behavioral-ai/core/core"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/resiliency/guidance"
)

type createAgent func(origin core.Origin, handler messaging.OpsAgent, global messaging.Dispatcher) messaging.Agent

func newFeedbackAgent(origin core.Origin, handler messaging.OpsAgent, global messaging.Dispatcher) messaging.Agent {
	return feedback.NewAgent(origin, handler, global)
}

//type newServiceAgent func(origin core.Origin, handler messaging.OpsAgent, dispatcher messaging.Dispatcher) messaging.Agent
//type newFeedbackAgent func(origin core.Origin, handler messaging.OpsAgent, dispatcher messaging.Dispatcher) messaging.Agent

func emissaryAttend(agent *caseOfficer, assignments *guidance.Assignments, newService createAgent, newFeedback createAgent) {
	paused := false
	createAssignments(agent, assignments, newService)
	addFeedback(agent, newFeedback)
	agent.startup()
	agent.dispatch(messaging.StartupEvent)

	for {
		select {
		case <-agent.ticker.C():
			if !paused {
				updateAssignments(agent, assignments, newService)
				agent.dispatch(messaging.TickEvent)
			}
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
				if msg.IsContentType(guidance.ContentTypeCalendar) {
					agent.serviceAgents.Broadcast(msg)
				}
			default:
				agent.Notify(messaging.EventErrorStatus(agent.Uri(), msg))
			}
			agent.dispatch(msg.Event())
		default:
		}
	}
}
