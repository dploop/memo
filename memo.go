package memo

import (
	"sync"

	"github.com/dploop/memo/stl/multidict"
	"github.com/dploop/memo/stl/types"
)

type Memo struct {
	mu    sync.Mutex
	o     options
	cache map[Key]*entry
	timer *multidict.Dict
}

type entry struct {
	mu       sync.Mutex
	value    Value
	err      error
	expireAt int64
}

func NewMemo(opts ...Option) *Memo {
	keyComp := func(x types.Data, y types.Data) bool {
		return x.(int64) < y.(int64)
	}

	m := &Memo{
		o:     newOptions(opts...),
		cache: make(map[Key]*entry),
		timer: multidict.New(keyComp),
	}

	return m
}

func (m *Memo) Get(key Key, opts ...Option) (Value, error) {
	o := m.o.newGetOptions(opts...)
	m.mu.Lock()
	now := m.o.clock.Now()
	m.cleanup(now)

	var expireAt int64
	if o.expiration != 0 {
		expireAt = now + int64(o.expiration)
	}

	e := m.cache[key]
	if e != nil {
		m.mu.Unlock()
		e.mu.Lock()
		defer e.mu.Unlock()

		return e.value, e.err
	}

	if o.loader == nil {
		m.mu.Unlock()

		return nil, ErrNotFound
	}

	e = &entry{expireAt: expireAt}
	m.cache[key] = e
	m.timerInsert(e.expireAt, key)

	e.mu.Lock()
	m.mu.Unlock()
	defer e.mu.Unlock()
	e.value, e.err = o.loader(key)

	return e.value, e.err
}

func (m *Memo) Set(key Key, value Value, opts ...Option) {
	o := m.o.newSetOptions(opts...)
	m.mu.Lock()
	now := m.o.clock.Now()
	m.cleanup(now)

	var expireAt int64
	if o.expiration != 0 {
		expireAt = now + int64(o.expiration)
	}

	e := m.cache[key]
	if e == nil {
		e = &entry{value: value, expireAt: expireAt}
		m.cache[key] = e
		m.timerInsert(e.expireAt, key)
		m.mu.Unlock()

		return
	}

	if e.expireAt != expireAt {
		m.timerErase(e.expireAt, key)
		e.expireAt = expireAt
		m.timerInsert(e.expireAt, key)
	}

	m.mu.Unlock()
	e.mu.Lock()
	defer e.mu.Unlock()
	e.value, e.err = value, nil
}

func (m *Memo) Del(key Key) {
	m.mu.Lock()
	defer m.mu.Unlock()
	now := m.o.clock.Now()
	m.cleanup(now)

	e := m.cache[key]
	if e == nil {
		return
	}

	delete(m.cache, key)
	m.timerErase(e.expireAt, key)
}

func (m *Memo) cleanup(now int64) {
	i, last := m.timer.Begin(), m.timer.End()
	for !i.ImplEqual(last) {
		pair := i.Read().(multidict.Value)
		if pair.Key.(int64) > now {
			break
		}

		delete(m.cache, pair.Mapped)
		i = m.timer.Erase(i)
	}
}

func (m *Memo) timerInsert(expireAt int64, key Key) {
	if expireAt == 0 {
		return
	}

	m.timer.Insert(expireAt, key)
}

func (m *Memo) timerErase(expireAt int64, key Key) {
	if expireAt == 0 {
		return
	}

	i, last := m.timer.EqualRange(expireAt)
	for ; !i.ImplEqual(last); i = i.ImplNext() {
		pair := i.Read().(multidict.Value)
		if pair.Mapped == key {
			break
		}
	}

	_ = m.timer.Erase(i)
}
