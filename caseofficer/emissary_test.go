package caseofficer

import (
	"fmt"
	"github.com/advanced-go/agents/service"
	"github.com/advanced-go/common/core"
	"github.com/advanced-go/common/messaging"
	"github.com/advanced-go/common/test"
	"github.com/advanced-go/resiliency/guidance"
	"time"
)

var (
	shutdown   = messaging.NewControlMessage("", "", messaging.ShutdownEvent)
	dataChange = messaging.NewControlMessage("", "", messaging.DataChangeEvent)
)

func ExampleEmissary() {
	ch := make(chan struct{})
	agent := newAgent(core.Origin{Region: "us-west"}, test.NewAgent("agent-test"), newTestDispatcher())
	dataChange.SetContent(guidance.ContentTypeCalendar, guidance.NewProcessingCalendar())

	go func() {
		go emissaryAttend(agent, guidance.Assign, service.NewAgent)
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
