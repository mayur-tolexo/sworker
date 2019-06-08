[![Build Status](https://travis-ci.org/mayur-tolexo/sworker.svg)](https://travis-ci.org/mayur-tolexo/sworker)
[![Godocs](https://img.shields.io/badge/golang-documentation-blue.svg)](https://www.godoc.org/github.com/mayur-tolexo/sworker/draught)
[![Go Report Card](https://goreportcard.com/badge/github.com/mayur-tolexo/sworker)](https://goreportcard.com/report/github.com/mayur-tolexo/sworker)
[![Release](https://img.shields.io/github/release/mayur-tolexo/sworker.svg?style=flat-square)](https://github.com/mayur-tolexo/sworker/releases)

# sworker
Easy worker setup for your code.  
Have some draught and let this repo manage your load using go routines.  
Checkout NSQ repo for msg queuing *-* [drift](https://github.com/mayur-tolexo/drift)

### install
```
go get github.com/mayur-tolexo/sworker/draught
```

### Benchmark
![Screenshot 2019-06-07 at 1 30 32 AM](https://user-images.githubusercontent.com/20511920/59062640-f744eb00-88c3-11e9-8701-48e51fe6f71d.png)

### Fetures
- [Recovery](#recovery)
- [Logger](#logger)
- [Error Pool](#error-pool)
- [Retries](#retries)
- [Set exponent base by which retry will be delayed](#set-exponent-base-of-retry)
- [Complete Pool Stats](#pool-stats)
	- [Total job count](#total-job-count)
	- [Successful job count](#successful-job-count)
	- [Failure job count](#failure-job-count)
	- [Retry count](#retry-count)
	- [Worker count](#worker-count)
- [Add Job](#add-job) _(Thread safe)_
- [Add Worker](#add-worker) _(Thread safe)_
- [Set Tag](#set-tag)
- [Profiler](#profiler) _(To print the current pool status after every specified interval)_
- Worker job process limit (Beyond which it will log worker current status)
- [Console log](#console-log)

### Example 1
Basic print job using two workers and 100 jobs
```
//print : function which worker will call to execute
func print(ctx context.Context, value ...interface{}) error {
	fmt.Println(value)
	return nil
}

//main function
func main() {
	handler := print                 //handler function which the go routine will call
	n := 100                          //no of jobs
	pool := draught.NewSimplePool(n) //new job pool created
	pool.AddWorker(2, handler, true) //adding 2 workers

	for i := 0; i < n; i++ {
		pool.AddJob(i) //adding jobs
	}
	pool.Close() //closed the job pool
}
```

### Example 2
```
//print : function which worker will call to execute
func print(ctx context.Context, value ...interface{}) error {
	fmt.Println(value)
	time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
	return nil
}

//main function
func main() {
	rand.Seed(time.Now().UnixNano())
	handler := print                 //handler function which the go routine will call
	n := 10                          //no of jobs
	pool := draught.NewSimplePool(n) //new job pool created
	pool.SetTag("PRINTER")           //added tag to the pool
	pool.SetProfiler(5)              //added profiler batch size
	pool.SetConsoleLog(true)         //enable console log
	pool.AddWorker(2, handler, true) //adding 2 workers

	for i := 0; i < n; i++ {
		pool.AddJob(i) //adding jobs
	}
	pool.Close() //closed the job pool
	pool.Stats() //pool stats
}
```
### Output
![Screenshot 2019-06-07 at 11 57 29 AM](https://user-images.githubusercontent.com/20511920/59085198-774a6f80-891b-11e9-903f-e3ac36fae790.png)

### Handler
```
//Handler function which will be called by the go routine
type Handler func(context.Context, ...interface{}) error

Here print is a handler function.  
Define your own handler and pass it to the workers and you are ready to go.
```

### Recovery
```
Each job contains it's own recovery.  
If any panic occured while processing the job  
then that panic will be logged and worker will continue doing next job.
```

### Logger
There are two ways to set logger in the pool.
- While creating the pool
- Using SetLogger() method after pool creation


```
type Logger struct{}
//Implementing Logger interface
func (l Logger)Print(pool *Pool, value []interface{}, err error){
}

// While creating the pool
NewPool(size int, tag string, Logger{})

// Using SetLogger() method after pool creation
pool.SetLogger(Logger{})
Console log will enable pool error and close notification logging in console
```

### Error Pool
```
pool.GetErrorPool()
This will return a channel of workerJob which contains error occured and job value.
```

### Retries
```
pool.SetMaxRetry(2)
To set maximum number of retires to be done if error occured while processing the job.  
Default is 0. Retry won't work if panic occured while processing the job.
```

### Set exponent base of retry
```
pool.SetRetryExponent(2)
If error occured then that job will be delayed exponentially.  
Default exponent base is 10. We can set to accordingly.
```

### Pool Stats
```
pool.Stats()
Pool's complete status.
```
#### Total job count
```
pool.TotalCount()
Total job processed by the pool workers.
```
#### Successful job count
```
pool.SuccessCount()
Successfully processed jobs count
```
#### Failure job count
```
pool.ErrorCount()
Error count while processing job
```
#### Retry count
```
pool.RetryCount()
Retry count while processing job
```
#### Worker count
```
pool.WorkerCount()
No of Worker added on pool
```

### Add job
```
for i := 1; i <= n; i++ {
	go func(){pool.AddJob(i)}()
}
You can add job in go routines as well.
```

### Add worker
```
pool.AddWorker(2, handler, true)
```

### Set tag
```
pool.SetTag("PRINTER")
You can set tag to a pool to identify its log.
```

### Profiler
```
pool.SetProfiler(1000)
To get the status of processed jobs after specified amount of work done.
```
#### Profiler Example 
```
PRINTER: Processed:3 jobs(total:10 success:1 error:2 retry:2) in 0.00179143 SEC
```

### Console Log
```
pool.SetConsoleLog(true)
Console log will enable pool error and close notification logging in console
```
