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
	workerPool  []*Worker
}

//errorLog model
type errorLog struct {
	logValue interface{}
	logTime  time.Time
}

//Job model
type Job struct {
	Runtime time.Time
	Value   []interface{}
}

//Handler function
type Handler func(value ...interface{}) error

//Worker model
type Worker struct {
	workerID int
	jobPool  *JobPool
	handler  Handler
}
