[![Godocs](https://img.shields.io/badge/golang-documentation-blue.svg)](https://www.godoc.org/github.com/mayur-tolexo/sworker/draught)
[![Go Report Card](https://goreportcard.com/badge/github.com/mayur-tolexo/sworker)](https://goreportcard.com/report/github.com/mayur-tolexo/sworker)
[![Release](https://img.shields.io/github/release/mayur-tolexo/sworker.svg?style=flat-square)](https://github.com/mayur-tolexo/sworker/releases)

# sworker
Easy worker setup for your code.
Checkout NSQ repo for msg queuing *-* [drift](https://github.com/mayur-tolexo/drift)

### install
```
go get github.com/mayur-tolexo/sworker/draught
```

### Benchmark
![Screenshot 2019-06-07 at 1 30 32 AM](https://user-images.githubusercontent.com/20511920/59062640-f744eb00-88c3-11e9-8701-48e51fe6f71d.png)

### Example 1
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

Here print is a handler function. Define your own handler and pass it in the jobpool and you are ready to go.
```
