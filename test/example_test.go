package example_test

import (
	"testing"

	"github.com/masakurapa/gocover-html/example"
)

func TestHoge(t *testing.T) {
	tc := []struct {
		i int
		e bool
	}{
		// {9, false},
		{10, true},
		{11, true},
	}

	for _, tt := range tc {
		t.Run("test", func(t *testing.T) {
			if example.Hoge(tt.i) != tt.e {
				t.Fatal("error!")
			}
		})
	}
}
