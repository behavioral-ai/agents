package common

import (
	"context"
	"github.com/behavioral-ai/core/core"
	"github.com/behavioral-ai/core/messaging"
	"github.com/behavioral-ai/log/timeseries"
	"time"
)

const (
	timeseriesDuration = time.Second * 2
)

// Observation - observation functions struct
type Observation struct {
	Timeseries func(h messaging.Notifier, origin core.Origin) (timeseries.Entry, *core.Status)
}

var Observe = func() *Observation {
	return &Observation{
		Timeseries: func(h messaging.Notifier, origin core.Origin) (timeseries.Entry, *core.Status) {
			ctx, cancel := context.WithTimeout(context.Background(), timeseriesDuration)
			defer cancel()
			e, status := timeseries.Query(ctx, origin)
			if !status.OK() && !status.NotFound() {
				h.Notify(status)
			}
			return e, status
		},
	}
}()
