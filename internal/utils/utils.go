package utils

import (
	"fmt"
	"reflect"
	"unsafe"
)

func getPackedOffset(st reflect.Type, recursive bool, outerOffsets ...[]uintptr) []uintptr {
	var offsets []uintptr
	if len(outerOffsets) != 0 {
		offsets = outerOffsets[0]
	}

	var baseOffset uintptr

	if offsets != nil {
		baseOffset = offsets[len(offsets)-1]
	}

	fieldNum := st.NumField()

	for i := 0; i < fieldNum; i++ {
		f := st.Field(i)
		if recursive && f.Type.Kind() == reflect.Struct {
			offsets = getPackedOffset(f.Type, recursive, offsets)
		}

		baseOffset += f.Type.Size()

		offsets = append(offsets, baseOffset)
	}

	return offsets
}

// GetPackedOffset calculates and returns the offsets of fields within
// a struct without considering any padding by summing up the sizes of
// each field's type using the reflect package. If the `recursive` is
// set to true, the offsets for fields within nested structs are retrieved,
// relative to the outermost struct.
func GetPackedOffset(st any, recursive bool) ([]uintptr, error) {
	t := reflect.TypeOf(st)
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("not a struct")
	}

	return getPackedOffset(t, recursive), nil
}
