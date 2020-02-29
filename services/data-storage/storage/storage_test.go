package storage

import (
	data "github.com/angrymuskrat/event-monitoring-system/services/proto"
	"reflect"
	"testing"
)

func Test_formSelectGrids(t *testing.T) {
	type args struct {
		ids []int64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test",
			args: args{ids: []int64{14, 88, 345}},
			want: "SELECT id, blob FROM grids WHERE id IN (14,88,345);",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := formSelectGrids(tt.args.ids); got != tt.want {
				t.Errorf("formSelectGrids() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_putEvent(t *testing.T) {
	type args struct {
		e   data.Event
		evs []data.Event
	}
	tests := []struct {
		name string
		args args
		want []data.Event
	}{
		{
			name: "nil slice",
			args: args{
				e:   data.Event{Start: 1},
				evs: nil,
			},
			want: []data.Event{{Start: 1}},
		},
		{
			name: "empty slice",
			args: args{
				e:   data.Event{Start: 1},
				evs: []data.Event{},
			},
			want: []data.Event{{Start: 1}},
		},
		{
			name: "insert to the beginning",
			args: args{
				e:   data.Event{Start: 2},
				evs: []data.Event{{Start: 3}, {Start: 5}},
			},
			want: []data.Event{{Start: 2}, {Start: 3}, {Start: 5}},
		},
		{
			name: "insert to the middle",
			args: args{
				e:   data.Event{Start: 4},
				evs: []data.Event{{Start: 3}, {Start: 5}},
			},
			want: []data.Event{{Start: 3}, {Start: 4}, {Start: 5}},
		},
		{
			name: "insert th the end",
			args: args{
				e:   data.Event{Start: 7},
				evs: []data.Event{{Start: 3}, {Start: 5}},
			},
			want: []data.Event{{Start: 3}, {Start: 5}, {Start: 7}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := putEvent(tt.args.e, tt.args.evs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("putEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}
