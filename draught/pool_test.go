package draught

import (
	"context"
	"fmt"
	"math/rand"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

//I job struct
type I struct {
	priority int
	value    interface{}
}

func (i I) String() string {
	return fmt.Sprintf("Priority %v Value %v", i.priority, i.value)
}

func printIT(ctx context.Context, value ...interface{}) (err error) {
	return nil
}

func TestPool(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	handler := printIT
	pool := NewPool(2, "", nil)
	pool.AddWorker(2, handler, true)
	for i := 0; i < 20; i++ {
		pool.AddJob(I{value: i, priority: rand.Intn(3)})
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
			pool.AddJob(I{value: 1})
		}
	})
	pool.Close()
	// pool.Stats()
}
