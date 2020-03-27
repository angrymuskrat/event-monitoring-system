package detection

import data "github.com/angrymuskrat/event-monitoring-system/services/proto"

type eventHolder struct {
	users    map[string]bool
	posts    map[string]data.Post
	postTags map[string][]string
	tags     map[string]int
}
