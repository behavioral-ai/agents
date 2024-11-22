package service

import (
	"fmt"
	"github.com/advanced-go/common/messaging"
)

type emissaryT struct{ channel string }

func newTestEmissaryDispatcher() dispatcher {
	d := new(emissaryT)
	d.channel = messaging.EmissaryChannel
	return d
}

func (d *emissaryT) setup(_ *service, _ string) {}

func (d *emissaryT) dispatch(agent *service, event string) {
	switch event {
	case messaging.DataChangeEvent:
		agent.handler.Trace(agent, d.channel, event, "Broadcast() -> calendar data change event")
	}
}

func ExampleTestDispatcher() {
	fmt.Printf("test: TestDispatch() \n")

	//Output:
	//test: TestDispatch()

}
