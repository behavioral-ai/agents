package service

import (
	"fmt"
	"github.com/behavioral-ai/core/core"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/core/test"
)

var (
	masterShutdown = messaging.NewControlMessage(messaging.MasterChannel, "", messaging.ShutdownEvent)
	observationMsg = messaging.NewControlMessage(messaging.MasterChannel, "", messaging.ObservationEvent)
)

func ExampleMaster() {
	ch := make(chan struct{})
	traceDispatch := messaging.NewTraceDispatcher(nil, "")
	agent := newAgent(core.Origin{Region: "us-west"}, test.NewAgent("agent-test"), traceDispatch, newMasterDispatcher(true), newEmissaryDispatcher(true))

	go func() {
		go masterAttend(agent)
		//agent.Message(observationMsg)
		agent.Message(masterShutdown)
		fmt.Printf("test: masterAttend() -> [finalized:%v]\n", agent.master.isFinalized())
		ch <- struct{}{}
	}()
	<-ch
	close(ch)

	//Output:
	//test: masterAttend() -> [finalized:true]

}
