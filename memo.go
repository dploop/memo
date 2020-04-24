package memo

import (
	"sync"
)

type Memo struct {
	sync.Mutex
	Options
	cache map[Key]*item
}

type item struct {
	sync.Mutex
	value    Value
	err      error
	expireAt int64
}

func NewMemo(opts ...Option) *Memo {
	m := &Memo{
		Options: newOptions(opts...),
		cache:   make(map[Key]*item),
	}
	return m
}

func (m *Memo) Get(key Key, opts ...Option) (value Value, err error) {
	o := m.newGetOptions(opts...)
	m.Lock()
	i := m.cache[key]
	if o.Loader == nil {
		if i == nil {
			m.Unlock()
		} else {
			m.Unlock()
			i.Lock()
			defer i.Unlock()
			if i.expireAt == 0 || i.expireAt > m.Clock.Now().UnixNano() {
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
			if i.expireAt == 0 || i.expireAt > m.Clock.Now().UnixNano() {
				return i.value, i.err
			}
		}
		value, err = o.Loader(key)
		i.value, i.err = value, err
		if o.Expiration == NoExpire {
			i.expireAt = 0
		} else {
			i.expireAt = m.Clock.Now().Add(o.Expiration).UnixNano()
		}
	}
	return value, err
}

func (m *Memo) Set(key Key, value Value, opts ...Option) {
	o := m.newSetOptions(opts...)
	m.Lock()
	i := m.cache[key]
	if i == nil {
		i = &item{}
		m.cache[key] = i
	}
	i.value, i.err = value, nil
	if o.Expiration == NoExpire {
		i.expireAt = 0
	} else {
		i.expireAt = m.Clock.Now().Add(o.Expiration).UnixNano()
	}
	m.Unlock()
}
