// Test in memory database implementation
package test

/*
import (
	"fmt"
	. "github.com/go-yaaf/yaaf-common/database"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
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
	skipCI(t)
	db, fe := getInitializedDb()
	assert.Nil(t, fe, "error initializing DB")

	hero, fe := db.Get(NewHero, "2", "")

	assert.Nil(t, fe, "error")
	assert.NotNilf(t, hero, "hero is nil")
	return
}

func TestInMemoryDatabase_List(t *testing.T) {
	skipCI(t)
	db, fe := getInitializedDb()
	assert.Nil(t, fe, "error initializing DB")

	heroes, fe := db.List(NewHero, []string{"1", "2", "3", "4"}, "")

	assert.Nil(t, fe, "error")
	assert.NotNilf(t, heroes, "heroes is nil")
	assert.Equal(t, 4, len(heroes), "count should be 4")

	return
}

func TestInMemoryDatabase_Like_Suffix(t *testing.T) {
	skipCI(t)
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
	skipCI(t)
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

func TestInMemoryDatabase_IsEmpty(t *testing.T) {
	skipCI(t)
	db, fe := getInitializedDb()
	assert.Nil(t, fe, "error initializing DB")

	// Add hero with null and one with empty string
	h1 := NewHero()
	h1.(*Hero).Id = "40"
	h1.(*Hero).Key = 40
	_, _ = db.Insert(h1)
	_, _ = db.Insert(NewHero1("41", 41, ""))

	heroes, total, fe := db.Query(NewHero).Filter(F("name").IsEmpty()).Find()

	assert.Equal(t, 2, total, "total should be 2")
	assert.Nil(t, fe, "error")
	assert.NotNilf(t, heroes, "heroes is nil")
	assert.Equal(t, 2, len(heroes), "count should be 4")

	return
}

func TestInMemoryDatabase_Backup(t *testing.T) {
	skipCI(t)
	db, fe := getInitializedDb()
	assert.Nil(t, fe, "error initializing DB")

	inMemDb, ok := db.(*InMemoryDatabase)
	if !ok {
		fmt.Printf("Not in-memory database\n")
	}

	path, err := os.Getwd()
	assert.Nil(t, err, "error getting working directory")

	path = filepath.Join(path, "backup.json")
	//path = filepath.Join(path, "backup.bin")

	err = inMemDb.Backup(path)
	assert.Nil(t, err, "error backup directory")

	return
}

func TestInMemoryDatabase_Restore(t *testing.T) {
	skipCI(t)

	db, fe := NewInMemoryDatabase()
	assert.Nil(t, fe, "error initializing DB")

	inMemDb, ok := db.(*InMemoryDatabase)
	if !ok {
		fmt.Printf("Not in-memory database\n")
	}

	path, err := os.Getwd()
	assert.Nil(t, err, "error getting working directory")

	path = filepath.Join(path, "backup.json")
	err = inMemDb.Restore(path)
	assert.Nil(t, err, "error backup directory")

	return
}
*/
