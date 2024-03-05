package binary

import (
	"unsafe"
)

// const notEnoughBytesLayout = "not enough bytes available to decode <%T>, needed %d and has an available %d"

func getStringBytes(str *string) *[]byte {
	return (*[]byte)(unsafe.Pointer(str))
}

func getStringFromBytes(bs []byte) string {
	return *((*string)(unsafe.Pointer(&bs)))
}

//func newNotEnoughBytesError(target interface{}, needed, remaining int) (err error) {
//	err = fmt.Errorf(notEnoughBytesLayout, target, needed, remaining)
//	return
//}
//
//type reader interface {
//	io.Reader
//	io.ByteReader
//}

func expandSlice(bs *[]byte, sz int) {
	if *bs != nil && cap(*bs) >= sz {
		*bs = (*bs)[:sz]
		return
	}

	*bs = make([]byte, sz)
}
