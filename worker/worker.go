package worker

import (
	"fmt"
	"time"
)

//startHandler : call handler of the current worker
func (w *Worker) startHandler(job Job) {
	defer w.jobPool.wg.Done()

	sTime := time.Now()
	defer func(sTime time.Time) {
		if rec := recover(); rec != nil {
			w.log(errorLog{logValue: rec, logTime: sTime})
		}
	}(sTime)

	if w.jobPool.workDisplay {
		fmt.Printf("Worker: %d STARTED at %v:%v:%v\n", w.workerID,
			sTime.Hour(), sTime.Minute(), sTime.Second())
	}
	w.handler(job.Value...)
	if w.jobPool.workDisplay {
		fmt.Printf("Worker: %d END in %v SEC\n\n", w.workerID, time.Since(sTime).Seconds())
	}
}

//Start worker
func (w *Worker) start() {
	go func() {
		for job := range w.jobPool.job {
			w.startHandler(job)
		}
	}()
}

//log start logging
func (w *Worker) log(log errorLog) {
	fmt.Println("HERE", log.logValue)
}
