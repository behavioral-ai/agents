package service

import (
	"fmt"
	"github.com/advanced-go/common/core"
	"github.com/advanced-go/common/messaging"
	"github.com/advanced-go/common/test"
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
		fmt.Printf("test: masterAttend() -> [finalized:%v]\n", agent.isFinalizedMaster())
		ch <- struct{}{}
	}()
	<-ch
	close(ch)

	//Output:
	//test: masterAttend() -> [finalized:true]

}
