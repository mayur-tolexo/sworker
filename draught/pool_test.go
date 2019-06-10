package draught

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func print(ctx context.Context, value ...interface{}) (err error) {
	return nil
}

func TestPool(t *testing.T) {

	handler := print
	pool := NewPool(2, "", nil)
	pool.AddWorker(2, handler, true)
	for i := 0; i < 20; i++ {
		go pool.AddJob(i)
	}
	pool.Close()
	assert := assert.New(t)
	assert.Equal(pool.TotalCount(), pool.SuccessCount()+pool.ErrorCount(),
		"Total job should be equal to success + error job count")
}

type logger struct {
	t *testing.T
}

func (l logger) Print(pool *Pool, value []interface{}, err error) {
	assert := assert.New(l.t)
	assert.Contains(err.Error(), "Panic Print")
}

func panicPrint(ctx context.Context, value ...interface{}) (err error) {
	if true {
		panic("Panic Print")
	}
	return nil
}

func TestLogger(t *testing.T) {
	pool := NewPool(1, "", logger{t})
	pool.AddWorker(2, panicPrint, true)
	pool.AddJob(1)
	pool.Close()
}

type logger2 struct {
}

func (l logger2) Print(pool *Pool, value []interface{}, err error) {
}

func errorPrint(ctx context.Context, value ...interface{}) (err error) {
	return fmt.Errorf("Error Print")
}

func TestErrorPool(t *testing.T) {
	pool := NewPool(1, "", logger2{})
	pool.AddWorker(1, errorPrint, true)
	pool.AddJob(1)
	ep := pool.GetErrorPool()
	pool.Close()
	if wj, open := <-ep; open {
		assert := assert.New(t)
		assert.NotNil(wj)
		assert.EqualError(wj.GetError()[0], "Error Print")
	}
}

func TestErrorCount(t *testing.T) {
	pool := NewPool(1, "", logger2{})
	pool.AddWorker(1, errorPrint, true)
	for i := 0; i < 10; i++ {
		pool.AddJob(i)
	}
	pool.Close()
	assert := assert.New(t)
	assert.Equal(10, pool.ErrorCount())
}
