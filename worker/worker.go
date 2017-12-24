package worker

import (
	"fmt"
	"time"
)

//startHandler : call handler of the current worker
func (w *worker) startHandler(job Job) {
	defer w.jobPool.wg.Done()

	sTime := time.Now()
	defer func(sTime time.Time) {
		if rec := recover(); rec != nil {
			if w.jobPool.log {
				w.log(errorLog{logValue: rec, logTime: sTime})
			} else {
				fmt.Println("Panic Recovered:\n", rec)
			}
		}
	}(sTime)
	if w.jobPool.workDisplay {
		fmt.Printf("worker: %d STARTED at %v:%v:%v\n", w.workerID,
			sTime.Hour(), sTime.Minute(), sTime.Second())
	}
	if err := w.handler(job.Value...); err != nil {
		if w.jobPool.log {
			w.log(errorLog{logValue: err.Error(), logTime: sTime})
		} else {
			fmt.Println("Error while processing handler:\n", err.Error())
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
