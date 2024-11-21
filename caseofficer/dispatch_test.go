package caseofficer

import (
	"fmt"
	"github.com/advanced-go/common/messaging"
)

type dispatchT struct{}

func newTestDispatcher() dispatcher {
	d := new(dispatchT)
	return d
}

func (d *dispatchT) setup(_ *caseOfficer, _, _ string) {}

func (d *dispatchT) dispatch(agent *caseOfficer, channel, event string) {
	switch event {
	case messaging.DataChangeEvent:
		agent.handler.Trace(agent, channel, event, "Broadcast() -> calendar data change event")
	}
}

func ExampleTestDispatcher() {
	fmt.Printf("test: TestDispatch() \n")

	//Output:
	//test: TestDispatch()

}
