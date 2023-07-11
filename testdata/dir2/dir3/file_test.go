package dir3

import (
	"testing"
)

func TestFunc3(t *testing.T) {
	s := Struct2{}
	s.Func1()
}

func TestFunc4(t *testing.T) {
	s := Struct2{}
	s.Func4("c")
}
