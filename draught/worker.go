package draught

import (
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
			case job := <-w.jobPool.pool:
				w.processJob(job)
			default:
			}
		}
	}()
}

//processJob will process the job
func (w *Worker) processJob(wj workerJob) {
	defer w.jobPool.wg.Done()
	//if timer is set then check if timeout is done or not
	if wj.timer != nil {
		select {
		case <-wj.timer.C:
			break //if timeout is done then process the job
		default:
			w.jobPool.retryJob(wj)
			return
		}
	}

	if len(wj.value) == 0 {
		return
	}
	if err := w.handler(w.ctx, wj.value...); err == nil {
		w.jobPool.counterPool <- 1
	} else {
		//if logger is set
		if w.jobPool.logger != nil {
			w.jobPool.logger.Print(w.jobPool, wj.value, err)
		} else {
			log.Println(err)
		}
		w.retry(wj, err)
		w.jobPool.counterPool <- 0
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
		w.jobPool.counterPool <- 2
		w.jobPool.retryJob(wj)
	}
}
