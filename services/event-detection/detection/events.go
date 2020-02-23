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
			//tagPosts := map[string]map[string]bool{}
			posts := make([]data.Post, len(tree.Points))
			evHolders := []eventHolder{}
			for i, item := range tree.Points {
				tags := []string{}
				post := item.Content.(data.Post)
				posts[i] = post
				//tags := []string{}
				r, _ := regexp.Compile("#[^\\s|\\#|\\n|!|\\.|\\?]+")
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
				r, _ = regexp.Compile("@[^\\s|\\#|\\n|!|\\.|\\?]+")
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
				if len(tags) == 0 {
					continue
				}
				exists := false
				for _, tag := range tags {
					for _, h := range evHolders {
						if _, ok := h.tags[tag]; ok {
							h.posts[post.Shortcode] = true
							for _, t := range tags {
								if _, ok := h.tags[t]; ok {
									h.tags[t]++
								} else {
									h.tags[t] = 1
								}
							}
							exists = true
							continue
						}
					}
					if exists {
						break
					}
				}
				if !exists {
					h := eventHolder{
						posts: map[string]bool{},
						tags:  map[string]int{},
					}
					h.posts[post.Shortcode] = true
					for _, tag := range tags {
						h.tags[tag] = 1
					}
					evHolders = append(evHolders, h)
				}
				//if len(tags) > 0 {
				//for _, tag := range tags {
				//	if _, ok := tagPosts[tag]; !ok {
				//		tagPosts[tag] = map[string]bool{}
				//	}
				//	tagPosts[tag][post.Shortcode] = true
				//}
				//}
			}
			//for tag, tPosts := range tagPosts {
			//	added := false
			//	for tCode := range tPosts {
			//		for idx, holder := range evHolders {
			//			if _, ok := holder.posts[tCode]; ok {
			//				added = true
			//				for eCode := range tPosts {
			//					holder.posts[eCode] = true
			//				}
			//				holder.tags[tag]++
			//				evHolders[idx] = holder
			//			}
			//		}
			//	}
			//	if !added {
			//		holder := eventHolder{
			//			posts: map[string]bool{},
			//			tags:  map[string]int{},
			//		}
			//		for code := range tPosts {
			//			holder.posts[code] = true
			//		}
			//		holder.tags[tag] = 1
			//		evHolders = append(evHolders, holder)
			//	}
			//}
			for _, e := range evHolders {
				if len(e.posts) >= maxPoints/2 {
					for k, v := range e.tags {
						if v < len(e.posts)/4 {
							delete(e.tags, k)
						}
					}
					//eventPosts := map[string]bool{}
					//for tag := range e.tags {
					//	for p := range tagPosts[tag] {
					//		eventPosts[p] = true
					//	}
					//}
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
					event := data.Event{
						Center:    eventCenter(e.posts, posts),
						PostCodes: postCodes,
						Tags:      tags,
						Title:     maxTag,
						Start:     start,
						Finish:    finish,
					}
					if len(postCodes) > 0 {
						result = append(result, event)
					}

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
