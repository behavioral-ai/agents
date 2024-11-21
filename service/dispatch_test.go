package service

import (
	"fmt"
	"github.com/advanced-go/common/messaging"
)

type masterT struct{}

func newTestMasterDispatcher() dispatcher {
	d := new(masterT)
	return d
}

func (d *masterT) setup(_ *service, _ string) {}

func (d *masterT) dispatch(agent *service, event string) {
	switch event {
	case messaging.DataChangeEvent:
		agent.handler.Trace(agent, messaging.MasterChannel, event, "Broadcast() -> calendar data change event")
	}
}

type emissaryT struct{}

func newTestEmissaryDispatcher() dispatcher {
	d := new(emissaryT)
	return d
}

func (d *emissaryT) setup(_ *service, _ string) {}

func (d *emissaryT) dispatch(agent *service, event string) {
	switch event {
	case messaging.DataChangeEvent:
		agent.handler.Trace(agent, messaging.EmissaryChannel, event, "Broadcast() -> calendar data change event")
	}
}

func ExampleTestDispatcher() {
	fmt.Printf("test: TestDispatch() \n")

	//Output:
	//test: TestDispatch()

}
