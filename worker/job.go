package worker

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"
	"strings"
	"time"

	"github.com/fatih/color"
)

//NewJobPool create new job pool
func NewJobPool(bufferSize int) *JobPool {
	jp := &JobPool{
		job:        make(chan Job, bufferSize),
		workerPool: make(map[int]*worker),
	}

	jp.ticker = time.NewTicker(getSlowDuration(jp))
	return jp
}

//AddJob new job in job pool
func (jobPool *JobPool) AddJob(value ...interface{}) {
	jobPool.wg.Add(1)
	jobPool.total++
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
	// close(jobPool.job)
	jobPool.wg.Wait()
	// jobPool.KillWorker(jobPool.WorkerCount())
}

//KClose the job pool and wait until all the jobs are completed
//after complition will kill all routines
func (jobPool *JobPool) KClose() {
	jobPool.wg.Wait()
	close(jobPool.job)
	jobPool.Closed = true
	jobPool.KillWorker(jobPool.WorkerCount())
	jobPool.ticker.Stop()
	jobPool.counterWG.Wait()
	if jobPool.lastPrintCount != jobPool.sc {
		d := color.New(color.FgGreen, color.Bold)
		d.Printf("%d %s JOBs DONE IN %.8f SEC\n", jobPool.sc,
			jobPool.Tag, time.Since(jobPool.startTime).Seconds())
		d.Printf("--- %s POOL CLOSED ---\n\n", jobPool.Tag)
	}
}

//SetWorkDisplay : enable or disable work display of worker
func (jobPool *JobPool) SetWorkDisplay(wd bool) {
	jobPool.workDisplay = wd
}

//SetLogger : enable or disable error logger
func (jobPool *JobPool) SetLogger(log bool, path string) {
	jobPool.log = log
	jobPool.LogPath = path
	if log {
		jobPool.initErrorLog()
	}
}

//SetStackTrace : enable or disable error stackTrace
func (jobPool *JobPool) SetStackTrace(st bool) {
	jobPool.stackTrace = st
}

//SetSlowDuration : will set slow duration for job pool
func (jobPool *JobPool) SetSlowDuration(d time.Duration) {
	jobPool.slowDuration = d
	jobPool.ticker = time.NewTicker(getSlowDuration(jobPool))
}

//StartWorker : start worker
func (jobPool *JobPool) StartWorker(noOfWorker int, handler Handler) {
	sTime := time.Now()
	jobPool.startTime = sTime
	jobPool.lastPrint = sTime
	jobPool.scPool = make(chan bool, noOfWorker)
	jobPool.ecPool = make(chan bool, noOfWorker)
	jobPool.startCounter()
	jobPool.handler = handler
	jobPool.AddWorker(noOfWorker)
}

//AddWorker will add worker in job pool default handler
func (jobPool *JobPool) AddWorker(noOfWorker int) {
	for i := 1; i <= noOfWorker; i++ {
		sTime := time.Now()
		w := &worker{
			workerID: i + sTime.Nanosecond(),
			jobPool:  jobPool,
			quit:     make(chan int, 2),
			handler:  jobPool.handler,
		}
		jobPool.workerPool[w.workerID] = w
		w.start()
	}
}

func (jobPool *JobPool) startCounter() {
	go func() {
		jobPool.counterWG.Add(1)
		defer jobPool.counterWG.Done()
		for {
			select {
			case <-jobPool.scPool:
				jobPool.sc++
				if jobPool.batchSize != 0 && jobPool.sc%jobPool.batchSize == 0 {
					if jobPool.sc != jobPool.lastPrintCount {
						fmt.Printf("%d %s JOBs DONE IN %.8f SEC\n", jobPool.sc,
							jobPool.Tag, time.Since(jobPool.startTime).Seconds())
					}
					jobPool.lastPrint = time.Now()
					jobPool.lastPrintCount = jobPool.sc
				}
			case <-jobPool.ecPool:
				jobPool.ec++
			case <-jobPool.ticker.C:
				if jobPool.lastPrint.Before(time.Now().Add(-1 * getSlowDuration(jobPool))) {
					d := color.New(color.FgHiBlue, color.Bold)
					d.Printf("SLOW PROFILER - %d %s JOBs DONE IN %.8f SEC\n", jobPool.sc,
						jobPool.Tag, time.Since(jobPool.startTime).Seconds())
					jobPool.Stats()
					jobPool.WorkerJobs()
				}
				jobPool.lastPrint = time.Now()
				jobPool.lastPrintCount = jobPool.sc
			default:
				if jobPool.Closed {
					return
				}
			}
		}
	}()
}

//WorkerJobs will print worker current jobs
func (jobPool *JobPool) WorkerJobs() {
	for _, w := range jobPool.workerPool {
		if w.inProcess {
			fmt.Printf("%v JOB VALUE: %v\n", w.jobPool.Tag, w.job.Value)
		}
	}
	fmt.Println()
}

//Stats will print cur stats
func (jobPool *JobPool) Stats() {
	count := 0
	wCount := 0
	for _, w := range jobPool.workerPool {
		wCount++
		if w.inProcess {
			count++
		}
	}
	fmt.Printf("%v STATS - WORKER %d JOB: %d PENDING: %d IN-PROCESS: %d PROCESSED: %d ERROR: %d\n",
		jobPool.Tag, wCount, jobPool.total, len(jobPool.job), count, jobPool.sc, jobPool.ec)
}

//GetWorkers return the worker of the current jobpool
func (jobPool *JobPool) GetWorkers() map[int]*worker {
	return jobPool.workerPool
}

//SuccessCount return the successful job count
func (jobPool *JobPool) SuccessCount() int {
	return jobPool.sc
}

//ErrorCount return the worker error count
func (jobPool *JobPool) ErrorCount() int {
	return jobPool.ec
}

//GetBufferSize return the job buffer count
func (jobPool *JobPool) GetBufferSize() int {
	return cap(jobPool.job)
}

//ResetCounter will reset the job counter
func (jobPool *JobPool) ResetCounter() {
	jobPool.sc = 0
	jobPool.ec = 0
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
	if killCount > total {
		killCount = total
	}
	count := 0
	for workerID, w := range jobPool.workerPool {
		if count == killCount {
			break
		}
		w.quit <- 1
		delete(jobPool.workerPool, workerID)
		count++
	}
}

//initErrorLog will initialize logger
func (jobPool *JobPool) initErrorLog() {
	var fileErr error
	sTime := time.Now()
	if jobPool.LogPath == "" {
		jobPool.LogPath = jobPool.Tag + ".error_log"
	}
	path := fmt.Sprintf("%s_%d-%d-%d.log", strings.TrimSuffix(jobPool.LogPath, ".log"),
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
		logger.Printf("\nERROR:%v %v\nJOB VALUE: %v\nSTACK TRACE:\n%v", jobPool.Tag, err.logValue, err.jobValue,
			string(debug.Stack()))
	} else {
		logger.Printf("\nERROR:%v %v\nJOB VALUE: %v", jobPool.Tag, err.logValue, err.jobValue)
	}
}
