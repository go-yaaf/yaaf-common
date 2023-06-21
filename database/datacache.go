// Copyright (c) 2022. Motty Cohen

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

// IDataCache DataCache interface
type IDataCache interface {

	// Closer includes method Close()
	io.Closer

	// Ping tests connectivity for retries number of time with time interval (in seconds) between retries
	Ping(retries uint, intervalInSeconds uint) error

	// Clone Returns a clone (copy) of the instance
	Clone() (IDataCache, error)

	// region Key actions ----------------------------------------------------------------------------------------------

	// Get the value of a key
	Get(factory EntityFactory, key string) (Entity, error)

	// GetRaw gets the raw value of a key
	GetRaw(key string) ([]byte, error)

	// GetKeys Get the value of all the given keys
	GetKeys(factory EntityFactory, keys ...string) ([]Entity, error)

	// GetRawKeys gets the raw value of all the given keys
	GetRawKeys(keys ...string) ([]Tuple[string, []byte], error)

	// Set value of key with optional expiration
	Set(key string, entity Entity, expiration ...time.Duration) error

	// SetRaw sets the raw value of key with optional expiration
	SetRaw(key string, bytes []byte, expiration ...time.Duration) error

	// SetNX sets value of key only if it is not exist with optional expiration, return false if the key exists
	SetNX(key string, entity Entity, expiration ...time.Duration) (bool, error)

	// SetRawNX sets the raw value of key only if it is not exist with optional expiration, return false if the key exists
	SetRawNX(key string, bytes []byte, expiration ...time.Duration) (bool, error)

	// Add Set the value of a key only if the key does not exist
	Add(key string, entity Entity, expiration time.Duration) (bool, error)

	// AddRaw sets the raw value of a key only if the key does not exist
	AddRaw(key string, bytes []byte, expiration time.Duration) (bool, error)

	// Del Delete keys
	Del(keys ...string) (err error)

	// Rename a key
	Rename(key string, newKey string) (err error)

	// Exists Check if key exists
	Exists(key string) (result bool, err error)

	// Scan keys from the provided cursor
	Scan(from uint64, match string, count int64) (keys []string, cursor uint64, err error)

	// endregion

	// region Hash actions ---------------------------------------------------------------------------------------------

	// HGet gets the value of a hash field
	HGet(factory EntityFactory, key, field string) (result Entity, err error)

	// HGetRaw gets the raw value of a hash field
	HGetRaw(key, field string) ([]byte, error)

	// HKeys gets all the fields in a hash
	HKeys(key string) (fields []string, err error)

	// HGetAll gets all the fields and values in a hash
	HGetAll(factory EntityFactory, key string) (result map[string]Entity, err error)

	// HGetRawAll gets all the fields and raw values in a hash
	HGetRawAll(key string) (result map[string][]byte, err error)

	// HSet sets the value of a hash field
	HSet(key, field string, entity Entity) (err error)

	// HSetRaw sets the raw value of a hash field
	HSetRaw(key, field string, bytes []byte) (err error)

	// HSetNX Set value of key only if it is not exist with optional expiration, return false if the key exists
	HSetNX(key, field string, entity Entity) (result bool, err error)

	// HSetRawNX sets the raw value of key only if it is not exist with optional expiration, return false if the key exists
	HSetRawNX(key, field string, bytes []byte) (result bool, err error)

	// HDel delete one or more hash fields
	HDel(key string, fields ...string) (err error)

	// HAdd sets the value of a key only if the key does not exist
	HAdd(key, field string, entity Entity) (result bool, err error)

	// HAddRaw sets the raw value of a key only if the key does not exist
	HAddRaw(key, field string, bytes []byte) (result bool, err error)

	// HExists Check if key exists
	HExists(key, field string) (result bool, err error)

	// endregion

	// region List actions ---------------------------------------------------------------------------------------------

	// RPush Append one or multiple values to a list
	RPush(key string, value ...Entity) (err error)

	// LPush Prepend one or multiple values to a list
	LPush(key string, value ...Entity) (err error)

	// RPop Remove and get the last element in a list
	RPop(factory EntityFactory, key string) (entity Entity, err error)

	// LPop Remove and get the first element in a list
	LPop(factory EntityFactory, key string) (entity Entity, err error)

	// BRPop Remove and get the last element in a list or block until one is available
	BRPop(factory EntityFactory, timeout time.Duration, keys ...string) (key string, entity Entity, err error)

	// BLPop Remove and get the first element in a list or block until one is available
	BLPop(factory EntityFactory, timeout time.Duration, keys ...string) (key string, entity Entity, err error)

	// LRange Get a range of elements from list
	LRange(factory EntityFactory, key string, start, stop int64) (result []Entity, err error)

	// LLen Get the length of a list
	LLen(key string) (result int64)

	// endregion

	// region List actions ---------------------------------------------------------------------------------------------

	// ObtainLocker tries to obtain a new lock using a key with the given TTL
	ObtainLocker(key string, ttl time.Duration) (ILocker, error)

	// endregion
}

// ILocker represents distributed lock
type ILocker interface {
	// Key returns the locker key
	Key() string

	// Token returns the token value set by the lock.
	Token() string

	// TTL returns the remaining time-to-live. Returns 0 if the lock has expired.
	TTL(ctx context.Context) (time.Duration, error)

	// Refresh extends the lock with a new TTL.
	Refresh(ctx context.Context, ttl time.Duration) error

	// Release manually releases the lock.
	Release(ctx context.Context) error
}
