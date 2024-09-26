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
