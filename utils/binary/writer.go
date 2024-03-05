package binary

import (
	"bytes"
	"github.com/go-yaaf/yaaf-common/entity"
	"io"
	"math"
	"unsafe"
)

// NewWriter will initialize a new instance of writer
func NewWriter() *Writer {
	return &Writer{
		buffer:  bytes.NewBuffer(nil),
		bs:      nil,
		written: 0,
	}
}

// Writer manages the writing of the output
type Writer struct {
	buffer  *bytes.Buffer
	bs      []byte
	written int64
}

// Uint will encode unsigned int value
func (w *Writer) Uint(v uint) *Writer {
	return w.Uint64(uint64(v))
}

// Uint8 will encode unsigned int 8 bit value (0 .. 255)
func (w *Writer) Uint8(v uint8) *Writer {
	w.bs = append(w.bs, v)
	w.flush()
	return w
}

// Uint16 will encode unsigned int 16 bit value (0 .. 65,535)
func (w *Writer) Uint16(v uint16) *Writer {
	return w.Uint64(uint64(v))
}

// Uint32 will encode unsigned int 32 bit value (0 .. 4,294,967,295)
func (w *Writer) Uint32(v uint32) *Writer {
	return w.Uint64(uint64(v))
}

// Uint64 will encode unsigned int 64 bits value (0 .. 18,446,744,073,709,551,615)
func (w *Writer) Uint64(v uint64) *Writer {
	w.varInt(v)
	w.flush()
	return w
}

// varInt will put variable length integer in the temp buffer
func (w *Writer) varInt(v uint64) {

	switch {
	case v < 1<<7-1:
		w.bs = append(w.bs, byte(v))
	case v < 1<<14-1:
		w.bs = append(w.bs, byte(v)|0x80, byte(v>>7))
	case v < 1<<21-1:
		w.bs = append(w.bs, byte(v)|0x80, byte(v>>7)|0x80, byte(v>>14))
	case v < 1<<28-1:
		w.bs = append(w.bs, byte(v)|0x80, byte(v>>7)|0x80, byte(v>>14)|0x80, byte(v>>21))
	case v < 1<<35-1:
		w.bs = append(w.bs, byte(v)|0x80, byte(v>>7)|0x80, byte(v>>14)|0x80, byte(v>>21)|0x80, byte(v>>28))
	case v < 1<<42-1:
		w.bs = append(w.bs, byte(v)|0x80, byte(v>>7)|0x80, byte(v>>14)|0x80, byte(v>>21)|0x80, byte(v>>28)|0x80, byte(v>>35))
	case v < 1<<49-1:
		w.bs = append(w.bs, byte(v)|0x80, byte(v>>7)|0x80, byte(v>>14)|0x80, byte(v>>21)|0x80, byte(v>>28)|0x80, byte(v>>35)|0x80, byte(v>>42))
	case v < 1<<56-1:
		w.bs = append(w.bs, byte(v)|0x80, byte(v>>7)|0x80, byte(v>>14)|0x80, byte(v>>21)|0x80, byte(v>>28)|0x80, byte(v>>35)|0x80, byte(v>>42)|0x80, byte(v>>49))
	default:
		w.bs = append(w.bs, byte(v)|0x80, byte(v>>7)|0x80, byte(v>>14)|0x80, byte(v>>21)|0x80, byte(v>>28)|0x80, byte(v>>35)|0x80, byte(v>>42)|0x80, byte(v>>49)|0x80, byte(v>>56))
	}
}

// Int will encode int value
func (w *Writer) Int(v int) *Writer {
	return w.Int64(int64(v))
}

// IntArray will encode variable length array of int values
func (w *Writer) IntArray(v []int) *Writer {
	// Write array sized
	w.varInt(uint64(len(v)))
	for _, val := range v {
		w.varInt(uint64(val))
	}
	w.flush()
	return w
}

// Int8 will encode unsigned int 8 bit value (-128 .. 127)
func (w *Writer) Int8(v int8) *Writer {
	return w.Uint8(*(*uint8)(unsafe.Pointer(&v)))
}

// Int16 will encode int 16 bit value (-32,768 .. 32,767)
func (w *Writer) Int16(v int16) *Writer {
	return w.Int64(int64(v))
}

// Int32 will encode int 32 bit value (-2,147,483,648 .. 2,147,483,647)
func (w *Writer) Int32(v int32) *Writer {
	return w.Int64(int64(v))
}

// Int64 will encode int 64 bit value (-9,223,372,036,854,775,808 .. 9,223,372,036,854,775,807)
func (w *Writer) Int64(v int64) *Writer {
	return w.Uint64(*(*uint64)(unsafe.Pointer(&v)))
}

// Float32 will encode float 32 bit value (single)
func (w *Writer) Float32(v float32) *Writer {
	return w.Uint32(math.Float32bits(v))
}

// Float32Array will encode variable length array of float32 values
func (w *Writer) Float32Array(v []float32) *Writer {
	// Write array sized
	w.varInt(uint64(len(v)))
	for _, val := range v {
		w.varInt(uint64(math.Float32bits(val)))
	}
	w.flush()
	return w
}

// Float64 will encode float 64 bit value (double)
func (w *Writer) Float64(v float64) *Writer {
	return w.Uint64(math.Float64bits(v))
}

// Float64Array will encode variable length array of float64 values
func (w *Writer) Float64Array(v []float64) *Writer {
	// Write array sized
	w.varInt(uint64(len(v)))
	for _, val := range v {
		w.varInt(math.Float64bits(val))
	}
	w.flush()
	return w
}

// Bool will encode a boolean value
func (w *Writer) Bool(v bool) *Writer {
	if v {
		return w.Uint8(1)
	} else {
		return w.Uint8(0)
	}
}

// String will encode a variable length string
func (w *Writer) String(v string) *Writer {
	bsp := getStringBytes(&v)
	return w.Object(bsp)
}

// StringArray will encode variable length array of strings
func (w *Writer) StringArray(v []string) *Writer {
	// Write array sized
	w.varInt(uint64(len(v)))
	for _, val := range v {
		bsp := getStringBytes(&val)
		w.varInt(uint64(len(*bsp)))
		w.bs = append(w.bs, *bsp...)
	}
	w.flush()
	return w
}

// Object will encode an arbitrary object represented as variable length byte array
func (w *Writer) Object(v *[]byte) *Writer {
	w.varInt(uint64(len(*v)))
	w.bs = append(w.bs, *v...)
	w.flush()
	return w
}

// ObjectArray will encode variable length array of arbitrary objects
func (w *Writer) ObjectArray(v *[][]byte) *Writer {
	// Write array size
	w.varInt(uint64(len(*v)))

	// for each item, write the size of the item and then its content
	for _, val := range *v {
		w.varInt(uint64(len(val)))
		w.bs = append(w.bs, val...)
	}
	w.flush()
	return w
}

// Timestamp will encode a timestamp (int64) type
func (w *Writer) Timestamp(v entity.Timestamp) *Writer {
	return w.Int64(int64(v))
}

// IP will encode an IPv4 or IPv6 to byte array, to distinguish between IP types, we need a small uint8 header:
// 1: IP represented as string, 4: IP represented as IPv4 int (uint32), 6: IP represented as IPv6 bigInt (2 * uint64)
func (w *Writer) IP(v string) *Writer {
	if ip, version, err := parseIP(v); err != nil {
		return w.Uint8(1).String(v)
	} else {
		if version == 4 {
			if ipv4, er := IPv4ToInt(ip); er != nil {
				return w.Uint8(1).String(v)
			} else {
				return w.Uint8(4).Uint32(ipv4)
			}
		} else {
			if ipv6, er := IPv6ToInt(ip); er != nil {
				return w.Uint8(1).String(v)
			} else {
				return w.Uint8(6).Uint64(ipv6[0]).Uint64(ipv6[1])
			}
		}
	}
}

// IPArray will encode a list of IPv4 or IPv6 to byte array, each IP is stored as defined in the IP() method
func (w *Writer) IPArray(v []string) *Writer {
	// Write array size
	w.varInt(uint64(len(v)))

	for _, ip := range v {
		w.IP(ip)
	}

	w.flush()
	return w
}

// Reset will reset the underlying bytes of the Encoder
func (w *Writer) Reset() {
	w.bs = w.bs[:0]
}

// WriteTo will write to an io.Writer
func (w *Writer) WriteTo(dest io.Writer) (int64, error) {
	if written, err := dest.Write(w.GetBytes()); err != nil {
		return 0, err
	} else {
		n := int64(written)
		w.Reset()
		return n, nil
	}
}

// GetBytes will expose the underlying bytes
func (w *Writer) GetBytes() []byte {
	return w.bs
}

// Written will return the total number of bytes written
func (w *Writer) Written() int64 {
	return w.written
}

// Close will close the writer
func (w *Writer) Close() (err error) {
	w.bs = nil
	return
}

// Flash temp bytes array (bs) to the buffer
func (w *Writer) flush() {
	if n, err := w.buffer.Write(w.bs); err == nil {
		w.written += int64(n)
		w.bs = w.bs[:]
		//return nil
	} else {
		panic(err)
		//w.bs = w.bs[:]
		//return err
	}
}
