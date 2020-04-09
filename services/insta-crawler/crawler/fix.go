package crawler

import (
	"github.com/angrymuskrat/event-monitoring-system/services/insta-crawler/crawler/data"
)

type Fixer struct {
	Init bool
	loc  map[string]data.Location
}

func NewFixer(ls []data.Location) (Fixer, error) {
	loc := map[string]data.Location{}
	for _, l := range ls {
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
