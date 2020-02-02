package data

import (
	"encoding/json"
)

type Location struct {
	ID    string
	Title string
	Lat   float64
	Lon   float64
	Slug  string
}

func (l *Location) Marshal() ([]byte, error) {
	b, err := json.Marshal(l)
	return b, err
}

func (l *Location) Unmarshal(d []byte) error {
	return json.Unmarshal(d, l)
}
