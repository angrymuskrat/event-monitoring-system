package service

import data "github.com/angrymuskrat/event-monitoring-system/services/proto"

type HeatmapRequest struct {
	City        string     `json:"city"`
	TopLeft     data.Point `json:"top-left"`
	BottomRight data.Point `json:"bottom-right"`
	Hour        int64      `json:"hour"`
}

type TimelineRequest struct {
	City   string `json:"city"`
	Start  int64  `json:"start"`
	Finish int64  `json:"finish"`
}

type EventsRequest struct {
	City        string     `json:"city"`
	TopLeft     data.Point `json:"top-left"`
	BottomRight data.Point `json:"bottom-right"`
	Hour        int64      `json:"hour"`
}

type SearchRequest struct {
	City    string   `json:"city"`
	Keytags []string `json:"tags"`
	Start   int64    `json:"start"`
	Finish  int64    `json:"finish"`
}
