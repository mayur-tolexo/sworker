package draught

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/fatih/color"
)

//NewPool will create new pool
func NewPool(size int, tag string, logger Logger) *Pool {
	ctx, cancel := context.WithCancel(context.Background())
	p := Pool{
		Tag:        tag,
		pool:       make(chan workerJob, size),
		ePool:      make(chan workerJob, size/2), //as their is 1/2 probability of success or failure
		ctx:        ctx,
		cancel:     cancel,
		logger:     logger,
		workerPool: make(map[int]*Worker),
		exponent:   10.0,
	}
	p.counterPool = make(chan int, size/2) //as their is 1/2 probability of success or failure
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

//SetMaxRetry will set max retries for jobs.
//Default value is 0
func (p *Pool) SetMaxRetry(n int) {
	p.maxRetry = n
}

//SetRetryExponent will set retry exponential base value
func (p *Pool) SetRetryExponent(n float64) {
	p.exponent = n
}

//SetConsoleLog will enable/disable console logging
func (p *Pool) SetConsoleLog(enable bool) {
	p.consoleLog = enable
}

//SetProfiler will set profiler
//will fill print/log the job done in given batch size
func (p *Pool) SetProfiler(batchSize int) {
	p.profiler = batchSize
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
			if p.profiler != 0 {
				p.profile(p.totalCount, p.successCount, p.errCount, p.retryCount)
			}
		}
	}()
}

func (p *Pool) profile(total, success, errorCount, retry int) {
	processed := success + errorCount
	if processed != p.lastProfile && (processed)%p.profiler == 0 {
		p.lastProfile = processed
		tag := p.Tag
		if tag == "" {
			tag = "Stats"
		}

		msg := fmt.Sprintf("%v: Processed:%d jobs(total:%d success:%d error:%d retry:%d) in\t%.8f SEC\n",
			tag, processed, total, success, errorCount, retry, time.Since(p.sTime).Seconds())
		if p.consoleLog {
			d := color.New(color.FgHiBlue, color.Bold)
			d.Print(msg)
		} else {
			fmt.Print(msg)
		}
	}
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
		}
		p.mtx.Lock()
		p.workerPool[w.ID] = w
		p.wCount++
		p.mtx.Unlock()
		if len(start) > 0 && start[0] {
			w.start()
			p.sTime = sTime
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
	p.sTime = time.Now()
}

//AddJob will enqueue job in the pool
func (p *Pool) AddJob(value ...interface{}) {
	if p.closed == false {
		p.wg.Add(1)
		p.counterPool <- 3
		p.pool <- workerJob{value: value}
	}
}

//retryJob will add job again in the pool
func (p *Pool) retryJob(job workerJob) {
	if p.closed == false {
		p.wg.Add(1)
		p.pool <- job
	}
}

//Close will close the pool
func (p *Pool) Close() {
	p.wg.Wait()          //waiting for all job to be done
	p.closed = true      //marking pool as closed
	p.cancel()           //cancel all worker (go routines)
	close(p.counterPool) //close counter pool
	p.countWG.Wait()     //waiting for counter to complete the count
	if p.consoleLog {
		d := color.New(color.FgGreen, color.Bold)
		d.Printf("--- %s POOL CLOSED ---\n", p.Tag)
	}
}

//Stats will print pool stats
func (p *Pool) Stats() {
	tag := p.Tag
	if tag == "" {
		tag = "Stats"
	}
	msg := fmt.Sprintf("\n%v: Woker %d Jobs: Total %d Success %d Error %d Retry %d\n",
		tag, p.wCount, p.totalCount, p.successCount, p.errCount, p.retryCount)
	if p.consoleLog {
		d := color.New(color.FgHiMagenta, color.Bold)
		d.Print(msg)
	} else {
		log.Print(msg)
	}
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

//WorkerCount will return worker count
func (p *Pool) WorkerCount() int {
	return p.wCount
}
