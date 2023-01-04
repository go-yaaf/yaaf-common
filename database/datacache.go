// Copyright (c) 2022. Motty Cohen

// Package database
//
// General interface for distributed cache and data structure store (e.g. Redis)
package database

import (
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

	// region Key actions ----------------------------------------------------------------------------------------------

	// Get the value of a key
	Get(factory EntityFactory, key string) (result Entity, err error)

	// GetKeys Get the value of all the given keys
	GetKeys(factory EntityFactory, keys ...string) (results []Entity, err error)

	// Set value of key with optional expiration
	Set(key string, entity Entity, expiration ...time.Duration) (err error)

	// Add Set the value of a key only if the key does not exist
	Add(key string, entity Entity, expiration time.Duration) (result bool, err error)

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

	// HKeys gets all the fields in a hash
	HKeys(key string) (fields []string, err error)

	// HGetAll gets all the fields and values in a hash
	HGetAll(factory EntityFactory, key string) (result map[string]Entity, err error)

	// HSet sets the value of a hash field
	HSet(key, field string, entity Entity) (err error)

	// HDel delete one or more hash fields
	HDel(key string, fields ...string) (err error)

	// HAdd sets the value of a key only if the key does not exist
	HAdd(key, field string, entity Entity) (result bool, err error)

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

}
