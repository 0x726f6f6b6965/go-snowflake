package snowflake

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}

var start time.Time

func setup() {
	start = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	fmt.Printf("\033[1;33m%s\033[0m", "> Setup completed\n")
}

func teardown() {
	fmt.Printf("\033[1;33m%s\033[0m", "> Teardown completed")
	fmt.Printf("\n")
}

func TestNextMonotonic(t *testing.T) {
	gen, _ := NewGenerator(10, start)
	out := make([]string, 100000)

	for i := range out {
		seq, _ := gen.Next()
		out[i] = seq.String()
	}

	// ensure they are all distinct and increasing
	for i := range out[1:] {
		if out[i] >= out[i+1] {
			t.Fatal("bad entries:", out[i], out[i+1])
		}
	}
}

func TestMultiCall(t *testing.T) {
	gen, _ := NewGenerator(3, start)
	c := make(chan uint64)
	times := rand.Intn(100000) + 1000
	num := rand.Intn(20) + 5
	go func() {
		defer close(c)
		var wg sync.WaitGroup
		defer wg.Wait()
		for i := 0; i < num; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for j := 0; j < times; j++ {
					seq, _ := gen.Next()
					c <- seq.Uint64()
				}
			}()
		}
	}()
	show := map[uint64]bool{}
	for v := range c {
		if show[v] {
			t.Fatal("get repeat squence")
		}
		show[v] = true
	}
	assert.Equal(t, times*num, len(show), "the count sequence should be equal")
}

func TestErrInitGenerator(t *testing.T) {
	var err error

	_, err = NewGenerator(maxNode+1, start)
	assert.ErrorIs(t, err, ErrInvalidNode)

	_, err = NewGenerator(5, time.Time{})
	assert.ErrorIs(t, err, ErrStartZero)

	_, err = NewGenerator(5, time.Now().Add(time.Hour))
	assert.ErrorIs(t, err, ErrStartFuture)

	_, err = NewGenerator(5, time.Date(1954, 1, 1, 0, 0, 0, 0, time.UTC))
	assert.ErrorIs(t, err, ErrStartExceed)

}

func TestClose(t *testing.T) {
	gen, _ := NewGenerator(7, start)
	gen.Close()
	_, err := gen.Next()
	assert.ErrorIs(t, err, ErrGeneratorClosed)
}

func TestNextExceed(t *testing.T) {
	date := time.Date(1954, 1, 1, 0, 0, 0, 0, time.UTC)
	g := generator{
		nodeID:    7,
		baseEpoch: date.UnixMilli(),
		sequence:  0, // sequence is now an int64
	}
	assert.NotNil(t, g)
	// g.sequence <- 0 // This line is no longer needed as sequence is not a channel
	_, err := g.Next()
	assert.ErrorIs(t, err, ErrStartExceed, err)
}

func BenchmarkNext(b *testing.B) {
	gen, err := NewGenerator(0, time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC))
	if err != nil {
		b.Fatalf("Failed to create generator: %v", err)
	}
	defer gen.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := gen.Next()
		if err != nil {
			b.Fatalf("Failed to get next ID: %v", err)
		}
	}
}
