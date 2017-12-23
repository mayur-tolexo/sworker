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
	WorkerPool  []*Worker
)

const (
	BUFFER_SIZE = 100
)

func init() {
	logPool = make(chan interface{}, BUFFER_SIZE)
}

//NewWorker new worker creation
func NewWorker(noOfWorker int, jobPool *JobPool, handler Handler) {
	for i := 1; i <= noOfWorker; i++ {
		w := &Worker{
			workerID:    i,
			jobPool:     jobPool,
			logPool:     logPool,
			handler:     handler,
			log:         true,
			workDisplay: workDisplay,
		}
		WorkerPool = append(WorkerPool, w)
		w.Start()
	}
	return
}

//SetWorkDisplay enable or disable work display of worker
func SetWorkDisplay(wd bool) {
	workDisplay = wd
}

//starthandler call handler
func (w *Worker) starthandler(job Job) {
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
func (w *Worker) Start() {
	go func() {
		for job := range w.jobPool.job {
			w.starthandler(job)
		}
	}()
}
