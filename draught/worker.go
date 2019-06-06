package draught

import (
	"fmt"
	"log"
	"math"
	"time"
)

//start will start the worker
func (w *Worker) start() {
	go func() {
		for {
			select {
			case <-w.ctx.Done():
				return
			case job, open := <-w.jobPool.pool:
				if open == false {
					return
				}
				w.working = true
				w.processJob(job)
			default: //non blocking
			}
		}
	}()
}

//processJob will process the job
func (w *Worker) processJob(wj workerJob) {
	defer w.jobPool.wg.Done()
	defer func() {
		w.working = false
		if rec := recover(); rec != nil {
			w.jobPool.counterPool <- 0             //error
			w.log(wj.value, fmt.Errorf("%v", rec)) //logged the panic
		}
	}()

	if wj.timer != nil { //if timer is set then check if timeout is done or not
		select {
		case <-wj.timer.C:
			break //if timeout is done then process the job
		default:
			w.jobPool.retryJob(wj) //requeue the job
			return
		}
	}

	//calling the handler
	if err := w.handler(w.ctx, wj.value...); err == nil {
		w.jobPool.counterPool <- 1 //success
	} else {
		w.log(wj.value, err)       //logging the error
		w.retry(wj, err)           //adding job again to retry if possible
		w.jobPool.counterPool <- 0 //error
	}
}

func (w *Worker) log(value []interface{}, err error) {
	if w.jobPool.logger != nil { //if logger is set
		w.jobPool.logger.Print(w.jobPool, value, err)
	} else {
		log.Println(err)
	}
}

//retry will retry job
func (w *Worker) retry(wj workerJob, err error) {
	if wj.retry < w.jobPool.maxRetry {
		wj.err = append(wj.err, err)
		wj.retry++

		//retrying exponentially
		dur := int(math.Pow(w.jobPool.exponent, float64(wj.retry)))
		wj.timer = time.NewTimer(time.Duration(dur) * time.Millisecond)

		log.Printf("Retrying after: %v ms Job: %v\n", dur, wj.value)
		w.jobPool.counterPool <- 2 //retry
		w.jobPool.retryJob(wj)
	}
}
