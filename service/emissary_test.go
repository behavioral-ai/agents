package service

import (
	"fmt"
	"github.com/behavioral-ai/core/core"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/core/test"
	"github.com/behavioral-ai/resiliency/guidance"
)

var (
	emissaryShutdown = messaging.NewControlMessage(messaging.EmissaryChannel, "", messaging.ShutdownEvent)
	dataChange       = messaging.NewControlMessageWithBody("", "", messaging.DataChangeEvent, guidance.NewProcessingCalendar())
)

func ExampleEmissary() {
	ch := make(chan struct{})
	traceDispatch := messaging.NewTraceDispatcher(nil, "")
	agent := newAgent(core.Origin{Region: "us-west"}, test.NewAgent("agent-test"), traceDispatch, newMasterDispatcher(true), newEmissaryDispatcher(true))
	dataChange.SetContent(guidance.ContentTypeCalendar, guidance.NewProcessingCalendar())

	go func() {
		go emissaryAttend(agent, nil)
		agent.Message(dataChange)
		agent.Message(emissaryShutdown)
		fmt.Printf("test: emissaryAttend() -> [finalized:%v]\n", agent.emissary.isFinalized())
		ch <- struct{}{}
	}()
	<-ch
	close(ch)

	//Output:
	//test: Trace() -> service:us-west.. : emissary event:data-change Broadcast() -> calendar data change event]
	//test: emissaryAttend() -> [finalized:true]

}
