// Entity tests

package test

import (
	"fmt"
	"github.com/go-yaaf/yaaf-common/entity"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEntityIDs(t *testing.T) {

	sid := entity.ShortID()
	fmt.Printf("%-8s %s\n", "ShortID", sid)
	assert.Equal(t, 6, len(sid), "short id should be 6 digits length")

	lid := entity.ID()
	fmt.Printf("%-8s %s\n", "ID", lid)
	assert.Equal(t, 10, len(lid), "long id should be 10 digits length")

	nid := entity.NanoID()
	fmt.Printf("%-8s %s\n", "NanoID", nid)
	assert.Equal(t, 21, len(nid), "long id should be 21 digits length")

	gid := entity.GUID()
	fmt.Printf("%-8s %s\n", "GUID", gid)
	assert.Equal(t, 36, len(gid), "guid should be 36 digits length")

	fmt.Printf("\n\n")
}
