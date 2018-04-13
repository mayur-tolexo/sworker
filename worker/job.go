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

//Close the job pool and wait until all the jobs are completed
func (jobPool *JobPool) Close() {
	close(jobPool.job)
	jobPool.wg.Wait()
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

//GetWorkers return the worker of the current jobpool
func (jobPool *JobPool) GetWorkers() []*worker {
	return jobPool.workerPool
}

//WorkerCount return the worker count
func (jobPool *JobPool) WorkerCount() int {
	return len(jobPool.workerPool)
}

//KillWorker will kill worker
func (jobPool *JobPool) KillWorker() {
	total := jobPool.WorkerCount()
	if total > 1 {
		jobPool.workerPool[0].quit <- 1
		jobPool.workerPool = jobPool.workerPool[1:]
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
