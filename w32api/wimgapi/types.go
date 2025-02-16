package wimgapi

import (
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/Snshadow/winimg/w32api"
)

type LPARAM w32api.LONG_PTR
type WPARAM w32api.ULONG_PTR

// WIMCreateFile
const (
	WIM_GENERIC_READ  = windows.GENERIC_READ
	WIM_GENERIC_WRITE = windows.GENERIC_WRITE
	WIM_GENERIC_MOUNT = windows.GENERIC_EXECUTE

	WIM_CREATE_NEW    = windows.CREATE_NEW
	WIM_CREATE_ALWAYS = windows.CREATE_ALWAYS
	WIM_OPEN_EXISTING = windows.OPEN_EXISTING
	WIM_OPEN_ALWAYS   = windows.OPEN_ALWAYS
)

// Compression type of .wim file
type WimCompressionType uint32

const (
	WIM_COMPRESS_NONE WimCompressionType = iota
	WIM_COMPRESS_XPRESS
	WIM_COMPRESS_LZX
	WIM_COMPRESS_LZMS
)

const (
	WIM_CREATED_NEW = iota
	WIM_OPENED_EXISTING
)

// WIMCreateFile, WIMCaptureImage, WIMApplyImage, WIMMountImageHandle flags
const (
	WIM_FLAG_RESERVED       = 0x00000001
	WIM_FLAG_VERIFY         = 0x00000002
	WIM_FLAG_INDEX          = 0x00000004
	WIM_FLAG_NO_APPLY       = 0x00000008
	WIM_FLAG_NO_DIRACL      = 0x00000010
	WIM_FLAG_NO_FILEACL     = 0x00000020
	WIM_FLAG_SHARE_WRITE    = 0x00000040
	WIM_FLAG_FILEINFO       = 0x00000080
	WIM_FLAG_NO_RP_FIX      = 0x00000100
	WIM_FLAG_MOUNT_READONLY = 0x00000200
	WIM_FLAG_MOUNT_FAST     = 0x00000400
	WIM_FLAG_MOUNT_LEGACY   = 0x00000800
	WIM_FLAG_APPLY_CI_EA    = 0x00001000
	WIM_FLAG_WIM_BOOT       = 0x00002000
	WIM_FLAG_APPLY_COMPACT  = 0x00004000
	WIM_FLAG_SUPPORT_EA     = 0x00008000
)

// WIMGetMountedImageList flags
const (
	WIM_MOUNT_FLAG_MOUNTED           = 0x00000001
	WIM_MOUNT_FLAG_MOUNTING          = 0x00000002
	WIM_MOUNT_FLAG_REMOUNTABLE       = 0x00000004
	WIM_MOUNT_FLAG_INVALID           = 0x00000008
	WIM_MOUNT_FLAG_NO_WIM            = 0x00000010
	WIM_MOUNT_FLAG_NO_MOUNTDIR       = 0x00000020
	WIM_MOUNT_FLAG_MOUNTDIR_REPLACED = 0x00000040
	WIM_MOUNT_FLAG_READWRITE         = 0x00000100
)

// WIMCommitImageHandle flags
const (
	WIM_COMMIT_FLAG_APPEND = 0x00000001
)

// WIMSetReferenceFile
const (
	WIM_REFERENCE_APPEND  = 0x00010000
	WIM_REFERENCE_REPLACE = 0x00020000
)

// WIMExportImage
const (
	WIM_EXPORT_ALLOW_DUPLICATES   = 0x00000001
	WIM_EXPORT_ONLY_RESOURCES     = 0x00000002
	WIM_EXPORT_ONLY_METADATA      = 0x00000004
	WIM_EXPORT_VERIFY_SOURCE      = 0x00000008
	WIM_EXPORT_VERIFY_DESTINATION = 0x00000010
)

// WimRegisterMessageCallback
const (
	INVALID_CALLBACK_VALUE = 0xffffffff
)

// WIMCopyFile
const (
	WIM_COPY_FILE_RETRY = 0x01000000
)

// WIMDeleteImageMounts
const (
	WIM_DELETE_MOUNTS_ALL = 0x00000001
)

// WIMRegisterLogFile
const (
	WIM_LOGFILE_UTF8 = 0x00000001
)

// WIMMessageCallback Notifications
//
// https://learn.microsoft.com/en-us/windows-hardware/manufacture/desktop/wim/dd851929(v=msdn.10)
const (
	WIM_MSG = w32api.WM_APP + 0x1476 + iota
	WIM_MSG_TEXT
	WIM_MSG_PROGRESS
	WIM_MSG_PROCESS
	WIM_MSG_SCANNING
	WIM_MSG_SETRANGE
	WIM_MSG_SETPOS
	WIM_MSG_STEPIT
	WIM_MSG_COMPRESS
	WIM_MSG_ERROR
	WIM_MSG_ALIGNMENT
	WIM_MSG_RETRY
	WIM_MSG_SPLIT
	WIM_MSG_FILEINFO
	WIM_MSG_INFO
	WIM_MSG_WARNING
	WIM_MSG_CHK_PROCESS
	WIM_MSG_WARNING_OBJECTID
	WIM_MSG_STALE_MOUNT_DIR
	WIM_MSG_STALE_MOUNT_FILE
	WIM_MSG_MOUNT_CLEANUP_PROGRESS
	WIM_MSG_CLEANUP_SCANNING_DRIVE
	WIM_MSG_IMAGE_ALREADY_MOUNTED
	WIM_MSG_CLEANUP_UNMOUNTING_IMAGE
	WIM_MSG_QUERY_ABORT
	WIM_MSG_IO_RANGE_START_REQUEST_LOOP
	WIM_MSG_IO_RANGE_END_REQUEST_LOOP
	WIM_MSG_IO_RANGE_REQUEST
	WIM_MSG_IO_RANGE_RELEASE
	WIM_MSG_VERIFY_PROGRESS
	WIM_MSG_COPY_BUFFER
	WIM_MSG_METADATA_EXCLUDE
	WIM_MSG_GET_APPLY_ROOT
	WIM_MSG_MDPAD
	WIM_MSG_STEPNAME
	WIM_MSG_PERFILE_COMPRESS
	WIM_MSG_CHECK_CI_EA_PREREQUISITE_NOT_MET
	WIM_MSG_JOURNALING_ENABLED
)

// WIMMessageCallback Return Codes
const (
	WIM_MSG_SUCCESS     = 0x00000000 // ERROR_SUCCESS
	WIM_MSG_DONE        = 0xFFFFFFF0
	WIM_MSG_SKIP_ERROR  = 0xFFFFFFFE
	WIM_MSG_ABORT_IMAGE = 0xFFFFFFFF
)

// WIM_INFO flags values
const (
	WIM_ATTRIBUTE_NORMAL        = 0x00000000
	WIM_ATTRIBUTE_RESOURCE_ONLY = 0x00000001
	WIM_ATTRIBUTE_METADATA_ONLY = 0x00000002
	WIM_ATTRIBUTE_VERIFY_DATA   = 0x00000004
	WIM_ATTRIBUTE_RP_FIX        = 0x00000008
	WIM_ATTRIBUTE_SPANNED       = 0x00000010
	WIM_ATTRIBUTE_READONLY      = 0x00000020
)

// An abstract type implemented by the caller when using File I/O callbacks
type PFILEIOCALLBACK_SESSION unsafe.Pointer

// used by WIMGetAttributes
type WIM_INFO struct {
	WimPath         [windows.MAX_PATH]uint16
	Guid            windows.GUID
	ImageCount      uint32
	CompressionType uint32
	PartNumber      uint16
	TotalParts      uint16
	BootIndex       uint32
	WimAttributes   uint32
	WimFlagsAndAttr uint32
}

// used for getting the list of mounted images
type WIM_MOUNT_LIST struct {
	WimPath      [windows.MAX_PATH]uint16
	MountPath    [windows.MAX_PATH]uint16
	ImageIndex   uint32
	MountedForRW int32 // BOOL
}

type WIM_MOUNT_INFO_LEVEL0 WIM_MOUNT_LIST

// new structure with additional data
type WIM_MOUNT_INFO_LEVEL1 struct {
	WimPath    [windows.MAX_PATH]uint16
	MountPath  [windows.MAX_PATH]uint16
	ImageIndex uint32
	MountFlags uint32
}

// enumeration for WIMGetMountedImageInfo
type MOUNTED_IMAGE_INFO_LEVELS uint32

const (
	MountedImageLevel0 MOUNTED_IMAGE_INFO_LEVELS = iota
	MountedImageLevel1
	MountedImageLevelInvalid
)

type WIM_IO_RANGE_CALLBACK struct {
	// the callback session that corresponds to the file that is being queried
	Session PFILEIOCALLBACK_SESSION
	// LARGE_INTEGER ; filled in by WIMGAPI for both messages
	Offset, Size int64
	// BOOL ; filled in by the callback for WIM_MSG_IO_RANGE_REQUEST (set to TRUE to
	// indicate data in the specified range is available, and FALSE to indicate
	// it is not yet available)
	Available int32
}

type WIM_FIND_DATA struct {
	w32api.WIN32_FIND_DATAW
	Hash               [20]byte
	SecurityDescriptor *windows.SECURITY_DESCRIPTOR
	// double-null terminated, read from *AlternateStreamNames with [windows.UTF16ToString] & [unsafe.Add]
	AlternateStreamNames  **uint16
	PbReparseData         *byte
	CbReparseData         uint32
	ResourceSize          uint64 // ULARGE_INTEGER
	Resourceoffset        int64  // LARGE_INTEGER
	ResourceReferencePath *uint16
}

// File I/O callback prototypes, used with [windows.NewCallback]

type FileIOCallbackOpenFile func(fileName *uint16) PFILEIOCALLBACK_SESSION
type FileIOCallbackCloseFile func(file PFILEIOCALLBACK_SESSION) /* BOOL */ uintptr
type FileIOCallbackReadFile func(file PFILEIOCALLBACK_SESSION, buffer unsafe.Pointer, numberOfBytesToRead uint32, numberOfBytesRead *uint32, overlapped *windows.Overlapped) /* BOOL */ uintptr
type FileIOCallbackSetFilePointer func(file PFILEIOCALLBACK_SESSION, distanceToMove int64, newFilePointer *int64, moveMethod uint32) /* BOOL */ uintptr
type FileIOCallbackGetFileSize func(file windows.Handle, fileSize *int64) /* BOOL */ uintptr

type SFileIOCallbackInfo struct {
	OpenFile       uintptr // FileIOCallbackOpenFile
	CloseFile      uintptr // FileIOCallbackCloseFile
	ReadFile       uintptr // FileIOCallbackReadFile
	SetFilePointer uintptr // FileIOCallbackSetFilePointer
	GetFileSize    uintptr // FileIOCallbackGetFileSize
}

// Abstract (opaque) type for WIM files used with
// WIMEnumImageFiles API
type PWIM_ENUM_FILE unsafe.Pointer

// API for fast enumeration for image files
type WIMEnumImageFilesCallback func(
	findFileData *WIM_FIND_DATA,
	enumFile PWIM_ENUM_FILE,
	enumContext unsafe.Pointer,
) /* HRESULT */ uintptr
