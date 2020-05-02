package iterators

type MutableRandomAccessIterator interface {
	MutableIterator
	RandomAccessIterator
}
