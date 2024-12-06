package caseofficer

import (
	"github.com/behavioral-ai/core/core"
)

func addFeedback(agent *caseOfficer, newAgent createAgent) {
	a := newAgent(agent.origin, agent, agent.dispatcher)
	err := agent.serviceAgents.Register(a)
	if err != nil {
		agent.Notify(core.NewStatusError(core.StatusInvalidArgument, err))
	} else {
		a.Run()
	}
}
