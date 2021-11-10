package positional

import (
	data "github.com/angrymuskrat/event-monitoring-system/services/proto"
	rand "github.com/angrymuskrat/event-monitoring-system/utils/rand"
	"math"
)

type Rand struct {
	rand   *rand.SimpleRand
	center data.Point
	delta  data.Point
}

func New(firstCorner, secondCorner data.Point) *Rand {
	simpleRand := rand.New()
	center := data.Point{Lat: (firstCorner.Lat + secondCorner.Lat) / 2, Lon: (firstCorner.Lon + secondCorner.Lon) / 2}
	delta := data.Point{Lat: math.Abs(center.Lat - firstCorner.Lat), Lon: math.Abs(center.Lon - firstCorner.Lat)}
	r := Rand{rand: simpleRand, center: center, delta: delta}
	return &r
}

func NewByDelta(center, delta data.Point) *Rand {
	simpleRand := rand.New()
	r := Rand{rand: simpleRand, center: center, delta: delta}
	return &r
}

func (r *Rand) Point() data.Point {
	return r.rand.Point(r.center, r.delta)
}

func (r *Rand) Post(minTimestamp, maxTimestamp int64) *data.Post {
	p := r.Point()
	return &data.Post{
		ID:            r.rand.FixString(20),
		Shortcode:     r.rand.FixString(10),
		ImageURL:      r.rand.String(50, 300),
		IsVideo:       r.rand.Bool(),
		Caption:       r.rand.String(0, 500),
		CommentsCount: r.rand.AbsInt64(0, 1000),
		Timestamp:     r.rand.AbsInt64(minTimestamp, maxTimestamp),
		LikesCount:    r.rand.AbsInt64(0, 10000),
		IsAd:          r.rand.Bool(),
		AuthorID:      r.rand.FixString(15),
		LocationID:    r.rand.FixString(15),
		Lat:           p.Lat,
		Lon:           p.Lon,
	}
}

func (r *Rand) Event(minTimestamp, maxTimestamp int64) *data.Event {
	lenTags := int(r.rand.AbsInt64(1, 10))
	tags := make([]string, lenTags)
	for i := 0; i < lenTags; i++ {
		tags[i] = r.rand.String(1, 6)
	}
	lenCodes := int(r.rand.AbsInt64(1, 10))
	codes := make([]string, lenCodes)
	for i := 0; i < lenCodes; i++ {
		codes[i] = r.rand.FixString(10)
	}

	startTime := r.rand.AbsInt64(minTimestamp, maxTimestamp)
	finishTime := startTime + r.rand.AbsInt64(0, 3600-startTime%3600)
	return &data.Event{
		Center:    data.Point{},
		PostCodes: codes,
		Tags:      tags,
		Title:     r.rand.String(2, 20),
		Start:     startTime,
		Finish:    finishTime,
	}
}

func (r *Rand) Posts(length int, minTimestamp, maxTimestamp int64) []data.Post {
	posts := make([]data.Post, length)
	for i := 0; i < length; i++ {
		posts[i] = *r.Post(minTimestamp, maxTimestamp)
	}
	return posts
}

func (r *Rand) Events(length int, minTimestamp, maxTimestamp int64) []data.Event {
	events := make([]data.Event, length)
	for i := 0; i < length; i++ {
		events[i] = *r.Event(minTimestamp, maxTimestamp)
	}
	return events
}

func (r *Rand) ShortPost(minTimestamp, maxTimestamp int64) *data.ShortPost {
	p := r.Point()
	return &data.ShortPost{
		Shortcode:     r.rand.FixString(10),
		Caption:       r.rand.String(0, 500),
		CommentsCount: r.rand.AbsInt64(0, 1000),
		Timestamp:     r.rand.AbsInt64(minTimestamp, maxTimestamp),
		LikesCount:    r.rand.AbsInt64(0, 10000),
		Lat:           p.Lat,
		Lon:           p.Lon,
	}
}
