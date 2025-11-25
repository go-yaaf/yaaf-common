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

// Marshal returns the JSON encoding of v.
// It is a wrapper around json.Marshal.
func Marshal(v any) ([]byte, error) {
	return json.Marshal(&v)
}

// Unmarshal parses the JSON-encoded data and stores the result in the value pointed to by v.
// It is a wrapper around json.Unmarshal.
func Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

// JsonMarshal converts any Go value to a generic Json map (map[string]any).
// It first marshals the value to JSON bytes, then unmarshals it back into a map.
// This is useful for converting structs to maps.
func JsonMarshal(v any) (Json, error) {
	bytes, err := json.Marshal(&v)
	if err != nil {
		return nil, err
	}

	result := Json{}
	if er := Unmarshal(bytes, &result); er != nil {
		return nil, er
	} else {
		return result, nil
	}
}

// JsonUnmarshal converts a generic Json map (map[string]any) to a specific Go value.
// It first marshals the map to JSON bytes, then unmarshals it into the target value.
func JsonUnmarshal(data Json, v any) error {
	bytes, err := json.Marshal(&data)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, v)
}
