package service

import (
	"github.com/angrymuskrat/event-monitoring-system/services/insta-crawler/crawler"
	"github.com/angrymuskrat/event-monitoring-system/services/insta-crawler/crawler/data"
)

type CrawlerService interface {
	New(p crawler.Parameters) (string, error)
	Status(id string) (crawler.OutStatus, error)
	Stop(id string) (bool, error)
	Entities(id string) ([]data.Entity, error)
	Posts(id, cursor string, num int) ([]data.Post, string, error)
}
