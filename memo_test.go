package memo_test

import (
	"testing"
	"time"

	"github.com/dploop/memo"
)

func TestMemo_Simple(t *testing.T) {
	// Case: a is not in m
	m := memo.NewMemo()
	a, err := m.Get("a")
	if err != memo.ErrNotFound {
		t.Errorf("err(%v) is expected to be not found", err)
	}
	if a != nil {
		t.Errorf("a(%v) is expected to be nil", a)
	}
	// Case: a is in m
	m.Set("a", 1)
	a, err = m.Get("a")
	if err != nil {
		t.Errorf("err(%v) is expected to be nil", err)
	}
	if a.(int) != 1 {
		t.Errorf("a(%v) is expected to be 1", a)
	}
}

func TestMemo_Expiration(t *testing.T) {
	// Case: a is not in m
	fakeClock := memo.NewFakeClock()
	m := memo.NewMemo(
		memo.WithClock(fakeClock),
		memo.WithDefaultExpiration(10*time.Minute),
	)
	a, err := m.Get("a")
	if err != memo.ErrNotFound {
		t.Errorf("err(%v) is expected to be not found", err)
	}
	if a != nil {
		t.Errorf("a(%v) is expected to be nil", a)
	}
	// Case: a is in m
	m.Set("a", 1)
	a, err = m.Get("a")
	if err != nil {
		t.Errorf("err(%v) is expected to be nil", err)
	}
	if a.(int) != 1 {
		t.Errorf("a(%v) is expected to be 1", a)
	}
	// Case: a is in m
	fakeClock.Advance(9 * time.Minute)
	a, err = m.Get("a")
	if err != nil {
		t.Errorf("err(%v) is expected to be nil", err)
	}
	if a.(int) != 1 {
		t.Errorf("a(%v) is expected to be 1", a)
	}
	// Case: a is not in m
	fakeClock.Advance(1 * time.Minute)
	a, err = m.Get("a")
	if err != memo.ErrNotFound {
		t.Errorf("err(%v) is expected to be not found", err)
	}
	if a != nil {
		t.Errorf("a(%v) is expected to be nil", a)
	}
	// Case: a is in m
	m.SetWithExpiration("a", 1, 5*time.Minute)
	a, err = m.Get("a")
	if err != nil {
		t.Errorf("err(%v) is expected to be nil", err)
	}
	if a.(int) != 1 {
		t.Errorf("a(%v) is expected to be 1", a)
	}
	// Case: a is in m
	fakeClock.Advance(4 * time.Minute)
	a, err = m.Get("a")
	if err != nil {
		t.Errorf("err(%v) is expected to be nil", err)
	}
	if a.(int) != 1 {
		t.Errorf("a(%v) is expected to be 1", a)
	}
	// Case: a is not in m
	fakeClock.Advance(1 * time.Minute)
	a, err = m.Get("a")
	if err != memo.ErrNotFound {
		t.Errorf("err(%v) is expected to be not found", err)
	}
	if a != nil {
		t.Errorf("a(%v) is expected to be nil", a)
	}
}

func TestMemo_Loader(t *testing.T) {
	// Case: a is not in m
	fakeClock := memo.NewFakeClock()
	loader1 := func(_ interface{}) (interface{}, error) {
		return 1, nil
	}
	loader2 := func(_ interface{}) (interface{}, error) {
		return 2, nil
	}
	m := memo.NewMemo(
		memo.WithClock(fakeClock),
		memo.WithDefaultLoader(loader1),
	)
	a, err := m.GetWithLoader("a", nil)
	if err != memo.ErrNotFound {
		t.Errorf("err(%v) is expected to be not found", err)
	}
	if a != nil {
		t.Errorf("a(%v) is expected to be nil", a)
	}
	// Case: a is in m
	a, err = m.Get("a")
	if err != nil {
		t.Errorf("err(%v) is expected to be nil", err)
	}
	if a.(int) != 1 {
		t.Errorf("a(%v) is expected to be 1", a)
	}
	a, err = m.GetWithLoader("a", loader2)
	if err != nil {
		t.Errorf("err(%v) is expected to be nil", err)
	}
	if a.(int) != 1 {
		t.Errorf("a(%v) is expected to be 1", a)
	}
	// Case: a is not in m
	m.SetWithExpiration("a", 1, 5*time.Minute)
	fakeClock.Advance(5 * time.Minute)
	a, err = m.GetWithLoader("a", nil)
	if err != memo.ErrNotFound {
		t.Errorf("err(%v) is expected to be not found", err)
	}
	if a != nil {
		t.Errorf("a(%v) is expected to be nil", a)
	}
	// Case: a is in m
	a, err = m.GetWithLoaderExpiration("a", loader1, 5*time.Minute)
	if err != nil {
		t.Errorf("err(%v) is expected to be nil", err)
	}
	if a.(int) != 1 {
		t.Errorf("a(%v) is expected to be 1", a)
	}
	// Case: a is not in m
	fakeClock.Advance(5 * time.Minute)
	a, err = m.GetWithLoader("a", nil)
	if err != memo.ErrNotFound {
		t.Errorf("err(%v) is expected to be not found", err)
	}
	if a != nil {
		t.Errorf("a(%v) is expected to be nil", a)
	}
}
