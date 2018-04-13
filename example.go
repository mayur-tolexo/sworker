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
	jp := worker.NewJobPool(3, handler)

	// jp.KillWorker()

	//job added in jobpool
	jp.SetWorkDisplay(true)
	jp.AddJob("Hello", "Mayur")
	jp.AddJob("World", 12345)
	jp.AddJob(1001)

	//close the jobpool
	jp.Close()
}
