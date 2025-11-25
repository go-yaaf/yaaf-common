package binary

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"unsafe"

	"github.com/go-yaaf/yaaf-common/entity"
)

const (
	byte1Subtractor = (1 << 7)
	byte2Subtractor = (1<<7 + 1<<14)
	byte3Subtractor = (1<<7 + 1<<14 + 1<<21)
	byte4Subtractor = (1<<7 + 1<<14 + 1<<21 + 1<<28)
	byte5Subtractor = (1<<7 + 1<<14 + 1<<21 + 1<<28 + 1<<35)
	byte6Subtractor = (1<<7 + 1<<14 + 1<<21 + 1<<28 + 1<<35 + 1<<42)
	byte7Subtractor = (1<<7 + 1<<14 + 1<<21 + 1<<28 + 1<<35 + 1<<42 + 1<<49)
	byte8Subtractor = (1<<7 + 1<<14 + 1<<21 + 1<<28 + 1<<35 + 1<<42 + 1<<49 + 1<<56)
)

// NewReader will initialize a new instance of writer
func NewReader(data []byte) *Reader {
	rd := bytes.NewReader(data)
	return &Reader{reader: rd}
}

// Reader manages the reading of binary data
type Reader struct {
	reader *bytes.Reader
}

// Uint read unsigned int value
func (r *Reader) Uint() (uint, error) {
	if u64, err := r.Uint64(); err != nil {
		return 0, err
	} else {
		return uint(u64), nil
	}
}

// Uint8 read unsigned int 8 bit value
func (r *Reader) Uint8() (uint8, error) {
	return r.reader.ReadByte()
}

// Uint16 read unsigned int 16 bit value
func (r *Reader) Uint16() (uint16, error) {
	if u64, err := r.Uint64(); err != nil {
		return 0, err
	} else {
		return uint16(u64), nil
	}
}

// Uint32 read unsigned int 32 bit value
func (r *Reader) Uint32() (uint32, error) {
	if u64, err := r.Uint64(); err != nil {
		return 0, err
	} else {
		return uint32(u64), nil
	}
}

// Uint64 read unsigned int 64 bit value
func (r *Reader) Uint64() (v uint64, err error) {
	var b byte
	if b, err = r.reader.ReadByte(); err != nil {
		return
	}

	if v = uint64(b); v < ceiling {
		return
	}

	// Read next byte
	if b, err = r.reader.ReadByte(); err != nil {
		return
	}

	v += uint64(b) << 7

	if b < ceiling {
		v -= byte1Subtractor
		return
	}

	// Read next byte
	if b, err = r.reader.ReadByte(); err != nil {
		return
	}

	v += uint64(b) << 14
	if b < ceiling {
		v -= byte2Subtractor
		return
	}

	// Read next byte
	if b, err = r.reader.ReadByte(); err != nil {
		return
	}

	v += uint64(b) << 21
	if b < ceiling {
		v -= byte3Subtractor
		return
	}

	// Read next byte
	if b, err = r.reader.ReadByte(); err != nil {
		return
	}

	v += uint64(b) << 28
	if b < ceiling {
		v -= byte4Subtractor
		return
	}

	// Read next byte
	if b, err = r.reader.ReadByte(); err != nil {
		return
	}

	v += uint64(b) << 35
	if b < ceiling {
		v -= byte5Subtractor
		return
	}

	// Read next byte
	if b, err = r.reader.ReadByte(); err != nil {
		return
	}

	v += uint64(b) << 42
	if b < ceiling {
		v -= byte6Subtractor
		return
	}

	// Read next byte
	if b, err = r.reader.ReadByte(); err != nil {
		return
	}

	v += uint64(b) << 49
	if b < ceiling {
		v -= byte7Subtractor
		return
	}

	// Read next byte
	if b, err = r.reader.ReadByte(); err != nil {
		return
	}

	v += uint64(b) << 56
	v -= byte8Subtractor
	return
}

// Int read integer value
func (r *Reader) Int() (int, error) {
	if u64, err := r.Int64(); err != nil {
		return 0, err
	} else {
		return int(u64), nil
	}
}

// IntArray read variable length array of int values
func (r *Reader) IntArray() ([]int, error) {
	// Read array sized
	size, err := r.Int()
	if err != nil {
		return nil, err
	}

	result := make([]int, 0)
	for i := 0; i < size; i++ {
		if val, e := r.Int(); e != nil {
			return nil, e
		} else {
			result = append(result, val)
		}
	}

	return result, nil
}

// Int8 read int 8 bit value
func (r *Reader) Int8() (int8, error) {
	if u8, err := r.Uint8(); err != nil {
		return 0, err
	} else {
		v := *(*int8)(unsafe.Pointer(&u8))
		return v, nil
	}
}

// Int16 read int 16 bit value
func (r *Reader) Int16() (int16, error) {
	if i64, err := r.Int64(); err != nil {
		return 0, err
	} else {
		return int16(i64), nil
	}
}

// Int32 read int 32 bit value
func (r *Reader) Int32() (int32, error) {
	if i64, err := r.Int64(); err != nil {
		return 0, err
	} else {
		return int32(i64), nil
	}
}

// Int64 read int 64 bit value
func (r *Reader) Int64() (int64, error) {
	if u64, err := r.Uint64(); err != nil {
		return 0, err
	} else {
		v := *(*int64)(unsafe.Pointer(&u64))
		return v, nil
	}
}

// Float32 read float 32 bit value (single)
func (r *Reader) Float32() (float32, error) {
	if u32, err := r.Uint32(); err != nil {
		return 0, err
	} else {
		return math.Float32frombits(u32), nil
	}
}

// Float64 read float 64 bit value (double)
func (r *Reader) Float64() (float64, error) {
	if u64, err := r.Uint64(); err != nil {
		return 0, err
	} else {
		return math.Float64frombits(u64), nil
	}
}

// Object read an arbitrary byte array representing an object
func (r *Reader) Object() (result []byte, err error) {
	var bsLength int
	if bsLength, err = r.Int(); err != nil {
		err = fmt.Errorf("error decoding bytes length: %v", err)
		return
	}

	expandSlice(&result, bsLength)

	if bsLength == 0 {
		// We do not have any bytes to read, return
		return
	}

	_, err = io.ReadAtLeast(r.reader, result, bsLength)
	return
}

// ObjectArray read variable length array of arbitrary objects
func (r *Reader) ObjectArray() ([][]byte, error) {
	// Read array size
	size, err := r.Int()
	if err != nil {
		return nil, fmt.Errorf("error reading array length: %v", err)
	}

	result := make([][]byte, 0)
	for i := 0; i < size; i++ {
		if data, e := r.Object(); e != nil {
			return nil, e
		} else {
			result = append(result, data)
		}
	}

	return result, nil
}

// String read string value
func (r *Reader) String() (string, error) {
	if bs, err := r.Object(); err != nil {
		return "", err
	} else {
		return getStringFromBytes(bs), nil
	}
}

// StringArray read array of strings
func (r *Reader) StringArray() ([]string, error) {
	// Read array size
	size, err := r.Int()
	if err != nil {
		return nil, err
	}

	result := make([]string, 0)
	for i := 0; i < size; i++ {
		if val, e := r.String(); e != nil {
			return nil, e
		} else {
			result = append(result, val)
		}
	}

	return result, nil
}

// IP read encoded IP value that can be represented as uint32 (IPv4), bigInt (IPv6) or string (dns)
// In order to parse header correctly, we need to read the header first:
// 1: IP represented as string, 4: IP represented as IPv4 int (uint32), 6: IP represented as IPv6 bigInt (2 * uint64)
func (r *Reader) IP() (string, error) {
	// Read IP header
	hdr, err := r.Uint8()
	if err != nil {
		return "", err
	}
	// If header is 4, it is IPv4 represented as Uint32
	if hdr == 4 {
		if val, er := r.Uint32(); er != nil {
			return "", err
		} else {
			ip := IntToIPv4(val)
			return ip.String(), nil
		}
	}
	// If header is 6, it is IPv6 represented as two Uint64 values
	if hdr == 6 {
		high, er := r.Uint64()
		if er != nil {
			return "", err
		}
		low, er := r.Uint64()
		if er != nil {
			return "", err
		}
		ip := IntToIPv6(high, low)
		return ip.String(), nil
	}
	// If header is not 4 or 6, read the IP as string
	return r.String()
}

// IPArray will encode a list of IPv4 or IPv6 to byte array, each IP is stored as defined in the IP() method
func (r *Reader) IPArray() ([]string, error) {
	// Read array size
	size, err := r.Int()
	if err != nil {
		return nil, err
	}

	result := make([]string, 0)
	for i := 0; i < size; i++ {
		if val, e := r.IP(); e != nil {
			return nil, e
		} else {
			result = append(result, val)
		}
	}
	return result, nil
}

// Bool read boolean value
func (r *Reader) Bool() (bool, error) {
	if u8, err := r.Uint8(); err != nil {
		return false, err
	} else {
		return u8 == 1, nil
	}
}

// Timestamp read int 64 bit value and return it as timestamp
func (r *Reader) Timestamp() (entity.Timestamp, error) {
	if u64, err := r.Uint64(); err != nil {
		return 0, err
	} else {
		v := *(*int64)(unsafe.Pointer(&u64))
		return entity.Timestamp(v), nil
	}
}

// Close will close the reader
func (r *Reader) Close() (err error) {
	return
}
