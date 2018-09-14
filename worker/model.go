package worker

import (
	"os"
	"sync"
	"time"
)

//JobPool contain job pool and wait group
type JobPool struct {
	Tag              string
	job              chan Job
	jobCounter       int
	wErrorCounter    int
	jobCounterPool   chan bool
	errorCounterPool chan bool
	batchSize        int
	startTime        time.Time
	lastPrint        time.Time
	lastPrintCount   int
	slowDuration     time.Duration
	ticker           *time.Ticker
	wg               sync.WaitGroup
	workDisplay      bool
	log              bool
	LogPath          string
	stackTrace       bool
	errorFP          *os.File
	workerPool       []*worker
	Closed           bool
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
	workerID int
	jobPool  *JobPool
	quit     chan int
	handler  Handler
}
