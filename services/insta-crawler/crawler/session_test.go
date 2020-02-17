package crawler

import "testing"

func Test_convertSession(t *testing.T) {
	type args struct {
		sess    Session
		rootDir string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "simple test",
			args: args{
				sess: Session{
					ID: "test session",
					Params: Parameters{
						CityID:   "spb",
						Entities: []string{"123", "456"},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := convertSession(tt.args.sess, tt.args.rootDir); (err != nil) != tt.wantErr {
				t.Errorf("convertSession() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
