package multidict

import (
	base "github.com/dploop/memo/stl/avl"
	"github.com/dploop/memo/stl/constraints"
	"github.com/dploop/memo/stl/iterators"
	"github.com/dploop/memo/stl/types"
)

var _ iterators.MutableBidirectionalIterator = Iterator{}

type Iterator struct {
	base base.Iterator
}

func (i Iterator) Write(data types.Data) {
	i.base.Write(data)
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
	i.base = i.base.ImplNext()

	return i
}

func (i Iterator) Equal(other constraints.EqualityComparable) bool {
	return i.ImplEqual(other.(Iterator))
}

func (i Iterator) ImplEqual(other Iterator) bool {
	return i == other
}

func (i Iterator) Read() types.Data {
	return i.base.Read()
}

func (i Iterator) Prev() constraints.Decrementable {
	return i.ImplPrev()
}

func (i Iterator) ImplPrev() Iterator {
	i.base = i.base.ImplPrev()

	return i
}
