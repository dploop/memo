package iterators

import (
	"github.com/dploop/memo/stl/constraints"
)

type MutableIterator interface {
	constraints.Writeable
}
