package worker

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/fatih/color"
)

//startHandler : call handler of the current worker
func (w *worker) startHandler(job Job) {
	defer w.jobPool.wg.Done()
	if len(job.Value) == 0 {
		return
	}
	w.inProcess = true
	w.job = job

	sTime := time.Now()
	d := color.New(color.FgHiRed)
	defer func(jobValue interface{}) {
		w.inProcess = false
		if rec := recover(); rec != nil {
			w.jobPool.errorCounterPool <- true
			if w.jobPool.log {
				w.log(errorLog{logValue: rec, jobValue: jobValue})
			} else {
				d.Printf("\nPANIC RECOVERED:%v %v\n%v\nJOB VALUE: %v\n", w.jobPool.Tag, rec, string(debug.Stack()), jobValue)
			}
		}
	}(w.job.Value)
	if w.jobPool.workDisplay {
		fmt.Printf("worker: %d STARTED at %v:%v:%v\n", w.workerID,
			sTime.Hour(), sTime.Minute(), sTime.Second())
	}
	if err := w.handler(w.job.Value...); err != nil {
		w.jobPool.errorCounterPool <- true
		if w.jobPool.log {
			w.log(errorLog{logValue: err, jobValue: w.job.Value})
		} else {
			d.Printf("\nERROR IN PROCESSING HANDLER:%v %v\nJOB VALUE: %v\n", w.jobPool.Tag, err, w.job.Value)
			// w.jobPool.Stats()
		}
	} else {
		w.jobPool.jobCounterPool <- true
	}
	if w.jobPool.workDisplay {
		fmt.Printf("worker: %d END in %v SEC\n\n", w.workerID, time.Since(sTime).Seconds())
	}
}

//Start worker
func (w *worker) start() {
	go func() {
		for job := range w.jobPool.job {
			var quit bool
			select {
			case <-w.quit:
				quit = true
			default:
				quit = false
			}
			w.startHandler(job)
			if quit {
				break
			}
		}
	}()
}

//log start logging
func (w *worker) log(log errorLog) {
	defer func() {
		if rec := recover(); rec != nil {
			fmt.Printf("Error while logging: %v\n%v", rec, string(debug.Stack()))
		}
	}()
	w.jobPool.logError(log)
}
