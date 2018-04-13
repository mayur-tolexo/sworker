package worker

import (
	"fmt"
	"time"
)

//startHandler : call handler of the current worker
func (w *worker) startHandler(job Job) {
	defer w.jobPool.wg.Done()

	sTime := time.Now()
	defer func(jobValue interface{}) {
		if rec := recover(); rec != nil {
			if w.jobPool.log {
				w.log(errorLog{logValue: rec, jobValue: jobValue})
			} else {
				fmt.Printf("\nPANIC RECOVERED: %v\nJOB VALUE: %v\n", rec, jobValue)
			}
		}
	}(job.Value)
	if w.jobPool.workDisplay {
		fmt.Printf("worker: %d STARTED at %v:%v:%v\n", w.workerID,
			sTime.Hour(), sTime.Minute(), sTime.Second())
	}
	if err := w.handler(job.Value...); err != nil {
		if w.jobPool.log {
			w.log(errorLog{logValue: err.Error(), jobValue: job.Value})
		} else {
			fmt.Printf("\nERROR IN PROCESSING HANDLER: %v\nJOB VALUE: %v\n", err.Error(), job.Value)
		}
	}
	if w.jobPool.workDisplay {
		fmt.Printf("worker: %d END in %v SEC\n\n", w.workerID, time.Since(sTime).Seconds())
	}
}

//Start worker
func (w *worker) start() {
	go func() {
		for {
			select {
			case <-w.quit:
				return
			case job := <-w.jobPool.job:
				w.startHandler(job)
			default:
			}
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
