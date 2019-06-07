package main

import (
	"context"
	"fmt"

	"github.com/mayur-tolexo/sworker/draught"
)

//printAll : function which worker will call to execute
func printAll(ctx context.Context, value ...interface{}) error {
	fmt.Println(value)
	return nil
}

//main function
func main() {
	handler := printAll              //handler function which the go routine will call
	n := 200                         //no of jobs
	pool := draught.NewSimplePool(n) //new job pool created
	pool.SetTag("Printer")           //added tag to the pool
	pool.SetProfiler(5)              //added profiler batch size
	pool.SetConsoleLog(true)         //enable console log
	pool.AddWorker(2, handler, true) //adding 2 workers
	for i := 0; i < n; i++ {
		pool.AddJob(i) //adding jobs
	}
	pool.Close() //closed the job pool
}
