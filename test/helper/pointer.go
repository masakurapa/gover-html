package helper

// P returns pointer to value
func P[T any](v T) *T {
	return &v
}
