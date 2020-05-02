package iterators

type MutableBidirectionalIterator interface {
	MutableIterator
	BidirectionalIterator
}
