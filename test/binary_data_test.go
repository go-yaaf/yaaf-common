// Entity tests

package test

import (
	"github.com/go-yaaf/yaaf-common/entity"
	"github.com/go-yaaf/yaaf-common/utils/binary"
)

type SampleObjectNotBinary struct {
	Timestamp   entity.Timestamp
	SrcIP       string
	DstIPs      []string
	IntValue    int
	Int32Value  int32
	Int64Value  int64
	IntArray    []int
	StringValue string
	StringArray []string
}

type SampleObject struct {
	Timestamp   entity.Timestamp
	SrcIP       string
	DstIPs      []string
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
	w.IP(s.SrcIP)
	w.IPArray(s.DstIPs)
	w.Timestamp(s.Timestamp).Int(s.IntValue).Int32(s.Int32Value).Int64(s.Int64Value).IntArray(s.IntArray).String(s.StringValue).StringArray(s.StringArray)
	return w.GetBytes(), nil
}

// UnmarshalBinary reads a wire-format byte array to fill the current structure
func (s *SampleObject) UnmarshalBinary(data []byte) (e error) {
	r := binary.NewReader(data)
	if s.SrcIP, e = r.IP(); e != nil {
		return e
	}
	if s.DstIPs, e = r.IPArray(); e != nil {
		return e
	}
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

type ComplexObject struct {
	SrcIP         string
	DstIPs        []string
	IntValue      int
	IntArray      []int
	StringValue   string
	StringArray   []string
	ComplexValue  SampleObject
	ComplexValues []SampleObject
}

// MarshalBinary convert current structure to a minimal wire-format byte array
func (c *ComplexObject) MarshalBinary() (data []byte, err error) {
	w := binary.NewWriter()
	w.IP(c.SrcIP)
	w.IPArray(c.DstIPs)
	w.Int(c.IntValue).IntArray(c.IntArray).String(c.StringValue).StringArray(c.StringArray)

	// Marshal complex value
	if cvd, e := c.ComplexValue.MarshalBinary(); e != nil {
		return nil, e
	} else {
		w.Object(&cvd)
	}

	// Marshal array of complex values
	objects := make([][]byte, 0)
	for _, cv := range c.ComplexValues {
		if cvd, e := cv.MarshalBinary(); e != nil {
			return nil, e
		} else {
			objects = append(objects, cvd)
		}
	}
	w.ObjectArray(&objects)

	return w.GetBytes(), nil
}

// UnmarshalBinary reads a wire-format byte array to fill the current structure
func (c *ComplexObject) UnmarshalBinary(data []byte) (e error) {
	r := binary.NewReader(data)
	if c.SrcIP, e = r.IP(); e != nil {
		return e
	}
	if c.DstIPs, e = r.IPArray(); e != nil {
		return e
	}
	if c.IntValue, e = r.Int(); e != nil {
		return e
	}
	if c.IntArray, e = r.IntArray(); e != nil {
		return e
	}
	if c.StringValue, e = r.String(); e != nil {
		return e
	}
	if c.StringArray, e = r.StringArray(); e != nil {
		return e
	}

	// Unmarshal objects into ComplexValue field
	if bytes, err := r.Object(); err != nil {
		return err
	} else {
		if er := c.ComplexValue.UnmarshalBinary(bytes); er != nil {
			return er
		}
	}

	// Unmarshal objects into ComplexValues field array
	c.ComplexValues = make([]SampleObject, 0)

	if mBytes, err := r.ObjectArray(); err != nil {
		return err
	} else {
		for _, objData := range mBytes {
			obj := SampleObject{}
			if er := obj.UnmarshalBinary(objData); er != nil {
				return er
			} else {
				c.ComplexValues = append(c.ComplexValues, obj)
			}
		}
	}

	return nil
}
