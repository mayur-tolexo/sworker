package worker

import (
	"time"
)

//NewJobPool create new job pool
func NewJobPool(bufferSize int) JobPool {
	return JobPool{
		job: make(chan Job, bufferSize),
	}
}

//AddJob new job in job pool
func (jp JobPool) AddJob(value ...interface{}) {
	jp.job <- Job{
		Runtime: time.Now(),
		Value:   value,
	}
}
