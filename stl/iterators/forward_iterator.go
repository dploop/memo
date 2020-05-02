package iterators

type ForwardIterator interface {
	InputIterator
	// Multiple passes allowed.
}
