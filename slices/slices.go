// Package slices defines various functions useful with slices of any type.
// Unless otherwise specified, these functions all apply to the elements
// of a slice at index 0 <= i < len(s).
// This package is based on slices package proposal: https://github.com/golang/go/issues/45955
//
package slices

import "github.com/syumai/go-generics/constraints" // See https://github.com/golang/go/issues/45458

// Equal reports whether two slices are equal: the same length and all
// elements equal. If the lengths are different, Equal returns false.
// Otherwise, the elements are compared in index order, and the
// comparison stops at the first unequal pair.
// Floating point NaNs are not considered equal.
func Equal[T comparable](s1, s2 []T) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := 0; i < len(s1); i++ {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}

// EqualFunc reports whether two slices are equal using a comparison
// function on each pair of elements. If the lengths are different,
// EqualFunc returns false. Otherwise, the elements are compared in
// index order, and the comparison stops at the first index for which
// eq returns false.
func EqualFunc[T1, T2 any](s1 []T1, s2 []T2, eq func(T1, T2) bool) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := 0; i < len(s1); i++ {
		if !eq(s1[i], s2[i]) {
			return false
		}
	}
	return true
}

// Compare compares the elements of s1 and s2.
// The elements are compared sequentially starting at index 0,
// until one element is not equal to the other. The result of comparing
// the first non-matching elements is the result of the comparison.
// If both slices are equal until one of them ends, the shorter slice is
// considered less than the longer one
// The result will be 0 if s1==s2, -1 if s1 < s2, and +1 if s1 > s2.
func Compare[T constraints.Ordered](s1, s2 []T) int {
	maxLen := len(s1)
	if maxLen > len(s2) {
		maxLen = len(s2)
	}

	for i := 0; i < maxLen; i++ {
		if s1[i] == s2[i] {
			continue
		}
		if s1[i] < s2[i] {
			return -1
		}
		return 1
	}

	// all elements are equal until maxLen
	if len(s1) == len(s2) {
		return 0
	}
	if len(s1) < len(s2) {
		return -1
	}
	return 1
}

// CompareFunc is like Compare, but uses a comparison function
// on each pair of elements. The elements are compared in index order,
// and the comparisons stop after the first time cmp returns non-zero.
// The result will be the first non-zero result of cmp; if cmp always
// returns 0 the result is 0 if len(s1) == len(s2), -1 if len(s1) < len(s2),
// and +1 if len(s1) > len(s2).
func CompareFunc[T any](s1, s2 []T, cmp func(T, T) int) int {
	maxLen := len(s1)
	if maxLen > len(s2) {
		maxLen = len(s2)
	}

	for i := 0; i < maxLen; i++ {
		if v := cmp(s1[i], s2[i]); v != 0 {
			return v
		}
	}

	// all elements are equal until maxLen
	if len(s1) == len(s2) {
		return 0
	}
	if len(s1) < len(s2) {
		return -1
	}
	return 1
}

// Index returns the index of the first occurrence of v in s, or -1 if not present.
func Index[T comparable](s []T, v T) int {
	for i := 0; i < len(s); i++ {
		if s[i] == v {
			return i
		}
	}
	return -1
}

// IndexFunc returns the index into s of the first element
// satisfying f(c), or -1 if none do.
func IndexFunc[T any](s []T, f func(T) bool) int {
	for i := 0; i < len(s); i++ {
		if f(s[i]) {
			return i
		}
	}
	return -1
}

// Contains reports whether v is present in s.
func Contains[T comparable](s []T, v T) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == v {
			return true
		}
	}
	return false
}

// Insert inserts the values v... into s at index i, returning the modified slice.
// In the returned slice r, r[i] == the first v.  Insert panics if i is out of range.
//
// This implementation is copied from SliceTricks' InsertVector.
// - https://github.com/golang/go/wiki/SliceTricks#insertvector
//
// The following verbose way only copies elements
// in s[i:] once and allocates at most once.
// But, as of Go toolchain 1.16, due to lacking of
// optimizations to avoid elements clearing in the
// "make" call, the verbose way is not always faster.
//
// Future compiler optimizations might implement
// both in the most efficient ways.
func Insert[S constraints.Slice[T], T any](s S, i int, v ...T) S {
	if n := len(s) + len(vs); n <= cap(s) {
		s2 := s[:n]
		copy(s2[i+len(vs):], s[i:])
		copy(s2[i:], vs)
		return s2
	}
	s2 := make([]int, len(s) + len(vs))
	copy(s2, s[:i])
	copy(s2[i:], vs)
	copy(s2[i+len(vs):], s[i:])
	return s2
}

// Delete removes the elements s[i:j] from s, returning the modified slice.
// Delete panics if s[i:j] is not a valid slice of s.
// Delete modifies the contents of the slice s; it does not create a new slice.
// Delete is O(len(s)-(j-i)), so if many items must be deleted, it is better to
// make a single call deleting them all together than to delete one at a time.
func Delete[S constraints.Slice[T], T any](s S, i, j int) S {
	return s[:i+copy(s[i:], s[j:])]
}

// Clone returns a copy of the slice.
// The elements are copied using assignment, so this is a shallow clone.
func Clone[S constraints.Slice[T], T any](s S) S {
	s2 := make(S, len(s), cap(s))
	copy(s2, s)
	return s2
}

// Compact replaces consecutive runs of equal elements with a single copy.
// This is like the uniq command found on Unix.
// Compact modifies the contents of the slice s; it does not create a new slice.
func Compact[S constraints.Slice[T], T comparable](s S) S {
	if len(s) == 0 || len(s) == 1 {
		return s
	}
	j := 1
	for i := 1; i < len(s); i++ {
		if s[i] == s[j-1] {
			continue
		}
		s[j] = s[i]
		j++
	}
	return s[:j]
}

// CompactFunc is like Compact, but uses a comparison function.
func CompactFunc[S constraints.Slice[T], T any](s S, cmp func(T, T) bool) S {
	if len(s) == 0 || len(s) == 1 {
		return s
	}
	j := 1
	for i := 1; i < len(s); i++ {
		if cmp(s[i], s[j-1]) {
			continue
		}
		s[j] = s[i]
		j++
	}
	return s[:j]
}

// Grow grows the slice's capacity, if necessary, to guarantee space for
// another n elements. After Grow(n), at least n elements can be appended
// to the slice without another allocation. If n is negative or too large to
// allocate the memory, Grow will panic.
func Grow[S constraints.Slice[T], T any](s S, n int) S {
	return append(s[:cap(s)], make([]T, n, n)...)[:len(s)]
}

// Clip removes unused capacity from the slice, returning s[:len(s):len(s)].
func Clip[S constraints.Slice[T], T any](s S) S {
	return s[:len(s):len(s)]
}
