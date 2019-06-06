package draught

import (
	"context"
	"fmt"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func printIT(ctx context.Context, value ...interface{}) (err error) {
	return nil
}

func TestPool(t *testing.T) {

	handler := printIT
	pool := NewPool(2, "", nil)
	pool.AddWorker(2, handler, true)
	for i := 0; i < 20; i++ {
		pool.AddJob(i)
	}
	pool.Close()
	assert := assert.New(t)
	assert.Equal(pool.TotalCount(), pool.SuccessCount()+pool.ErrorCount(),
		"Total job should be equal to success + error job count")
}

func BenchmarkWorker(b *testing.B) {
	for i := 1; i <= runtime.NumCPU()+1; i++ {
		b.Run(fmt.Sprintf("%v-Worker", i), func(b *testing.B) {
			runBenchmark(b, i)
		})
	}
}

func runBenchmark(b *testing.B, wCount int) {
	handler := printIT
	pool := NewPool(b.N, "", nil)
	pool.AddWorker(wCount, handler, true)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			pool.AddJob(1)
		}
	})
	pool.Close()
	// pool.Stats()
}
