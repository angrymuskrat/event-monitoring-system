package service

import (
	"github.com/angrymuskrat/event-monitoring-system/services/insta-crawler/crawler"
	"github.com/angrymuskrat/event-monitoring-system/services/insta-crawler/crawler/data"
)

type newEpResponse struct {
	ID    string `json:"id"`
	Error string `json:"error,omitempty"`
}

type statusEpResponse struct {
	Status crawler.OutStatus `json:"status"`
	Error  string            `json:"error,omitempty"`
}

type stopEpResponse struct {
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
