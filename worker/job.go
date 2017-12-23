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
func (jobPool *JobPool) AddJob(value ...interface{}) {
	jobPool.wg.Add(1)
	jobPool.job <- Job{
		Runtime: time.Now(),
		Value:   value,
	}
}

//Close the job pool
func (jobPool *JobPool) Close() {
	close(jobPool.job)
	jobPool.wg.Wait()
}

//SetWorkDisplay : enable or disable work display of worker
func (jobPool *JobPool) SetWorkDisplay(wd bool) {
	jobPool.workDisplay = wd
}

//StartWorker : start worker
func (jobPool *JobPool) StartWorker(noOfWorker int, handler Handler) {
	sTime := time.Now().Nanosecond()
	for i := 1; i <= noOfWorker; i++ {
		w := &Worker{
			workerID:    i + sTime,
			jobPool:     jobPool,
			logPool:     logPool,
			handler:     handler,
			log:         true,
			workDisplay: jobPool.workDisplay,
		}
		WorkerPool = append(WorkerPool, w)
		w.Start()
	}
}
