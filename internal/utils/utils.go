package utils

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

// HresultToError parses HRESULT to error value for comparision
//
// see https://learn.microsoft.com/en-us/openspecs/windows_protocols/ms-erref/0642cb2f-2075-4469-918c-4441e69c548a
func HresultToError(err error) error {
	errno, ok := err.(windows.Errno)
	if !ok {
		return err
	}

	hresVal := int32(errno)
	if windows.Errno(hresVal) != errno {
		return err
	}

	if hresVal&0x10000000 != 0 { // this is an NTStatus value
		return windows.NTStatus(hresVal)
	}

	if hresVal&0x20000000 == 0 && hresVal&0x07ff0000 == windows.FACILITY_WIN32 {
		// this is an undecorated error code
		return windows.Errno(hresVal & 0xffff)
	}

	return err
}

func StrSliceToUtf16PtrArr(strSlice []string) ([]*uint16, error) {
	var u16Arr []*uint16
	var count uint32

	for _, str := range strSlice {
		u16Ptr, err := windows.UTF16PtrFromString(str)
		if err != nil {
			return nil, err
		}

		u16Arr = append(u16Arr, u16Ptr)
		count++
	}

	return u16Arr, nil
}

func PZZWSTRToStrings(pzzwstr **uint16) []string {
	result := make([]string, 0)
	bufPtr := *pzzwstr

	for *bufPtr != 0 {
		result = append(result, windows.UTF16PtrToString(bufPtr))

		bufPtr = (*uint16)(unsafe.Add(unsafe.Pointer(bufPtr), (len(result)+1)*2))
	}

	return result
}
