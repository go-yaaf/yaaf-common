package entity

// This functions replacing the standard encoding/json library with 10x faster json decoder
// Using go get github.com/json-iterator/go instead of the standard encoding/json library

import (
	jsoniter "github.com/json-iterator/go"
)

var fastjson = jsoniter.ConfigCompatibleWithStandardLibrary

// Marshal returns the JSON encoding of v
func Marshal(v any) ([]byte, error) {
	return fastjson.Marshal(&v)
}

// Unmarshal returns the struct from JSON byte array
func Unmarshal(data []byte, v any) error {
	return fastjson.Unmarshal(data, v)
}

// UnmarshalFromString returns the struct from JSON string
func UnmarshalFromString(data string, v any) error {
	return fastjson.UnmarshalFromString(data, v)
}
