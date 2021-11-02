package rand

import (
	"math"
	"math/rand"
	"time"

	data "github.com/angrymuskrat/event-monitoring-system/services/proto"
)
const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type GenConfig struct {
	Center     data.Point
	DeltaPoint data.Point
	StartTime  int64
	FinishTime int64
}

type Rand struct {
	seeded *rand.Rand
}

func New() *Rand {
	r := Rand{}
	r.seeded = rand.New(rand.NewSource(time.Now().UnixNano()))
	return &r
}

func (r *Rand) FixString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[r.seeded.Intn(len(charset))]
	}
	return string(b)
}

func (r *Rand) Uint64(min, max int64) int64 {
	return int64(r.seeded.Uint64()>>1)%(max-min) + min
}

func (r *Rand) String(min, max int) string {
	length := int(r.Uint64(int64(min), int64(max)))
	return r.FixString(length)
}

func (r *Rand) Bool() bool {
	return r.seeded.Int()%2 == 1
}

func (r *Rand) Sign() float64 {
	return float64(1 + -2*r.seeded.Int()%2)
}

func (r *Rand) Double() float64 {
	return r.seeded.Float64()
}

func (r *Rand) DeltaDouble(delta float64) float64 {
	return r.Sign() * r.Double() * delta
}

func (r *Rand) Point(point data.Point, delta data.Point) data.Point {
	return data.Point{Lat: point.Lat + r.DeltaDouble(delta.Lat), Lon: point.Lon + r.DeltaDouble(delta.Lon)}
}

func (r *Rand) Post(conf GenConfig) *data.Post {
	p := r.Point(conf.Center, conf.DeltaPoint)
	return &data.Post{
		ID:            r.FixString(20),
		Shortcode:     r.FixString(10),
		ImageURL:      r.String(50, 300),
		IsVideo:       r.Bool(),
		Caption:       r.String(0, 500),
		CommentsCount: r.Uint64(0, 1000),
		Timestamp:     r.Uint64(conf.StartTime, conf.FinishTime),
		LikesCount:    r.Uint64(0, 10000),
		IsAd:          r.Bool(),
		AuthorID:      r.FixString(15),
		LocationID:    r.FixString(15),
		Lat:           p.Lat,
		Lon:           p.Lon,
	}
}

func (r *Rand) Event(conf GenConfig) *data.Event {
	lenTags := int(r.Uint64(1, 10))
	tags := make([]string, lenTags)
	for i := 0; i < lenTags; i++ {
		tags[i] = r.String(1, 6)
	}
	lenCodes := int(r.Uint64(1, 10))
	codes := make([]string, lenCodes)
	for i := 0; i < lenCodes; i++ {
		codes[i] = r.FixString(10)
	}

	startTime := r.Uint64(conf.StartTime, conf.FinishTime)
	finishTime := startTime + r.Uint64(0, 3600-startTime%3600)
	return &data.Event{
		Center:    data.Point{},
		PostCodes: codes,
		Tags:      tags,
		Title:     r.String(2, 20),
		Start:     startTime,
		Finish:    finishTime,
	}
}

func (r *Rand) Posts(length int, conf GenConfig) []data.Post {
	posts := make([]data.Post, length)
	for i := 0; i < length; i++ {
		posts[i] = *r.Post(conf)
	}
	return posts
}

func (r *Rand) Events(length int, conf GenConfig) []data.Event {
	events := make([]data.Event, length)
	for i := 0; i < length; i++ {
		events[i] = *r.Event(conf)
	}
	return events
}

func (r *Rand) ShortPost(conf GenConfig) *data.ShortPost {
	p := r.Point(conf.Center, conf.DeltaPoint)
	return &data.ShortPost{
		Shortcode:     r.FixString(10),
		Caption:       r.String(0, 500),
		CommentsCount: r.Uint64(0, 1000),
		Timestamp:     r.Uint64(conf.StartTime, conf.FinishTime),
		LikesCount:    r.Uint64(0, 10000),
		Lat:           p.Lat,
		Lon:           p.Lon,
	}
}

func MakeGenConfigByCorner(topLeft data.Point, botRight data.Point, start, finish int64) GenConfig {
	center := data.Point{ Lat: (topLeft.Lat + botRight.Lat) / 2, Lon: (topLeft.Lon + botRight.Lon) / 2}
	delta := data.Point{ Lat: math.Abs(center.Lat - botRight.Lat), Lon: math.Abs(center.Lon - botRight.Lon) }
	genConfig := GenConfig{Center: center, DeltaPoint: delta, StartTime: start, FinishTime: finish}
	return genConfig
}