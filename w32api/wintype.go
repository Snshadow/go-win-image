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

// Some common undocumented error values from wimgapi or DISM API
// functions. Defined for showing error message which would be
// shown in Dism.exe instead of showing "winapi error #(num)"
type WinimgInternalErr uint32

// TODO need to add more?
const (
	ErrDirNotExist       WinimgInternalErr = 0xc1420114
	ErrMountDirInUse     WinimgInternalErr = 0xc1420117
	ErrCommitUnavailable WinimgInternalErr = 0xc142011d
	ErrAlreadyMounted    WinimgInternalErr = 0xc1420127
)

func (e WinimgInternalErr) Error() string {
	switch e {
	case ErrDirNotExist:
		return "The user attempted to mount to a directory that does not exist. This is not supported."
	case ErrMountDirInUse:
		return "The directory could not be completely unmounted. This is usually due to applications that still have files opened withing the directory. Close these files and unmount again to complete the unmount process."
	case ErrCommitUnavailable:
		return "The specified mounted image cannot be committed back into the WIM. This occurs when an image has been through a partial unmount or when an image is still being mounted. If this image was unmounted with commit earlier, then the commit probably succeeded. Please validate that this is the case and then unmount without commit."
	case ErrAlreadyMounted:
		return "The specified image in the specified wim is already mounted for read and write access."
	}

	return ""
}

// WrapInternalErr tries to wrap error value if it is a
// known undocumented error value to show error message.
func WrapInternalErr(e error) error {
	errno, ok := e.(windows.Errno)
	if !ok {
		return e
	}

	lookupErr := WinimgInternalErr(errno)

	if lookupErr.Error() != "" {
		return lookupErr
	}

	return e
}
