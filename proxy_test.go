package twitch

import (
	"testing"
)

func TestListFromAPI(t *testing.T) {
	tests := []struct {
		name    string
		want    int
		wantErr bool
	}{
		{
			name:    "get list from api",
			want:    100,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ListFromAPI()

			if (err != nil) != tt.wantErr {
				t.Errorf("ListFromAPI() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got.Proxies) != tt.want {
				t.Errorf("ListFromAPI() got = %v, want %v", got, tt.want)
			}

			t.Log(got.Info.GetNextRefresh())
		})
	}
}
