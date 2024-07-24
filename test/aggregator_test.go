package test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/go-yaaf/yaaf-common/entity"
	. "github.com/go-yaaf/yaaf-common/utils/Aggregator"
)

func TestAggregator(t *testing.T) {
	skipCI(t)

	agg := NewAggregator[entity.Entity](10, 10*time.Second)

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

	for i := 0; i < 100; i++ {

		h := NewHero1(fmt.Sprintf("hero-%d", i), i, "Hero")
		agg.Add(h)

		sec := rand.Int31n(11)
		time.Sleep(time.Duration(sec) * time.Second)
	}
}
