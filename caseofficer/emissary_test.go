package caseofficer

import (
	"fmt"
	"github.com/advanced-go/common/core"
	"github.com/advanced-go/common/messaging"
	"github.com/advanced-go/common/test"
	"github.com/advanced-go/resiliency/guidance"
)

var (
	shutdownMsg   = messaging.NewControlMessage("", "", messaging.ShutdownEvent)
	dataChangeMsg = messaging.NewControlMessage("", "", messaging.DataChangeEvent)
)

func init() {
	dataChangeMsg.SetContent(guidance.ContentTypeCalendar, guidance.NewProcessingCalendar())
}

func ExampleEmissary() {
	ch := make(chan struct{})
	agent := newAgent(core.Origin{Region: "us-west"}, test.NewAgent("agent-test"), newTestDispatcher())

	go func() {
		go emissaryAttend(agent, nil, nil, nil)
		//agent.Message(dataChangeMsg)
		agent.Message(shutdownMsg)
		fmt.Printf("test: emissaryAttend() -> [finalized:%v]\n", agent.IsFinalized())
		ch <- struct{}{}
	}()
	<-ch
	close(ch)

	//Output:
	//test: emissaryAttend() -> [finalized:true]

}
