package objpool

import (
	"fmt"
	"sync/atomic"
	"testing"
)

type Obj struct {
	name string
}

func TestPool(t *testing.T) {
	var count int64
	pool := NewPool(func(s string) (interface{}, error) {
		atomic.AddInt64(&count, 1)
		fmt.Println(atomic.LoadInt64(&count))
		return &Obj{name: s}, nil
	})

	t.Parallel()

	for {
		go func() {
			for i := 0; i < 100; i++ {
				v, err := pool.Get(fmt.Sprint("hello", i))
				fmt.Println(v, err)
			}
		}()
	}

	select {}
}
