package draught

import (
	"context"
	"fmt"
	"sync"
	"time"
)

//Logger will be called if error occured.
//@params: pool, job value, error
type Logger interface {
	Print(*Pool, []interface{}, error)
}

//Pool contains the jobs
type Pool struct {
	Tag        string //tag used to identify a pool
	pool       chan *WorkerJob
	wg         sync.WaitGroup
	logger     Logger
	ctx        context.Context
	cancel     context.CancelFunc
	workerPool map[int]*Worker
	mtx        sync.RWMutex
	maxRetry   int
	exponent   float64 //retry exponent
	closed     bool
	counter
	flags
	pError
}

type pError struct {
	ePool       chan *WorkerJob
	ePoolEnable bool
}

//flags contains pool flags
type flags struct {
	sTime       time.Time    //pool start time
	consoleLog  bool         //this flag will enable console logging
	profiler    int          //if this flag is set then profiler will be activated
	lastProfile int          //contains last profiler count
	ticker      *time.Ticker //time profiler ticker
	tickerCount int          //time profiler count
	tickerLimit int          //time profiler count limit
}

//counter of the pool
type counter struct {
	errCount     int            //error count
	successCount int            //success count
	retryCount   int            //retry count
	wCount       int            //worker count
	totalCount   int            //total count
	counterPool  chan int       //0:error 1:success 2:retry 3:total
	countWG      sync.WaitGroup //counter wait group
}

//Worker will perform the job
type Worker struct {
	ID      int        //worker ID
	jobPool *Pool      //common job pool
	handler Handler    //handler to call
	job     *WorkerJob //current job
	working bool       //flag to check worker is idel or not
	once    sync.Once  //worker will start once
}

//WorkerJob : job assigned to a worker
type WorkerJob struct {
	value []interface{} //job value
	retry int           //number of retries done
	timer *time.Timer   //timer set if job fails
	err   []error       //all the retries error at their respective indices
}

//GetValue will return job value
func (wj *WorkerJob) GetValue() []interface{} {
	return wj.value
}

//GetError will return errors
func (wj *WorkerJob) GetError() []error {
	return wj.err
}

func (wj *WorkerJob) String() string {
	return fmt.Sprintf("Job value:%v error:%v", wj.value, wj.err)
}

//Handler function which will be called by the go routine
type Handler func(context.Context, ...interface{}) error
