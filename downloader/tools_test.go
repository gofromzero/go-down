package downloader

import "testing"

func Test_calcLength(t *testing.T) {
	type args struct {
		L int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "calcLength KB",
			args: args{
				L: 1024,
			},
			want: "1 KB",
		},
		{
			name: "calcLength MB",
			args: args{
				L: 1024 * 1024,
			},
			want: "1 MB",
		},
		{
			name: "calcLength GB",
			args: args{
				L: 1024 * 1024 * 1024,
			},
			want: "1 GB",
		},
		{
			name: "calcLength PB",
			args: args{
				L: 1024 * 1024 * 1024 * 1024,
			},
			want: "1 PB",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calcLength(tt.args.L); got != tt.want {
				t.Errorf("calcLength() = %v, want %v", got, tt.want)
			}
		})
	}
}
