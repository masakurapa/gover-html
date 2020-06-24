package example

import "testing"

func TestExample(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "10",
			args: args{s: "10"},
			want: "hello",
		},
		{
			name: "not integer",
			args: args{s: "str"},
			want: "error!!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Example(tt.args.s); got != tt.want {
				t.Errorf("Example() = %v, want %v", got, tt.want)
			}
		})
	}
}
