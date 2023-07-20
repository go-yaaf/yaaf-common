// Entity tests

package test

import (
	"fmt"
	"github.com/go-yaaf/yaaf-common/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestEntityIDs(t *testing.T) {
	skipCI(t)
	sid := entity.ShortID()
	fmt.Printf("%-8s %s\n", "ShortID", sid)
	assert.Equal(t, 6, len(sid), "short id should be 6 digits length")

	sidn := entity.ShortIDN()
	fmt.Printf("%-8s %s\n", "ShortIDN", sidn)
	assert.Equal(t, 10, len(sidn), "short idn should be 6 digits length")

	lid := entity.ID()
	fmt.Printf("%-8s %s\n", "ID", lid)
	assert.Equal(t, 10, len(lid), "long id should be 10 digits length")

	did := entity.IDN()
	fmt.Printf("%-8s %s\n", "NumID", did)
	assert.Equal(t, 16, len(did), "number id should be 16 digits length")

	nid := entity.NanoID()
	fmt.Printf("%-8s %s\n", "NanoID", nid)
	assert.Equal(t, 21, len(nid), "long id should be 21 digits length")

	gid := entity.GUID()
	fmt.Printf("%-8s %s\n", "GUID", gid)
	assert.Equal(t, 36, len(gid), "guid should be 36 digits length")

	fmt.Printf("\n\n")
}

func TestFastJson(t *testing.T) {
	skipCI(t)

	hero := NewHero1("spider-man", 234, "Spider Man")

	bytes, err := entity.Marshal(hero)
	require.Nil(t, err)

	result := string(bytes)
	fmt.Println(result)

	// Test Unmarshal
	expected := NewHero()
	err = entity.Unmarshal(bytes, expected)
	require.Nil(t, err)
	fmt.Println(expected.NAME())

	fmt.Printf("Done \n\n")
}
