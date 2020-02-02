package data

import "encoding/json"

type Post struct {
	ID            string
	Shortcode     string
	ImageURL      string
	IsVideo       bool
	TaggedUsers   []UserTag
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

type UserTag struct {
	Username string
	X        float64
	Y        float64
}

func (p Post) Marshal() ([]byte, error) {
	b, err := json.Marshal(p)
	return b, err
}

func (p *Post) Unmarshal(d []byte) error {
	return json.Unmarshal(d, p)
}
