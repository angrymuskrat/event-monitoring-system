package data

import "testing"

func TestCrawlingType_UnmarshalJSON(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		s       CrawlingType
		args    args
		wantErr bool
		want    CrawlingType
	}{
		{
			name:    "locations test",
			args:    args{[]byte(`"locations"`)},
			wantErr: false,
			want:    LocationsType,
		},
		{
			name:    "profiles test",
			args:    args{[]byte(`"profiles"`)},
			wantErr: false,
			want:    ProfilesType,
		},
		{
			name:    "internal profiles test",
			args:    args{[]byte(`"profiles-internal"`)},
			wantErr: false,
			want:    InternalProfilesType,
		},
		{
			name:    "stories test",
			args:    args{[]byte(`"stories"`)},
			wantErr: false,
			want:    StoriesType,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.UnmarshalJSON(tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.s != tt.want {
				t.Errorf("UnmarshalJSON() got = %v, want %v", tt.s, tt.want)
			}
		})
	}
}
