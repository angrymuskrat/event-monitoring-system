package crawler

import (
	"github.com/angrymuskrat/event-monitoring-system/services/insta-crawler/crawler/data"
	protodata "github.com/angrymuskrat/event-monitoring-system/services/proto"
)

type Parameters struct {
	CityID          string
	TopLeft         protodata.Point
	BottomRight     protodata.Point
	Type            data.CrawlingType
	Description     string
	Entities        []string
	FinishTimestamp int64
	DetailedPosts   bool
	LoadMedia       bool
	FixLocations    []Location
	Checkpoints     map[string]string
}
