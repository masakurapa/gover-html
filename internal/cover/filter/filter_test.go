package filter_test

import (
	"testing"

	"github.com/masakurapa/gover-html/internal/cover/filter"
)

func TestFilter_IsOutputTarget(t *testing.T) {
	type args struct {
		include *string
		exclude *string
	}
	tests := []struct {
		name string
		args args
		path string
		want bool
	}{
		{
			name: "include, excludeがnilの場合はtrueが返る",
			args: args{
				include: nil,
				exclude: nil,
			},
			path: "path/to/dir1",
			want: true,
		},
		{
			name: "include, excludeが空文字の場合はtrueが返る",
			args: args{
				include: stringP(""),
				exclude: stringP(""),
			},
			path: "path/to/dir1",
			want: true,
		},
		{
			name: "pathが絶対パス指定の場合はfalseが返る",
			args: args{
				include: stringP(""),
				exclude: stringP(""),
			},
			path: "/path/to/dir1",
			want: false,
		},

		// validate "include"
		{
			name: "include = path の場合はtrueが返る",
			args: args{
				include: stringP("path/to/dir1"),
			},
			path: "path/to/dir1",
			want: true,
		},
		{
			name: "include != path の場合はfalseが返る",
			args: args{
				include: stringP("path/to/dir1"),
			},
			path: "path/to/dir2",
			want: false,
		},
		{
			name: "includeの接頭語(./)を除いた値 = path の場合はtrueが返る",
			args: args{
				include: stringP("./path/to/dir1"),
			},
			path: "path/to/dir1",
			want: true,
		},
		{
			name: "includeの接頭語(./)を除いた値 != path の場合はfalseが返る",
			args: args{
				include: stringP("./path/to/dir1"),
			},
			path: "path/to/dir2",
			want: false,
		},
		{
			name: "includeの接尾語(/)を除いた値 = path の場合はtrueが返る",
			args: args{
				include: stringP("path/to/dir1/"),
			},
			path: "path/to/dir1",
			want: true,
		},
		{
			name: "includeの接尾語(/)を除いた値 != path の場合はfalseが返る",
			args: args{
				include: stringP("path/to/dir1/"),
			},
			path: "path/to/dir2",
			want: false,
		},
		{
			name: "include = pathの接頭語(./)を除いた値 の場合はtrueが返る",
			args: args{
				include: stringP("path/to/dir1"),
			},
			path: "./path/to/dir1",
			want: true,
		},
		{
			name: "include != pathの接頭語(./)を除いた値 の場合はfalseが返る",
			args: args{
				include: stringP("path/to/dir1"),
			},
			path: "./path/to/dir2",
			want: false,
		},
		{
			name: "include = pathの接尾語(/)を除いた値 の場合はtrueが返る",
			args: args{
				include: stringP("path/to/dir1"),
			},
			path: "./path/to/dir1/",
			want: true,
		},
		{
			name: "include != pathの接尾語(/)を除いた値 の場合はfalseが返る",
			args: args{
				include: stringP("path/to/dir1"),
			},
			path: "./path/to/dir2/",
			want: false,
		},
		{
			name: "include = pathの接頭語 の場合はtrueが返る",
			args: args{
				include: stringP("path/to"),
			},
			path: "path/to/dir1",
			want: true,
		},
		{
			name: "include != pathの接頭語 の場合はfalseが返る",
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
			name: "includeがカンマ区切りで複数あり、いずれかと同じパスにもマッチしない場合はfalseが返る",
			args: args{
				include: stringP("path/to/dir2, path/to/dir3,path/to/dir1"),
			},
			path: "path/to/dir4",
			want: false,
		},

		// validate "exclude"
		{
			name: "exclude = path の場合はfalseが返る",
			args: args{
				exclude: stringP("path/to/dir1"),
			},
			path: "path/to/dir1",
			want: false,
		},
		{
			name: "exclude != path の場合はtrueが返る",
			args: args{
				exclude: stringP("path/to/dir1"),
			},
			path: "path/to/dir2",
			want: true,
		},
		{
			name: "excludeの接頭語(./)を除いた値 = path の場合はfalseが返る",
			args: args{
				exclude: stringP("./path/to/dir1"),
			},
			path: "path/to/dir1",
			want: false,
		},
		{
			name: "excludeの接頭語(./)を除いた値 != path の場合はtrueが返る",
			args: args{
				exclude: stringP("./path/to/dir1"),
			},
			path: "path/to/dir2",
			want: true,
		},
		{
			name: "excludeの接尾語(/)を除いた値 = path の場合はfalseが返る",
			args: args{
				exclude: stringP("path/to/dir1/"),
			},
			path: "path/to/dir1",
			want: false,
		},
		{
			name: "excludeの接尾語(/)を除いた値 != path の場合はtrueが返る",
			args: args{
				exclude: stringP("path/to/dir1/"),
			},
			path: "path/to/dir2",
			want: true,
		},
		{
			name: "exclude = pathの接頭語(./)を除いた値 の場合はfalseが返る",
			args: args{
				exclude: stringP("path/to/dir1"),
			},
			path: "./path/to/dir1",
			want: false,
		},
		{
			name: "exclude != pathの接頭語(./)を除いた値 の場合はtrueが返る",
			args: args{
				exclude: stringP("path/to/dir1"),
			},
			path: "./path/to/dir2",
			want: true,
		},
		{
			name: "exclude = pathの接尾語(/)を除いた値 の場合はfalseが返る",
			args: args{
				exclude: stringP("path/to/dir1"),
			},
			path: "./path/to/dir1/",
			want: false,
		},
		{
			name: "exclude != pathの接尾語(/)を除いた値 の場合はtrueが返る",
			args: args{
				exclude: stringP("path/to/dir1"),
			},
			path: "./path/to/dir2/",
			want: true,
		},
		{
			name: "exclude = pathの接頭語 の場合はfalseが返る",
			args: args{
				exclude: stringP("path/to"),
			},
			path: "path/to/dir1",
			want: false,
		},
		{
			name: "exclude != pathの接頭語 の場合はtrueが返る",
			args: args{
				exclude: stringP("path/to"),
			},
			path: "path/tooo/dir1",
			want: true,
		},
		{
			name: "excludeがカンマ区切りで複数あり、いずれかと同じパスの場合はfalseが返る",
			args: args{
				exclude: stringP("path/to/dir2, path/to/dir3,path/to/dir1"),
			},
			path: "path/to/dir1",
			want: false,
		},
		{
			name: "excludeがカンマ区切りで複数あり、いずれかと同じパスにもマッチしない場合はtrueが返る",
			args: args{
				exclude: stringP("path/to/dir2, path/to/dir3,path/to/dir1"),
			},
			path: "path/to/dir4",
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := filter.New(tt.args.include, tt.args.exclude)
			if got := f.IsOutputTarget(tt.path); got != tt.want {
				t.Errorf("filter.IsOutputTarget() = %v, want %v", got, tt.want)
			}
		})
	}
}

func stringP(s string) *string {
	return &s
}
