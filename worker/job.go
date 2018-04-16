package worker

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/rightjoin/aero/conf"
)

//NewJobPool create new job pool
func NewJobPool(bufferSize int) *JobPool {
	jp := &JobPool{
		job: make(chan Job, bufferSize),
		log: true,
	}
	if jp.log {
		jp.initErrorLog()
	}
	return jp
}

//AddJob new job in job pool
func (jobPool *JobPool) AddJob(value ...interface{}) {
	jobPool.wg.Add(1)
	jobPool.job <- Job{
		Runtime: time.Now(),
		Value:   value,
	}
}

//CurrentBuffSize returns the current data size in buffer
func (jobPool *JobPool) CurrentBuffSize() (n int) {
	n = len(jobPool.job)
	return
}

//Close the job pool and wait until all the jobs are completed
func (jobPool *JobPool) Close() {
	close(jobPool.job)
	jobPool.wg.Wait()
	jobPool.KillWorker(jobPool.WorkerCount())
}

//SetWorkDisplay : enable or disable work display of worker
func (jobPool *JobPool) SetWorkDisplay(wd bool) {
	jobPool.workDisplay = wd
}

//SetLogger : enable or disable error logger
func (jobPool *JobPool) SetLogger(log bool) {
	jobPool.log = log
}

//SetStackTrace : enable or disable error stackTrace
func (jobPool *JobPool) SetStackTrace(st bool) {
	jobPool.stackTrace = st
}

//StartWorker : start worker
func (jobPool *JobPool) StartWorker(noOfWorker int, handler Handler) {
	sTime := time.Now()
	jobPool.startTime = sTime
	jobPool.jobCounterPool = make(chan bool, noOfWorker)
	jobPool.startCounter()

	for i := 1; i <= noOfWorker; i++ {
		w := &worker{
			workerID: i + sTime.Nanosecond(),
			jobPool:  jobPool,
			quit:     make(chan int, 2),
			handler:  handler,
		}
		jobPool.workerPool = append(jobPool.workerPool, w)
		w.start()
	}
}

func (jobPool *JobPool) startCounter() {
	go func() {
		for {
			select {
			case <-jobPool.jobCounterPool:
				jobPool.jobCounter++
				if jobPool.batchSize != 0 && jobPool.jobCounter%jobPool.batchSize == 0 {
					fmt.Printf("%d\tJOBs DONE IN\t%.8f SEC\n", jobPool.jobCounter, time.Since(jobPool.startTime).Seconds())
				}
			default:
			}
		}
	}()
}

//GetWorkers return the worker of the current jobpool
func (jobPool *JobPool) GetWorkers() []*worker {
	return jobPool.workerPool
}

//BatchSize will set profiling batch size for the counter
func (jobPool *JobPool) BatchSize(n int) {
	jobPool.batchSize = n
}

//WorkerCount return the worker count
func (jobPool *JobPool) WorkerCount() int {
	return len(jobPool.workerPool)
}

//KillWorker will kill worker
func (jobPool *JobPool) KillWorker(n ...int) {
	killCount := 1
	if len(n) > 0 {
		killCount = n[0]
	}
	total := jobPool.WorkerCount()
	if (total - 1) < killCount {
		killCount = total - 1
	}
	for i := 0; i < killCount; i++ {
		jobPool.workerPool[i].quit <- 1
	}
	if killCount == total {
		jobPool.workerPool = nil
	} else {
		jobPool.workerPool = jobPool.workerPool[killCount:]
	}
}

//initErrorLog will initialize logger
func (jobPool *JobPool) initErrorLog() {
	var fileErr error
	sTime := time.Now()
	path := conf.String("error_log", "logs.error_log")
	path = fmt.Sprintf("%s_%d-%d-%d.log", strings.TrimSuffix(path, ".log"),
		sTime.Day(), sTime.Month(), sTime.Year())
	if jobPool.errorFP, fileErr = os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND,
		0666); fileErr != nil {
		fmt.Println("Couldn't able to create the error log file", fileErr.Error())
	}
}

//logError will log given error
func (jobPool *JobPool) logError(err errorLog) {
	logger := log.New(jobPool.errorFP, "\n", 2)
	if jobPool.stackTrace {
		logger.Printf("\nERROR: %v\nJOB VALUE: %v\nSTACK TRACE:\n%v", err.logValue, err.jobValue,
			string(debug.Stack()))
	} else {
		logger.Printf("\nERROR: %v\nJOB VALUE: %v", err.logValue, err.jobValue)
	}
}
