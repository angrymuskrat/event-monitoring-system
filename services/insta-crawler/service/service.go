package service

import (
	"github.com/angrymuskrat/event-monitoring-system/services/insta-crawler/crawler"
)

type CrawlerService interface {
	New(p crawler.Parameters) (string, error)
	Status(id string) (crawler.OutStatus, error)
	Stop(id string) (bool, error)
}
