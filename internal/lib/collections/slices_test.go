package collections

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func generateRandomSlice(n int) []int {
	slice := make([]int, n)
	for i := 0; i < n; i++ {
		slice[i] = rand.Intn(n / 2) // This will ensure some duplicates
	}
	return slice
}

func BenchmarkUnique(b *testing.B) {
	sizes := []int{100000, 10000000, 100000000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			slice := generateRandomSlice(size)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				Unique(slice)
			}
		})
	}
}

func TestUnique(t *testing.T) {
	t.Run(fmt.Sprintf("size 100000000"), func(t *testing.T) {
		slice := generateRandomSlice(100000000)
		tm := time.Now()
		Unique(slice)
		t.Log(time.Since(tm))
	})
}

func TestUnique2(t *testing.T) {
	t.Run(fmt.Sprintf("size 100000000"), func(t *testing.T) {
		slice := generateRandomSlice(100000000)
		tm := time.Now()
		Unique2(slice)
		t.Log(time.Since(tm))
	})
}

func BenchmarkUnique2(b *testing.B) {
	sizes := []int{100000, 10000000, 100000000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			slice := generateRandomSlice(size)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				Unique2(slice)
			}
		})
	}
}
