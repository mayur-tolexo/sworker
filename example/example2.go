package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/mayur-tolexo/sworker/draught"
)

var n = 10 //no of jobs

//print : function which worker will call to execute
func print(ctx context.Context, value ...interface{}) error {
	if value[0].(int)%rand.Intn(n) == 0 {
		return fmt.Errorf("Random Error")
	}
	fmt.Println(value)
	time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)
	return nil
}

//main function
func main() {
	rand.Seed(time.Now().UnixNano())
	handler := print                 //handler function which the go routine will call
	pool := draught.NewSimplePool(n) //new job pool created
	ePool := pool.GetErrorPool()     //will return error pool
	pool.SetMaxRetry(1)              //setting max retry count
	pool.SetTag("PRINTER")           //added tag to the pool
	pool.SetProfiler(3)              //added profiler batch size
	pool.SetConsoleLog(true)         //enable console log
	pool.AddWorker(2, handler, true) //adding 2 workers

	for i := 1; i <= n; i++ {
		pool.AddJob(i) //adding jobs
	}
	pool.Close()  //closed the job pool
	wj := <-ePool //fetching only the first error
	fmt.Println(wj)
	pool.Stats() //pool stats
}
