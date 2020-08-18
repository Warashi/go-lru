package lru_test

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/Warashi/go-lru"
)

func TestCache(t *testing.T) {
	const size = 10
	c := lru.New(size)
	for i := 1; i <= size; i++ {
		if v := c.Store(i, i); v != i {
			t.Errorf("expected: %d, but %d\n", i, v)
		}
		if l := c.Len(); l != i {
			t.Errorf("expected: %d, but %d\n", i, l)
		}
	}

	for i := 1; i <= size; i++ {
		v, ok := c.Load(i)
		if !ok {
			t.Errorf("expectec: %t, but %t\n", true, ok)
		}
		if v != i {
			t.Errorf("expected: %d, but %d\n", i, v)
		}
	}

	for i := 1; i <= size; i++ {
		v, ok := c.Delete(i)
		if !ok {
			t.Errorf("expectec: %t, but %t\n", true, ok)
		}
		if v != i {
			t.Errorf("expected: %d, but %d\n", i, v)
		}
		if l := c.Len(); l != (size - i) {
			t.Errorf("expected: %d, but %d\n", size-i, l)
		}
	}

	for i := 1; i <= size; i++ {
		c.Store(i, i)
	}
	for i := 1; i <= size; i++ {
		c.Store(i, i)
		if l := c.Len(); l != size {
			t.Errorf("expected %d, but %d\n", size, l)
		}
	}
}

func BenchmarkCache_Store(b *testing.B) {
	for i := 0; i < 20; i++ {
		i := i
		b.Run(fmt.Sprintf("cap %d", i), func(b *testing.B) {
			size := 1 << i
			c := lru.New(size)
			b.RunParallel(func(pb *testing.PB) {
				rnd := rand.New(rand.NewSource(0))
				for pb.Next() {
					c.Store(rnd.Int(), rnd.Int())
				}
			})
		})
	}
}

func BenchmarkCache_Load(b *testing.B) {
	for i := 0; i < 20; i++ {
		i := i
		b.Run(fmt.Sprintf("cap %d", i), func(b *testing.B) {
			size := 1 << i
			c := lru.New(size)
			for j := 0; j < size; j++ {
				c.Store(j, j)
			}
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				rnd := rand.New(rand.NewSource(0))
				for pb.Next() {
					c.Load(rnd.Intn(size))
				}
			})
		})
	}
}
