package iterators

import (
	"github.com/dploop/memo/stl/constraints"
)

type InputIterator interface {
	Iterator
	constraints.EqualityComparable
	constraints.Readable
}
