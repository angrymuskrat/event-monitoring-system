package storage

import (
	"encoding/json"
	"os"

	"github.com/angrymuskrat/event-monitoring-system/services/insta-crawler/crawler/data"
)

type Location struct {
	ID  string
	Lat float64
	Lon float64
}

type Fixer struct {
	Init bool
	loc  map[string]Location
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
	return Fixer{Init: true, loc: loc}, nil
}

func (f Fixer) Fix(d []data.Post) []data.Post {
	res := make([]data.Post, 0, len(d))
	for _, p := range d {
		if p.Lat != 0 {
			res = append(res, p)
			continue
		}
		l, ok := f.loc[p.LocationID]
		if !ok {
			res = append(res, p)
			continue
		}
		p.Lat = l.Lat
		p.Lon = l.Lon
		res = append(res, p)
	}
	return res
}
