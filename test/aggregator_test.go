package test

import (
	"fmt"
	"github.com/go-yaaf/yaaf-common/entity"
	aggregator "github.com/go-yaaf/yaaf-common/utils/Aggregator"
	"golang.org/x/exp/rand"
	"testing"
	"time"
)

func TestAggregator(t *testing.T) {
	skipCI(t)

	agg := aggregator.NewAggregator[entity.Entity](10, 10*time.Second)

	bc := func(bulk []entity.Entity) {
		fmt.Println("---------------- BULK ----------------")
		for _, h := range bulk {
			fmt.Println(h.KEY(), h.ID())
		}
	}

	tc := func(bulk []entity.Entity) {
		fmt.Println("----- TIMEOUT ---------------------------")
		for _, h := range bulk {
			fmt.Println(h.KEY(), h.ID())
		}
	}

	agg.SetBulkCallback(bc)
	agg.SetTimeoutCallback(tc)

	rand.Seed(uint64(time.Now().UnixNano()))

	for i := 0; i < 100; i++ {

		h := NewHero1(fmt.Sprintf("hero-%d", i), i, "Hero")
		agg.Add(h)

		sec := rand.Int31n(11)
		time.Sleep(time.Duration(sec) * time.Second)
	}
}
