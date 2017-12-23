package worker

import (
	"time"
)

type JobPool struct {
	job chan Job
}

//Job model
type Job struct {
	Runtime time.Time
	Value   []interface{}
}

//Handler function
type Handler func(value ...interface{}) bool

//Worker model
type Worker struct {
	workerID    int
	jobPool     JobPool
	logPool     chan interface{}
	handler     Handler
	sync        bool
	log         bool
	workDisplay bool
}
