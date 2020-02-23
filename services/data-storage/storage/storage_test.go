package storage

import "testing"

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
			want: "SELECT id, blob FROM gridsWHERE id IN (14,88,345);",
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
