package detection

import data "github.com/angrymuskrat/event-monitoring-system/services/proto"

type eventStatusType int

const (
	oldEventStatus eventStatusType = iota
	newEventStatus
	updatedEventStatus
)

type eventWithStatus struct {
	event  data.Event
	status eventStatusType
}

type eventHolder struct {
	users    map[string]bool
	posts    map[string]data.Post
	postTags map[string][]string
	tags     map[string]int
}
