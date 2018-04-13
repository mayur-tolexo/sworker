package worker

import (
	"os"
	"sync"
	"time"
)

//JobPool contain job pool and wait group
type JobPool struct {
	job         chan Job
	wg          sync.WaitGroup
	workDisplay bool
	log         bool
	stackTrace  bool
	errorFP     *os.File
	workerPool  []*worker
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
