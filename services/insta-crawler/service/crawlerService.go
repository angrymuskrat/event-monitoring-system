package service

import (
	"github.com/angrymuskrat/event-monitoring-system/services/insta-crawler/crawler"
	"github.com/angrymuskrat/event-monitoring-system/services/insta-crawler/crawler/data"
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

func (s *crawlerService) Entities(id string) (data []data.Entity, err error) {
	data, err = s.crawler.Entities(id)
	return
}

func (s *crawlerService) Posts(id, offset string, num int) (posts []data.Post, cursor string, err error) {
	if offset == "none" {
		offset = ""
	}
	posts, cursor, err = s.crawler.Posts(id, offset, num)
	return
}
