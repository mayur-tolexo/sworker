package main

import (
	"fmt"

	"github.com/mayur-tolexo/sworker/worker"
)

//PrintAll : function which worker will call to execute
func PrintAll(value ...interface{}) error {
	fmt.Println(value)
	return nil
}

//main function
func main() {

	//handler for the worker created
	handler := PrintAll

	//jobpool created
	jp := worker.NewJobPool(3)

	//5 worker created
	worker.NewWorker(5, jp, handler)

	//job added in jobpool
	jp.AddJob("Hello", "Mayur")
	jp.AddJob("World")
	jp.AddJob("YOYOYO")
	jp.Close()
}
