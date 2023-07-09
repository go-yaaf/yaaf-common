# Binary marshaling utility

This utility is a set of helper functions to implement `MarshalBinary()` and `UnmarshalBinary()` methods in an easy and convenient way.
Users sometime would like to implement their own manual binary serialization without using a reflection.
This enables better performance (no reflection required) and smaller object size (no headers or schema is stored, only the pure data).
This implementation (like protobuf) is using the actual field value to determine the size of storage int8, int16, int32, int64) to store the data.
But, unlike protobuf, this implementation is using the byte size (int8) as the minimal storage unit. it will not save more than on field in a single byte
to avoid BigEndian/LittleEndian encoding issues.

The current implementation does not support versioning so the user **must keep the same order** in the `MarshalBinary()` and `UnmarshalBinary()` methods

## Marshal and Unmarshal Example

```go
type SampleObject struct {
	Timestamp   entity.Timestamp
	IntValue    int
	Int32Value  int32
	Int64Value  int64
	IntArray    []int
	StringValue string
	StringArray []string
}

// MarshalBinary convert current structure to a minimal wire-format byte array
func (s *SampleObject) MarshalBinary() (data []byte, err error) {
    w := binary.NewWriter()
    w.Timestamp(s.Timestamp).Int(s.IntValue).Int32(s.Int32Value).Int64(s.Int64Value).IntArray(s.IntArray).String(s.StringValue).StringArray(s.StringArray)
    return w.GetBytes(), nil
}

// UnmarshalBinary reads a wire-format byte array to fill the current structure
func (s *SampleObject) UnmarshalBinary(data []byte) (e error) {
    r := binary.NewReader(data)
    if s.Timestamp, e = r.Timestamp(); e != nil {
        return e
    }
    if s.IntValue, e = r.Int(); e != nil {
        return e
    }
    if s.Int32Value, e = r.Int32(); e != nil {
        return e
    }
    if s.Int64Value, e = r.Int64(); e != nil {
        return e
    }
    if s.IntArray, e = r.IntArray(); e != nil {
        return e
    }
    if s.StringValue, e = r.String(); e != nil {
        return e
    }
    if s.StringArray, e = r.StringArray(); e != nil {
        return e
    }
    return nil
}

```