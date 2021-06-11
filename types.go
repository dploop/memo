package memo

import (
	"errors"
	"time"
	"unsafe"

	"github.com/dploop/memo/clock"
)

var (
	// ErrNotFound represents not found error.
	ErrNotFound = errors.New("memo: not found")
	// ErrInvalidExpiration represents invalid expiration error.
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
		opt((*options)(noescape(unsafe.Pointer(&o))))
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

type getOptions struct {
	loader     Loader
	expiration Expiration
}

type GetOption func(*getOptions)

func (base *options) newGetOptions(opts ...GetOption) getOptions {
	o := getOptions{
		loader:     base.loader,
		expiration: base.expiration,
	}
	for _, opt := range opts {
		opt((*getOptions)(noescape(unsafe.Pointer(&o))))
	}

	if o.expiration < 0 {
		panic(ErrInvalidExpiration)
	}

	return o
}

func GetWithLoader(loader Loader) GetOption {
	return func(o *getOptions) {
		o.loader = loader
	}
}

func GetWithExpiration(expiration Expiration) GetOption {
	return func(o *getOptions) {
		o.expiration = expiration
	}
}

type setOptions struct {
	expiration Expiration
}

type SetOption func(*setOptions)

func (base *options) newSetOptions(opts ...SetOption) setOptions {
	o := setOptions{
		expiration: base.expiration,
	}

	for _, opt := range opts {
		opt((*setOptions)(noescape(unsafe.Pointer(&o))))
	}

	if o.expiration < 0 {
		panic(ErrInvalidExpiration)
	}

	return o
}

func SetWithExpiration(expiration Expiration) SetOption {
	return func(o *setOptions) {
		o.expiration = expiration
	}
}

//go:nosplit
func noescape(p unsafe.Pointer) unsafe.Pointer {
	x := uintptr(p)

	return unsafe.Pointer(x ^ 0) //nolint:staticcheck
}
