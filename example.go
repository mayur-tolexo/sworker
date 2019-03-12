package main

import (
	"fmt"
	"time"

	"github.com/mayur-tolexo/sworker/worker"
)

//PrintAll : function which worker will call to execute
func PrintAll(value ...interface{}) error {
	fmt.Println(value)
	// time.Sleep(5 * time.Second)
	return nil
}

//main function
func main() {

	sTime := time.Now()

	//handler to which worker will call
	handler := PrintAll

	//jobpool created
	jp := worker.NewJobPool(10)
	jp.SetSlowDuration(5 * time.Second)

	//5 worker started
	jp.StartWorker(5, handler)

	//job added in jobpool
	jp.AddJob("Hello", "Mayur")
	jp.AddJob("World")
	jp.AddJob(1001)
	jp.KillWorker()

	for i := 0; i < 5; i++ {
		jp.AddJob(i)
	}
	//close the jobpool
	jp.KClose()
	jp.Stats()
	fmt.Println(time.Since(sTime))
}
