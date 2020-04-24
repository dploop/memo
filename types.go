package memo

import (
	"errors"
	"time"

	"github.com/dploop/memo/clock"
)

const (
	DoExpire Expiration = -1
	NoExpire Expiration = 0
)

var (
	ErrNotFound = errors.New("memo: not found")
)

type (
	Key        = interface{}
	Value      = interface{}
	Clock      = clock.Clock
	Loader     func(Key) (Value, error)
	Expiration = time.Duration
)

type Options struct {
	Clock      Clock
	Loader     Loader
	Expiration Expiration
}

type Option func(*Options)

func newOptions(opts ...Option) Options {
	o := Options{
		Clock: clock.NewRealClock(),
	}
	for _, opt := range opts {
		opt(&o)
	}
	return o
}

func (base *Options) newGetOptions(opts ...Option) Options {
	o := Options{
		Loader:     base.Loader,
		Expiration: base.Expiration,
	}
	for _, opt := range opts {
		opt(&o)
	}
	return o
}

func (base *Options) newSetOptions(opts ...Option) Options {
	o := Options{
		Expiration: base.Expiration,
	}
	for _, opt := range opts {
		opt(&o)
	}
	return o
}

func WithClock(clock Clock) Option {
	return func(o *Options) {
		o.Clock = clock
	}
}

func WithLoader(loader Loader) Option {
	return func(o *Options) {
		o.Loader = loader
	}
}

func WithExpiration(expiration Expiration) Option {
	return func(o *Options) {
		o.Expiration = expiration
	}
}
