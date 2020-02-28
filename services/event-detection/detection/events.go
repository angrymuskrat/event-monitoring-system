package detection

import (
	"regexp"
	"strings"

	data "github.com/angrymuskrat/event-monitoring-system/services/proto"
	convtree "github.com/visheratin/conv-tree"
)

func FindEvents(histGrid convtree.ConvTree, posts []data.Post, maxPoints int, filterTags map[string]bool, start, finish int64) ([]data.Event, bool) {
	candGrid, wasFound := findCandidates(&histGrid, posts, maxPoints)
	if !wasFound {
		return nil, false
	}
	splitGrid(candGrid, maxPoints)
	events := treeEvents(candGrid, maxPoints, filterTags, start, finish)
	if len(events) == 0 {
		return nil, false
	}
	return events, true
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

func treeEvents(tree *convtree.ConvTree, maxPoints int, filterTags map[string]bool, start, finish int64) []data.Event {
	if tree.IsLeaf {
		result := []data.Event{}
		if len(tree.Points) >= maxPoints {
			evHolders, posts := eventHolders(tree.Points, filterTags)
			for _, e := range evHolders {
				event, ok := checkEvent(e, maxPoints, posts, start, finish)
				if ok {
					result = append(result, event)
				}

			}
		}
		return result
	} else {
		result := []data.Event{}
		result = append(result, treeEvents(tree.ChildBottomLeft, maxPoints, filterTags, start, finish)...)
		result = append(result, treeEvents(tree.ChildBottomRight, maxPoints, filterTags, start, finish)...)
		result = append(result, treeEvents(tree.ChildTopLeft, maxPoints, filterTags, start, finish)...)
		result = append(result, treeEvents(tree.ChildTopRight, maxPoints, filterTags, start, finish)...)
		return result
	}
}

func extractTags(post data.Post, filterTags map[string]bool) []string {
	tags := []string{}
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
		tags = append(tags, tag)
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
		tags = append(tags, tag)
	}
	return tags
}

func eventHolders(d []convtree.Point, filterTags map[string]bool) ([]eventHolder, []data.Post) {
	evHolders := []eventHolder{}
	posts := make([]data.Post, len(d))
	for pi, p := range d {
		post := p.Content.(data.Post)
		posts[pi] = post
		tags := extractTags(post, filterTags)
		if len(tags) == 0 {
			continue
		}
		h := eventHolder{
			users: map[string]bool{},
			posts: map[string]bool{},
			tags:  map[string]int{},
		}
		h.users[post.AuthorID] = true
		h.posts[post.Shortcode] = true
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
	res := []eventHolder{}
	for i := 0; i < len(evHolders); i++ {
		eh1 := evHolders[i]
		for j := i + 1; j < len(evHolders); j++ {
			eh2 := evHolders[j]
			found := false
			for t1 := range eh1.tags {
				for t2 := range eh2.tags {
					if t1 == t2 {
						eh1 = combineHolders(eh1, eh2)
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
	return res, posts
}

func combineHolders(eh1, eh2 eventHolder) eventHolder {
	res := eventHolder{
		users: map[string]bool{},
		posts: map[string]bool{},
		tags:  map[string]int{},
	}
	for u := range eh1.users {
		res.users[u] = true
	}
	for u := range eh2.users {
		res.users[u] = true
	}
	for p := range eh1.posts {
		res.posts[p] = true
	}
	for p := range eh2.posts {
		res.posts[p] = true
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

func checkEvent(e eventHolder, maxPoints int, posts []data.Post, start, finish int64) (data.Event, bool) {
	if len(e.users) < 2 {
		return data.Event{}, false
	}
	if len(e.users) < maxPoints/2 {
		return data.Event{}, false
	}
	for k, v := range e.tags {
		if v < len(e.posts)/3 {
			delete(e.tags, k)
		}
	}
	postCodes := []string{}
	for k := range e.posts {
		postCodes = append(postCodes, k)
	}
	tags := []string{}
	var max int
	var maxTag string
	for k := range e.tags {
		if e.tags[k] > max {
			max = e.tags[k]
			maxTag = k
		}
		tags = append(tags, k)
	}
	return data.Event{
		Center:    eventCenter(e.posts, posts),
		PostCodes: postCodes,
		Tags:      tags,
		Title:     maxTag,
		Start:     start,
		Finish:    finish,
	}, true
}

func eventCenter(codes map[string]bool, posts []data.Post) data.Point {
	points := map[convtree.Point]int{}
	for _, post := range posts {
		if _, ok := codes[post.Shortcode]; ok {
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
	}
	lat, lon := 0.0, 0.0
	for p, w := range points {
		lon += p.X * float64(w) / float64(len(codes))
		lat += p.Y * float64(w) / float64(len(codes))
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
