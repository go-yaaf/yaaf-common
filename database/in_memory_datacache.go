package database

import (
	"container/list"
	"fmt"
	"regexp"
	"sync"
	"time"

	. "github.com/go-yaaf/yaaf-common/entity"
	"github.com/go-yaaf/yaaf-common/logger"
	"github.com/go-yaaf/yaaf-common/utils"
	"github.com/go-yaaf/yaaf-common/utils/cache"
	"github.com/go-yaaf/yaaf-common/utils/collections"
)

// region Database store definitions -----------------------------------------------------------------------------------

// InMemoryDataCache represent in memory data cache
type InMemoryDataCache struct {
	keys   *cache.Cache
	lists  map[string]*list.List
	queues map[string]collections.Queue

	mu sync.RWMutex
}

// endregion

// region Factory and connectivity methods for Database ----------------------------------------------------------------

// NewInMemoryDataCache is a factory method for DB store
func NewInMemoryDataCache() (dc IDataCache, err error) {
	return &InMemoryDataCache{
		keys:   cache.NewTtlCache(),
		lists:  make(map[string]*list.List),
		queues: make(map[string]collections.Queue),
	}, nil
}

// Ping tests connectivity for retries number of time with time interval (in seconds) between retries
func (dc *InMemoryDataCache) Ping(retries uint, interval uint) error {
	return nil
}

// Close cache and free resources
func (dc *InMemoryDataCache) Close() error {
	logger.Debug("In memory data-cache closed")
	return nil
}

//endregion

// region Key actions ----------------------------------------------------------------------------------------------

// Get the value of a key
func (dc *InMemoryDataCache) Get(factory EntityFactory, key string) (result Entity, err error) {

	defer utils.RecoverAll(func(er any) {
		if er != nil {
			err = fmt.Errorf("%v", er)
		}
	})

	if value, ok := dc.keys.Get(key); ok {
		return value.(Entity), nil
	} else {
		return nil, fmt.Errorf("key %s not found", key)
	}
}

// GetKeys Get the value of all the given keys
func (dc *InMemoryDataCache) GetKeys(factory EntityFactory, keys ...string) (results []Entity, err error) {
	results = make([]Entity, 0)
	for _, key := range keys {
		if entity, fe := dc.Get(factory, key); fe == nil {
			results = append(results, entity)
		}
	}
	return
}

// Set value of key with optional expiration
func (dc *InMemoryDataCache) Set(key string, entity Entity, expiration ...time.Duration) (err error) {
	dc.keys.Set(key, entity)
	return nil
}

// Add Set the value of a key only if the key does not exist
func (dc *InMemoryDataCache) Add(key string, entity Entity, expiration time.Duration) (result bool, err error) {
	if _, fe := dc.Get(nil, key); fe != nil {
		_ = dc.Set(key, entity, expiration)
		return true, nil
	} else {
		return false, nil
	}
}

// Del Delete keys
func (dc *InMemoryDataCache) Del(keys ...string) (err error) {
	for _, key := range keys {
		dc.keys.Delete(key)
	}
	return nil
}

// Rename a key
func (dc *InMemoryDataCache) Rename(key string, newKey string) (err error) {
	exists, fe := dc.Exists(newKey)
	if fe != nil {
		return fe
	}
	if exists {
		return fmt.Errorf("key %s already exists", newKey)
	}

	if entity, fe := dc.Get(nil, key); fe != nil {
		return fe
	} else {
		_ = dc.Set(newKey, entity)
		_ = dc.Del(key)
		return nil
	}
}

// Exists checks if key exists
func (dc *InMemoryDataCache) Exists(key string) (result bool, err error) {
	_, exists := dc.keys.Get(key)
	return exists, nil
}

// Scan keys from the provided cursor
func (dc *InMemoryDataCache) Scan(from uint64, match string, count int64) (keys []string, cursor uint64, err error) {
	rex, _ := regexp.Compile(match)

	keys = make([]string, 0)
	cb := func(k, v any) bool {
		if rex != nil {
			if rex.MatchString(fmt.Sprintf("%v", k)) {
				keys = append(keys, fmt.Sprintf("%v", k))
			}
		} else {
			keys = append(keys, fmt.Sprintf("%v", k))
		}
		return true
	}
	dc.keys.Range(cb)
	return
}

// endregion

// region Hash actions ---------------------------------------------------------------------------------------------

// HGet gets the value of a hash field
func (dc *InMemoryDataCache) HGet(factory EntityFactory, key, field string) (result Entity, err error) {
	hKey := fmt.Sprintf("%s@%s", key, field)
	return dc.Get(factory, hKey)
}

// HKeys get all the fields in a hash
func (dc *InMemoryDataCache) HKeys(key string) (fields []string, err error) {
	keys, _, fe := dc.Scan(0, key, 0)
	return keys, fe
}

// HGetAll gets all the fields and values in a hash
func (dc *InMemoryDataCache) HGetAll(factory EntityFactory, key string) (result map[string]Entity, err error) {
	result = make(map[string]Entity)
	keys, _, err := dc.Scan(0, key, 0)
	for _, k := range keys {
		if entity, fe := dc.Get(factory, k); fe == nil {
			result[k] = entity
		}
	}
	return
}

// HSet sets the value of a hash field
func (dc *InMemoryDataCache) HSet(key, field string, entity Entity) (err error) {
	hKey := fmt.Sprintf("%s@%s", key, field)
	return dc.Set(hKey, entity)
}

// HDel delete one or more hash fields
func (dc *InMemoryDataCache) HDel(key string, fields ...string) (err error) {
	keys := make([]string, 0)
	for _, field := range fields {
		hKey := fmt.Sprintf("%s@%s", key, field)
		keys = append(keys, hKey)
	}

	return dc.Del(keys...)
}

// HAdd sets the value of a key only if the key does not exist
func (dc *InMemoryDataCache) HAdd(key, field string, entity Entity) (result bool, err error) {
	hKey := fmt.Sprintf("%s@%s", key, field)
	return dc.Add(hKey, entity, 0)
}

// HExists Check if key exists
func (dc *InMemoryDataCache) HExists(key, field string) (result bool, err error) {
	hKey := fmt.Sprintf("%s@%s", key, field)
	return dc.Exists(hKey)
}

// endregion

// region List actions ---------------------------------------------------------------------------------------------

// RPush append (add to the right) one or multiple values to a list
func (dc *InMemoryDataCache) RPush(key string, value ...Entity) (err error) {
	// Ensure list exists
	lst, ok := dc.lists[key]
	if !ok {
		lst = list.New()
		lst.Init()
		dc.lists[key] = lst
	}

	for _, val := range value {
		lst.PushBack(val)
	}
	return nil
}

// LPush Prepend (add to the left) one or multiple values to a list
func (dc *InMemoryDataCache) LPush(key string, value ...Entity) (err error) {
	// Ensure list exists
	lst, ok := dc.lists[key]
	if !ok {
		lst = list.New()
		lst.Init()
		dc.lists[key] = lst
	}

	for _, val := range value {
		lst.PushFront(val)
	}
	return nil
}

// RPop Remove and get the last element in a list
func (dc *InMemoryDataCache) RPop(factory EntityFactory, key string) (entity Entity, err error) {
	// Ensure list exists
	if lst, ok := dc.lists[key]; !ok {
		return nil, fmt.Errorf("list %s not exists", key)
	} else {
		if e := lst.Back(); e != nil {
			return nil, fmt.Errorf("end of list")
		} else {
			entity = e.Value.(Entity)
			lst.Remove(e)
		}
	}
	return entity, nil
}

// LPop Remove and get the first element in a list
func (dc *InMemoryDataCache) LPop(factory EntityFactory, key string) (entity Entity, err error) {
	// Ensure list exists
	if lst, ok := dc.lists[key]; !ok {
		return nil, fmt.Errorf("list %s not exists", key)
	} else {
		if e := lst.Front(); e != nil {
			return nil, fmt.Errorf("end of list")
		} else {
			entity = e.Value.(Entity)
			lst.Remove(e)
		}
	}
	return entity, nil
}

// BRPop Remove and get the last element in a list or block until one is available
func (dc *InMemoryDataCache) BRPop(factory EntityFactory, timeout time.Duration, keys ...string) (key string, value Entity, err error) {
	select {
	case _ = <-time.Tick(time.Millisecond):
		if k, v, exists := dc.brPop(keys...); exists {
			return k, v, nil
		}
	case <-time.After(timeout):
		return "", nil, fmt.Errorf("timeout")
	}

	return "", nil, fmt.Errorf("not exists")
}

// Internal implementation of brPop
func (dc *InMemoryDataCache) brPop(keys ...string) (key string, entity Entity, exists bool) {
	for _, k := range keys {
		if v, _ := dc.RPop(nil, k); v != nil {
			return k, v, true
		}
	}
	return "", nil, false
}

// BLPop Remove and get the first element in a list or block until one is available
func (dc *InMemoryDataCache) BLPop(factory EntityFactory, timeout time.Duration, keys ...string) (key string, entity Entity, err error) {
	select {
	case _ = <-time.Tick(time.Millisecond):
		if k, v, exists := dc.blPop(keys...); exists {
			return k, v, nil
		}
	case <-time.After(timeout):
		return "", nil, fmt.Errorf("timeout")
	}

	return "", nil, fmt.Errorf("not exists")
}

// Internal implementation of blPop
func (dc *InMemoryDataCache) blPop(keys ...string) (key string, entity Entity, exists bool) {
	for _, k := range keys {
		if v, _ := dc.LPop(nil, k); v != nil {
			return k, v, true
		}
	}
	return "", nil, false
}

// LRange Get a range of elements from list
func (dc *InMemoryDataCache) LRange(factory EntityFactory, key string, start, stop int64) (result []Entity, err error) {
	result = make([]Entity, 0)

	index := int64(-1)

	// Ensure list exists
	if lst, ok := dc.lists[key]; !ok {
		return nil, fmt.Errorf("key %s not found", key)
	} else {
		for e := lst.Front(); e != nil; e = e.Next() {
			index += 1
			if index < start {
				continue
			}
			if stop > 0 && index > stop {
				continue
			}
			result = append(result, e.Value.(Entity))
		}
		return
	}
}

// LLen Get the length of a list
func (dc *InMemoryDataCache) LLen(key string) (result int64) {
	// Ensure list exists
	if lst, ok := dc.lists[key]; !ok {
		return 0
	} else {
		return int64(lst.Len())
	}
}

// endregion
