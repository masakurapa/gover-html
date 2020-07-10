package ex

import "strconv"

func Example(s string) string {
	n, err := strconv.Atoi(s)
	if err != nil {
		return "error!!"
	}

	if n <= 10 {
		return "hello"
	} else if n <= 20 {
		return "world"
	}
	return "ninja!"
}
