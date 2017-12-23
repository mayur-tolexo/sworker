package worker

import (
	"time"
)

//NewJobPool create new job pool
func NewJobPool(bufferSize int) *JobPool {
	return &JobPool{
		job: make(chan Job, bufferSize),
	}
}

//AddJob new job in job pool
func (jp *JobPool) AddJob(value ...interface{}) {
	jp.wg.Add(1)
	jp.job <- Job{
		Runtime: time.Now(),
		Value:   value,
	}
}

//Close the job pool
func (jp *JobPool) Close() {
	close(jp.job)
	jp.wg.Wait()
}
