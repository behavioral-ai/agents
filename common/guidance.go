package common

import (
	"github.com/advanced-go/common/core"
)

// Guidance - guidance functions struct, a nod to Linus Torvalds and plain C
type Guidance struct {
	Assignments    func(h core.ErrorHandler, origin core.Origin) ([]HostEntry, *core.Status)
	NewAssignments func(h core.ErrorHandler, origin core.Origin) ([]HostEntry, *core.Status)
}

var Guide = func() *Guidance {
	return &Guidance{
		Assignments: func(h core.ErrorHandler, origin core.Origin) ([]HostEntry, *core.Status) {
			e, status := GetEntry(origin)
			if !status.OK() {
				h.Handle(status)
			}
			return []HostEntry{e}, status
		},
		NewAssignments: func(h core.ErrorHandler, origin core.Origin) ([]HostEntry, *core.Status) {
			return nil, core.StatusNotFound()
		},
	}
}()
