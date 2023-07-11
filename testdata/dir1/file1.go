package dir1

import (
	"fmt"
)

func Func3() {
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

func Func4(s string) bool {
	switch s {
	case "a", "b":
		return true
	}
	return false
}
