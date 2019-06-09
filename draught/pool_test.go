package draught

import (
	"context"
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
		pool.AddJob(i)
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
	panic("Panic Print")
	return nil
}

func TestLogger(t *testing.T) {
	pool := NewPool(1, "", logger{t})
	pool.AddWorker(2, panicPrint, true)
	for i := 0; i < 20; i++ {
		pool.AddJob(i)
	}
	pool.Close()
}
