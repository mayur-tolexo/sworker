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

func BenchmarkSworker(b *testing.B) {
	for i := 1; i <= runtime.NumCPU()+1; i++ {
		// i := runtime.NumCPU()
		b.Run(fmt.Sprintf("%v-Worker", i), func(b *testing.B) {
			runDraughtBenchmark(b, i)
		})
	}
}

func runDraughtBenchmark(b *testing.B, wCount int) {
	handler := print
	pool := draught.NewPool(b.N, "", nil)
	pool.DisableCounter()
	pool.AddWorker(wCount, handler, true)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			pool.AddJob(1)
		}
	})
	pool.Close()
}

// func executeServe(p *slaves.Pool, rounds int) {
// 	for i := 0; i < rounds; i++ {
// 		p.Serve(i)
// 	}
// }

// func BenchmarkSlavePool(b *testing.B) {
// 	ch := make(chan int, b.N)

// 	sp := slaves.NewPool(0, func(obj interface{}) {
// 		ch <- obj.(int)
// 	})

// 	go executeServe(&sp, b.N)

// 	i := 0
// 	for i < b.N {
// 		select {
// 		case <-ch:
// 			i++
// 		}
// 	}
// 	close(ch)
// 	sp.Close()
// }

// func BenchmarkGrPool(b *testing.B) {
// 	n := runtime.NumCPU()
// 	b.Run(fmt.Sprintf("%v-Worker", n), func(b *testing.B) {
// 		// number of workers, and size of job queue
// 		pool := grpool.NewPool(n, b.N)
// 		defer pool.Release()

// 		// how many jobs we should wait
// 		pool.WaitCount(b.N)

// 		// submit one or more jobs to pool
// 		for i := 0; i < b.N; i++ {

// 			pool.JobQueue <- func() {
// 				// say that job is done, so we can know how many jobs are finished
// 				defer pool.JobDone()
// 			}
// 		}

// 		// wait until we call JobDone for all jobs
// 		pool.WaitAll()
// 	})
// }

// func BenchmarkTunny(b *testing.B) {
// 	pool := tunny.NewFunc(10, func(in interface{}) interface{} {
// 		intVal := in.(int)
// 		return intVal * 2
// 	})
// 	defer pool.Close()

// 	b.ResetTimer()

// 	for i := 0; i < b.N; i++ {
// 		ret := pool.Process(10)
// 		if exp, act := 20, ret.(int); exp != act {
// 			b.Errorf("Wrong result: %v != %v", act, exp)
// 		}
// 	}
// }

// func BenchmarkWorkerpool(b *testing.B) {
// 	n := runtime.NumCPU()
// 	b.Run(fmt.Sprintf("%v-Worker", n), func(b *testing.B) {
// 		wp := workerpool.New(n)
// 		defer wp.Stop()
// 		releaseChan := make(chan struct{})

// 		b.ResetTimer()

// 		// Start workers, and have them all wait on a channel before completing.
// 		for i := 0; i < b.N; i++ {
// 			wp.Submit(func() { <-releaseChan })
// 		}
// 		close(releaseChan)
// 	})
// }
