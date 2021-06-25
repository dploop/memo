# Memo [![License](https://img.shields.io/badge/license-BSD%202--Clause-green.svg)](https://opensource.org/licenses/BSD-2-Clause) [![Build Status](https://github.com/dploop/memo/actions/workflows/build.yml/badge.svg)](https://github.com/dploop/memo/actions?query=branch%3Amaster)

## How to use

### Use as a simple map
```golang
var cache = memo.NewMemo()

func main() {
	cache.Set("x", 1)
	fmt.Println(cache.Get("x"))
}
```

### Concurrently safe
```golang
var cache = memo.NewMemo()

func main() {
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000000; j++ {
				cache.Set("x", j)
			}
		}()
	}
	wg.Wait()
	fmt.Println(cache.Get("x"))
}
```

### Support expiration
```golang
var cache = memo.NewMemo()

func main() {
	cache.Set("x", 1, memo.SetWithExpiration(time.Second))
	fmt.Println(cache.Get("x"))
	time.Sleep(time.Second)
	fmt.Println(cache.Get("x"))
}
```

### Support loader
```golang
var cache = memo.NewMemo()

func main() {
	fmt.Println(cache.Get("x"))
	loader := func(k memo.Key) (memo.Value, error) {
		return len(k.(string)), nil
	}
	fmt.Println(cache.Get("x", memo.GetWithLoader(loader)))
}
```

### Load exactly once per key
```golang
var cache = memo.NewMemo()

func main() {
	fmt.Println(cache.Get("x"))

	var loadCounter, loadSum int64
	loader := func(k memo.Key) (memo.Value, error) {
		time.Sleep(time.Second)
		atomic.AddInt64(&loadCounter, 1)
		return int64(1), nil
	}

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000000; j++ {
				value, _ := cache.Get("x", memo.GetWithLoader(loader))
				atomic.AddInt64(&loadSum, value.(int64))
			}
		}()
	}
	wg.Wait()
	fmt.Println("counter", atomic.LoadInt64(&loadCounter))
	fmt.Println("sum", atomic.LoadInt64(&loadSum))
	fmt.Println(cache.Get("x"))
}
```

