package detection

import (
	data "github.com/angrymuskrat/event-monitoring-system/services/proto"
	convtree "github.com/visheratin/conv-tree"
	"regexp"
	"strings"
)

func FindEvents(histGrid convtree.ConvTree, posts []data.Post, maxPoints int, filterTags map[string]bool, start, finish int64) ([]data.Event, error) {
	candGrid, err := findCandidates(histGrid, posts, maxPoints)
	if err != nil {
		return nil, err
	}
	splitGrid(&candGrid, maxPoints)
	events := iterateGrid(&candGrid, maxPoints, filterTags, start, finish)
	return events, nil
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

func iterateGrid(tree *convtree.ConvTree, maxPoints int, filterTags map[string]bool, start, finish int64) []data.Event {
	if tree.IsLeaf {
		result := []data.Event{}
		if len(tree.Points) >= maxPoints {
			tagPosts := map[string]map[string]bool{}
			posts := make([]data.Post, len(tree.Points))
			for i, item := range tree.Points {
				post := item.Content.(data.Post)
				posts[i] = post
				tags := []string{}
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
					exists := false
					for _, t := range tags {
						if t == tag {
							exists = true
							break
						}
					}
					if !exists {
						tags = append(tags, tag)
					}
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
					exists := false
					for _, t := range tags {
						if t == tag {
							exists = true
							break
						}
					}
					if !exists {
						tags = append(tags, tag)
					}
				}
				if len(tags) > 0 {
					for _, tag := range tags {
						if _, ok := tagPosts[tag]; !ok {
							tagPosts[tag] = map[string]bool{}
						}
						tagPosts[tag][post.Shortcode] = true
					}
				}
			}
			evHolders := []eventHolder{}
			for tag, tPosts := range tagPosts {
				added := false
				for tCode := range tPosts {
					for idx, holder := range evHolders {
						if _, ok := holder.posts[tCode]; ok {
							added = true
							for eCode := range tPosts {
								holder.posts[eCode] = true
							}
							holder.tags[tag]++
							evHolders[idx] = holder
						}
					}
				}
				if !added {
					holder := eventHolder{
						posts: map[string]bool{},
						tags:  map[string]int{},
					}
					for code := range tPosts {
						holder.posts[code] = true
					}
					holder.tags[tag] = 1
					evHolders = append(evHolders, holder)
				}
			}
			for _, e := range evHolders {
				if len(e.posts) >= maxPoints/2 {
					for k, v := range e.tags {
						if v < len(e.posts)/2 {
							delete(e.tags, k)
						}
					}
					eventPosts := map[string]bool{}
					for tag := range e.tags {
						for p := range tagPosts[tag] {
							eventPosts[p] = true
						}
					}
					postCodes := []string{}
					for k := range eventPosts {
						postCodes = append(postCodes, k)
					}
					tags := []string{}
					for k := range e.tags {
						tags = append(tags, k)
					}
					event := data.Event{
						Center:    eventCenter(eventPosts, posts),
						PostCodes: postCodes,
						Tags:      tags,
						Title:     "",
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
		result = append(result, iterateGrid(tree.ChildBottomLeft, maxPoints, filterTags, start, finish)...)
		result = append(result, iterateGrid(tree.ChildBottomRight, maxPoints, filterTags, start, finish)...)
		result = append(result, iterateGrid(tree.ChildTopLeft, maxPoints, filterTags, start, finish)...)
		result = append(result, iterateGrid(tree.ChildTopRight, maxPoints, filterTags, start, finish)...)
		return result
	}
}

func eventCenter(codes map[string]bool, posts []data.Post) data.Point {
	points := map[convtree.Point]int{}
	for _, post := range posts {
		if _, ok := codes[post.Shortcode]; ok {
			p := convtree.Point{
				X: post.Lat,
				Y: post.Lon,
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
		lat += p.X * float64(w) / float64(len(codes))
		lon += p.Y * float64(w) / float64(len(codes))
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
