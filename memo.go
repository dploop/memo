package memo

import (
	"sync"
	"time"
)

type Memo interface {
	Get(key interface{}) (interface{}, error)
	GetWithLoader(key interface{}, loader Loader) (interface{}, error)
	GetWithLoaderExpiration(key interface{}, loader Loader, expiration time.Duration) (interface{}, error)
	Set(key interface{}, value interface{})
	SetWithExpiration(key interface{}, value interface{}, expiration time.Duration)
}

type memo struct {
	sync.Mutex
	defaultExpiration time.Duration
	defaultLoader     Loader
	clock             Clock
	cache             map[interface{}]*item
}

type item struct {
	sync.Mutex
	value    interface{}
	err      error
	expireAt int64
}

func NewMemo(options ...Option) Memo {
	m := &memo{
		defaultExpiration: NoExpire,
		defaultLoader:     nil,
		clock:             NewRealClock(),
		cache:             make(map[interface{}]*item),
	}
	for _, option := range options {
		option(m)
	}
	return m
}

func (m *memo) Get(key interface{}) (interface{}, error) {
	return m.GetWithLoader(key, m.defaultLoader)
}

func (m *memo) GetWithLoader(key interface{}, loader Loader) (interface{}, error) {
	return m.GetWithLoaderExpiration(key, loader, m.defaultExpiration)
}

func (m *memo) GetWithLoaderExpiration(key interface{}, loader Loader, expiration time.Duration) (value interface{}, err error) {
	m.Lock()
	i := m.cache[key]
	if loader == nil {
		if i == nil {
			m.Unlock()
		} else {
			m.Unlock()
			i.Lock()
			defer i.Unlock()
			if i.expireAt == 0 || i.expireAt > m.clock.Now().UnixNano() {
				return i.value, i.err
			}
		}
		value, err = nil, ErrNotFound
	} else {
		if i == nil {
			i = &item{}
			m.cache[key] = i
			m.Unlock()
			i.Lock()
			defer i.Unlock()
		} else {
			m.Unlock()
			i.Lock()
			defer i.Unlock()
			if i.expireAt == 0 || i.expireAt > m.clock.Now().UnixNano() {
				return i.value, i.err
			}
		}
		value, err = loader(key)
		i.value, i.err = value, err
		if expiration == NoExpire {
			i.expireAt = 0
		} else {
			i.expireAt = m.clock.Now().Add(expiration).UnixNano()
		}
	}
	return value, err
}

func (m *memo) Set(key interface{}, value interface{}) {
	m.SetWithExpiration(key, value, m.defaultExpiration)
}

func (m *memo) SetWithExpiration(key interface{}, value interface{}, expiration time.Duration) {
	m.Lock()
	i := m.cache[key]
	if i == nil {
		i = &item{}
		m.cache[key] = i
	}
	i.value, i.err = value, nil
	if expiration == NoExpire {
		i.expireAt = 0
	} else {
		i.expireAt = m.clock.Now().Add(expiration).UnixNano()
	}
	m.Unlock()
}
