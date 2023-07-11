package dir2

import (
	"testing"
)

func TestFunc1(t *testing.T) {
	s := Struct1{}
	s.Func1()
}

func TestFunc2(t *testing.T) {
	s := Struct1{}
	s.Func2("a")
	s.Func2("b")
	s.Func2("d")
}
