package caseofficer

import (
	"errors"
	"github.com/advanced-go/common/core"
	"github.com/advanced-go/resiliency/guidance"
)

func createAssignments(agent *caseOfficer, guide *guidance.Guidance, newAgent newServiceAgent) {
	if newAgent == nil {
		agent.Notify(core.NewStatusError(core.StatusInvalidArgument, errors.New("error: create assignments newAgent is nil")))
		return
	}
	entry, status := guide.Assignments(agent.handler, agent.origin)
	if status.OK() {
		for _, e := range entry {
			addAssignment(agent, e, newAgent)
		}
		return
	}
	if !status.NotFound() {
		agent.Notify(status)
	}
}

func updateAssignments(agent *caseOfficer, guide *guidance.Guidance, newAgent newServiceAgent) {
	if newAgent == nil {
		agent.Notify(core.NewStatusError(core.StatusInvalidArgument, errors.New("error: update assignments newAgent is nil")))
		return
	}
	entry, status := guide.NewAssignments(agent.handler, agent.origin)
	if status.OK() {
		for _, e := range entry {
			addAssignment(agent, e, newAgent)
		}
		return
	}
	if !status.NotFound() {
		agent.Notify(status)
	}
}

func addAssignment(agent *caseOfficer, e guidance.HostEntry, newAgent newServiceAgent) {
	a := newAgent(e.Origin, agent)
	err := agent.serviceAgents.Register(a)
	if err != nil {
		agent.Notify(core.NewStatusError(core.StatusInvalidArgument, err))
	} else {
		a.Run()
	}
}
