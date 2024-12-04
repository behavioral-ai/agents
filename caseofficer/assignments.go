package caseofficer

import (
	"errors"
	"github.com/behavioral-ai/core/core"
	"github.com/behavioral-ai/resiliency/guidance"
)

func createAssignments(agent *caseOfficer, assignments *guidance.Assignments, newAgent newServiceAgent) {
	if newAgent == nil {
		agent.Notify(core.NewStatusError(core.StatusInvalidArgument, errors.New("error: create assignments newAgent is nil")))
		return
	}
	entry, status := assignments.All(agent.handler, agent.origin)
	if status.OK() {
		addAssignments(agent, entry, newAgent)
	}
	if !status.NotFound() {
		agent.Notify(status)
	}
}

func updateAssignments(agent *caseOfficer, assignments *guidance.Assignments, newAgent newServiceAgent) {
	if newAgent == nil {
		agent.Notify(core.NewStatusError(core.StatusInvalidArgument, errors.New("error: update assignments newAgent is nil")))
		return
	}
	entry, status := assignments.New(agent.handler, agent.origin)
	if status.OK() {
		addAssignments(agent, entry, newAgent)
	}
	if !status.NotFound() {
		agent.Notify(status)
	}
}

func addAssignments(agent *caseOfficer, entry []guidance.HostEntry, newAgent newServiceAgent) {
	for _, e := range entry {
		a := newAgent(e.Origin, agent, agent.dispatcher)
		err := agent.serviceAgents.Register(a)
		if err != nil {
			agent.Notify(core.NewStatusError(core.StatusInvalidArgument, err))
		} else {
			a.Run()
		}
	}
}
