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

### Handler
```
//Handler function which will be called by the go routine
type Handler func(context.Context, ...interface{}) error

Here PrintAll is a handler function. Define your own handler and pass it in the jobpool and you are ready to go.
```

### Example
```
//printAll : function which worker will call to execute
func printAll(ctx context.Context, value ...interface{}) error {
	fmt.Println(value)
	return nil
}

//main function
func main() {
	handler := printAll              //handler function which the go routine will call
	n := 20                          //no of jobs
	pool := draught.NewSimplePool(n) //new job pool created
	pool.AddWorker(2, handler, true) //adding 2 workers
	for i := 0; i < n; i++ {
		pool.AddJob(i) //adding jobs
	}
	pool.Close() //closed the job pool
}
```
### SetStackTrace
```
  To log complete stacktrace of the error
  By default it's false but to activate it
  jp.SetStackTrace(true)
```
### SetLogger
```
  To log all error
  By default it's true but to deactivate it
  jb.SetLogger(false)
``` 
### SetWorkDisplay
```
  To Print worker start time and end time while processing handler
  By default it's false but to activate it
  jb.SetWorkDisplay(true)
```

### To change error log path
```
create config.yaml file in package
add field logs: error_logs: YOUR_ERROR_PATH (as mentioned in config.yaml file in this repo)
```

### Worker inside worker example
```
//ChildHandler : second handler
func ChildHandler(value ...interface{}) error {
	fmt.Println("CHILD", value)
	return nil
}

//PrintAll : function which worker will call to execute
func PrintAll(value ...interface{}) error {
	fmt.Println("PARENT", value)
	jp := worker.NewJobPool(1)
	jp.AddJob("World")
	jp.StartWorker(3, ChildHandler)
	jp.Close()
	return nil
}

//main function
func main() {

	//handler to which worker will call
	handler := PrintAll

	//jobpool created
	jp := worker.NewJobPool(3)

	//job added in jobpool
	jp.AddJob("Hello", "Hello")
	jp.AddJob("Mayur")
	jp.AddJob(1001)

	//5 worker started
	jp.StartWorker(5, handler)

	//close the jobpool
	jp.Close()
}
```
