package util

func Map[T any, R any](input []T, f func(T) R) []R {
	output := make([]R, len(input))
	for i, e := range input {
		output[i] = f(e)
	}
	return output
}
