package constraints

type Incrementable interface {
	Next() Incrementable
}
