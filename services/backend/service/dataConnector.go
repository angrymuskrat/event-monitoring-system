package service

import (
	"context"
	service "github.com/angrymuskrat/event-monitoring-system/services/data-storage"
	data "github.com/angrymuskrat/event-monitoring-system/services/proto"
	"github.com/visheratin/unilog"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type DataConnector struct {
	dsClient service.GrpcService
}

func NewDataConnector(storageAddress string) (DataConnector, error) {
	conn, err := grpc.Dial(storageAddress, grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(service.MaxMsgSize)))
	if err != nil {
		unilog.Logger().Error("unable to connect to data strorage", zap.Error(err))
		return DataConnector{}, err
	}
	svc := service.NewGRPCClient(conn)
	return DataConnector{dsClient: svc}, nil
}

func (c DataConnector) HeatmapPosts(city string, topLeft, botRight data.Point, hour int64) ([]data.AggregatedPost, error) {
	posts, err := c.dsClient.SelectAggrPosts(context.Background(), city,
		data.SpatioHourInterval{
			Hour: hour,
			Area: data.Area{
				TopLeft:  &topLeft,
				BotRight: &botRight,
			}})
	if err != nil {
		unilog.Logger().Error("unable to get aggregated posts", zap.Error(err))
		return nil, err
	}
	return posts, nil
}

func (c DataConnector) Timeline(city string, start, finish int64) (Timeline, error) {
	tl, err := c.dsClient.PullTimeline(context.Background(), city, start, finish)
	if err != nil {
		unilog.Logger().Error("unable to get timeline", zap.Error(err))
		return nil, err
	}
	return tl, nil
}

func (c DataConnector) Events(city string, topLeft, botRight data.Point, hour int64) ([]data.Event, error) {
	evs, err := c.dsClient.PullEvents(context.Background(), city,
		data.SpatioHourInterval{
			Hour: hour,
			Area: data.Area{
				TopLeft:  &topLeft,
				BotRight: &botRight,
			}})
	if err != nil {
		unilog.Logger().Error("unable to get events", zap.Error(err))
		return nil, err
	}
	for i := range evs {
		evs[i].Tags = filterTags(evs[i].Tags, 5)
	}
	return evs, nil
}

func (c DataConnector) EventsByTags(city string, keytags []string, start, finish int64) ([]data.Event, error) {
	evs, err := c.dsClient.PullEventsTags(context.Background(), city, keytags, start, finish)
	if err != nil {
		unilog.Logger().Error("unable to search events by tags", zap.Error(err))
		return nil, err
	}
	for i := range evs {
		evs[i].Tags = filterTags(evs[i].Tags, 5)
	}
	return evs, nil
}

func (c DataConnector) ShortPostsInInterval(city string, shortcodes []string, start, finish int64) ([]data.ShortPost, error) {
	evs, err := c.dsClient.PullShortPostInInterval(context.Background(), city, shortcodes, start, finish)
	if err != nil {
		unilog.Logger().Error("unable to extract short posts in interval", zap.Error(err))
		return nil, err
	}
	return evs, nil
}

func (c DataConnector) SingleShortPost(city, shortcode string) (*data.ShortPost, error) {
	evs, err := c.dsClient.PullSingleShortPost(context.Background(), city, shortcode)
	if err != nil {
		unilog.Logger().Error("unable to extract single post", zap.Error(err))
		return nil, err
	}
	return evs, nil
}

func filterTags(tags []string, max int) []string {
	l := max
	if l > (len(tags) - 1) {
		l = len(tags)
	}
	tags = tags[0:l]
	return tags
}
