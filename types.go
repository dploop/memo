package memo

import (
	"errors"
	"time"

	"github.com/dploop/memo/clock"
)

var (
	ErrNotFound          = errors.New("memo: not found")
	ErrInvalidExpiration = errors.New("memo: invalid expiration")
)

type (
	Key        = interface{}
	Value      = interface{}
	Clock      = clock.Clock
	Loader     func(Key) (Value, error)
	Expiration = time.Duration
)

type options struct {
	clock      Clock
	loader     Loader
	expiration Expiration
}

type Option func(*options)

func newOptions(opts ...Option) options {
	o := options{
		clock: clock.NewRealClock(),
	}
	for _, opt := range opts {
		opt(&o)
	}

	if o.expiration < 0 {
		panic(ErrInvalidExpiration)
	}

	return o
}

func (base *options) newGetOptions(opts ...Option) options {
	o := options{
		loader:     base.loader,
		expiration: base.expiration,
	}
	for _, opt := range opts {
		opt(&o)
	}

	if o.expiration < 0 {
		panic(ErrInvalidExpiration)
	}

	return o
}

func (base *options) newSetOptions(opts ...Option) options {
	o := options{
		expiration: base.expiration,
	}
	for _, opt := range opts {
		opt(&o)
	}

	if o.expiration < 0 {
		panic(ErrInvalidExpiration)
	}

	return o
}

func WithClock(clock Clock) Option {
	return func(o *options) {
		o.clock = clock
	}
}

func WithLoader(loader Loader) Option {
	return func(o *options) {
		o.loader = loader
	}
}

func WithExpiration(expiration Expiration) Option {
	return func(o *options) {
		o.expiration = expiration
	}
}
