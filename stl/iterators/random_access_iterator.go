package iterators

import (
	"github.com/dploop/memo/stl/constraints"
	"github.com/dploop/memo/stl/types"
)

type RandomAccessIterator interface {
	BidirectionalIterator
	constraints.LessThanComparable
	At(diff types.Size) types.Data
	Advance(diff types.Size) RandomAccessIterator
	Distance(other RandomAccessIterator) types.Size
}
