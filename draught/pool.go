package draught

import (
	"context"
	"fmt"
	"time"
)

//NewPool will create new pool
func NewPool(size int, tag string, logger Logger) *Pool {
	ctx, cancel := context.WithCancel(context.Background())
	p := Pool{
		Tag:         tag,
		pool:        make(chan workerJob, size),
		counterPool: make(chan int, size/2), //as their is 1/2 probability of success or failure
		ctx:         ctx,
		cancel:      cancel,
		logger:      logger,
		workerPool:  make(map[int]*Worker),
		exponent:    10.0,
	}
	p.startCount()
	return &p
}

//NewSimplePool will create new pool without any logger and tag
func NewSimplePool(size int) *Pool {
	return NewPool(size, "", nil)
}

//SetTag will set tag in the pool to identify the pool
func (p *Pool) SetTag(tag string) {
	p.Tag = tag
}

//SetMaxRetry will set max retries for jobs
func (p *Pool) SetMaxRetry(n int) {
	p.maxRetry = n
}

//SetRetryExponent will set retry exponential base value
func (p *Pool) SetRetryExponent(n float64) {
	p.exponent = n
}

//startCount will start counter on job pool
func (p *Pool) startCount() {
	p.countWG.Add(1) //one job added for counter to complete
	go func() {
		defer p.countWG.Done()
		for val := range p.counterPool {
			switch val {
			case 0:
				p.errCount++
			case 1:
				p.successCount++
			case 2:
				p.retryCount++
			case 3:
				p.totalCount++
			}
		}
	}()
}

//AddWorker will add worker in the pool.
//If start value is true then it will immediately start the worker as well
func (p *Pool) AddWorker(n int, handler Handler, start ...bool) {
	sTime := time.Now()
	for i := 1; i <= n; i++ {
		w := &Worker{
			ID:      i + sTime.Nanosecond(),
			jobPool: p,
			handler: handler,
			ctx:     p.ctx,
			cancel:  p.cancel,
		}
		p.mtx.Lock()
		p.workerPool[w.ID] = w
		p.wCount++
		p.mtx.Unlock()
		if len(start) > 0 && start[0] {
			w.start()
		}
	}
}

//Start will start all the workers added in the pool
func (p *Pool) Start() {
	p.mtx.RLock()
	defer p.mtx.RUnlock()
	for _, w := range p.workerPool {
		w.start()
	}
}

//AddJob will enqueue job in the pool
func (p *Pool) AddJob(value ...interface{}) {
	p.wg.Add(1)
	p.counterPool <- 3
	p.pool <- workerJob{value: value}
}

//retryJob will add job again in the pool
func (p *Pool) retryJob(job workerJob) {
	p.wg.Add(1)
	p.pool <- job
}

//Close will close the pool
func (p *Pool) Close() {
	p.wg.Wait()
	p.cancel()
	close(p.counterPool)
	p.countWG.Wait()
}

//Stats will print pool stats
func (p *Pool) Stats() {
	tag := p.Tag
	if tag == "" {
		tag = "Stats"
	}
	fmt.Printf("\n%v: Woker %d Job: Total %d Retry %d Success %d Error %d\n",
		tag, p.wCount, p.totalCount, p.retryCount, p.successCount, p.errCount)
}

//SuccessCount will return success count
func (p *Pool) SuccessCount() int {
	return p.successCount
}

//ErrorCount will return error count
func (p *Pool) ErrorCount() int {
	return p.errCount
}

//TotalCount will return total count
func (p *Pool) TotalCount() int {
	return p.totalCount
}

//RetryCount will return retry count
func (p *Pool) RetryCount() int {
	return p.retryCount
}
