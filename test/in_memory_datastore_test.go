// Copyright 2022. Motty Cohen
//
// Test in memory datastore implementation tests
//
package test

import (
	. "github.com/go-yaaf/yaaf-common/database"
	"github.com/stretchr/testify/assert"
	"testing"
)

// region Init DB ------------------------------------------------------------------------------------------------------

func getInitializedDs() (ds IDatastore, err error) {
	ds, err = NewInMemoryDatastore()
	if err != nil {
		return nil, err
	}

	_, err = ds.BulkInsert(list_of_heroes)
	return ds, err
}

// endregion

func TestInMemoryDatastore_Get(t *testing.T) {

	ds, fe := getInitializedDs()
	assert.Nil(t, fe, "error initializing Datastore")

	hero, fe := ds.Get(NewHero, "2")

	assert.Nil(t, fe, "error")
	assert.NotNilf(t, hero, "hero is nil")
	return
}

func TestInMemoryDatastore_List(t *testing.T) {

	ds, fe := getInitializedDs()
	assert.Nil(t, fe, "error initializing Datastore")

	heroes, fe := ds.List(NewHero, []string{"1", "2", "3", "4"})

	assert.Nil(t, fe, "error")
	assert.NotNilf(t, heroes, "heroes is nil")
	assert.Equal(t, 4, len(heroes), "count should be 4")

	return
}

func TestInMemoryDatastore_Like_Suffix(t *testing.T) {

	ds, fe := getInitializedDs()
	assert.Nil(t, fe, "error initializing Datastore")

	filter := Filter("name")
	filter = filter.Like("Black*")
	heroes, count, fe := ds.Query(NewHero).Filter(filter).Find()

	assert.Nil(t, fe, "error")
	assert.NotNilf(t, heroes, "heroes is nil")
	assert.Equal(t, count, int64(2), "count should be 2")

	return
}

func TestInMemoryDatastore_Like_Prefix(t *testing.T) {

	ds, fe := getInitializedDs()
	assert.Nil(t, fe, "error initializing Datastore")

	filter := Filter("name")
	filter = filter.Like("*man")
	heroes, count, fe := ds.Query(NewHero).Filter(filter).Find()

	assert.Nil(t, fe, "error")
	assert.NotNilf(t, heroes, "heroes is nil")
	assert.Equal(t, count, int64(6), "count should be 6")

	return
}
