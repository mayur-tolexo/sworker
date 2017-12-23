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

	//handler to which worker will call
	handler := PrintAll

	//jobpool created
	jp := worker.NewJobPool(4)
	//job added in jobpool
	jp.AddJob("Hello", "Mayur")
	jp.AddJob("World")
	jp.AddJob("YOYOYO")

	//5 worker started
	jp.StartWorker(5, handler)

	//close the jobpool
	jp.Close()
}
