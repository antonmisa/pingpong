package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func foo_ch(n int, in chan struct{}, out chan struct{}, wg *sync.WaitGroup, f func(string)) {
	defer wg.Done()

	for i := 0; i < n; i++ {
		f("Foo")
		in <- struct{}{}
		<-out
	}
}

func bar_ch(n int, in chan struct{}, out chan struct{}, wg *sync.WaitGroup, f func(string)) {
	defer wg.Done()

	for i := 0; i < n; i++ {
		<-out
		f("Bar\n")
		in <- struct{}{}
	}
}

func pingpong_ch(n int, f func(string)) {
	var wg sync.WaitGroup
	in, out := make(chan struct{}), make(chan struct{})

	wg.Add(2)
	go foo_ch(n, in, out, &wg, f)
	go bar_ch(n, out, in, &wg, f)

	wg.Wait()
	close(in)
	close(out)
}

func foo_lock(n int, c *sync.Cond, b *bool, wg *sync.WaitGroup, f func(string)) {
	defer wg.Done()

	for i := 0; i < n; i++ {
		c.L.Lock()
		for *b {
			c.Wait()
		}
		f("Foo")
		*b = true
		c.Signal()
		c.L.Unlock()
	}
}

func bar_lock(n int, c *sync.Cond, b *bool, wg *sync.WaitGroup, f func(string)) {
	defer wg.Done()

	for i := 0; i < n; i++ {
		c.L.Lock()
		for !*b {
			c.Wait()
		}
		f("Bar\n")
		*b = false
		c.Signal()
		c.L.Unlock()
	}
}

func pingpong_lock(n int, f func(string)) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var c *sync.Cond = sync.NewCond(&mu)
	b := false

	wg.Add(2)
	go foo_lock(n, c, &b, &wg, f)
	go bar_lock(n, c, &b, &wg, f)

	wg.Wait()
}

func foo_atomic(n int, bfoo *atomic.Bool, bbar *atomic.Bool, wg *sync.WaitGroup, f func(string)) {
	defer wg.Done()

	for i := 0; i < n; i++ {
		for {
			if bfoo.CompareAndSwap(false, true) {
				break
			}
		}
		f("Foo")
		bbar.Store(false)
	}
}

func bar_atomic(n int, bfoo *atomic.Bool, bbar *atomic.Bool, wg *sync.WaitGroup, f func(string)) {
	defer wg.Done()

	for i := 0; i < n; i++ {
		for {
			if bbar.CompareAndSwap(false, true) {
				break
			}
		}
		f("Bar\n")
		bfoo.Store(false)
	}
}

func pingpong_atomic(n int, f func(string)) {
	var wg sync.WaitGroup
	var bfoo atomic.Bool
	bfoo.Store(false)
	var bbar atomic.Bool
	bbar.Store(true)

	wg.Add(2)
	go foo_atomic(n, &bfoo, &bbar, &wg, f)
	go bar_atomic(n, &bfoo, &bbar, &wg, f)

	wg.Wait()
}

func main() {
	var f = func(s string) {
		fmt.Print(s)
	}

	pingpong_ch(10, f)
	fmt.Println("------------------------------")
	pingpong_lock(10, f)
	fmt.Println("------------------------------")
	pingpong_atomic(10, f)
}
