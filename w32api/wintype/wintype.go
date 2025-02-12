package wintype

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	WM_APP = 0x8000
)

type WIN32_FIND_DATAW struct {
	FileAttributes    uint32
	CreationFile      windows.Filetime
	LastAccessTime    windows.Filetime
	LastWriteTime     windows.Filetime
	FileSizeHigh      uint32
	FileSizeLow       uint32
	_                 uint32 // Reserved0
	_                 uint32 // Reserved1
	FileName          [windows.MAX_PATH]uint16
	AlternateFileName [14]uint16
	FileType          uint32 // Deprecated: Obsolete. Do not use
	CreatorType       uint32 // Deprecated: Obsolete. Do not use
	FinderFlags       uint16 // Deprecated: Obsolete. Do not use
}

type LONG_PTR int
type ULONG_PTR uintptr

// value of callbackReason from [LPPROGRESS_ROUTINE]
const (
	CALLBACK_CHUNK_FINISHED = 0x0
	CALLBACK_STREAM_SWITCH  = 0x1
)

// return value from [LPPROGRESS_ROUTINE]
const (
	PROGRESS_CONTINUE = 0
	PROGRESS_CANCEL   = 1
	PROGRESS_STOP     = 2
	PROGRESS_QUIET    = 3
)

// application-defined callback function, use with [windows.NewCallback]
type LPPROGRESS_ROUTINE func(
	totalFileSize int64, // LARGE_INTEGER
	totalBytesTransferred int64, // LARGE_INTEGER
	streamSize int64, // LARGE_INTEGER
	streamBytesTransferred int64, // LARGE_INTEGER
	streamNumber uint32,
	callbackReason uint32,
	sourceFile windows.Handle,
	destinationFile windows.Handle,
	data unsafe.Pointer,
) /* uint32 */ uintptr
