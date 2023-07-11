package dir2

type Struct1 struct{}

func (*Struct1) Func1() string {
	return ""
}

func (*Struct1) Func2(s string) int {
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
