//go:build windows
// +build windows

package dismapi

import (
	"bytes"
	"reflect"
	"unsafe"
)

type OneBytePacked[T any] interface {
	ToStruct() T
	ToBuffer(st T) []byte
}

// types for handling convert raw byte buffer into struct
type DismFeatureBuf [12]byte
type DismCapabilityBuf [12]byte
type DismPackageInfoBuf [172]byte
type DismPackageInfoExBuf [180]byte
type DismFeatureInfoBuf [44]byte
type DismCapabilityInfoBuf [36]byte
type DismImageInfoBuf [140]byte
type DismMountedImageInfoBuf [28]byte
type DismDriverPackageBuf [100]byte
type DismDriverBuf [52]byte
type DismAppxPackageBuf [68]byte

func (b *DismFeatureBuf) ToStruct() DismFeature {
	return *(*DismFeature)(unsafe.Pointer(&b[0]))
}

func (b *DismCapabilityBuf) ToStruct() DismCapability {
	return *(*DismCapability)(unsafe.Pointer(&b[0]))
}

func (b *DismPackageInfoBuf) ToStruct() DismPackageInfo {
	
}


