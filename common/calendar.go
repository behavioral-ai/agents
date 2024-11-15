package common

import (
	"github.com/advanced-go/common/core"
	"github.com/advanced-go/common/messaging"
)

type ProcessingCalendar struct {
	week [7][24]string
}

func NewProcessingCalendar() *ProcessingCalendar {
	c := new(ProcessingCalendar)
	return c
}

func GetCalendar(h core.ErrorHandler, agentId string, msg *messaging.Message) *ProcessingCalendar {
	if !msg.IsContentType(ContentTypeCalendar) {
		return nil
	}
	if p, ok := msg.Body.(*ProcessingCalendar); ok {
		return p
	}
	h.Handle(ProfileTypeErrorStatus(agentId, msg.Body))
	return nil
}
