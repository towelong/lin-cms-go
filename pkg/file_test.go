package pkg

import (
	"testing"
)

func TestCreateDirAndFileForCurrentTime(t *testing.T) {
	type args struct {
		fileDir string
		format  string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "生成目录",
			args: args{
				fileDir: "tt1",
				format:  "2006-01-02",
			},
			want: "/Users/welong/Project/lin-cms-go/tt1/2021-11-18",
		},
		{
			name: "生成目录2",
			args: args{
				fileDir: "tt2",
				format:  "2006-01-02",
			},
			want: "/Users/welong/Project/lin-cms-go/tt2/2021-11-18",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := CreateDirAndFileForCurrentTime(tt.args.fileDir, tt.args.format)
			if got != tt.want {
				t.Errorf("CreateDirAndFileForCurrentTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
