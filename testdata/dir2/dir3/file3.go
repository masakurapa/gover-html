package dir3

import (
	"fmt"
)

type Struct2 struct{}

func (s *Struct2) Func1() {
	for i := 1; i <= 10; i++ {
		if i%3 == 0 && i%5 == 0 {
			fmt.Println("FizzBuzz")
			continue
		}
		if i%3 == 0 {
			fmt.Println("Fizz")
			continue
		}
		if i%5 == 0 {
			fmt.Println("Buzz")
			continue
		}
		fmt.Println(i)
	}
}

func (s *Struct2) Func4(ss string) bool {
	switch ss {
	case "a", "b":
		return true
	}
	return false
}
