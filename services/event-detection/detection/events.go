package detection

import (
	"regexp"
	"sort"
	"strings"

	convtree "github.com/angrymuskrat/conv-tree"
	data "github.com/angrymuskrat/event-monitoring-system/services/proto"
)

func FindEvents(histGrid convtree.ConvTree, posts []data.Post, maxPoints float64, minUser int,
	filterTags map[string]bool, start, finish int64, anomalyMode bool) ([]data.Event, bool) {
	candidateGrid, wasFound := findCandidates(&histGrid, posts, maxPoints)
	if !wasFound {
		return nil, false
	}
	splitGrid(candidateGrid, maxPoints)
	events := treeEvents(candidateGrid, maxPoints, minUser, filterTags, start, finish, anomalyMode)
	if len(events) == 0 {
		return nil, false
	}
	return events, true
}

func splitGrid(tree *convtree.ConvTree, maxPoints float64) {
	if tree.IsLeaf {
		if sumWeightPoints(tree.Points) >= maxPoints {
			tree.Check()
		}
	} else {
		splitGrid(tree.ChildBottomLeft, maxPoints)
		splitGrid(tree.ChildBottomRight, maxPoints)
		splitGrid(tree.ChildTopLeft, maxPoints)
		splitGrid(tree.ChildTopRight, maxPoints)
	}
}

func treeEvents(tree *convtree.ConvTree, maxPoints float64, minUser int, filterTags map[string]bool,
	start, finish int64, anomalyMode bool) []data.Event {
	if tree.IsLeaf {
		var result []data.Event
		if sumWeightPoints(tree.Points) >= maxPoints {
			eventHolders, posts := eventHolders(tree.Points, filterTags, anomalyMode)
			for _, holder := range eventHolders {
				event, ok := checkEvent(holder, minUser, posts, start, finish)
				if ok {
					result = append(result, event)
				}

			}
		}
		return result
	} else {
		var result []data.Event
		result = append(result, treeEvents(tree.ChildBottomLeft, maxPoints, minUser, filterTags, start, finish, anomalyMode)...)
		result = append(result, treeEvents(tree.ChildBottomRight, maxPoints, minUser, filterTags, start, finish, anomalyMode)...)
		result = append(result, treeEvents(tree.ChildTopLeft, maxPoints, minUser, filterTags, start, finish, anomalyMode)...)
		result = append(result, treeEvents(tree.ChildTopRight, maxPoints, minUser, filterTags, start, finish, anomalyMode)...)
		return result
	}
}

func extractTags(post data.Post, filterTags map[string]bool) []string {
	var tags []string
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

func eventHolders(leafPoints []convtree.Point, filterTags map[string]bool, anomalyMode bool) ([]eventHolder, []data.Post) {
	var holders []eventHolder
	posts := make([]data.Post, len(leafPoints))
	for pointInd, point := range leafPoints {
		post := point.Content.(data.Post)
		posts[pointInd] = post
		tags := extractTags(post, filterTags)
		if (len(tags) == 0) && !anomalyMode {
			continue
		}
		holder := eventHolder{
			users: map[string]bool{},
			posts: map[string]bool{},
			tags:  map[string]int{},
		}
		holder.users[post.AuthorID] = true
		holder.posts[post.Shortcode] = true
		for _, tag := range tags {
			holder.tags[tag] = 1
		}
		found := false
		for i := range holders {
			if _, ok := holders[i].users[post.AuthorID]; ok {
				holder = combineHolders(holders[i], holder)
				holders[i] = holder
				found = true
			}
		} // каждый holders - публикации от одного пользователя в данной ячейке
		if !found {
			holders = append(holders, holder)
		}
	}
	result, wasUnion := uniteHolders(holders, anomalyMode)
	for wasUnion {
		result, wasUnion = uniteHolders(result, anomalyMode)
	}
	return result, posts
}

func uniteHolders(holders []eventHolder, anomalyMode bool) ([]eventHolder, bool) {
	united := false
	var res []eventHolder
	for i := 0; i < len(holders); i++ {
		holder1 := holders[i]
		for j := i + 1; j < len(holders)-1; j++ {
			/*if j == (len(holders) - 1) { in for was j < len(holders) - What the fuck was that?
				continue
			}*/
			holder2 := holders[j]
			found := false
			for tag1 := range holder1.tags {
				for tag2 := range holder2.tags {
					if (tag1 == tag2) || anomalyMode {
						holder1 = combineHolders(holder1, holder2)
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
				holders = append(holders[:j], holders[j+1:]...)
				j--
			}
		}
		res = append(res, holder1)
	}
	return res, united
}

func combineHolders(holder1, holder2 eventHolder) eventHolder {
	res := eventHolder{
		users: map[string]bool{},
		posts: map[string]bool{},
		tags:  map[string]int{},
	}
	for userId := range holder1.users {
		res.users[userId] = true
	}
	for userId := range holder2.users {
		res.users[userId] = true
	}
	for postCode := range holder1.posts {
		res.posts[postCode] = true
	}
	for postCode := range holder2.posts {
		res.posts[postCode] = true
	}
	for tag, count := range holder1.tags {
		if _, ok := res.tags[tag]; ok {
			res.tags[tag] += count
		} else {
			res.tags[tag] = count
		}
	}
	for tag, count := range holder2.tags {
		if _, ok := res.tags[tag]; ok {
			res.tags[tag] += count
		} else {
			res.tags[tag] = count
		}
	}
	return res
}

func checkEvent(e eventHolder, minUser int, posts []data.Post, start, finish int64) (data.Event, bool) {
	if len(e.users) < 2 {
		return data.Event{}, false
	}
	if len(e.users) < minUser {
		return data.Event{}, false
	}
	tags := sortTags(e.tags)
	var postCodes []string
	for k := range e.posts {
		postCodes = append(postCodes, k)
	}
	return data.Event{
		Center:    eventCenter(e.posts, posts),
		PostCodes: postCodes,
		Tags:      tags,
		Title:     tags[0],
		Start:     start,
		Finish:    finish,
	}, true
}

func sortTags(tags map[string]int) []string {
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
	var res []string
	for _, k := range keys {
		res = append(res, rev[k]...)
	}
	return res
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

func sumWeightPoints(points []convtree.Point) float64 {
	sumWeight := float64(0)
	for _, point := range points {
		sumWeight += point.Weight
	}
	return sumWeight
}
