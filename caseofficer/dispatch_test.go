package caseofficer

import (
	"fmt"
	"github.com/advanced-go/common/messaging"
)

type dispatchT struct{ channel string }

func newTestDispatcher() dispatcher {
	d := new(dispatchT)
	d.channel = messaging.EmissaryChannel
	return d
}

func (d *dispatchT) setup(_ *caseOfficer, _ string) {}

func (d *dispatchT) dispatch(agent *caseOfficer, event string) {
	switch event {
	case messaging.StartupEvent:
		agent.handler.Trace(agent, d.channel, event, fmt.Sprintf("count:%v", agent.serviceAgents.Count()))
	case messaging.DataChangeEvent:
		agent.handler.Trace(agent, d.channel, event, "Broadcast() -> calendar data change event")
	case messaging.ShutdownEvent:
		agent.handler.Trace(agent, d.channel, event, "")
	case messaging.TickEvent:
		agent.handler.Trace(agent, d.channel, event, "")
	}
}

func ExampleTestDispatcher() {
	fmt.Printf("test: TestDispatch() \n")

	//Output:
	//test: TestDispatch()

}
