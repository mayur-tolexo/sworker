package worker

import (
	"sync"
	"time"
)

//JobPool contain job pool and wait group
type JobPool struct {
	job         chan Job
	wg          sync.WaitGroup
	workDisplay bool
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
	workerID    int
	jobPool     *JobPool
	logPool     chan interface{}
	handler     Handler
	log         bool
	workDisplay bool
}
