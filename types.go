package memo

import (
	"errors"
	"time"
)

var (
	ErrNotFound = errors.New("memo: not found")
)

type Option func(*memo)

const (
	NoExpire time.Duration = 0
)

func WithDefaultExpiration(defaultExpiration time.Duration) Option {
	return func(m *memo) {
		m.defaultExpiration = defaultExpiration
	}
}

type Loader func(interface{}) (interface{}, error)

func WithDefaultLoader(defaultLoader Loader) Option {
	return func(m *memo) {
		m.defaultLoader = defaultLoader
	}
}

func WithClock(clock Clock) Option {
	return func(m *memo) {
		m.clock = clock
	}
}
