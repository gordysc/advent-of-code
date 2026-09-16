// Package set provides a generic hash set built on a map.
//
// A map with an empty struct value uses no memory for the values, which is the
// idiomatic way to build a set in Go.
package set

import "iter"

// Set is an unordered collection of unique values.
type Set[T comparable] map[T]struct{}

// New returns an empty set.
func New[T comparable]() Set[T] {
	return make(Set[T])
}

// Of returns a set containing the given items.
func Of[T comparable](items ...T) Set[T] {
	s := make(Set[T], len(items))
	s.Add(items...)

	return s
}

// From returns a set with every element of a slice.
func From[T comparable](items []T) Set[T] {
	return Of(items...)
}

// Add inserts one or more items. Adding an existing item is a no-op.
func (s Set[T]) Add(items ...T) {
	for _, item := range items {
		s[item] = struct{}{}
	}
}

// Has reports whether the item is in the set.
func (s Set[T]) Has(item T) bool {
	_, ok := s[item]

	return ok
}

// Remove deletes the item. Removing an absent item is a no-op.
func (s Set[T]) Remove(item T) {
	delete(s, item)
}

// Len returns the number of items.
func (s Set[T]) Len() int {
	return len(s)
}

// Items returns the members as a slice in no particular order.
func (s Set[T]) Items() []T {
	out := make([]T, 0, len(s))

	for item := range s {
		out = append(out, item)
	}

	return out
}

// All yields every member. Use it as: for item := range s.All() { ... }.
// Plain `for item := range s` also works because Set is a map.
func (s Set[T]) All() iter.Seq[T] {
	return func(yield func(T) bool) {
		for item := range s {
			if !yield(item) {
				return
			}
		}
	}
}

// Union returns a new set with the members of both sets.
func (s Set[T]) Union(other Set[T]) Set[T] {
	out := make(Set[T], len(s)+len(other))

	for item := range s {
		out.Add(item)
	}
	for item := range other {
		out.Add(item)
	}

	return out
}

// Intersect returns a new set with the members present in both sets.
func (s Set[T]) Intersect(other Set[T]) Set[T] {
	out := New[T]()

	for item := range s {
		if other.Has(item) {
			out.Add(item)
		}
	}

	return out
}

// Difference returns a new set with the members of s that are not in other.
func (s Set[T]) Difference(other Set[T]) Set[T] {
	out := New[T]()

	for item := range s {
		if !other.Has(item) {
			out.Add(item)
		}
	}

	return out
}

// Clone returns a shallow copy.
func (s Set[T]) Clone() Set[T] {
	return s.Union(New[T]())
}
