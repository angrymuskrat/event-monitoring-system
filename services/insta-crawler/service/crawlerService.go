package service

import (
	"github.com/angrymuskrat/event-monitoring-system/services/insta-crawler/crawler"
)

type crawlerService struct {
	crawler *crawler.Crawler
}

func (s *crawlerService) New(p crawler.Parameters) (id string, err error) {
	id, err = s.crawler.NewSession(p)
	return
}

func (s *crawlerService) Status(id string) (status crawler.OutStatus, err error) {
	status, err = s.crawler.Status(id)
	return
}

func (s *crawlerService) Stop(id string) (ok bool, err error) {
	ok, err = s.crawler.Stop(id)
	return
}
