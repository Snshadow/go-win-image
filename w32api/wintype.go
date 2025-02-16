package w32api

import (
	"strconv"
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
	Reserved0         uint32
	Reserved1         uint32
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

// Windows processor architectures
//
// https://learn.microsoft.com/en-us/windows/win32/api/sysinfoapi/ns-sysinfoapi-system_info#members
type Architecture uint32

const (
	Intel       Architecture = 0 // x86
	MIPS        Architecture = 1
	ALPHA       Architecture = 2
	PPC         Architecture = 3
	SHX         Architecture = 4
	ARM         Architecture = 5
	IA64        Architecture = 6
	ALPHA64     Architecture = 7
	MSIL        Architecture = 8
	AMD64       Architecture = 9 // x64
	IA32OnAMD64 Architecture = 10
	Neutral     Architecture = 11
	ARM64       Architecture = 12
	Unknown     Architecture = 0xffff
)

const archStr = "x86MIPSALPHAPPCSHXARMIA-64ALPHA64MSILx86_64IA32OnAMD64NeutralARM64Unknown"

func (a Architecture) String() string {
	switch a {
	case Intel:
		return archStr[:3]
	case MIPS:
		return archStr[3:7]
	case ALPHA:
		return archStr[7:12]
	case PPC:
		return archStr[12:15]
	case SHX:
		return archStr[15:18]
	case ARM:
		return archStr[18:21]
	case IA64:
		return archStr[21:26]
	case ALPHA64:
		return archStr[26:33]
	case MSIL:
		return archStr[33:37]
	case AMD64:
		return archStr[37:43]
	case IA32OnAMD64:
		return archStr[43:54]
	case Neutral:
		return archStr[54:61]
	case ARM64:
		return archStr[61:66]
	case Unknown:
		return archStr[66:]
	}

	return "Architecture(" + strconv.FormatUint(uint64(a), 10) + ")"
}
