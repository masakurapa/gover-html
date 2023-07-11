package dir1

func Func1() string {
	return "a"
}

func Func2(s string) int {
	switch s {
	case "a":
		return 1
	case "b":
		return 1
	case "c":
		return 3
	}
	return 0
}
