package multidict

import (
	base "github.com/dploop/memo/stl/avl"
	"github.com/dploop/memo/stl/constraints"
	"github.com/dploop/memo/stl/types"
)

type Dict struct {
	base *base.Tree
	comp constraints.LessThan
}

type Value struct {
	Key    types.Data
	Mapped types.Data
}

func New(keyComp constraints.LessThan) *Dict {
	valueComp := func(x types.Data, y types.Data) bool {
		return keyComp(x.(Value).Key, y.(Value).Key)
	}

	return &Dict{
		base: base.New(valueComp),
		comp: keyComp,
	}
}

func (d *Dict) Size() types.Size {
	return d.base.Size()
}

func (d *Dict) Empty() bool {
	return d.base.Empty()
}

func (d *Dict) Begin() Iterator {
	return Iterator{base: d.base.Begin()}
}

func (d *Dict) End() Iterator {
	return Iterator{base: d.base.End()}
}

func (d *Dict) ReverseBegin() Iterator {
	return Iterator{base: d.base.ReverseBegin()}
}

func (d *Dict) ReverseEnd() Iterator {
	return Iterator{base: d.base.ReverseEnd()}
}

func (d *Dict) Count(k types.Data) types.Size {
	return d.base.CountMulti(Value{Key: k})
}

func (d *Dict) Find(k types.Data) Iterator {
	return Iterator{base: d.base.Find(Value{Key: k})}
}

func (d *Dict) Contains(k types.Data) bool {
	return d.base.Contains(Value{Key: k})
}

func (d *Dict) EqualRange(k types.Data) (Iterator, Iterator) {
	lb, ub := d.base.EqualRangeMulti(Value{Key: k})

	return Iterator{base: lb}, Iterator{base: ub}
}

func (d *Dict) LowerBound(k types.Data) Iterator {
	return Iterator{base: d.base.LowerBound(Value{Key: k})}
}

func (d *Dict) UpperBound(k types.Data) Iterator {
	return Iterator{base: d.base.UpperBound(Value{Key: k})}
}

func (d *Dict) Clear() {
	d.base.Clear()
}

func (d *Dict) Insert(k types.Data, m types.Data) Iterator {
	i := d.base.InsertMulti(Value{Key: k, Mapped: m})

	return Iterator{base: i}
}

func (d *Dict) Erase(i Iterator) Iterator {
	return Iterator{base: d.base.Delete(i.base)}
}
