package blitra

// NOTE: The names of these utilities are as such. The first letter includes
// a transform to a pointer or a value. If no transform is needed then skip
// it. The bit is an action. Or indicates a fallback to a default. Map indicates
// a transformation via a function. The last letter indicates the type of the
// default value which must be provided only if the first letter is omitted.

// A helper function which will take a value and return a pointer to it.
// P stands for Pointer.
func P[T any](v T) *T {
	return &v
}

// A helper function which will take a pointer and returns it if not nil.
// If it is then the default pointer is returned instead.
// OrP stands for Or Pointer.
func OrP[T any](p *T, defaultP *T) *T {
	if p == nil {
		return defaultP
	}
	return p
}

// VOr takes a pointer and dereferences it, returning the value.
// If the pointer is nil, then the default value is returned instead.
// VOr stands for Value Or.
func VOr[T any](p *T, defaultV T) T {
	if p == nil {
		return defaultV
	}
	return *p
}

// V takes a pointer and dereferences it, returning the value.
// If the pointer is nil, then the zero value of the type is returned instead.
// V stands for Value.
func V[T any](p *T) T {
	var v T
	return VOr(p, v)
}

// VMap takes a pointer and a map function. The return value of the map function
// is returned. If the pointer is nil, then the zero value of the return type
// is returned instead.
// VMap stands for Value Map.
func VMap[T any, U any](p *T, fn func(T) U) U {
	if p == nil {
		var v U
		return v
	}
	return fn(*p)
}
