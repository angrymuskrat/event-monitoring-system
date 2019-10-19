package dbsvc

import (
	"context"
	"fmt"

	"github.com/angrymuskrat/event-monitoring-system/services/dbsvc/pb"
)


type Repository interface {
	PushPosts(ctx context.Context, posts []pd.Post) error
	SelectPosts(ctx context.Context, interval SpatialTemporalInterval) ([]Post, error)
}

func (i SpatialTemporalInterval) String() string {
	return fmt.Sprintf("Time: %v-%v, lat: %v-%v, lon: %v-%v", i.MinTime, i.MaxTime, i.MinLat, i.MaxLat, i.MinLon, i.MaxLon)
}
