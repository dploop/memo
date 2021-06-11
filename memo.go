package memo

import (
	"container/heap"
	"sync"
)

type Memo struct {
	mu sync.Mutex
	o  options
	c  *cache
}

func NewMemo(opts ...Option) *Memo {
	return &Memo{o: newOptions(opts...), c: newCache()}
}

func (m *Memo) Get(key Key, opts ...GetOption) (Value, error) {
	o := m.o.newGetOptions(opts...)
	now := m.o.clock.Now()

	var expireAt int64
	if o.expiration != 0 {
		expireAt = now + int64(o.expiration)
	}

	m.mu.Lock()
	m.cleanup(now)

	e := m.c.dict[key]
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

	e = newEntry()
	m.c.dict[key] = e

	if expireAt != zeroExpireAt {
		heap.Push(m.c, node{key: key, expireAt: expireAt})
	}

	e.mu.Lock()
	m.mu.Unlock()
	defer e.mu.Unlock()
	e.value, e.err = o.loader(key)

	return e.value, e.err
}

func (m *Memo) Set(key Key, value Value, opts ...SetOption) {
	o := m.o.newSetOptions(opts...)
	now := m.o.clock.Now()

	var expireAt int64
	if o.expiration != 0 {
		expireAt = now + int64(o.expiration)
	}

	m.mu.Lock()
	m.cleanup(now)

	e := m.c.dict[key]
	if e == nil {
		e = newEntry()
		e.value = value
		m.c.dict[key] = e

		if expireAt != zeroExpireAt {
			heap.Push(m.c, node{key: key, expireAt: expireAt})
		}

		m.mu.Unlock()

		return
	}

	switch {
	case e.position == zeroPosition && expireAt != zeroExpireAt:
		heap.Push(m.c, node{key: key, expireAt: expireAt})
	case e.position != zeroPosition && expireAt == zeroExpireAt:
		heap.Remove(m.c, e.position)
	case e.position != zeroPosition && m.c.heap[e.position].expireAt != expireAt:
		m.c.heap[e.position].expireAt = expireAt
		heap.Fix(m.c, e.position)
	}

	m.mu.Unlock()
	e.mu.Lock()
	defer e.mu.Unlock()
	e.value, e.err = value, nil
}

func (m *Memo) Del(key Key) {
	now := m.o.clock.Now()
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cleanup(now)

	e := m.c.dict[key]
	if e == nil {
		return
	}

	if e.position != zeroPosition {
		heap.Remove(m.c, e.position)
	}

	delete(m.c.dict, key)
}

func (m *Memo) cleanup(now int64) {
	for m.c.heapSize != 0 {
		top := m.c.heap[0]
		if top.expireAt > now {
			break
		}

		_ = heap.Pop(m.c)
		delete(m.c.dict, top.key)
	}
}

type cache struct {
	dict     map[Key]*entry
	heap     []node
	heapSize int
}

func newCache() *cache {
	return &cache{dict: make(map[Key]*entry)}
}

const zeroPosition = -1

type entry struct {
	mu       sync.Mutex
	position int
	value    Value
	err      error
}

func newEntry() *entry {
	return &entry{position: zeroPosition}
}

const zeroExpireAt = 0

type node struct {
	key      Key
	expireAt int64
}

func (c *cache) Len() int {
	return c.heapSize
}

func (c *cache) Less(i, j int) bool {
	return c.heap[i].expireAt < c.heap[j].expireAt
}

func (c *cache) Swap(i, j int) {
	if i != j {
		c.heap[i], c.heap[j] = c.heap[j], c.heap[i]
		c.dict[c.heap[i].key].position = i
		c.dict[c.heap[j].key].position = j
	}
}

func (c *cache) Push(x interface{}) {
	if c.heapSize == len(c.heap) {
		c.heap = append(c.heap, node{})
	}

	c.heap[c.heapSize] = x.(node)
	c.dict[c.heap[c.heapSize].key].position = c.heapSize
	c.heapSize++
}

func (c *cache) Pop() interface{} {
	c.heapSize--
	c.dict[c.heap[c.heapSize].key].position = zeroPosition

	return c.heap[c.heapSize]
}
