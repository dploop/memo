package iterators

import (
	"github.com/dploop/memo/stl/constraints"
)

type Iterator interface {
	constraints.Cloneable
	constraints.Incrementable
}
