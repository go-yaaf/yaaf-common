package entity

// This functions replacing the standard encoding/json library with 10x faster json decoder
// Using go get github.com/json-iterator/go instead of the standard encoding/json library

import (
	"encoding/json"
)

//var fastjson = jsoniter.ConfigCompatibleWithStandardLibrary

// BinaryMarshal returns the best wire format encoding of v. If v implements encoding.BinaryMarshaler it will return the binary format, otherwise it will apply JSON encoding
//func BinaryMarshal(v any) ([]byte, error) {
//	if bm, ok := v.(encoding.BinaryMarshaler); ok {
//		return bm.MarshalBinary()
//	} else {
//		return json.Marshal(&v)
//	}
//}

// BinaryUnmarshal will try to unmarshal binary data if v implements encoding.BinaryUnmarshaler, otherwise it will use JSON unmarshal
//func BinaryUnmarshal(data []byte, v any) error {
//	if bm, ok := v.(encoding.BinaryUnmarshaler); ok {
//		return bm.UnmarshalBinary(data)
//	} else {
//		return json.Unmarshal(data, v)
//	}
//}

// Marshal returns the JSON encoding of v
func Marshal(v any) ([]byte, error) {
	return json.Marshal(&v)
}

// Unmarshal returns the struct from JSON byte array
func Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

// UnmarshalFromString returns the struct from JSON string
//func UnmarshalFromString(data string, v any) error {
//	return fastjson.UnmarshalFromString(data, v)
//}
