package service

import (
	data "github.com/angrymuskrat/event-monitoring-system/services/proto"
	postrand "github.com/angrymuskrat/event-monitoring-system/utils/rand/positional"
	"math/rand"
	"strconv"
)

type MockConnector struct{}

func (c MockConnector) HeatmapPosts(city string, topLeft, botRight data.Point, hour int64) ([]data.AggregatedPost, error) {
	res := []data.AggregatedPost{}
	for x := botRight.Lat; x <= topLeft.Lat; x += 0.001 {
		for y := topLeft.Lon; y <= botRight.Lon; y += 0.001 {
			c := int64(rand.Float64() * 1000)
			p := data.AggregatedPost{
				Center: data.Point{
					Lat: x + 0.0005,
					Lon: y + 0.0005,
				},
				Count: c,
			}
			res = append(res, p)
		}
	}
	return res, nil
}

func (c MockConnector) Timeline(city string, start, finish int64) (Timeline, error) {
	res := Timeline{}
	for t := start; t <= finish; t += 3600 {
		ts := data.Timestamp{
			Time:         t,
			PostsNumber:  int64(rand.Float64() * 1000),
			EventsNumber: int64(rand.Float64() * 50),
		}
		res = append(res, ts)
	}
	return res, nil
}

func (c MockConnector) Events(city string, topLeft, botRight data.Point, hour int64) ([]data.Event, error) {
	res := []data.Event{}
	for i := 0; i < 15; i++ {
		lat := botRight.Lat + rand.Float64()*(topLeft.Lat-botRight.Lat)
		lon := topLeft.Lon + rand.Float64()*(botRight.Lon-topLeft.Lon)
		e := data.Event{
			Center: data.Point{
				Lat: lat,
				Lon: lon,
			},
			PostCodes: []string{"B3Hdj8En6eR", "B3MGSDdnLtm", "B2ZNlFEnL0F"},
			Tags:      []string{"#testtag1", "#testtag2", "#quitelongtesttag", "#shorttag"},
			Title:     "Event " + strconv.Itoa(i),
			Start:     hour,
			Finish:    hour,
		}
		res = append(res, e)
	}
	return res, nil
}

var topLeft = data.Point{
	Lat: 60.12,
	Lon: 29.91,
}

var botRight = data.Point{
	Lat: 59.75,
	Lon: 30.69,
}

func (c MockConnector) EventsByTags(city string, keytags []string, start, finish int64) ([]data.Event, error) {
	res := []data.Event{}
	for i := 0; i < 15; i++ {
		lat := botRight.Lat + rand.Float64()*(topLeft.Lat-botRight.Lat)
		lon := topLeft.Lon + rand.Float64()*(botRight.Lon-topLeft.Lon)
		t := start + int64(rand.Float64()*float64(finish-start))
		e := data.Event{
			Center: data.Point{
				Lat: lat,
				Lon: lon,
			},
			PostCodes: []string{"B3Hdj8En6eR", "B3MGSDdnLtm", "B2ZNlFEnL0F"},
			Tags:      []string{"#testtag1", "#testtag2", "#quitelongtesttag", "#shorttag"},
			Title:     "Event " + strconv.Itoa(i),
			Start:     t,
			Finish:    t + 3600,
		}
		res = append(res, e)
	}
	return res, nil
}

func (c MockConnector) ShortPostsInInterval(city string, shortcodes []string, start, end int64) ([]data.ShortPost, error) {
	var res []data.ShortPost
	if start > end {
		return res, nil
	}
	generator := postrand.New(topLeft, botRight)
	for _, code := range shortcodes {
		post := generator.ShortPost(start, end)
		post.Shortcode = code
		res = append(res, *post)
	}
	return res, nil
}

const StartTimestamp, EndTimestamp = 1546300800, 1577836800

func (c MockConnector) SingleShortPost(city, shortcode string) (*data.ShortPost, error) {
	generator := postrand.New(topLeft, botRight)
	post := generator.ShortPost(StartTimestamp, EndTimestamp)
	post.Shortcode = shortcode
	return post, nil
}
