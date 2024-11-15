package common

import (
	"errors"
	"fmt"
	"github.com/advanced-go/common/core"
	"github.com/advanced-go/common/messaging"
)

func MessageEventErrorStatus(agentId string, msg *messaging.Message) *core.Status {
	err := errors.New(fmt.Sprintf("error: message event:%v is invalid for agent:%v", msg.Event(), agentId))
	return core.NewStatusError(core.StatusInvalidArgument, err)
}

func MessageContentTypeErrorStatus(agentId string, msg *messaging.Message) *core.Status {
	err := errors.New(fmt.Sprintf("error: message content:%v is invalid for agent:%v and event:%v", msg.ContentType(), agentId, msg.Event()))
	return core.NewStatusError(core.StatusInvalidArgument, err)
}
