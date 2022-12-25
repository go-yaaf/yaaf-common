// Copyright 2022. Motty Cohen
//
// Hash utilities

package utils

import (
	"crypto/sha1"
	"encoding/base64"
	"hash/fnv"
	"sync"
)

// region Singleton pattern --------------------------------------------------------------------------------------------
var doOnceForHashUtils sync.Once

type hashUtils struct{}

var hashUtilsSingleton *hashUtils = nil

// HashUtils is a factory method that acts as a static member
func HashUtils() *hashUtils {
	doOnceForHashUtils.Do(func() {
		hashUtilsSingleton = &hashUtils{}
	})
	return hashUtilsSingleton
}

// endregion

// region Hash functions -----------------------------------------------------------------------------------------------

// Hash hashes a string using FNV hash
func (t *hashUtils) Hash(s string) uint32 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(s))
	return h.Sum32()
}

// HashStringToString returns a string containing the base64-encoded SHA-1 hash
// of the input string.
func (t *hashUtils) HashStringToString(s string) string {
	h := sha1.New()
	_, _ = h.Write([]byte(s))
	return base64.RawStdEncoding.EncodeToString(h.Sum(nil))
}

// endregion
