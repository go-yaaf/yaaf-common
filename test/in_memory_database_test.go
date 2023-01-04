// Test in memory database implementation
package test

import (
	. "github.com/go-yaaf/yaaf-common/database"
	"github.com/stretchr/testify/assert"
	"testing"
)

// region Init DB ------------------------------------------------------------------------------------------------------

func getInitializedDb() (dbs IDatabase, err error) {
	dbs, err = NewInMemoryDatabase()
	if err != nil {
		return nil, err
	}

	_, err = dbs.BulkInsert(list_of_heroes)
	return dbs, err
}

// endregion

func TestInMemoryDatabase_Get(t *testing.T) {

	db, fe := getInitializedDb()
	assert.Nil(t, fe, "error initializing DB")

	hero, fe := db.Get(NewHero, "2", "")

	assert.Nil(t, fe, "error")
	assert.NotNilf(t, hero, "hero is nil")
	return
}

func TestInMemoryDatabase_List(t *testing.T) {

	db, fe := getInitializedDb()
	assert.Nil(t, fe, "error initializing DB")

	heroes, fe := db.List(NewHero, []string{"1", "2", "3", "4"}, "")

	assert.Nil(t, fe, "error")
	assert.NotNilf(t, heroes, "heroes is nil")
	assert.Equal(t, 4, len(heroes), "count should be 4")

	return
}

func TestInMemoryDatabase_Like_Suffix(t *testing.T) {

	db, fe := getInitializedDb()
	assert.Nil(t, fe, "error initializing DB")

	filter := Filter("name")
	filter = filter.Like("Black*")
	heroes, count, fe := db.Query(NewHero).Filter(filter).Find()

	assert.Nil(t, fe, "error")
	assert.NotNilf(t, heroes, "heroes is nil")
	assert.Equal(t, count, int64(2), "count should be 2")

	return
}

func TestInMemoryDatabase_Like_Prefix(t *testing.T) {

	db, fe := getInitializedDb()
	assert.Nil(t, fe, "error initializing DB")

	filter := Filter("name")
	filter = filter.Like("*man")
	heroes, count, fe := db.Query(NewHero).Filter(filter).Find()

	assert.Nil(t, fe, "error")
	assert.NotNilf(t, heroes, "heroes is nil")
	assert.Equal(t, count, int64(6), "count should be 6")

	return
}
