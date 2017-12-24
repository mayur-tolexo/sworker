# sworker
Easy worker setup for your code.

# install
```
go get github.com/mayur-tolexo/sworker/worker
```

# Handler
```
handler is a function to which the worker will call. it is a function like:
func(value ...interface{}) error

Here PrintAll is a handler function. Define your own handler and pass it in the jobpool and you are ready to go.
```

# Example
```

//PrintAll : function which worker will call to execute
func PrintAll(value ...interface{}) error {
	fmt.Println(value)
	return nil
}

//main function
func main() {

	//handler to which worker will call
	handler := PrintAll

	//jobpool created
	jp := worker.NewJobPool(3)

	//job added in jobpool
	jp.AddJob("Hello", "Mayur")
	jp.AddJob("World")
	jp.AddJob(1001)

	//5 worker started
	jp.StartWorker(5, handler)

	//close the jobpool
	jp.Close()
}
```
# SetStackTrace
```
  To log complete stacktrace of the error
  By default it's false but to activate it
  jp.SetStackTrace(true)
```
# SetLogger
```
  To log all error
  By default it's true but to deactivate it
  jb.SetLogger(false)
``` 
# SetWorkDisplay
```
  To Print worker start time and end time while processing handler
  By default it's false but to activate it
  jb.SetWorkDisplay(true)
```

# To change error log path
```
create config.yaml file in package
add field logs: error_logs: YOUR_ERROR_PATH (as mentioned in config.yaml file in this repo)
```
