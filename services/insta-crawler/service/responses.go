package service

import (
	"github.com/angrymuskrat/event-monitoring-system/services/insta-crawler/crawler"
	"github.com/angrymuskrat/event-monitoring-system/services/insta-crawler/crawler/data"
)

type NewEpResponse struct {
	ID    string `json:"id"`
	Error string `json:"error,omitempty"`
}

type StatusEpResponse struct {
	Status crawler.OutStatus `json:"status"`
	Error  string            `json:"error,omitempty"`
}

type StopEpResponse struct {
	Ok    bool   `json:"stopped"`
	Error string `json:"error,omitempty"`
}

type entitiesEpResponse struct {
	Entities []data.Entity `json:"entities"`
	Error    string        `json:"error,omitempty"`
}

type postsEpResponse struct {
	Posts  []data.Post `json:"posts"`
	Cursor string      `json:"cursor"`
	Error  string      `json:"error,omitempty"`
}
