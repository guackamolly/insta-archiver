package model

func GroupBy[K comparable, T any](
	values []T,
	idf func(T) K,
) map[K][]T {
	m := map[K][]T{}

	for _, v := range values {
		k := idf(v)
		if s, ok := m[k]; ok {
			m[k] = append(s, v)
		} else {
			m[k] = []T{v}
		}
	}

	return m
}

func Filter[T any](
	values []T,
	pred func(T) bool,
) []T {
	var r []T

	for _, v := range values {
		if pred(v) {
			r = append(r, v)
		}
	}

	return r
}

func Map[T any, R any](
	values []T,
	mapf func(T) R,
) []R {
	var r []R

	for _, v := range values {
		m := mapf(v)
		r = append(r, m)
	}

	return r
}

func MapFilter[T any, M any](
	values []T,
	mapf func(T) (M, error),
	pred func(T) bool,
) []M {
	var r []M

	for _, v := range values {
		if pred(v) {
			m, err := mapf(v)

			if err == nil {
				r = append(r, m)
			}
		}
	}

	return r
}

func Find[T any](
	values []T,
	pred func(T) bool,
) *T {
	for _, v := range values {
		if pred(v) {
			return &v
		}
	}

	return nil
}
