package collections

import (
	"fmt"
	"math/rand"
	"slices"
	"testing"

	"github.com/pkg/profile"

	"github.com/brianvoe/gofakeit/v6"
)

func TestFilter(t *testing.T) {

	numbers := intSlice(1_000_000_00)

	filtered := Filter(Iter(numbers), func(i int) bool {
		return i%2 == 0
	})

	s := slices.Collect(filtered)

	t.Log(len(s))

}

func TestFilter2(t *testing.T) {

	numbers := intSlice(1_000_000_00)

	var out []int
	for _, n := range numbers {
		if n%2 == 0 {
			out = append(out, n)
		}
	}

	t.Log(len(out))

}

func intSlice(n int) []int {
	out := make([]int, n)
	for i := 0; i < n; i++ {
		out[i] = i
	}
	return out
}

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
	defer profile.Start(
		profile.MemProfile,
		profile.ProfilePath("./mem.out"),
		profile.MemProfileRate(1),
		profile.NoShutdownHook,
	).Stop()
	sizes := []int{
		100000,
		// 10000000,
		// 100000000,
	}
	for _, size := range sizes {
		b.Run(fmt.Sprintf("unique_int_size_%d", size), func(b *testing.B) {
			slice := randomIntSlice(size)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				Unique(slice)
			}
		})
		// b.Run(fmt.Sprintf("unique_string_size_%d", size), func(b *testing.B) {
		// 	slice := randomStringSlice(size)
		// 	b.ResetTimer()
		// 	for i := 0; i < b.N; i++ {
		// 		Unique(slice)
		// 	}
		// })
	}
}

// func TestUnique(t *testing.T) {

// 	defer profile.Start(
// 		profile.MemProfile,
// 		profile.ProfilePath("./mem.out"),
// 		profile.MemProfileRate(1),
// 		profile.NoShutdownHook,
// 	).Stop()

// 	slice := randomIntSlice(1000000)
// 	// Run the Unique function
// 	result := Unique(slice)

// 	// Use the result to prevent compiler optimization
// 	if len(result) == 0 {
// 		t.Fatal("Unexpected empty result")
// 	}

// }

// func TestUniqueIntMemProfile(t *testing.T) {
// 	// Skip this test unless explicitly run with -test.run=TestUniqueIntMemProfile
// 	if testing.Short() {
// 		t.Skip("Skipping memory profile test in short mode")
// 	}

// 	// Start memory profiling
// 	defer profile.Start(
// 		profile.MemProfile,
// 		profile.ProfilePath("."),  // Save in current directory
// 		profile.MemProfileRate(1), // Record all allocations
// 		profile.NoShutdownHook,
// 	).Stop()

// 	// Create a large slice of integers
// 	slice := randomIntSlice(10000000) // 10 million integers

// 	// Run the Unique function
// 	result := Unique(slice)

// 	// Use the result to prevent compiler optimization
// 	if len(result) == 0 {
// 		t.Fatal("Unexpected empty result")
// 	}
// }
