// Entity tests

package test

import (
	"fmt"
	"github.com/go-yaaf/yaaf-common/entity"
	"github.com/stretchr/testify/require"
	"testing"
)

// Test binary format of writer or reader
func TestBinaryOfSampleObject(t *testing.T) {
	skipCI(t)

	so := createSampleObject(16)
	data, err := so.MarshalBinary()
	require.NoError(t, err)

	expected := SampleObject{}
	err = expected.UnmarshalBinary(data)
	require.NoError(t, err)

	require.Equal(t, so.StringValue, expected.StringValue)

	require.Equal(t, so, expected)

	fmt.Printf("\n\n")
}

// Test binary format of writer or reader
func TestBinaryOfComplexObject(t *testing.T) {
	skipCI(t)

	co := createComplexObject(16, "complex")

	data, err := co.MarshalBinary()
	require.NoError(t, err)

	expected := ComplexObject{}
	err = expected.UnmarshalBinary(data)
	require.NoError(t, err)

	require.Equal(t, co.StringValue, expected.StringValue)

	require.Equal(t, co, expected)

	fmt.Printf("\n\n")
}

func createSampleObject(intVal int) SampleObject {
	return SampleObject{
		Timestamp:   entity.Now(),
		IntValue:    intVal,
		Int32Value:  18728836,
		Int64Value:  1523324432323,
		IntArray:    []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
		StringValue: "my string",
		StringArray: []string{"label_1", "label_2", "label_3", "label_4"},
	}
}

func createComplexObject(intVal int, strVal string) ComplexObject {
	co := ComplexObject{
		IntValue:    intVal,
		IntArray:    []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		StringValue: strVal,
		StringArray: []string{"complex_1", "complex_2"},
		ComplexValue: SampleObject{
			Timestamp:   entity.Now(),
			IntValue:    123,
			Int32Value:  3422,
			Int64Value:  5411321,
			IntArray:    []int{1, 2, 3},
			StringValue: "label_0",
			StringArray: []string{"l0", "l1", "l2"},
		},
		ComplexValues: nil,
	}

	for i := 0; i < 10; i++ {
		so := createSampleObject(i)
		co.ComplexValues = append(co.ComplexValues, so)
	}
	return co
}
