package memo_test

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/dploop/memo"
	"github.com/dploop/memo/clock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Memo", func() {
	expect := func(v memo.Value, err error) func(memo.Value, error) {
		return func(expectedV memo.Value, expectedErr error) {
			if expectedV == nil {
				Expect(v).To(BeNil())
			} else {
				Expect(v).To(Equal(expectedV))
			}
			if expectedErr == nil {
				Expect(err).NotTo(HaveOccurred())
			} else {
				Expect(err).To(MatchError(expectedErr))
			}
		}
	}

	Describe("Simple Memo", func() {
		It("can get nothing when empty", func() {
			m := memo.NewMemo()
			expect(m.Get("k"))(nil, memo.ErrNotFound)
		})
		It("can get something when set", func() {
			m := memo.NewMemo()
			m.Set("k", "v")
			expect(m.Get("k"))("v", nil)
		})
		It("can get nothing when set but deleted", func() {
			m := memo.NewMemo()
			m.Set("k", "v")
			expect(m.Get("k"))("v", nil)
			m.Del("k")
			expect(m.Get("k"))(nil, memo.ErrNotFound)
			m.Del("k")
			expect(m.Get("k"))(nil, memo.ErrNotFound)
			m.Set("k", "v", memo.SetWithExpiration(time.Minute))
			m.Del("k")
			expect(m.Get("k"))(nil, memo.ErrNotFound)
		})
	})

	vLoader := func(memo.Key) (memo.Value, error) {
		return "v", nil
	}
	uLoader := func(memo.Key) (memo.Value, error) {
		return "u", nil
	}
	Describe("Loader Memo", func() {
		It("can get something with specified loader", func() {
			m := memo.NewMemo()
			expect(m.Get("k", memo.GetWithLoader(vLoader)))("v", nil)
		})
		It("can get something with default loader", func() {
			m := memo.NewMemo(memo.WithLoader(vLoader))
			expect(m.Get("k"))("v", nil)
		})
		It("specified loader override default loader", func() {
			m := memo.NewMemo(memo.WithLoader(vLoader))
			expect(m.Get("k", memo.GetWithLoader(uLoader)))("u", nil)
		})
		It("loader will only take effect when needed", func() {
			m := memo.NewMemo()
			expect(m.Get("k", memo.GetWithLoader(vLoader)))("v", nil)
			expect(m.Get("k", memo.GetWithLoader(uLoader)))("v", nil)
		})
	})

	Describe("Expiration Memo", func() {
		var fakeClock *clock.FakeClock
		BeforeEach(func() {
			fakeClock = clock.NewFakeClock()
		})
		It("can get something when not expired", func() {
			m := memo.NewMemo(memo.WithClock(fakeClock))
			m.Set("k", "v", memo.SetWithExpiration(time.Minute))
			fakeClock.Advance(time.Minute - time.Second)
			expect(m.Get("k"))("v", nil)
		})
		It("expire something through specified expiration", func() {
			m := memo.NewMemo(memo.WithClock(fakeClock))
			m.Set("k", "v", memo.SetWithExpiration(time.Minute))
			fakeClock.Advance(time.Minute)
			expect(m.Get("k"))(nil, memo.ErrNotFound)
		})
		It("expire something through default expiration", func() {
			m := memo.NewMemo(memo.WithClock(fakeClock), memo.WithExpiration(time.Minute))
			m.Set("k", "v")
			fakeClock.Advance(time.Minute)
			expect(m.Get("k"))(nil, memo.ErrNotFound)
		})
		It("specified expiration will override default expiration", func() {
			m := memo.NewMemo(memo.WithClock(fakeClock), memo.WithExpiration(time.Hour))
			m.Set("k", "v", memo.SetWithExpiration(time.Minute))
			fakeClock.Advance(time.Minute)
			expect(m.Get("k"))(nil, memo.ErrNotFound)
		})
		It("specified expiration will override specified expiration", func() {
			m := memo.NewMemo(memo.WithClock(fakeClock))
			m.Set("k", "v")
			m.Set("k", "v", memo.SetWithExpiration(time.Minute))
			fakeClock.Advance(time.Minute)
			expect(m.Get("k"))(nil, memo.ErrNotFound)
			m.Set("k", "v", memo.SetWithExpiration(time.Minute))
			m.Set("k", "v")
			fakeClock.Advance(time.Minute)
			expect(m.Get("k"))("v", nil)
			m.Set("k", "v", memo.SetWithExpiration(1*time.Minute))
			m.Set("k", "v", memo.SetWithExpiration(2*time.Minute))
			fakeClock.Advance(time.Minute)
			expect(m.Get("k"))("v", nil)
			fakeClock.Advance(time.Minute)
			expect(m.Get("k"))(nil, memo.ErrNotFound)
		})
		It("expiration can be accurate to the key level", func() {
			m := memo.NewMemo(memo.WithClock(fakeClock))
			m.Set("k1", "v1", memo.SetWithExpiration(1*time.Minute))
			m.Set("k2", "v2", memo.SetWithExpiration(2*time.Minute))
			expect(m.Get("k1"))("v1", nil)
			expect(m.Get("k2"))("v2", nil)
			fakeClock.Advance(time.Minute)
			expect(m.Get("k1"))(nil, memo.ErrNotFound)
			expect(m.Get("k2"))("v2", nil)
			fakeClock.Advance(time.Minute)
			expect(m.Get("k1"))(nil, memo.ErrNotFound)
			expect(m.Get("k2"))(nil, memo.ErrNotFound)
		})
		It("should panic when expiration is invalid", func() {
			func() {
				defer func() {
					Expect(recover()).To(MatchError(memo.ErrInvalidExpiration))
				}()
				_ = memo.NewMemo(memo.WithExpiration(-1))
			}()
			func() {
				defer func() {
					Expect(recover()).To(MatchError(memo.ErrInvalidExpiration))
				}()
				m := memo.NewMemo()
				_, _ = m.Get("k", memo.GetWithExpiration(-1))
			}()
			func() {
				defer func() {
					Expect(recover()).To(MatchError(memo.ErrInvalidExpiration))
				}()
				m := memo.NewMemo()
				m.Set("k", "v", memo.SetWithExpiration(-1))
			}()
		})
	})

	Describe("Complicated Memo", func() {
		var fakeClock *clock.FakeClock
		BeforeEach(func() {
			fakeClock = clock.NewFakeClock()
		})
		It("can get nothing when loaded but expired", func() {
			m := memo.NewMemo(memo.WithClock(fakeClock))
			expect(m.Get("k", memo.GetWithLoader(vLoader), memo.GetWithExpiration(time.Minute)))("v", nil)
			fakeClock.Advance(time.Minute)
			expect(m.Get("k"))(nil, memo.ErrNotFound)
		})
		It("access a key concurrently will only load once", func() {
			var counter int32
			sLoader := func(memo.Key) (memo.Value, error) {
				time.Sleep(time.Second)
				atomic.AddInt32(&counter, 1)
				return "s", nil
			}
			m := memo.NewMemo(memo.WithClock(fakeClock), memo.WithLoader(sLoader))
			start := time.Now()
			var wg sync.WaitGroup
			for i := 0; i < 10000; i++ {
				k := i % 100
				wg.Add(1)
				go func() {
					defer wg.Done()
					_, _ = m.Get(k)
				}()
			}
			wg.Wait()
			Expect(counter).To(BeEquivalentTo(100))
			Expect(time.Since(start)).To(
				BeNumerically("~", time.Second, time.Second/10),
			)
		})
	})
})

func BenchmarkMemo_Get(b *testing.B) {
	m := memo.NewMemo()
	m.Set("k", "v")
	for i := 0; i < b.N; i++ {
		_, _ = m.Get("k")
	}
}

func BenchmarkMemo_Set(b *testing.B) {
	m := memo.NewMemo()
	for i := 0; i < b.N; i++ {
		m.Set("k", "v")
	}
	_, _ = m.Get("k")
}

func BenchmarkMemo_GetConcurrent(b *testing.B) {
	m := memo.NewMemo()
	m.Set("k", "v")
	var wg sync.WaitGroup
	concurrency := runtime.NumCPU()
	each := b.N / concurrency
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < each; j++ {
				_, _ = m.Get("k")
			}
		}()
	}
	wg.Wait()
}

func BenchmarkMemo_SetConcurrent(b *testing.B) {
	m := memo.NewMemo()
	var wg sync.WaitGroup
	concurrency := runtime.NumCPU()
	each := b.N / concurrency
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < each; j++ {
				m.Set("k", "v")
			}
		}()
	}
	wg.Wait()
}
