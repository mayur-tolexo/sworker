package draught

import (
	"context"
	"fmt"
)

type myLogger struct {
}

func (l myLogger) Print(pool *Pool, value []interface{}, err error) {
	fmt.Println(value, err)
}

func egPrint(ctx context.Context, value ...interface{}) (err error) {
	return nil
}

func ExampleNewPool() {
	pool := NewPool(1, "Tag", myLogger{})
	pool.AddWorker(2, egPrint, true)
	pool.AddJob(1)
	pool.Close()
	fmt.Println(pool.GetStats())
	// Output: Tag: Woker 2 Jobs: Total 1 Success 1 Error 0 Retry 0
}
