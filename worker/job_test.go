package worker

import (
	"runtime"
	"testing"
)

//PrintAll : function which worker will call to execute
func PrintAll(value ...interface{}) error {
	// fmt.Println(value)
	// time.Sleep(5 * time.Second)
	return nil
}

func BenchmarkWorker1(b *testing.B) {
	runBenchmark(b, 1)
}

func BenchmarkWorker100(b *testing.B) {
	runBenchmark(b, 100)
}

func BenchmarkWorkerNumCPU(b *testing.B) {
	runBenchmark(b, runtime.NumCPU()+1)
}

func runBenchmark(b *testing.B, wCount int) {
	handler := PrintAll
	jp := NewJobPool(b.N)
	jp.StartWorker(wCount, handler)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			jp.AddJob(1)
		}
	})
	jp.KClose()
	jp.Stats()
}
