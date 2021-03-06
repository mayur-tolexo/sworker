package draught

import (
	"context"
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

//NewPool will create new pool
func NewPool(size int, tag string, logger Logger) *Pool {
	ctx, cancel := context.WithCancel(context.Background())
	p := Pool{
		Tag:        tag,
		pool:       make(chan *WorkerJob, size),
		ctx:        ctx,
		cancel:     cancel,
		logger:     logger,
		workerPool: make(map[int]*Worker),
		exponent:   10.0,
	}
	p.counterPool = make(chan int, size/2) //as their is 1/2 probability of success or failure
	p.tickerLimit = 3
	p.startCount()
	return &p
}

//NewSimplePool will create new pool without any logger and tag
func NewSimplePool(size int) *Pool {
	return NewPool(size, "", nil)
}

//SetLogger will set logger
func (p *Pool) SetLogger(logger Logger) {
	p.logger = logger
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

//SetBatchProfiler will set profiler by job processed batch
//will fill print/log the job done in given batch size
func (p *Pool) SetBatchProfiler(batchSize int) {
	p.profiler = batchSize
}

//SetTimeProfiler will set profiler by time
func (p *Pool) SetTimeProfiler(dur time.Duration) {
	p.ticker = time.NewTicker(dur)
}

//GetErrorPool will return error pool
//if any error occurred then worker will push that error on error pool
func (p *Pool) GetErrorPool() <-chan *WorkerJob {
	p.ePool = make(chan *WorkerJob, cap(p.pool))
	p.ePoolEnable = true
	return p.ePool
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
			p.logProfile(p.totalCount, p.successCount, p.errCount, p.retryCount)
		}
	}()
}

func (p *Pool) getProfilerMsg(total, success, errorCount, retry int) string {
	processed := success + errorCount
	return fmt.Sprintf("%v: Processed:%d jobs(total:%d success:%d error:%d retry:%d) in %.8f SEC\n",
		p.getTag(), processed, total, success, errorCount, retry, time.Since(p.sTime).Seconds())
}

func (p *Pool) logProfile(total, success, errorCount, retry int) {
	processed := success + errorCount
	//if batch profiler is enabled
	if p.profiler != 0 && processed%p.profiler == 0 {
		p.profile(total, success, errorCount, retry, false)
	}
	//if time profiler is enabled
	if p.ticker != nil {
		select {
		case _, open := <-p.ticker.C:
			if open {
				if p.tickerCount == p.tickerLimit {
					p.tickerCount = 0
					p.workerStatus()
				}
				p.profile(total, success, errorCount, retry, true)
			} else {
				p.ticker = nil
			}
		default:
		}
	}
}

func (p *Pool) profile(total, success, errorCount, retry int, timeProfile bool) {
	processed := success + errorCount
	if processed != p.lastProfile {
		p.lastProfile = processed
		p.tickerCount = 0
		msg := p.getProfilerMsg(total, success, errorCount, retry)
		if p.consoleLog {
			var d *color.Color
			d = color.New(color.FgHiBlue, color.Bold)
			if timeProfile {
				d = color.New(color.FgBlack, color.Bold)
			}
			d.Print(msg)
		} else {
			logrus.Print(msg)
		}
	} else if timeProfile {
		p.tickerCount++
	}
}

//workerStatus will print worker current status
func (p *Pool) workerStatus() {
	p.mtx.RLock()
	defer p.mtx.RUnlock()

	d := color.New(color.FgBlack, color.Bold)
	msg := fmt.Sprintf("---%v WORKER STATUS---\n", p.getTag())
	if p.consoleLog {
		d.Print(msg)
	} else {
		logrus.Print(msg)
	}

	for _, w := range p.workerPool {
		if w.working {
			msg = fmt.Sprintf("Value %v Error %v\n",
				w.job.GetValue(), w.job.GetError())
			if p.consoleLog {
				d = color.New(color.FgBlack)
				d.Print(msg)
			} else {
				logrus.Print(msg)
			}
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
			quite:   make(chan struct{}),
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
		p.pool <- &WorkerJob{value: value}
	}
}

//retryJob will add job again in the pool
func (p *Pool) retryJob(job *WorkerJob) {
	if p.closed == false {
		p.wg.Add(1)
		p.pool <- job
	}
}

//Close will close the pool
func (p *Pool) Close() {
	p.wg.Wait()          //waiting for all job to be done
	p.closed = true      //marking pool as closed
	close(p.pool)        //closing pool after all work is done
	close(p.counterPool) //close counter pool
	p.countWG.Wait()     //waiting for counter to complete the count

	if p.ePoolEnable { //if error pool is enable
		p.ePoolEnable = false //disabling flag
		close(p.ePool)        //closing the error pool
	}

	if p.ticker != nil { //if time profiler enabled
		p.ticker.Stop() //stoping ticker
		p.ticker = nil  //disabling time profiler
	}

	if p.consoleLog {
		d := color.New(color.FgGreen, color.Bold)
		if p.lastProfile != (p.successCount + p.errCount) {
			msg := p.getProfilerMsg(p.totalCount, p.successCount, p.errCount, p.retryCount)
			d.Print(msg)
		}
		d.Printf("--- %s POOL CLOSED ---\n", p.Tag)
	}
}

//Stats will print pool stats
func (p *Pool) Stats() {
	msg := fmt.Sprintf("\n%v: Woker %d Jobs: Total %d Success %d Error %d Retry %d\n",
		p.getTag(), p.wCount, p.totalCount, p.successCount, p.errCount, p.retryCount)
	if p.consoleLog {
		d := color.New(color.FgBlack, color.Bold)
		d.Print(msg)
	} else {
		logrus.Print(msg)
	}
}

//GetStats will return pool stats
func (p *Pool) GetStats() string {
	return fmt.Sprintf("\n%v: Woker %d Jobs: Total %d Success %d Error %d Retry %d\n",
		p.getTag(), p.wCount, p.totalCount, p.successCount, p.errCount, p.retryCount)
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

//PoolCap will return pool capacity
func (p *Pool) PoolCap() int {
	return cap(p.pool)
}

//PoolLen will return pool length
func (p *Pool) PoolLen() int {
	return len(p.pool)
}

func (p *Pool) getTag() string {
	tag := p.Tag
	if tag == "" {
		tag = "Pool"
	}
	return tag
}

// CloseGracefully will terminate all worker gracefully
func (p *Pool) CloseGracefully() {
	for _, w := range p.workerPool {
		w.quite <- struct{}{}
		<-w.quite
	}
}
