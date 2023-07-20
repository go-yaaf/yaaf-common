// Test in memory datastore implementation tests
package test

import (
	"fmt"
	. "github.com/go-yaaf/yaaf-common/database"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// region Init DB ------------------------------------------------------------------------------------------------------

func getInitializedCache() (dc IDataCache, err error) {
	dc, err = NewInMemoryDataCache()
	if err != nil {
		return nil, err
	}

	// fill keys
	for _, h := range list_of_heroes {
		_ = dc.Set(h.ID(), h)
	}
	return dc, nil
}

// endregion

func TestInMemoryDataCache_Get(t *testing.T) {
	skipCI(t)

	dc, fe := getInitializedCache()
	assert.Nil(t, fe, "error initializing DataCache")

	hero, fe := dc.Get(NewHero, "1")

	assert.Nil(t, fe, "error")
	assert.NotNilf(t, hero, "hero is nil")
	return
}

func TestInMemoryDataCache_SetWithTTL(t *testing.T) {
	skipCI(t)

	dc, fe := getInitializedCache()
	assert.Nil(t, fe, "error initializing DataCache")

	// Set value with TTL
	hero := NewHero()
	hero.(*Hero).Id = "test"
	hero.(*Hero).Name = "test_hero"
	fe = dc.Set("item_with_ttl", hero, time.Second)
	assert.Nil(t, fe, "error")

	// Ensure key is there
	result, err := dc.Get(NewHero, "item_with_ttl")
	fmt.Println("result", result, "error", err)
	assert.Nil(t, err, "error")
	assert.NotNilf(t, result, "result is nil")

	// Sleep for 2 minutes
	time.Sleep(time.Second * 2)

	result, err = dc.Get(NewHero, "item_with_ttl")
	fmt.Println("result", result, "error", err)

	return
}
