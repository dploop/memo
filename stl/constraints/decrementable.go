package constraints

type Decrementable interface {
	Prev() Decrementable
}
