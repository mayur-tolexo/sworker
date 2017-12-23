package main

import (
	"fmt"

	"github.com/mayur-tolexo/sworker/worker"
)

func main() {

	//handler created
	handler := func(value ...interface{}) bool {
		fmt.Println(value)
		return true
	}

	//jobpool created
	jp := worker.NewJobPool(2)

	//5 worker created
	worker.NewWorker(5, jp, handler)

	//job added in jobpool
	jp.AddJob("Hello", "Mayur")
	jp.AddJob("World")
	jp.AddJob("YOYOYO")

	//wait until all jobs of jobpool are not completed
	jp.Wait()
}
