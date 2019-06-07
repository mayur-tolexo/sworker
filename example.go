package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/mayur-tolexo/sworker/draught"
)

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
