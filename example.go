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
	n := 20                          //no of jobs
	pool := draught.NewSimplePool(n) //new job pool created
	pool.AddWorker(2, handler, true) //adding 2 workers
	for i := 0; i < n; i++ {
		pool.AddJob(i) //adding jobs
	}
	pool.Close() //closed the job pool
}
