package detection

import (
	"regexp"
	"sort"
	"strings"

	data "github.com/angrymuskrat/event-monitoring-system/services/proto"
	convtree "github.com/visheratin/conv-tree"
)

const closeDistanse = 0.005

func FindEvents(histGrid convtree.ConvTree, posts []data.Post, maxPoints int, filterTags map[string]bool, start, finish int64, existingEvents []data.Event) (newEvents, updatedEvents, deletedEvents []data.Event, found bool) {
	candGrid, wasFound := findCandidates(&histGrid, posts, maxPoints)
	if !wasFound {
		return
	}
	splitGrid(candGrid, maxPoints)

	existingEventsWithStatuses := make([]eventWithStatus, len(existingEvents))
	for i, existingEvent := range existingEvents {
		existingEventsWithStatuses[i] = eventWithStatus{
			event:  existingEvent,
			status: oldEventStatus,
		}
	}
	treeEvents(candGrid, maxPoints, filterTags, start, finish, &existingEventsWithStatuses)
	for _, existingEvent := range existingEventsWithStatuses {
		if existingEvent.status == newEventStatus {
			newEvents = append(newEvents, existingEvent.event)
		}
		if existingEvent.status == updatedEventStatus {
			updatedEvents = append(updatedEvents, existingEvent.event)
		}
		if existingEvent.status == deletedEventStatus {
			deletedEvents = append(deletedEvents, existingEvent.event)
		}
	}
	if len(newEvents)+len(updatedEvents) != 0 {
		found = true
	}
	return
}

func splitGrid(tree *convtree.ConvTree, maxPoints int) {
	if tree.IsLeaf {
		if len(tree.Points) >= maxPoints {
			tree.Check()
		}
	} else {
		splitGrid(tree.ChildBottomLeft, maxPoints)
		splitGrid(tree.ChildBottomRight, maxPoints)
		splitGrid(tree.ChildTopLeft, maxPoints)
		splitGrid(tree.ChildTopRight, maxPoints)
	}
}

func treeEvents(tree *convtree.ConvTree, maxPoints int, filterTags map[string]bool, start, finish int64, existingEvents *[]eventWithStatus) {
	if tree.IsLeaf {
		for _, baselineTag := range tree.BaselineTags {
			filterTags[baselineTag] = true
		}
		if len(tree.Points) >= maxPoints {
			evHolders := eventHolders(tree.Points, filterTags)
			for _, e := range evHolders {
				checkEvent(e, maxPoints, start, finish, tree.TopLeft, tree.BottomRight, existingEvents)
			}
		}
	} else {
		treeEvents(tree.ChildBottomLeft, maxPoints, filterTags, start, finish, existingEvents)
		treeEvents(tree.ChildBottomRight, maxPoints, filterTags, start, finish, existingEvents)
		treeEvents(tree.ChildTopLeft, maxPoints, filterTags, start, finish, existingEvents)
		treeEvents(tree.ChildTopRight, maxPoints, filterTags, start, finish, existingEvents)
	}
}

func extractTags(post data.Post, filterTags map[string]bool) []string {
	tags := []string{}
	mapTag := map[string]bool{}
	r, _ := regexp.Compile("#[a-zA-Z0-9а-яА-Я_]+")
	hashtags := r.FindAllString(post.Caption, -1)
	for idx := range hashtags {
		hashtags[idx] = strings.ToLower(hashtags[idx])
	}
	for _, tag := range hashtags {
		tag := strings.ToLower(tag)
		if filterTag(tag, filterTags) {
			continue
		}
		if _, ok := mapTag[tag]; !ok {
			tags = append(tags, tag)
			mapTag[tag] = true
		}
	}
	r, _ = regexp.Compile("@[a-zA-Z0-9а-яА-Я_]+")
	usernames := r.FindAllString(post.Caption, -1)
	for idx := range usernames {
		usernames[idx] = strings.ToLower(usernames[idx])
	}
	for _, tag := range usernames {
		if filterTag(tag, filterTags) {
			continue
		}
		if _, ok := mapTag[tag]; !ok {
			tags = append(tags, tag)
			mapTag[tag] = true
		}
	}
	return tags
}

func eventHolders(d []convtree.Point, filterTags map[string]bool) []eventHolder {
	evHolders := []eventHolder{}
	for _, p := range d {
		post := p.Content.(data.Post)

		tags := extractTags(post, filterTags)
		if len(tags) == 0 {
			continue
		}
		var duplicate bool
		for _, evHolder := range evHolders {
			for _, ehp := range evHolder.posts {
				if post.Caption == ehp.Caption && tags[0] != post.Caption {
					duplicate = true
					break
				}
			}
			if duplicate {
				break
			}
		}
		if duplicate {
			continue
		}

		h := eventHolder{
			users:    map[string]bool{},
			posts:    map[string]data.Post{},
			postTags: map[string][]string{},
			tags:     map[string]int{},
		}
		h.users[post.AuthorID] = true
		h.posts[post.Shortcode] = post
		h.postTags[post.Shortcode] = tags
		for _, tag := range tags {
			h.tags[tag] = 1
		}
		found := false
		for i := range evHolders {
			if _, ok := evHolders[i].users[post.AuthorID]; ok {
				h = combineHolders(evHolders[i], h)
				evHolders[i] = h
				found = true
			}
		}
		if !found {
			evHolders = append(evHolders, h)
		}
	}
	res, un := uniteHolders(evHolders)
	for un {
		res, un = uniteHolders(res)
	}
	return res
}

func uniteHolders(evHolders []eventHolder) ([]eventHolder, bool) {
	united := false
	res := []eventHolder{}
	for i := 0; i < len(evHolders); i++ {
		eh1 := evHolders[i]
		for j := i + 1; j < len(evHolders); j++ {
			if j == (len(evHolders) - 1) {
				continue
			}
			eh2 := evHolders[j]
			found := false
			for t1 := range eh1.tags {
				for t2 := range eh2.tags {
					if t1 == t2 {
						eh1 = combineHolders(eh1, eh2)
						united = true
						found = true
						break
					}
				}
				if found {
					break
				}
			}
			if found {
				evHolders = append(evHolders[:j], evHolders[j+1:]...)
				j--
			}
		}
		res = append(res, eh1)
	}
	return res, united
}

func combineHolders(eh1, eh2 eventHolder) eventHolder {
	res := eventHolder{
		users:    map[string]bool{},
		posts:    map[string]data.Post{},
		postTags: map[string][]string{},
		tags:     map[string]int{},
	}
	for u := range eh1.users {
		res.users[u] = true
	}
	for u := range eh2.users {
		res.users[u] = true
	}
	for p := range eh1.posts {
		res.posts[p] = eh1.posts[p]
		res.postTags[p] = eh1.postTags[p]

	}
	for p := range eh2.posts {
		if _, ok := res.posts[p]; !ok {
			res.posts[p] = eh2.posts[p]
			res.postTags[p] = eh2.postTags[p]
		}
	}
	for t, c := range eh1.tags {
		if _, ok := res.tags[t]; ok {
			res.tags[t] += c
		} else {
			res.tags[t] = c
		}
	}
	for t, c := range eh2.tags {
		if _, ok := res.tags[t]; ok {
			res.tags[t] += c
		} else {
			res.tags[t] = c
		}
	}
	return res
}

func checkEvent(e eventHolder, maxPoints int, start, finish int64, topLeft, bottomRight convtree.Point, existingEvents *[]eventWithStatus) {
	if len(e.users) < 2 {
		return
	}
	if len(e.users) < maxPoints/2 {
		return
	}

	updatedEventIndex := -1
	for i, oldEvent := range *existingEvents {
		if oldEvent.event.Center.Lat > bottomRight.Y-closeDistanse &&
			oldEvent.event.Center.Lat < topLeft.Y+closeDistanse &&
			oldEvent.event.Center.Lon > topLeft.X-closeDistanse &&
			oldEvent.event.Center.Lon < bottomRight.X+closeDistanse &&
			oldEvent.status != deletedEventStatus &&
			oldEvent.status != ignoredEventStatus {
		NEXTEVENT:
			for _, oldTag := range oldEvent.event.Tags {
				for newTag := range e.tags {
					if oldTag == newTag {
						if updatedEventIndex == -1 {
							event := combineHolderAndEvent(e, oldEvent, finish)
							(*existingEvents)[i] = event
							updatedEventIndex = i
						} else {
							event := combineTwoEvents((*existingEvents)[updatedEventIndex], oldEvent, e.posts)

							if oldEvent.status == newEventStatus {
								(*existingEvents)[i].status = ignoredEventStatus
								(*existingEvents)[updatedEventIndex] = event
							} else {
								if oldEvent.event.Start < (*existingEvents)[updatedEventIndex].event.Start {
									(*existingEvents)[updatedEventIndex].status = deletedEventStatus
									(*existingEvents)[i] = event
									updatedEventIndex = i
								} else {
									(*existingEvents)[i].status = deletedEventStatus
									(*existingEvents)[updatedEventIndex] = event
								}
							}
						}
						break NEXTEVENT
					}
				}
			}
		}
	}

	if updatedEventIndex == -1 {
		postCodes := []string{}
		for k := range e.posts {
			postCodes = append(postCodes, k)
		}
		tags, counts := sortTags(e.tags)

		event := eventWithStatus{
			event: data.Event{
				Center:    eventCenter(e.posts),
				PostCodes: postCodes,
				Tags:      tags,
				TagsCount: counts,
				Title:     tags[0],
				Start:     start,
				Finish:    finish,
			},
			status: newEventStatus,
		}
		*existingEvents = append(*existingEvents, event)
	}

}

func combineHolderAndEvent(e eventHolder, oldEvent eventWithStatus, finish int64) eventWithStatus {

	for i, ot := range oldEvent.event.Tags {
		if _, ok := e.tags[ot]; ok {
			e.tags[ot] += int(oldEvent.event.TagsCount[i])
		} else {
			e.tags[ot] = int(oldEvent.event.TagsCount[i])
		}
	}

	postCodes := []string{}
	mapPosts := make(map[string]bool)

	for _, oldPostCode := range oldEvent.event.PostCodes {
		postCodes = append(postCodes, oldPostCode)
		mapPosts[oldPostCode] = true
		if _, ok := e.posts[oldPostCode]; ok {
			for _, oldPostTag := range e.postTags[oldPostCode] {
				for _, oldEventTag := range oldEvent.event.Tags {
					if oldEventTag == oldPostTag {
						e.tags[oldPostTag]--
						break
					}
				}
			}
		}
	}

	newPosts := make(map[string]data.Post)
	for k, v := range e.posts {
		if _, ok := mapPosts[k]; !ok {
			postCodes = append(postCodes, k)
			newPosts[k] = v
		}
	}

	newPostsCenter := eventCenter(newPosts)
	tags, counts := sortTags(e.tags)

	var status eventStatusType
	if oldEvent.status == newEventStatus {
		status = newEventStatus
	} else {
		status = updatedEventStatus
	}
	return eventWithStatus{
		event: data.Event{
			ID:        oldEvent.event.ID,
			Center:    combineCenters(oldEvent.event.Center, newPostsCenter, len(oldEvent.event.PostCodes), len(newPosts)),
			PostCodes: postCodes,
			Tags:      tags,
			TagsCount: counts,
			Title:     tags[0],
			Start:     oldEvent.event.Start,
			Finish:    finish,
		},
		status: status,
	}
}

func combineTwoEvents(event1, event2 eventWithStatus, posts map[string]data.Post) eventWithStatus {
	var resultEvent eventWithStatus
	if event1.status == newEventStatus && event2.status == newEventStatus {
		resultEvent.status = newEventStatus
	} else {
		resultEvent.status = updatedEventStatus
	}

	if event2.event.Start < event1.event.Start {
		resultEvent.event.ID = event2.event.ID
		resultEvent.event.Start = event2.event.Start
	} else {
		resultEvent.event.ID = event1.event.ID
		resultEvent.event.Start = event1.event.Start
	}

	if event2.event.Finish > event1.event.Finish {
		resultEvent.event.Finish = event2.event.Finish
	} else {
		resultEvent.event.Finish = event1.event.Finish
	}

	mapPosts := make(map[string]bool)
	var postcodes []string
	for _, postcode := range event1.event.PostCodes {
		postcodes = append(postcodes, postcode)
		mapPosts[postcode] = true
	}

	center2 := event2.event.Center
	for _, postcode := range event2.event.PostCodes {
		if _, ok := mapPosts[postcode]; !ok {
			postcodes = append(postcodes, postcode)
		} else {
			if _, ok := posts[postcode]; ok {
				center2 = data.Point{
					Lat: (event2.event.Center.Lat*float64(len(event2.event.PostCodes)) - posts[postcode].Lat) / (float64(len(event2.event.PostCodes)) - 1),
					Lon: (event2.event.Center.Lon*float64(len(event2.event.PostCodes)) - posts[postcode].Lon) / (float64(len(event2.event.PostCodes)) - 1),
				}
			}
		}
	}
	resultEvent.event.PostCodes = postcodes
	resultEvent.event.Center = combineCenters(event1.event.Center, center2, len(event1.event.PostCodes), len(event2.event.PostCodes))

	tagsMap := make(map[string]int)
	for i, tag := range event1.event.Tags {
		tagsMap[tag] = int(event1.event.TagsCount[i])
	}
	for i, tag := range event2.event.Tags {
		if _, ok := tagsMap[tag]; ok {
			tagsMap[tag] += int(event2.event.TagsCount[i])
		} else {
			tagsMap[tag] = int(event2.event.TagsCount[i])
		}
	}
	tags, counts := sortTags(tagsMap)
	resultEvent.event.Tags = tags
	resultEvent.event.TagsCount = counts
	resultEvent.event.Title = tags[0]
	return resultEvent
}

func sortTags(tags map[string]int) ([]string, []int64) {
	rev := map[int][]string{}
	for t, c := range tags {
		if c != 1 {
			if _, ok := rev[c]; !ok {
				rev[c] = []string{t}
			} else {
				rev[c] = append(rev[c], t)
			}
		}
	}
	keys := make([]int, 0, len(rev))
	for k := range rev {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] > keys[j]
	})
	res := []string{}
	for _, k := range keys {
		res = append(res, rev[k]...)
	}
	counts := make([]int64, len(res))
	for ind, t := range res {
		counts[ind] = int64(tags[t])
	}
	return res, counts
}

func eventCenter(posts map[string]data.Post) data.Point {
	points := map[convtree.Point]int{}
	for _, post := range posts {
		p := convtree.Point{
			X: post.Lon,
			Y: post.Lat,
		}
		if _, ok := points[p]; ok {
			points[p]++
		} else {
			points[p] = 1
		}
	}
	lat, lon := 0.0, 0.0
	for p, w := range points {
		lon += p.X * float64(w) / float64(len(posts))
		lat += p.Y * float64(w) / float64(len(posts))
	}
	return data.Point{
		Lat: lat,
		Lon: lon,
	}
}

func combineCenters(center1, center2 data.Point, num1, num2 int) data.Point {
	return data.Point{
		Lat: (center1.Lat*float64(num1) + center2.Lat*float64(num2)) / float64((num1 + num2)),
		Lon: (center1.Lon*float64(num1) + center2.Lon*float64(num2)) / float64((num1 + num2)),
	}
}

func filterTag(tag string, filterTags map[string]bool) bool {
	_, ok := filterTags[tag]
	return ok
}
