// Package constraints defines a set of useful constraints to be used with type parameters.
// This package is based on these proposals.
// - constraints package: https://github.com/golang/go/issues/45458
// - type parameters proposal: https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md
package constraints

// Signed is a constraint that permits any signed integer type.
type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// Unsigned is a constraint that permits any unsigned integer type.
type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// Integer is a constraint that permits any integer type.
type Integer interface {
	Signed | Unsigned
}

// Float is a constraint that permits any floating-point type.
type Float interface {
	~float32 | ~float64
}

// Complex is a constraint that permits any complex numeric type.
type Complex interface {
	~complex64 | ~complex128
}

// Ordered is a constraint that permits any ordered type: any type that supports the operators < <= >= >.
type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 |
		~string
}

// Slice is a constraint that matches slices of any element type.
type Slice[Elem any] interface {
	~[]Elem
}

// Map is a constraint that matches maps of any element and value type.
type Map[Key comparable, Val any] interface {
	~map[Key]Val
}

// Chan is a constraint that matches channels of any element type.
type Chan[Elem any] interface {
	~chan Elem
}
