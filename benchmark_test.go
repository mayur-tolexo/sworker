package sworker

import (
	"context"
	"fmt"
	"runtime"
	"testing"

	"github.com/mayur-tolexo/sworker/draught"
)

func print(ctx context.Context, value ...interface{}) (err error) {
	return nil
}

func BenchmarkDraught(b *testing.B) {
	for i := 1; i <= runtime.NumCPU()+1; i++ {
		b.Run(fmt.Sprintf("%v-Worker", i), func(b *testing.B) {
			runBenchmark(b, i)
		})
	}
}

func runBenchmark(b *testing.B, wCount int) {
	handler := print
	pool := draught.NewPool(b.N, "", nil)
	pool.AddWorker(wCount, handler, true)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			pool.AddJob(1)
		}
	})
	pool.Close()
}
