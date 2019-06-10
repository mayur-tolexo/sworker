package worker

import (
	"os"
	"sync"
	"time"
)

//Logger will be called if error occurred
type Logger interface {
	Print(jp *JobPool, job interface{}, err interface{})
}

//JobPool contain job pool and wait group
type JobPool struct {
	Tag            string
	job            chan Job
	total          int
	sc             int
	scPool         chan bool
	ec             int
	ecPool         chan bool
	counterWG      sync.WaitGroup
	batchSize      int
	startTime      time.Time
	lastPrint      time.Time
	lastPrintCount int
	slowDuration   time.Duration
	ticker         *time.Ticker
	wg             sync.WaitGroup
	workDisplay    bool
	log            bool
	LogPath        string
	stackTrace     bool
	errorFP        *os.File
	workerPool     map[int]*worker
	Closed         bool
	handler        Handler
	logger         Logger
}

//errorLog model
type errorLog struct {
	logValue interface{}
	jobValue interface{}
}

//Job model
type Job struct {
	Runtime time.Time
	Value   []interface{}
}

//Handler function
type Handler func(value ...interface{}) error

//worker model
type worker struct {
	workerID  int
	jobPool   *JobPool
	quit      chan int
	handler   Handler
	job       Job
	inProcess bool
}
