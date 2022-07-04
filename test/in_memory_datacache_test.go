// Copyright 2022. Motty Cohen
//
// Test in memory datastore implementation tests
//
package test

import (
	. "github.com/mottyc/yaaf-common/database"
	"github.com/stretchr/testify/assert"
	"testing"
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

	dc, fe := getInitializedCache()
	assert.Nil(t, fe, "error initializing DataCache")

	hero, fe := dc.Get(NewHero, "1")

	assert.Nil(t, fe, "error")
	assert.NotNilf(t, hero, "hero is nil")
	return
}
