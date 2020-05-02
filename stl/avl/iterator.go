package avl

import (
	"github.com/dploop/memo/stl/constraints"
	"github.com/dploop/memo/stl/iterators"
	"github.com/dploop/memo/stl/types"
)

var _ iterators.MutableBidirectionalIterator = Iterator{}

type Iterator struct {
	n *node
}

func (i Iterator) Write(data types.Data) {
	i.n.data = data
}

func (i Iterator) Clone() constraints.Cloneable {
	return i.ImplClone()
}

func (i Iterator) ImplClone() Iterator {
	return i
}

func (i Iterator) Next() constraints.Incrementable {
	return i.ImplNext()
}

func (i Iterator) ImplNext() Iterator {
	i.n = successor(i.n)

	return i
}

func (i Iterator) Equal(other constraints.EqualityComparable) bool {
	return i.ImplEqual(other.(Iterator))
}

func (i Iterator) ImplEqual(other Iterator) bool {
	return i == other
}

func (i Iterator) Read() types.Data {
	return i.n.data
}

func (i Iterator) Prev() constraints.Decrementable {
	return i.ImplPrev()
}

func (i Iterator) ImplPrev() Iterator {
	i.n = predecessor(i.n)

	return i
}
