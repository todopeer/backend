package highorder

// would expect `github.com/Shopify/hoff` to give most of the HighOrder functions
// this package is to provde the not provided ones

func Uniq[T comparable](arr []T) []T {
	if len(arr) < 2 {
		return arr
	}
	res := []T{arr[0]}

	for i := 1; i < len(arr); i++ {
		if arr[i] == res[len(res)-1] {
			continue
		}

		res = append(res, arr[i])
	}
	return res
}
