package main

import (
	"fmt"
	"time"

	"github.com/mayur-tolexo/sworker/worker"
)

//PrintAll : function which worker will call to execute
func PrintAll(value ...interface{}) error {
	fmt.Println(value)
	return nil
}

//main function
func main() {

	sTime := time.Now()

	//handler to which worker will call
	handler := PrintAll

	//jobpool created
	jp := worker.NewJobPool(10)

	//5 worker started
	jp.StartWorker(2, handler)

	//job added in jobpool
	jp.AddJob("Hello", "Mayur")
	jp.AddJob("World")
	jp.AddJob(1001)
	jp.KillWorker()

	for i := 0; i < 100; i++ {
		jp.AddJob(i)
	}
	//close the jobpool
	jp.Close()
	fmt.Println(time.Since(sTime))
}
