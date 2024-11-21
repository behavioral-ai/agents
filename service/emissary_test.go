package service

import (
	"fmt"
	"github.com/advanced-go/common/core"
	"github.com/advanced-go/common/messaging"
	"github.com/advanced-go/common/test"
	"github.com/advanced-go/resiliency/guidance"
)

var (
	emissaryShutdown = messaging.NewControlMessage(messaging.EmissaryChannel, "", messaging.ShutdownEvent)
	dataChange       = messaging.NewControlMessageWithBody("", "", messaging.DataChangeEvent, guidance.NewProcessingCalendar())
)

func ExampleEmissary() {
	ch := make(chan struct{})
	agent := newAgent(core.Origin{Region: "us-west"}, test.NewAgent("agent-test"), newTestMasterDispatcher(), newTestEmissaryDispatcher())
	dataChange.SetContent(guidance.ContentTypeCalendar, guidance.NewProcessingCalendar())

	go func() {
		go emissaryAttend(agent, nil)
		agent.Message(dataChange)
		agent.Message(emissaryShutdown)
		fmt.Printf("test: emissaryAttend() -> [finalized:%v]\n", agent.isFinalizedEmissary())
		ch <- struct{}{}
	}()
	<-ch
	close(ch)

	//Output:
	//test: Trace() -> service:us-west.. : emissary event:data-change Broadcast() -> calendar data change event]
	//test: emissaryAttend() -> [finalized:true]

}
