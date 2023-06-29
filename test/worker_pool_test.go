package test

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"log"
	"math/rand"
	"testing"
	"time"

	"github.com/go-yaaf/yaaf-common/utils/pool"
)

var rnd = rand.New(rand.NewSource(time.Now().UnixMilli()))

type parseFileTask struct {
	file string
	num  int
	idx  int
}

func (p *parseFileTask) Run() int {
	// time.Sleep(time.Duration(rnd.Intn(8)) * time.Second)
	time.Sleep(time.Second)
	log.Println(fmt.Sprintf("Parse Finished: %s, (%d)", p.file, p.num))
	return rnd.Intn(1000)
}

func TestWorkPool(t *testing.T) {
	skipCI(t)
	start := time.Now().UnixMilli()

	wp := pool.NewWorkerPool[int](20, 50)
	err := wp.Start(nil)
	require.Nil(t, err, "failed to start worker pool")

	// submit tasks
	for i := 0; i < 30; i++ {
		wp.Submit(&parseFileTask{
			file: fmt.Sprintf("file: %d", i),
			num:  rnd.Intn(4),
			idx:  i,
		})
	}
	wp.Stop()

	duration := time.Now().UnixMilli() - start
	fmt.Println("Done within", duration, "milliseconds")
}

func TestWorkPoolWithResultsCallback(t *testing.T) {
	skipCI(t)
	start := time.Now().UnixMilli()

	cb := func(res int) {
		fmt.Println("Callback invoked:", res)
	}
	wp := pool.NewWorkerPool[int](20, 50)
	err := wp.Start(cb)
	require.Nil(t, err, "failed to start worker pool")

	// submit tasks
	for i := 0; i < 30; i++ {
		wp.Submit(&parseFileTask{
			file: fmt.Sprintf("file: %d", i),
			num:  rnd.Intn(4),
			idx:  i,
		})
	}
	wp.Stop()

	duration := time.Now().UnixMilli() - start
	fmt.Println("Done within", duration, "milliseconds")
}

func TestParallel(t *testing.T) {
	skipCI(t)
	start := time.Now().UnixMilli()

	wp := pool.NewParallel[parseFileTask](10, 50)
	err := wp.Start(parallelProcessor)
	require.Nil(t, err, "failed to start worker pool")

	// submit items
	for i := 0; i < 30; i++ {
		wp.Submit(parseFileTask{
			file: fmt.Sprintf("file: %d", i),
			num:  rnd.Intn(4),
			idx:  i,
		})
	}

	fmt.Println("Waiting for all tasks")
	wp.WaitAll()
	duration := time.Now().UnixMilli() - start
	fmt.Println("Done within", duration, "milliseconds")
}

func parallelProcessor(p parseFileTask) {
	// log.Println("process:", p.file, p.num, p.idx)
	time.Sleep(5 * time.Second)
	log.Println(fmt.Sprintf("%d Parse Finished: %s,", p.idx, p.file))
}
