package dbconnector

import "context"

type Post struct {
	ID            string
	Shortcode     string
	ImageURL      string
	IsVideo       bool
	Caption       string
	CommentsCount int
	Timestamp     int64
	LikesCount    int
	IsAd          bool
	AuthorID      string
	LocationID    string
	Lat           float64
	Lon           float64
}

type SpatialTemporalInterval struct {
	MinTime int64 `msg`
	MaxTime int64
	MinLat  float64
	MinLon  float64
	MaxLat  float64
	MaxLon  float64
}

type Repository interface {
	PushPosts (ctx context.Context, posts []Post) error
	SelectPosts (ctx context.Context, interval SpatialTemporalInterval) ([]Post, error)
}
