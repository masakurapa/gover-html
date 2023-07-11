package dir1

import (
	"testing"
)

func TestFunc1(t *testing.T) {
	Func1()
}

func TestFunc2(t *testing.T) {
	Func2("a")
	Func2("b")
	Func2("d")
}

func TestFunc3(t *testing.T) {
	Func3()
}

func TestFunc4(t *testing.T) {
	Func4("c")
}
