package util

func MapWithError[A, B any](arr []A, m func(A) (B, error)) ([]B, error) {
	res := make([]B, len(arr))

	var err error
	for i, a := range arr {
		res[i], err = m(a)
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func FindBy[A any](arr []A, f func(A) bool) int {
	for i, v := range arr {
		if f(v) {
			return i
		}
	}
	return -1
}