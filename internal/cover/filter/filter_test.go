package filter_test

import (
	"testing"

	"github.com/masakurapa/gover-html/internal/cover/filter"
)

func TestFilter_IsOutputTarget(t *testing.T) {
	type args struct {
		include *string
	}
	tests := []struct {
		name string
		args args
		path string
		want bool
	}{
		{
			name: "includeがnilの場合はtrueが返る",
			args: args{
				include: nil,
			},
			path: "path/to/dir1",
			want: true,
		},
		{
			name: "includeが空文字の場合はtrueが返る",
			args: args{
				include: stringP(""),
			},
			path: "path/to/dir1",
			want: true,
		},
		{
			name: "pathが絶対パス指定の場合はfalseが返る",
			args: args{
				include: stringP(""),
			},
			path: "/path/to/dir1",
			want: false,
		},
		{
			name: "includeと等しいパスの場合はtrueが返る",
			args: args{
				include: stringP("path/to/dir1"),
			},
			path: "path/to/dir1",
			want: true,
		},
		{
			name: "includeと等しくないパスの場合はfalseが返る",
			args: args{
				include: stringP("path/to/dir1"),
			},
			path: "path/to/dir2",
			want: false,
		},
		{
			name: "includeが./で始まるパスの場合はtrueが返る",
			args: args{
				include: stringP("./path/to/dir1"),
			},
			path: "path/to/dir1",
			want: true,
		},
		{
			name: "includeが./で始まっていて、マッチしないパスの場合はfalseが返る",
			args: args{
				include: stringP("./path/to/dir1"),
			},
			path: "path/to/dir2",
			want: false,
		},
		{
			name: "includeが/で始まるパスの場合はfalseが返る",
			args: args{
				include: stringP("/path/to/dir1"),
			},
			path: "path/to/dir1",
			want: false,
		},
		{
			name: "includeが/で終わるパスの場合はtrueが返る",
			args: args{
				include: stringP("path/to/dir1/"),
			},
			path: "path/to/dir1",
			want: true,
		},
		{
			name: "includeが/で終わっていて、マッチしないパスの場合はfalseが返る",
			args: args{
				include: stringP("path/to/dir1/"),
			},
			path: "path/to/dir2",
			want: false,
		},
		{
			name: "includeで始まるパスの場合はtrueが返る",
			args: args{
				include: stringP("path/to"),
			},
			path: "path/to/dir1",
			want: true,
		},
		{
			name: "includeで始まっていないパスの場合はfalseが返る",
			args: args{
				include: stringP("path/to"),
			},
			path: "path/tooo/dir1",
			want: false,
		},

		{
			name: "includeがカンマ区切りで複数あり、いずれかと同じパスの場合はtrueが返る",
			args: args{
				include: stringP("path/to/dir2, path/to/dir3,path/to/dir1"),
			},
			path: "path/to/dir1",
			want: true,
		},
		{
			name: "includeがカンマ区切りで複数あり、いずれかと同じパスにもマッチしない場合はtrueが返る",
			args: args{
				include: stringP("path/to/dir2, path/to/dir3,path/to/dir1"),
			},
			path: "path/to/dir4",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := filter.New(tt.args.include)
			if got := f.IsOutputTarget(tt.path); got != tt.want {
				t.Errorf("filter.IsOutputTarget() = %v, want %v", got, tt.want)
			}
		})
	}
}

func stringP(s string) *string {
	return &s
}
