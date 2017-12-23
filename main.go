package main

import (
	"fmt"

	"github.com/mayur-tolexo/sworker/worker"
)

func main() {
	jp := worker.NewJobPool(2)
	jp.AddJob("Hello", "Mayur")
	jp.AddJob("World")
	handler := func(value ...interface{}) bool {
		fmt.Println(value)
		return true
	}
	worker.NewWorker(5, jp, handler)
	return
}
