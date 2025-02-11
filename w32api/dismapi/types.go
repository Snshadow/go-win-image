//go:build windows
// +build windows

package dismapi

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

type DismSession uint32

// void DismProgressCallback(_In_ UINT Current,	_In_ UINT Total, _In_ PVOID UserData)
//
// return uintptr for windows.NewCallback requirement
type DismProgressCallback func(current uint32, total uint32, userData unsafe.Pointer) uintptr

// Dism error values
const (
	DISMAPI_S_RELOAD_IMAGE_SESSION_REQUIRED             = 0x00000001
	DISMAPI_E_DISMAPI_NOT_INITIALIZED                   = 0xc0040001
	DISMAPI_E_SHUTDOWN_IN_PROGRESS                      = 0xc0040002
	DISMAPI_E_OPEN_SESSION_HANDLES                      = 0xc0040003
	DISMAPI_E_INVALID_DISM_SESSION                      = 0xc0040004
	DISMAPI_E_INVALID_IMAGE_INDEX                       = 0xc0040005
	DISMAPI_E_INVALID_IMAGE_NAME                        = 0xc0040006
	DISMAPI_E_UNABLE_TO_UNMOUNT_IMAGE_PATH              = 0xc0040007
	DISMAPI_E_LOGGING_DISABLED                          = 0xc0040009
	DISMAPI_E_OPEN_HANDLES_UNABLE_TO_UNMOUNT_IMAGE_PATH = 0xc004000a
	DISMAPI_E_OPEN_HANDLES_UNABLE_TO_MOUNT_IMAGE_PATH   = 0xc004000b
	DISMAPI_E_OPEN_HANDLES_UNABLE_TO_REMOUNT_IMAGE_PATH = 0xc004000c
	DISMAPI_E_PARENT_FEATURE_DISABLED                   = 0xc004000d
	DISMAPI_E_MUST_SPECIFY_ONLINE_IMAGE                 = 0xc004000e
	DISMAPI_E_INVALID_PRODUCT_KEY                       = 0xc004000f
	DISMAPI_E_NEEDS_REMOUNT                             = 0xc1510114
	DISMAPI_E_UNKNOWN_FEATURE                           = 0x800f080c
	DISMAPI_E_BUSY                                      = 0x800f0902
)

var (
	DISM_ONLINE_IMAGE = [...]uint16{'D', 'I', 'S', 'M', '_', '{', '5', '3', 'B', 'F', 'A', 'E', '5', '2', '-', 'B', '1', '6', '7', '-', '4', 'E', '2', 'F', '-', 'A', '2', '5', '8', '-', '0', 'A', '3', '7', 'B', '5', '7', 'F', 'F', '8', '4', '5', '}', '\x00'} // L"DISM_{53BFAE52-B167-4E2F-A258-0A37B57FF845}"
)

// Dism constants
const (
	DISM_SESSION_DEFAULT           = 0
	DISM_MOUNT_READWRITE           = 0x00000000
	DISM_MOUNT_READONLY            = 0x00000001
	DISM_MOUNT_OPTIMIZE            = 0x00000002
	DISM_MOUNT_CHECK_INTEGRITY     = 0x00000004
	DISM_COMMIT_IMAGE              = 0x00000000
	DISM_DISCARD_IMAGE             = 0x00000001
	DISM_COMMIT_GENERATE_INTEGRITY = 0x00010000
	DISM_COMMIT_APPEND             = 0x00020000
	DISM_COMMIT_MASK               = 0xffff0000

	DISM_RESERVED_STORAGE_DISABLED = 0x00000000
	DISM_RESERVED_STORAGE_ENABLED  = 0x00000001
)

// Dism enums

type DismLogLevel uint32

const (
	DismLogErrors DismLogLevel = iota
	DismLogErrorsWarnings
	DismLogErrorsWarningsInfo
	DismLogErrorsWarningsIntoDebug
)

type DismImageIdentifier uint32

const (
	DismImageIndex DismImageIdentifier = iota
	DismImageName
	DismImageNone
)

type DismMountMode uint32

const (
	DismReadWrite DismMountMode = iota
	DismReadOnly
)

type DismImageType uint32

const (
	DismImageTypeUnsupported DismImageType = ^DismImageType(0) // -1
	DismImageTypeWim         DismImageType = 0
	DismImageTypeVhd         DismImageType = 1
)

type DismImageBootable uint32

const (
	DismImageBootableYes DismImageBootable = iota
	DismImageBootableNo
	DismImageBootableUnknown
)

type DismMountStatus uint32

const (
	DismMountStatusOk DismMountStatus = iota
	DismMountStatusNeedsRemount
	DismMountStatusInvalid
)

type DismImageHealthState uint32

const (
	DismImageHealthy DismImageHealthState = iota
	DismImageRepairable
	DismImageNonRepairable
)

type DismPackageIdentifier uint32

const (
	DismPackageNone DismPackageIdentifier = iota
	DismPackageName
	DismPackagePath
)

type DismPackageFeatureState uint32

const (
	DismStateNotPresent DismPackageFeatureState = iota
	DismStateUninstallPending
	DismStateStaged
	DismStateResolved, DismStateRemoved DismPackageFeatureState = iota, iota
	DismStateInstalled                  DismPackageFeatureState = iota
	DismStateInstallPending
	DismStateSuperseded
	DismStatePartiallyInstalled
)

type DismReleaseType uint32

const (
	DismReleaseTypeCriticalUpdate DismReleaseType = iota
	DismReleaseTypeDriver
	DismReleaseTypeFeaturePack
	DismReleaseTypeHotfix
	DismReleaseTypeSecurityUpdate
	DismReleaseTypeSoftwareUpdate
	DismReleaseTypeUpdate
	DismReleaseTypeUpdateRollup
	DismReleaseTypeLanguagePack
	DismReleaseTypeFoundation
	DismReleaseTypeServicePack
	DismReleaseTypeProduct
	DismReleaseTypeLocalPack
	DismReleaseTypeOther
	DismReleaseTypeOnDemandPack
)

type DismRestartType uint32

const (
	DismRestartNo DismRestartType = iota
	DismRestartPossible
	DismRestartRequired
)

type DismDriverSignature uint32

const (
	DismDriverSignatureUnknown DismDriverSignature = iota
	DismDriverSignatureUnsigned
	DismDriverSignatureSigned
)

type DismFullyOfflineInstallableType uint32

const (
	DismFullyOfflineInstallable DismFullyOfflineInstallableType = iota
	DismFullyOfflineNotInstallable
	DismFullyOfflineInstallableUndetermined
)

type DismStubPackageOption uint32

const (
	DismStubPackageOptionNone DismStubPackageOption = 0 + iota
	DismStubPackageOptionInstallFull
	DismStubPackageOptionInstallStub
)

// Dism structs (packed by 1 byte)

type DismPackage struct {
	PackageName  *uint16
	PackageState DismPackageFeatureState
	ReleaseType  DismReleaseType
	InstallTime  windows.Systemtime
}

type DismCustomProperty struct {
	Name  *uint16
	Value *uint16
	Path  *uint16
}

type DismFeature struct {
	FeatureName *uint16
	State       DismPackageFeatureState
}

type DismCapability struct {
	Name  *uint16
	State DismPackageFeatureState
}

type DismPackageInfo struct {
	PackageName         *uint16
	PackageState        DismPackageFeatureState
	ReleaseType         DismReleaseType
	InstallTime         windows.Systemtime
	Applicable          int32 // BOOL
	Copyright           *uint16
	Company             *uint16
	CreationTime        windows.Systemtime
	DisplayName         *uint16
	Description         *uint16
	InstallClient       *uint16
	InstallPackageName  *uint16
	LastUpdateTime      windows.Systemtime
	ProductName         *uint16
	ProductVersion      *uint16
	RestartRequired     DismRestartType
	FullyOffline        DismFullyOfflineInstallableType
	SupportInformation  *uint16
	CustomProperty      *DismCustomProperty
	CustomPropertyCount uint32
	Feature             *DismFeature
	FeatureCount        uint32
}

type DismPackageInfoEx struct {
	DismPackageInfo
	CapabilityId *uint16
}

type DismFeatureInfo struct {
	FeatureName         *uint16
	FeatureState        DismPackageFeatureState
	DisplayName         *uint16
	Description         *uint16
	RestartRequired     DismRestartType
	CustomProperty      *DismCustomProperty
	CustomPropertyCount uint32
}

type DismCapabilityInfo struct {
	Name         *uint16
	State        DismPackageFeatureState
	DisplayName  *uint16
	Description  *uint16
	DownloadSize uint32
	InstallSize  uint32
}

type DismString struct {
	Value *uint16
}

type DismLanguage DismString

type DismWimCustomizedInfo struct {
	Size           uint32
	DirectoryCount uint32
	FileCount      uint32
	CreateTime     windows.Systemtime
	ModifiedTime   windows.Systemtime
}

type DismImageInfo struct {
	ImageType            DismImageType
	ImageIndex           uint32
	ImageName            *uint16
	ImageDescription     *uint16
	ImageSize            uint64
	Architecture         uint32
	ProductName          *uint16
	EditionId            *uint16
	InstallationType     *uint16
	Hal                  *uint16
	ProductType          *uint16
	ProductSuite         *uint16
	MajorVersion         uint32
	MinorVersion         uint32
	Build                uint32
	SpBuild              uint32
	SpLevel              uint32
	Bootable             DismImageBootable
	SystemRoot           *uint16
	Language             *DismLanguage
	LanguageCount        uint32
	DefaultLanguageIndex uint32
	CustomizedInfo       unsafe.Pointer
}

type DismMountedImageInfo struct {
	MountPath     *uint16
	ImageFilePath *uint16
	ImageIndex    uint32
	MountMode     DismMountMode
	MountStatus   DismMountStatus
}

type DismDriverPackage struct {
	PublishedName    *uint16
	OriginalFileName *uint16
	Inbox            int32 // BOOL
	CatalogFile      *uint16
	ClassName        *uint16
	ClassGuid        *uint16
	ClassDescription *uint16
	BootCritical     int32 // BOOL
	DriverSignature  DismDriverSignature
	ProviderName     *uint16
	Date             windows.Systemtime
	MajorVersion     uint32
	MinorVersion     uint32
	Build            uint32
	Revision         uint32
}

type DismDriver struct {
	ManufacturerName    *uint16
	HardwareDescription *uint16
	HardwareId          *uint16
	Architecture        uint32
	ServiceName         *uint16
	CompatibleIds       *uint16
	ExcludeIds          *uint16
}

type DismAppxPackage struct {
	PackageName     *uint16
	DisplayName     *uint16
	PublisherId     *uint16
	MajorVersion    uint32
	MinorVersion    uint32
	Build           uint32
	RevisionNumber  uint32
	Architecture    uint32
	ResourceId      *uint16
	InstallLocation *uint16
	Region          *uint16
}
