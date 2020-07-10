package draught

import (
	"fmt"
	"math"
	"time"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

//start will start the worker
func (w *Worker) start() {
	w.once.Do(w.run)
}

func (w *Worker) run() {
	go func() {
		for w.job = range w.jobPool.pool {
			select {
			case <-w.quite:
				logrus.Debug("Closing worker", w.ID)
				w.quite <- struct{}{}
				return
			default:
				w.processJob()
			}
		}
	}()
}

//processJob will process the job
func (w *Worker) processJob() {
	var err error
	defer w.jobPool.wg.Done()
	defer func() {
		w.working = false
		if rec := recover(); rec != nil {
			err = fmt.Errorf("Panic: %v", rec)
			w.appendError(err) //appending error in worker job
			w.log(err)         //logged the panic
		}
	}()
	w.working = true

	if w.job.timer != nil { //if timer is set then check if timeout is done or not
		select {
		case <-w.job.timer.C:
			break //if timeout is done then process the job
		default:
			w.jobPool.retryJob(w.job) //requeue the job
			return
		}
	}

	//calling the handler
	if err = w.handler(w.jobPool.ctx, w.job.value...); err == nil {
		w.jobPool.counterPool <- 1 //success
	} else {
		w.appendError(err) //appending error in worker job
		w.log(err)         //logging the error
		w.retry(err)       //adding job again to retry if possible
	}
}

//appendError will addend error in worker job
//and enqueue in error pool if enabled
func (w *Worker) appendError(err error) {
	w.jobPool.counterPool <- 0 //error
	w.job.err = append(w.job.err, err)
	if w.jobPool.ePoolEnable {
		w.jobPool.ePool <- w.job
	}
}

func (w *Worker) log(err error) {
	if w.jobPool.logger != nil { //if logger is set
		w.jobPool.logger.Print(w.jobPool, w.job.value, err)
	} else {
		msg := fmt.Sprintf("\nERROR IN PROCESSING HANDLER:%v %v\nJOB VALUE: %v\n",
			w.jobPool.Tag, err, w.job.value)

		if w.jobPool.consoleLog {
			d := color.New(color.FgRed)
			d.Print(msg)
		} else {
			logrus.Println(msg)
		}
	}
}

//retry will retry job
func (w *Worker) retry(err error) {
	if w.job.retry < w.jobPool.maxRetry {
		w.job.retry++
		//retrying exponentially
		dur := int(math.Pow(w.jobPool.exponent, float64(w.job.retry)))
		w.job.timer = time.NewTimer(time.Duration(dur) * time.Millisecond)

		logrus.Printf("Retrying after: %v ms Job: %v\n", dur, w.job.value)
		w.jobPool.counterPool <- 2 //retry
		w.jobPool.retryJob(w.job)
	}
}
