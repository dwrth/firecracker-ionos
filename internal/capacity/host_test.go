package capacity

import "testing"

func Test_parseMemTotalMiB(t *testing.T) {
	tests := []struct {
		name    string
		meminfo []byte
		want    int64
		wantErr bool
	}{
		{
			name:    "happy path",
			meminfo: []byte("MemTotal: 1024000 kB"),
			want:    1024000 / 1024,
			wantErr: false,
		},
		{
			name:    "invalid meminfo line",
			meminfo: []byte("MemTotal: not-a-number kB"),
			want:    0,
			wantErr: true,
		},
		{
			name:    "missing MemTotal",
			meminfo: []byte("MemAvailable: 1024000 kB"),
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := parseMemTotalMiB(tt.meminfo)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("parseMemTotalMiB() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Errorf("parseMemTotalMiB() succeeded unexpectedly")
				return
			}
			if got != tt.want {
				t.Errorf("parseMemTotalMiB() = %v, want %v", got, tt.want)
			}
		})
	}
}
