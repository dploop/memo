package iterators

import (
	"github.com/dploop/memo/stl/constraints"
)

type BidirectionalIterator interface {
	ForwardIterator
	constraints.Decrementable
}
