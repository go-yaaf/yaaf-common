package utils

import (
	"crypto/sha1"
	"encoding/base64"
	"hash/fnv"
	"sync"
)

// region Singleton pattern --------------------------------------------------------------------------------------------

// HashUtilsStruct is a struct that provides hashing utility functions.
// It is used as a singleton to offer a centralized and efficient way to perform hashing operations.
type HashUtilsStruct struct{}

var (
	doOnceForHashUtils sync.Once
	hashUtilsSingleton *HashUtilsStruct
)

// HashUtils returns a singleton instance of HashUtilsStruct.
// This factory method ensures that the HashUtilsStruct is instantiated only once,
// providing a single point of access to the hashing utilities.
func HashUtils() *HashUtilsStruct {
	doOnceForHashUtils.Do(func() {
		hashUtilsSingleton = &HashUtilsStruct{}
	})
	return hashUtilsSingleton
}

// endregion

// region Hash functions -----------------------------------------------------------------------------------------------

// Hash computes a 32-bit FNV-1a hash of a given string.
// FNV (Fowler-Noll-Vo) is a non-cryptographic hash function known for its speed and low collision rate.
//
// Parameters:
//
//	s: The input string to hash.
//
// Returns:
//
//	A uint32 representing the FNV-1a hash of the string.
func (t *HashUtilsStruct) Hash(s string) uint32 {
	h := fnv.New32a()
	_, _ = h.Write([]byte(s))
	return h.Sum32()
}

// HashStringToString computes the SHA-1 hash of a string and returns it as a base64-encoded string.
// SHA-1 is a cryptographic hash function, though it is no longer considered secure for cryptographic purposes.
// It is still suitable for non-security-related applications like generating unique identifiers.
//
// Parameters:
//
//	s: The input string to hash.
//
// Returns:
//
//	A string containing the raw, standard base64-encoded SHA-1 hash.
func (t *HashUtilsStruct) HashStringToString(s string) string {
	h := sha1.New()
	_, _ = h.Write([]byte(s))
	return base64.RawStdEncoding.EncodeToString(h.Sum(nil))
}

// endregion
