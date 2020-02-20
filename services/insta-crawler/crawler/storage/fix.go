package storage

import (
	"encoding/json"
	"github.com/angrymuskrat/event-monitoring-system/services/insta-crawler/crawler/data"
	"os"
)

type Location struct {
	ID  string
	Lat float64
	Lon float64
}

type Fixer struct {
	loc map[string]Location
}

func NewFixer(path string) (Fixer, error) {
	f, err := os.Open(path)
	if err != nil {
		return Fixer{}, err
	}
	defer f.Close()
	var d []Location
	err = json.NewDecoder(f).Decode(&d)
	if err != nil {
		return Fixer{}, err
	}
	loc := map[string]Location{}
	for _, l := range d {
		loc[l.ID] = l
	}
	return Fixer{loc: loc}, nil
}

func (f Fixer) Fix(d []data.Post) []data.Post {
	res := make([]data.Post, len(d))
	for i, p := range d {
		res[i] = p
		l, ok := f.loc[p.LocationID]
		if !ok {
			continue
		}
		res[i].Lat = l.Lat
		res[i].Lon = l.Lon
	}
	return res
}
