package writer

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	service "github.com/angrymuskrat/event-monitoring-system/services/data-storage"
	convtree "github.com/visheratin/conv-tree"
	"google.golang.org/grpc"
	"io/ioutil"
	"time"
)

type Writer struct {
	storage service.Service
	conn    grpc.ClientConn
}

func New(url string) (*Writer, error) {
	var (
		svc service.Service
		err error
	)
	conn, err := grpc.Dial(url, grpc.WithInsecure(), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(service.MaxMsgSize)))
	if err != nil {
		return nil, err
	}
	svc = service.NewGRPCClient(conn)
	return &Writer{storage: svc, conn: *conn}, nil
}

func (w *Writer) WritePostsJson(ctx context.Context, cityId string, start, end int64, path string) error {
	res, _, err := w.storage.SelectPosts(ctx, cityId, start, end)
	if err != nil {
		return err
	}
	f, err := json.MarshalIndent(res, "", " ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(fmt.Sprintf("%vposts_%v_%v.json", path, start, end), f, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (w *Writer) WriteGridJson(ctx context.Context, cityId string, gridTime int64, path string) error {
	date := time.Unix(gridTime, 0).UTC()
	id := int64(date.Hour()) + int64(date.Month())*1000
	if date.Weekday().String() == "Sunday" || date.Weekday().String() == "Saturday" {
		id += 2 * 100
	} else {
		id += 1 * 100
	}
	res, err := w.storage.PullGrid(ctx, cityId, []int64{id})
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(res[id])
	dec := gob.NewDecoder(buf)

	var grid convtree.ConvTree

	if err := dec.Decode(&grid); err != nil {
		return err
	}
	f, err := json.MarshalIndent(grid, "", " ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(fmt.Sprintf("%vgrid%v.json", path, gridTime), f, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (w *Writer) Close() {
	w.conn.Close()
}
