package worker

import (
	"fmt"
	"sync"
)

var (
	logPool chan interface{}
	WG      sync.WaitGroup
)

const (
	BUFFER_SIZE = 100
)

func init() {
	logPool = make(chan interface{}, BUFFER_SIZE)
}

//NewWorker new worker creation
func NewWorker(noOfWorker int, jobPool JobPool, handler Handler) {
	for i := 1; i <= noOfWorker; i++ {
		w := &Worker{
			workerID:    i,
			jobPool:     jobPool,
			logPool:     logPool,
			handler:     handler,
			sync:        true,
			log:         true,
			workDisplay: true,
		}
		w.Start()
	}
	return
}

//starthandler call handler
func (w *Worker) starthandler(job Job) {
	if w.handler(job.Value...) {
		fmt.Println(w.workerID, "DONE")
	}
}

//Start worker
func (w *Worker) Start() {
	select {
	case job := <-w.jobPool.job:
		w.starthandler(job)
	default:
		fmt.Println("not found")
	}
}
