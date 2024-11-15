package common

import (
	"context"
	"github.com/advanced-go/common/core"
	"github.com/advanced-go/log/timeseries"
	"time"
)

const (
	timeseriesDuration = time.Second * 2
)

// Observation - observation functions struct, a nod to Linus Torvalds and plain C
type Observation struct {
	Timeseries func(h core.ErrorHandler, origin core.Origin) (timeseries.Entry, *core.Status)
}

var Observe = func() *Observation {
	return &Observation{
		Timeseries: func(h core.ErrorHandler, origin core.Origin) (timeseries.Entry, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), timeseriesDuration)
			defer cancel()
			e, status := timeseries.Get(ctx, origin)
			if !status.OK() && !status.NotFound() {
				h.Handle(status)
			}
			return e, status
		},
	}
}()
