package dbsvc

import (
	"context"
	"github.com/angrymuskrat/event-monitoring-system/services/dbsvc/pb"
)


type Repository interface {
	PushPosts(ctx context.Context, posts []pb.Post) error
	SelectPosts(ctx context.Context, interval pb.SpatialTemporalInterval) ([]pb.Post, error)
}
