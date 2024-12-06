package caseofficer

import (
	"fmt"
	"github.com/behavioral-ai/agents/service"
	"github.com/behavioral-ai/core/core"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/core/test"
	"github.com/behavioral-ai/resiliency/guidance"
	"time"
)

var (
	shutdown   = messaging.NewControlMessage("", "", messaging.ShutdownEvent)
	dataChange = messaging.NewControlMessage("", "", messaging.DataChangeEvent)
)

func ExampleEmissary() {
	ch := make(chan struct{})
	traceDispatch := messaging.NewTraceDispatcher([]string{messaging.StartupEvent, messaging.ShutdownEvent}, "")
	agent := newAgent(core.Origin{Region: guidance.WestRegion}, test.NewAgent("agent-test"), traceDispatch, newDispatcher(false))
	dataChange.SetContent(guidance.ContentTypeCalendar, guidance.NewProcessingCalendar())

	go func() {
		go emissaryAttend(agent, guidance.Assign, service.NewAgent, newFeedbackAgent)
		agent.Message(dataChange)
		time.Sleep(time.Minute * 1)
		agent.Message(shutdown)

		fmt.Printf("test: emissaryAttend() -> [finalized:%v]\n", agent.IsFinalized())
		ch <- struct{}{}
	}()
	<-ch
	close(ch)

	//Output:
	//test: emissaryAttend() -> [finalized:true]

}
