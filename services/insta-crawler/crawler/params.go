package crawler

import (
	"github.com/angrymuskrat/event-monitoring-system/services/insta-crawler/crawler/data"
	"github.com/angrymuskrat/event-monitoring-system/services/insta-crawler/crawler/storage"
)

type Parameters struct {
	CityID          string
	Type            data.CrawlingType
	Description     string
	Entities        []string
	FinishTimestamp int64
	DetailedPosts   bool
	LoadMedia       bool
	Reupload        bool
	FixLocations    []storage.Location
}
