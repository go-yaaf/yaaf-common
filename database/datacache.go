// Package database
//
// General interface for distributed cache and data structure store (e.g. Redis)
package database

import (
	"context"
	"io"
	"time"

	. "github.com/go-yaaf/yaaf-common/entity"
)

// IDataCache defines the interface for a distributed cache and data structure store (e.g., Redis).
// It supports key-value operations, hashes, lists, and distributed locking.
type IDataCache interface {

	// Closer includes method Close() to close the cache connection.
	io.Closer

	// Ping tests the cache connectivity.
	// It retries the connection 'retries' times with an 'intervalInSeconds' delay between attempts.
	Ping(retries uint, intervalInSeconds uint) error

	// CloneDataCache returns a clone (copy) of the cache instance.
	CloneDataCache() (IDataCache, error)

	// region Key actions ----------------------------------------------------------------------------------------------

	// Get retrieves the value of a key and unmarshals it into an Entity.
	Get(factory EntityFactory, key string) (Entity, error)

	// GetRaw retrieves the raw byte value of a key.
	GetRaw(key string) ([]byte, error)

	// GetKeys retrieves the values of multiple keys and unmarshals them into Entities.
	GetKeys(factory EntityFactory, keys ...string) ([]Entity, error)

	// GetRawKeys retrieves the raw byte values of multiple keys.
	GetRawKeys(keys ...string) ([]Tuple[string, []byte], error)

	// Set sets the value of a key with an optional expiration time.
	Set(key string, entity Entity, expiration ...time.Duration) error

	// SetRaw sets the raw byte value of a key with an optional expiration time.
	SetRaw(key string, bytes []byte, expiration ...time.Duration) error

	// SetNX sets the value of a key only if it does not exist.
	// Returns true if the key was set, false if it already exists.
	SetNX(key string, entity Entity, expiration ...time.Duration) (bool, error)

	// SetRawNX sets the raw byte value of a key only if it does not exist.
	// Returns true if the key was set, false if it already exists.
	SetRawNX(key string, bytes []byte, expiration ...time.Duration) (bool, error)

	// Add sets the value of a key only if it does not exist (alias for SetNX).
	Add(key string, entity Entity, expiration time.Duration) (bool, error)

	// AddRaw sets the raw byte value of a key only if it does not exist (alias for SetRawNX).
	AddRaw(key string, bytes []byte, expiration time.Duration) (bool, error)

	// Del deletes one or more keys.
	Del(keys ...string) (err error)

	// Rename changes the name of a key.
	Rename(key string, newKey string) (err error)

	// Exists checks if a key exists.
	Exists(key string) (result bool, err error)

	// Scan iterates over keys matching a pattern.
	Scan(from uint64, match string, count int64) (keys []string, cursor uint64, err error)

	// endregion

	// region Hash actions ---------------------------------------------------------------------------------------------

	// HGet retrieves the value of a hash field and unmarshals it into an Entity.
	HGet(factory EntityFactory, key, field string) (result Entity, err error)

	// HGetRaw retrieves the raw byte value of a hash field.
	HGetRaw(key, field string) ([]byte, error)

	// HKeys retrieves all field names in a hash.
	HKeys(key string) (fields []string, err error)

	// HGetAll retrieves all fields and values in a hash.
	HGetAll(factory EntityFactory, key string) (result map[string]Entity, err error)

	// HGetRawAll retrieves all fields and raw byte values in a hash.
	HGetRawAll(key string) (result map[string][]byte, err error)

	// HSet sets the value of a hash field.
	HSet(key, field string, entity Entity) (err error)

	// HSetRaw sets the raw byte value of a hash field.
	HSetRaw(key, field string, bytes []byte) (err error)

	// HSetNX sets the value of a hash field only if it does not exist.
	HSetNX(key, field string, entity Entity) (result bool, err error)

	// HSetRawNX sets the raw byte value of a hash field only if it does not exist.
	HSetRawNX(key, field string, bytes []byte) (result bool, err error)

	// HDel deletes one or more hash fields.
	HDel(key string, fields ...string) (err error)

	// HAdd sets the value of a hash field only if it does not exist (alias for HSetNX).
	HAdd(key, field string, entity Entity) (result bool, err error)

	// HAddRaw sets the raw byte value of a hash field only if it does not exist (alias for HSetRawNX).
	HAddRaw(key, field string, bytes []byte) (result bool, err error)

	// HExists checks if a hash field exists.
	HExists(key, field string) (result bool, err error)

	// endregion

	// region List actions ---------------------------------------------------------------------------------------------

	// RPush appends one or more values to the end of a list.
	RPush(key string, value ...Entity) (err error)

	// LPush prepends one or more values to the beginning of a list.
	LPush(key string, value ...Entity) (err error)

	// RPop removes and returns the last element of a list.
	RPop(factory EntityFactory, key string) (entity Entity, err error)

	// LPop removes and returns the first element of a list.
	LPop(factory EntityFactory, key string) (entity Entity, err error)

	// BRPop removes and returns the last element of a list, blocking if the list is empty.
	BRPop(factory EntityFactory, timeout time.Duration, keys ...string) (key string, entity Entity, err error)

	// BLPop removes and returns the first element of a list, blocking if the list is empty.
	BLPop(factory EntityFactory, timeout time.Duration, keys ...string) (key string, entity Entity, err error)

	// LRange retrieves a range of elements from a list.
	LRange(factory EntityFactory, key string, start, stop int64) (result []Entity, err error)

	// LLen returns the length of a list.
	LLen(key string) (result int64)

	// endregion

	// region Locker actions -------------------------------------------------------------------------------------------

	// ObtainLocker tries to obtain a distributed lock with the given key and TTL.
	ObtainLocker(key string, ttl time.Duration) (ILocker, error)

	// endregion
}

// ILocker defines the interface for a distributed lock.
type ILocker interface {
	// Key returns the key of the lock.
	Key() string

	// Token returns the unique token associated with the lock.
	Token() string

	// TTL returns the remaining time-to-live of the lock.
	// Returns 0 if the lock has expired.
	TTL(ctx context.Context) (time.Duration, error)

	// Refresh extends the lock's TTL.
	Refresh(ctx context.Context, ttl time.Duration) error

	// Release releases the lock.
	Release(ctx context.Context) error
}
