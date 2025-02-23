package w32api

import (
	"strconv"
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/Snshadow/winimg/internal/utils"
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

const archStr = "x86MIPSALPHAPPCSHXARMIA-64ALPHA64MSILx86_64IA-32OnAMD64NeutralARM64Unknown"

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
		return archStr[43:55]
	case Neutral:
		return archStr[55:62]
	case ARM64:
		return archStr[62:67]
	case Unknown:
		return archStr[67:]
	}

	return "Architecture(" + strconv.FormatUint(uint64(a), 10) + ")"
}

// Get error message from dll or collected string in
// advance.
//
// Defined for showing error message which would be shown
// in Dism.exe instead of showing "winapi error #(num)"
type WinimgInternalErr struct {
	errno  windows.Errno
	module uintptr
	str    string
}

func (e *WinimgInternalErr) Errno() windows.Errno {
	return e.errno
}

func (e *WinimgInternalErr) Error() string {
	if e.str != "" {
		return e.str
	}

	return utils.GetErrorMessage(uint32(e.errno), e.module)
}

func newInternalErr(e windows.Errno, mod uintptr, s string) *WinimgInternalErr {
	return &WinimgInternalErr{
		errno:  e,
		module: mod,
		str:    s,
	}
}

// WrapInternalErr tries to wrap error value if it was
// not handled by other error implementation to get error
// message from module or external function.
func WrapInternalErr(e error, module uintptr, str string) error {
	errno, ok := e.(windows.Errno)
	if !ok {
		return e
	}

	return newInternalErr(errno, module, str)
}
