package constraints

import (
	"github.com/dploop/memo/stl/types"
)

type EqualityComparable interface {
	Equal(other EqualityComparable) bool
}

type Equality func(types.Data, types.Data) bool

type LessThanComparable interface {
	Less(other LessThanComparable) bool
}

type LessThan func(types.Data, types.Data) bool
