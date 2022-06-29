// Copyright 2020. AgentVI Ltd.
//
// Test in memory database implementation
//
package test

import (
	. "github.com/mottyc/yaaf-common/database"
	. "github.com/mottyc/yaaf-common/entity"

	"github.com/stretchr/testify/assert"
	"testing"
)

// region Test Model ---------------------------------------------------------------------------------------------------
type Bird struct {
	BaseEntity
	Key  int    `json:"key"`  // Key
	Name string `json:"name"` // Name
}

func (a Bird) TABLE() string { return "bird" }
func (a Bird) NAME() string  { return a.Name }

func NewBird() Entity {
	return &Bird{}
}

func NewBird1(id string, key int, name string) Entity {
	return &Bird{
		BaseEntity: BaseEntity{Id: id, CreatedOn: Now(), UpdatedOn: Now()},
		Key:        key,
		Name:       name,
	}
}

// endregion

// region Init DB ------------------------------------------------------------------------------------------------------

func getInitializedDb() (dbs IDatabase, err error) {
	dbs, err = NewInMemoryDatabase()
	if err != nil {
		return nil, err
	}

	dbs.Insert(NewBird1("1", 1, "Blackbird"))
	dbs.Insert(NewBird1("2", 2, "Capercaillie"))
	dbs.Insert(NewBird1("3", 3, "Avocet"))
	dbs.Insert(NewBird1("4", 4, "Arctic Tern"))
	dbs.Insert(NewBird1("5", 5, "Bittern"))
	dbs.Insert(NewBird1("6", 6, "Blackcap"))
	dbs.Insert(NewBird1("7", 7, "Buzzard"))
	dbs.Insert(NewBird1("8", 8, "Chough"))
	dbs.Insert(NewBird1("9", 9, "Bullfinch"))
	dbs.Insert(NewBird1("10", 10, "Coot"))
	dbs.Insert(NewBird1("11", 11, "Chaffinch"))
	dbs.Insert(NewBird1("12", 12, "Crane"))
	dbs.Insert(NewBird1("13", 13, "Curlew"))

	return dbs, nil

}

// endregion

func TestInMemoryDatabase_Get(t *testing.T) {

	db, fe := getInitializedDb()
	assert.Nil(t, fe, "error initializing DB")

	bird, fe := db.Get(NewBird, "2")

	assert.Nil(t, fe, "error")
	assert.NotNilf(t, bird, "bird is nil")
	return
}

func TestInMemoryDatabase_List(t *testing.T) {

	db, fe := getInitializedDb()
	assert.Nil(t, fe, "error initializing DB")

	birds, fe := db.List(NewBird, []string{"1", "2", "3", "4"})

	assert.Nil(t, fe, "error")
	assert.NotNilf(t, birds, "birds is nil")
	assert.Equal(t, 4, len(birds), "count should be 4")

	return
}

func TestInMemoryDatabase_Like(t *testing.T) {

	db, fe := getInitializedDb()
	assert.Nil(t, fe, "error initializing DB")

	filter := F("name")
	filter = filter.Like("finch")
	birds, count, fe := db.Query(NewBird).Filter(filter).Find()

	assert.Nil(t, fe, "error")
	assert.NotNilf(t, birds, "birds is nil")
	assert.Equal(t, count, int64(2), "count should be 2")

	return
}
