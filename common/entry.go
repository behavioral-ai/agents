package common

import (
	"github.com/advanced-go/common/core"
	"time"
)

var (
	//safeEntry = common.NewSafe()
	entryData = []HostEntry{
		{Origin: core.Origin{Region: "us-west", Zone: "oregon", SubZone: "dc1", Host: "www.host1.com"}, CreatedTS: time.Date(2024, 6, 10, 7, 120, 35, 0, time.UTC)},
		{Origin: core.Origin{Region: "us-central", Zone: "iowa", SubZone: "dc2", Host: "www.host2.com"}, CreatedTS: time.Date(2024, 6, 10, 7, 120, 35, 0, time.UTC)},
	}
)

// HostEntry - host
type HostEntry struct {
	EntryId   int         `json:"entry-id"`
	CreatedTS time.Time   `json:"created-ts"`
	Origin    core.Origin `json:"origin"`
}

func GetEntry(origin core.Origin) (HostEntry, *core.Status) {
	if origin.Region == "us-west" {
		return entryData[0], core.StatusOK()
	}
	if origin.Region == "us-central" {
		return entryData[1], core.StatusOK()
	}
	return HostEntry{}, core.StatusNotFound()
}
