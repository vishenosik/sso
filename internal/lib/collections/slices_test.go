package collections

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
)

func randomIntSlice(n int) []int {
	slice := make([]int, n)
	for i := 0; i < n; i++ {
		slice[i] = rand.Intn(n / 2) // This will ensure some duplicates
	}
	return slice
}

func randomString(n int) []string {
	slice := make([]string, n)
	for i := 0; i < n; i++ {
		slice[i] = gofakeit.Email()
	}
	return slice
}

func randomStringSlice(n int) []string {
	random := randomString(n)
	slice := make([]string, n)
	for i := 0; i < n; i++ {
		slice[i] = gofakeit.RandomString(random)
	}
	return slice
}

func BenchmarkUnique(b *testing.B) {
	sizes := []int{
		100000,
		10000000,
		100000000,
	}
	for _, size := range sizes {
		b.Run(fmt.Sprintf("unique_int_size_%d", size), func(b *testing.B) {
			slice := randomIntSlice(size)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				Unique(slice)
			}
		})
		b.Run(fmt.Sprintf("unique_string_size_%d", size), func(b *testing.B) {
			slice := randomStringSlice(size)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				Unique(slice)
			}
		})
	}
}

// func TestUnique(t *testing.T) {
// 	t.Run(fmt.Sprintf("size 100000000"), func(t *testing.T) {
// 		slice := generateRandomSlice(100000000)
// 		tm := time.Now()
// 		Unique(slice)
// 		t.Log(time.Since(tm))
// 	})
// }
