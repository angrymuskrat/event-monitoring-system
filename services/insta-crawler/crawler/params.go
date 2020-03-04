package crawler

import "github.com/angrymuskrat/event-monitoring-system/services/insta-crawler/crawler/data"

type Parameters struct {
	CityID          string
	Type            data.CrawlingType
	Description     string
	Entities        []string
	FinishTimestamp int64
	DetailedPosts   bool
	LoadMedia       bool
}
