package detection

import (
	"regexp"
	"sort"
	"strings"

	data "github.com/angrymuskrat/event-monitoring-system/services/proto"
	convtree "github.com/visheratin/conv-tree"
)

const closeDistanse = 0.005

func FindEvents(histGrid convtree.ConvTree, posts []data.Post, maxPoints int, filterTags map[string]bool, start, finish int64, existingEvents []data.Event) (newEvents, updatedEvents []data.Event, found bool) {
	candGrid, wasFound := findCandidates(&histGrid, posts, maxPoints)
	if !wasFound {
		return
	}
	splitGrid(candGrid, maxPoints)
	newEvents, updatedEvents = treeEvents(candGrid, maxPoints, filterTags, start, finish, existingEvents)
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

func treeEvents(tree *convtree.ConvTree, maxPoints int, filterTags map[string]bool, start, finish int64, existingEvents []data.Event) (newEvents, updatedEvents []data.Event) {
	if tree.IsLeaf {
		if len(tree.Points) >= maxPoints {
			evHolders := eventHolders(tree.Points, filterTags)
			for _, e := range evHolders {
				event, isEvent, isNew := checkEvent(e, maxPoints, start, finish, existingEvents)
				if isEvent {
					if isNew {
						newEvents = append(newEvents, event)
					} else {
						updatedEvents = append(updatedEvents, event)
					}
				}

			}
		}
		return
	} else {
		var eventsBL []data.Event
		var eventsBR []data.Event
		var eventsTL []data.Event
		var eventsTR []data.Event
		del := tree.ChildBottomRight.TopLeft
		for _, event := range existingEvents {
			if event.Center.Lat > del.Y {
				if event.Center.Lon > del.X {
					eventsTR = append(eventsTR, event)
				} else {
					eventsTL = append(eventsTL, event)
				}
			} else {
				if event.Center.Lon > del.X {
					eventsBR = append(eventsBR, event)
				} else {
					eventsBL = append(eventsBL, event)
				}
			}
		}

		newEventsBL, updatedEventsBL := treeEvents(tree.ChildBottomLeft, maxPoints, filterTags, start, finish, eventsBL)
		newEvents = append(newEvents, newEventsBL...)
		updatedEvents = append(updatedEvents, updatedEventsBL...)

		border := tree.ChildBottomRight.TopLeft.X - closeDistanse
		for _, newEventBL := range newEventsBL {
			if newEventBL.Center.Lon > border {
				eventsBR = append(eventsBR, newEventBL)
			}
		}
		for _, updatedEventBL := range updatedEventsBL {
			if updatedEventBL.Center.Lon > border {
				eventsBR = append(eventsBR, updatedEventBL)
			}
		}
		newEventsBR, updatedEventsBR := treeEvents(tree.ChildBottomRight, maxPoints, filterTags, start, finish, eventsBR)
		newEvents = append(newEvents, newEventsBR...)
		updatedEvents = append(updatedEvents, updatedEventsBR...)

		border = tree.ChildTopLeft.BottomRight.Y - closeDistanse
		for _, newEventBL := range newEventsBL {
			if newEventBL.Center.Lat > border {
				eventsTL = append(eventsTL, newEventBL)
			}
		}
		for _, updatedEventBL := range updatedEventsBL {
			if updatedEventBL.Center.Lat > border {
				eventsTL = append(eventsTL, updatedEventBL)
			}
		}

		for _, newEventBR := range newEventsBR {
			if newEventBR.Center.Lat > border {
				eventsTL = append(eventsTL, newEventBR)
			}
		}
		for _, updatedEventBR := range updatedEventsBR {
			if updatedEventBR.Center.Lat > border {
				eventsTL = append(eventsTL, updatedEventBR)
			}
		}
		newEventsTL, updatedEventsTL := treeEvents(tree.ChildTopLeft, maxPoints, filterTags, start, finish, eventsTL)
		newEvents = append(newEvents, newEventsTL...)
		updatedEvents = append(updatedEvents, updatedEventsTL...)

		border = tree.ChildTopRight.BottomRight.Y - closeDistanse
		for _, newEventBL := range newEventsBL {
			if newEventBL.Center.Lat > border {
				eventsTR = append(eventsTR, newEventBL)
			}
		}
		for _, updatedEventBL := range updatedEventsBL {
			if updatedEventBL.Center.Lat > border {
				eventsTR = append(eventsTR, updatedEventBL)
			}
		}

		for _, newEventBR := range newEventsBR {
			if newEventBR.Center.Lat > border {
				eventsTR = append(eventsTR, newEventBR)
			}
		}
		for _, updatedEventBR := range updatedEventsBR {
			if updatedEventBR.Center.Lat > border {
				eventsTR = append(eventsTR, updatedEventBR)
			}
		}

		border = tree.ChildTopRight.TopLeft.X - closeDistanse
		for _, newEventTL := range newEventsTL {
			if newEventTL.Center.Lon > border {
				eventsTR = append(eventsTR, newEventTL)
			}
		}
		for _, updatedEventTL := range updatedEventsTL {
			if updatedEventTL.Center.Lon > border {
				eventsTR = append(eventsTR, updatedEventTL)
			}
		}
		newEventsTR, updatedEventsTR := treeEvents(tree.ChildTopRight, maxPoints, filterTags, start, finish, eventsTR)
		newEvents = append(newEvents, newEventsTR...)
		updatedEvents = append(updatedEvents, updatedEventsTR...)
		return
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

func checkEvent(e eventHolder, maxPoints int, start, finish int64, existingEvents []data.Event) (event data.Event, isEvent, isNew bool) {
	if len(e.users) < 2 {
		return
	}
	if len(e.users) < maxPoints/2 {
		return
	}

	resultStart := start
	isNew = true
	var id int64
SEARCH:
	for _, oldEvent := range existingEvents {
		for _, oldTag := range oldEvent.Tags {
			for newTag := range e.tags {
				if oldTag == newTag {
					for _, oldPostCode := range oldEvent.PostCodes {
						if _, ok := e.posts[oldPostCode]; !ok {
							for _, existingEventsPost := range oldEvent.Posts {
								if existingEventsPost.Shortcode == oldPostCode {
									e.posts[oldPostCode] = *existingEventsPost
									break
								}
							}
						} else {
							for _, oldPostTag := range e.postTags[oldPostCode] {
								e.tags[oldPostTag]--
							}
						}
					}

					for i, ot := range oldEvent.Tags {
						if _, ok := e.tags[ot]; ok {
							e.tags[ot] += int(oldEvent.TagsCount[i])
						} else {
							e.tags[ot] = int(oldEvent.TagsCount[i])
						}
					}

					id = oldEvent.ID
					resultStart = oldEvent.Start
					isNew = false
					break SEARCH
				}
			}
		}
	}

	postCodes := []string{}
	for k := range e.posts {
		postCodes = append(postCodes, k)
	}
	tags, counts := sortTags(e.tags, 5)

	event = data.Event{
		ID:        id,
		Center:    eventCenter(e.posts),
		PostCodes: postCodes,
		Tags:      tags,
		TagsCount: counts,
		Title:     tags[0],
		Start:     resultStart,
		Finish:    finish,
	}
	isEvent = true
	return
}

func sortTags(tags map[string]int, max int) ([]string, []int64) {
	rev := map[int][]string{}
	for t, c := range tags {
		if _, ok := rev[c]; !ok {
			rev[c] = []string{t}
		} else {
			rev[c] = append(rev[c], t)
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
		//if len(res) >= max {
		//	break
		//}
	}
	counts := make([]int64, len(res))
	for ind, t := range res {
		counts[ind] = int64(tags[t])
	}
	//l := max
	//if l > (len(res) - 1) {
	//	l = len(res)
	//}
	//res = res[0:l]
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

func filterTag(tag string, filterTags map[string]bool) bool {
	_, ok := filterTags[tag]
	return ok
}
