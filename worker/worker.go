package worker

import (
	"fmt"
	"sync"
	"time"
)

var (
	logPool     chan interface{}
	WG          sync.WaitGroup
	workDisplay bool
)

const (
	BUFFER_SIZE = 100
)

func init() {
	logPool = make(chan interface{}, BUFFER_SIZE)
}

//startHandler : call handler of the current worker
func (w *Worker) startHandler(job Job) {
	sTime := time.Now()
	if w.workDisplay {
		fmt.Printf("Worker: %d STARTED at %v:%v:%v\n", w.workerID, sTime.Hour(), sTime.Minute(), sTime.Second())
	}
	w.handler(job.Value...)
	if w.workDisplay {
		fmt.Printf("Worker: %d END in %v SEC\n\n", w.workerID, time.Since(sTime).Seconds())
	}
	w.jobPool.wg.Done()
}

//Start worker
func (w *Worker) start() {
	go func() {
		for job := range w.jobPool.job {
			w.startHandler(job)
		}
	}()
}
