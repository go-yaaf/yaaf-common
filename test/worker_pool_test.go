package test

import (
	"fmt"
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
	wp.Start(nil)

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
	wp.Start(cb)

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
