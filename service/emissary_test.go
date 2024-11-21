package service

import (
	"fmt"
	"github.com/advanced-go/common/core"
	"github.com/advanced-go/common/messaging"
	"github.com/advanced-go/common/test"
	"github.com/advanced-go/resiliency/guidance"
)

var (
	emissaryShutdownMsg = messaging.NewLeftChannelMessage("", "", messaging.ShutdownEvent, nil)
	dataChangeMsg       = messaging.NewControlMessage("", "", messaging.DataChangeEvent)
)

func init() {
	dataChangeMsg.SetContent(guidance.ContentTypeCalendar, guidance.NewProcessingCalendar())
}

func ExampleEmissary() {
	ch := make(chan struct{})
	agent := newAgent(core.Origin{Region: "us-west"}, test.NewAgent("agent-test"), newTestDispatcher())

	go func() {
		go emissaryAttend(agent, nil)
		//agent.Message(dataChangeMsg)
		agent.Message(emissaryShutdownMsg)
		fmt.Printf("test: emissaryAttend() -> [finalized:%v]\n", agent.isFinalizedEmissary())
		ch <- struct{}{}
	}()
	<-ch
	close(ch)

	//Output:
	//test: emissaryAttend() -> [finalized:true]

}
