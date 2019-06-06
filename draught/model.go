package draught

import (
	"context"
	"sync"
	"time"
)

//Logger will be called if error occured
type Logger interface {
	Print(pool *Pool, value []interface{}, err error)
}

//Pool contains the jobs
type Pool struct {
	Tag        string //tag used to identify a pool
	pool       chan workerJob
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
	ID      int                //worker ID
	jobPool *Pool              //common job pool
	handler Handler            //handler to call
	job     workerJob          //current job
	ctx     context.Context    //each worker context
	cancel  context.CancelFunc //context cancel function
	working bool               //flag to check worker is idel or not
}

//workerJob : job assigned to a worker
type workerJob struct {
	value []interface{} //job value
	retry int           //number of retries done
	timer *time.Timer   //timer set if job fails
	err   []error       //all the retries error at their respective indices
}

//Handler function which will be called by the go routines
type Handler func(context.Context, ...interface{}) error
