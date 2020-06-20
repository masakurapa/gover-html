package example_test

import (
	"testing"

	"github.com/masakurapa/go-cover/example2"
)

func TestHoge22(t *testing.T) {
	tc := []struct {
		i int
		e bool
	}{
		{9, false},
		// {10, true},
		// {11, true},
	}

	for _, tt := range tc {
		t.Run("test", func(t *testing.T) {
			if example2.Hoge(tt.i) != tt.e {
				t.Fatal("error!")
			}
			if example2.Hoge2(tt.i) != tt.e {
				t.Fatal("error!")
			}
		})
	}
}
