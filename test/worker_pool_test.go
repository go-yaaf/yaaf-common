// Base configuration utility tests

package test

import (
	"fmt"
	"github.com/go-yaaf/yaaf-common/utils/pool"
	"math/rand"
	"testing"
	"time"
)

type TestTask struct {
	data  int
	delay int
	index int
}

func (r *TestTask) Run() any {
	time.Sleep(time.Duration(r.delay) * time.Second)
	fmt.Println(fmt.Sprintf("%d: %d", r.index, r.data))
	return fmt.Sprintf("%d: %d", r.index, r.data)
}

func newTestTask(index, data, delay int) pool.Task {
	return &TestTask{
		index: index,
		data:  data,
		delay: delay,
	}
}

func TestWorkerPool_Execute(t *testing.T) {
	skipCI(t)

	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	results := make(chan any)

	// Create worker pool with 4 worker threads and queue size of 100000
	workerPool := pool.NewWorkerPool(4, 100000, results)

	for i := 1; i < 21; i++ {
		data := random.Intn(100000)
		delay := 1 + random.Intn(5)

		task := newTestTask(i, data, delay)
		if err := workerPool.Execute(task); err != nil {
			fmt.Println(i, "execute error", err.Error())
		} else {
			fmt.Println(i, "execute", data)
		}
	}

	//// Wait for all tasks in the pool
	//for res := range results {
	//	fmt.Println("result:", res)
	//}
	//workerPool.WaitForAll()

	workerPool.Close()
}
