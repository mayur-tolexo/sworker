package worker

import (
	"fmt"
	"time"
)

//startHandler : call handler of the current worker
func (w *worker) startHandler(job Job) {
	defer w.jobPool.wg.Done()

	sTime := time.Now()
	if w.jobPool.log {
		defer func(sTime time.Time) {
			if rec := recover(); rec != nil {
				w.log(errorLog{logValue: rec, logTime: sTime})
			}
		}(sTime)
	}
	if w.jobPool.workDisplay {
		fmt.Printf("worker: %d STARTED at %v:%v:%v\n", w.workerID,
			sTime.Hour(), sTime.Minute(), sTime.Second())
	}
	if err := w.handler(job.Value...); err != nil {
		if w.jobPool.log {
			w.log(errorLog{logValue: err.Error(), logTime: sTime})
		}
	}
	if w.jobPool.workDisplay {
		fmt.Printf("worker: %d END in %v SEC\n\n", w.workerID, time.Since(sTime).Seconds())
	}
}

//Start worker
func (w *worker) start() {
	go func() {
		for job := range w.jobPool.job {
			w.startHandler(job)
		}
	}()
}

//log start logging
func (w *worker) log(log errorLog) {
	defer func() {
		if rec := recover(); rec != nil {
			fmt.Println("Error while logging:\n", rec)
		}
	}()
	w.jobPool.logError(log)
}
