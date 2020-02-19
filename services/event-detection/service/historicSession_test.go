package service

import (
	"reflect"
	"testing"
)

func Test_getIntervals(t *testing.T) {
	type args struct {
		start  int64
		finish int64
		tz     string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string][][2]int64
		wantErr bool
	}{
		{
			"test case",
			args{start: 1483218000, finish: 1546290000, tz: "Europe/Moscow"},
			map[string][][2]int64{},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getIntervals(tt.args.start, tt.args.finish, tt.args.tz)
			if (err != nil) != tt.wantErr {
				t.Errorf("getIntervals() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getIntervals() = %v, want %v", got, tt.want)
			}
		})
	}
}
