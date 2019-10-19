package pb

import "fmt"

func (i SpatialTemporalInterval) ForLog() string {
	return fmt.Sprintf("Time: %v-%v, lat: %v-%v, lon: %v-%v", i.MinTime, i.MaxTime, i.MinLat, i.MaxLat, i.MinLon, i.MaxLon)
}

