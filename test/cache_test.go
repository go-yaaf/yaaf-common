package test

/*
import (
	"fmt"
	"github.com/go-yaaf/yaaf-common/entity"
	"github.com/go-yaaf/yaaf-common/utils/cache"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestGenericsCacheString(t *testing.T) {
	skipCI(t)

	ttlCache := cache.NewTtlCache[string, Hero]()
	h1 := Hero{
		BaseEntity: entity.BaseEntity{Id: "h1"},
		Key:        1,
		Name:       "Zeus",
	}
	h2 := Hero{
		BaseEntity: entity.BaseEntity{Id: "h2"},
		Key:        2,
		Name:       "Hara",
	}
	ttlCache.Set("1", h1)
	ttlCache.Set("2", h2)

	items := ttlCache.Count()
	require.Equal(t, 2, items)

	c1, ok := ttlCache.Get("1")
	require.True(t, ok)
	fmt.Println("1", c1.Name)

	c1, ok = ttlCache.Get("3")
	require.False(t, ok)

	ok = ttlCache.Remove("1")
	require.True(t, ok)

	ok = ttlCache.Remove("1")
	require.False(t, ok)

	items = ttlCache.Count()
	require.Equal(t, 1, items)
}

func TestGenericsCacheInt(t *testing.T) {
	skipCI(t)

	ttlCache := cache.NewTtlCache[int, Hero]()
	h1 := Hero{
		BaseEntity: entity.BaseEntity{Id: "h1"},
		Key:        1,
		Name:       "Zeus",
	}
	h2 := Hero{
		BaseEntity: entity.BaseEntity{Id: "h2"},
		Key:        2,
		Name:       "Hara",
	}
	ttlCache.Set(1, h1)
	ttlCache.Set(2, h2)

	items := ttlCache.Count()
	require.Equal(t, 2, items)

	c1, ok := ttlCache.Get(1)
	require.True(t, ok)
	fmt.Println("1", c1.Name)

	c1, ok = ttlCache.Get(3)
	require.False(t, ok)

	ok = ttlCache.Remove(1)
	require.True(t, ok)

	ok = ttlCache.Remove(1)
	require.False(t, ok)

	items = ttlCache.Count()
	require.Equal(t, 1, items)
}

func TestGenericsTtlCacheString(t *testing.T) {
	skipCI(t)

	ttlCache := cache.NewTtlCache[string, Hero]()
	h1 := Hero{
		BaseEntity: entity.BaseEntity{Id: "h1"},
		Key:        1,
		Name:       "Zeus",
	}
	h2 := Hero{
		BaseEntity: entity.BaseEntity{Id: "h2"},
		Key:        2,
		Name:       "Hara",
	}

	check := func(key string, value Hero) bool {
		fmt.Println("about to delete", key, "...")
		return true
	}
	ttlCache.SetCheckExpirationCallback(check)

	expirationCb := func(key string, value Hero) {
		fmt.Println(key, "deleted...")
	}
	ttlCache.SetExpirationCallback(expirationCb)

	ttlCache.SetWithTTL("1", h1, 5*time.Second)
	ttlCache.SetWithTTL("2", h2, 10*time.Second)

	count := ttlCache.Count()
	require.Equal(t, 2, count)
	fmt.Println("items:", count)

	// wait 6 seconds, only one item should be in the cache
	fmt.Println("wait 6 seconds ...")
	time.Sleep(6 * time.Second)
	require.Equal(t, 1, ttlCache.Count())
	fmt.Println("items:", ttlCache.Count())

	// wait 6 seconds, 0 item should be in the cache
	fmt.Println("wait 6 seconds ...")
	time.Sleep(6 * time.Second)
	require.Equal(t, 0, ttlCache.Count())
	fmt.Println("items:", ttlCache.Count())

}
*/
