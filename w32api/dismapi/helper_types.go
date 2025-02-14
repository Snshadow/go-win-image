package dismapi

import (
	"time"
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/Snshadow/winimg/w32api"
)

// helper structs with go friendly types

type GoDismPackage struct {
	PackageName  string
	PackageState DismPackageFeatureState
	ReleaseType  DismReleaseType
	InstallTime  time.Time
}

func (g *GoDismPackage) fill(st *DismPackage) {
	g.PackageName = windows.UTF16PtrToString(st.PackageName)
	g.PackageState = st.PackageState
	g.ReleaseType = st.ReleaseType
	g.InstallTime = time.Date(int(st.InstallTime.Year), time.Month(st.InstallTime.Month), int(st.InstallTime.Day), int(st.InstallTime.Hour), int(st.InstallTime.Minute), int(st.InstallTime.Second), int(st.InstallTime.Milliseconds)*1000000, time.Local)
}

type GoDismCustomProperty struct {
	Name  string
	Value string
	Path  string
}

func (g *GoDismCustomProperty) fill(st *DismCustomProperty) {
	g.Name = windows.UTF16PtrToString(st.Name)
	g.Value = windows.UTF16PtrToString(st.Value)
	g.Path = windows.UTF16PtrToString(st.Path)
}

type GoDismFeature struct {
	FeatureName string
	State       DismPackageFeatureState
}

func (g *GoDismFeature) fill(st *DismFeature) {
	g.FeatureName = windows.UTF16PtrToString(st.FeatureName)
	g.State = st.State
}

type GoDismCapability struct {
	Name  string
	State DismPackageFeatureState
}

func (g *GoDismCapability) fill(st *DismCapability) {
	g.Name = windows.UTF16PtrToString(st.Name)
	g.State = st.State
}

type GoDismPackageInfo struct {
	PackageName        string
	PackageState       DismPackageFeatureState
	ReleaseType        DismReleaseType
	InstallTime        time.Time
	Applicable         bool
	Copyright          string
	Company            string
	CreationTime       time.Time
	DisplayName        string
	Description        string
	InstallClient      string
	InstallPackageName string
	LastUpdateTime     time.Time
	ProductName        string
	ProductVersion     string
	RestartRequired    DismRestartType
	FullyOffline       DismFullyOfflineInstallableType
	SupportInformation string
	CustomProperty     []GoDismCustomProperty
	Feature            []GoDismFeature
}

func (g *GoDismPackageInfo) fill(st *DismPackageInfo) {
	g.PackageName = windows.UTF16PtrToString(st.PackageName)
	g.PackageState = st.PackageState
	g.ReleaseType = st.ReleaseType
	g.InstallTime = time.Date(int(st.InstallTime.Year), time.Month(st.InstallTime.Month), int(st.InstallTime.Day), int(st.InstallTime.Hour), int(st.InstallTime.Minute), int(st.InstallTime.Second), int(st.InstallTime.Milliseconds)*1000000, time.Local)
	g.Applicable = st.Applicable != 0
	g.Copyright = windows.UTF16PtrToString(st.Copyright)
	g.Company = windows.UTF16PtrToString(st.Company)
	g.CreationTime = time.Date(int(st.CreationTime.Year), time.Month(st.CreationTime.Month), int(st.CreationTime.Day), int(st.CreationTime.Hour), int(st.CreationTime.Minute), int(st.CreationTime.Second), int(st.CreationTime.Milliseconds)*1000000, time.Local)
	g.DisplayName = windows.UTF16PtrToString(st.DisplayName)
	g.Description = windows.UTF16PtrToString(st.Description)
	g.InstallClient = windows.UTF16PtrToString(st.InstallClient)
	g.InstallPackageName = windows.UTF16PtrToString(st.InstallPackageName)
	g.LastUpdateTime = time.Date(int(st.LastUpdateTime.Year), time.Month(st.LastUpdateTime.Month), int(st.LastUpdateTime.Day), int(st.LastUpdateTime.Hour), int(st.LastUpdateTime.Minute), int(st.LastUpdateTime.Second), int(st.LastUpdateTime.Milliseconds)*1000000, time.Local)
	g.ProductName = windows.UTF16PtrToString(st.ProductName)
	g.ProductVersion = windows.UTF16PtrToString(st.ProductVersion)
	g.RestartRequired = st.RestartRequired
	g.FullyOffline = st.FullyOffline
	g.SupportInformation = windows.UTF16PtrToString(st.SupportInformation)
	var customProperty []GoDismCustomProperty
	if st.CustomProperty != nil && st.CustomPropertyCount != 0 {
		stPtr := st.CustomProperty
		stSize := unsafe.Sizeof(*st.CustomProperty)

		for i := uint32(0); i < st.CustomPropertyCount; i++ {
			var goSt GoDismCustomProperty
			goSt.fill((*DismCustomProperty)(unsafe.Pointer(st.CustomProperty)))

			customProperty = append(customProperty, goSt)

			stPtr = (*DismCustomProperty)(unsafe.Add(unsafe.Pointer(stPtr), stSize))
		}
	}
	g.CustomProperty = customProperty
	var feature []GoDismFeature
	if st.Feature != nil && st.FeatureCount != 0 {
		stPtr := (*byte)(unsafe.Pointer(st.Feature))
		bufSize := GetPackedSize(*st.Feature)

		for i := uint32(0); i < st.FeatureCount; i++ {
			unpacked, _ := ToStruct[DismFeature](unsafe.Slice(stPtr, bufSize))

			var goSt GoDismFeature
			goSt.fill(&unpacked)

			feature = append(feature, goSt)

			stPtr = (*byte)(unsafe.Add(unsafe.Pointer(stPtr), bufSize))
		}
	}
	g.Feature = feature
}

type GoDismPackageInfoEx struct {
	GoDismPackageInfo
	CapabilityId string
}

func (g *GoDismPackageInfoEx) fill(st *DismPackageInfoEx) {
	g.GoDismPackageInfo.fill(&st.DismPackageInfo)
	g.CapabilityId = windows.UTF16PtrToString(st.CapabilityId)
}

type GoDismFeatureInfo struct {
	FeatureName     string
	FeatureState    DismPackageFeatureState
	DisplayName     string
	Description     string
	RestartRequired DismRestartType
	CustomProperty  []GoDismCustomProperty
}

func (g *GoDismFeatureInfo) fill(st *DismFeatureInfo) {
	g.FeatureName = windows.UTF16PtrToString(st.FeatureName)
	g.FeatureState = st.FeatureState
	g.DisplayName = windows.UTF16PtrToString(st.DisplayName)
	g.Description = windows.UTF16PtrToString(st.Description)
	g.RestartRequired = st.RestartRequired
	var customProperty []GoDismCustomProperty
	if st.CustomProperty != nil && st.CustomPropertyCount != 0 {
		stPtr := st.CustomProperty
		stSize := unsafe.Sizeof(*st.CustomProperty)

		for i := uint32(0); i < st.CustomPropertyCount; i++ {
			var goSt GoDismCustomProperty
			goSt.fill((*DismCustomProperty)(unsafe.Pointer(st.CustomProperty)))

			customProperty = append(customProperty, goSt)

			stPtr = (*DismCustomProperty)(unsafe.Add(unsafe.Pointer(stPtr), stSize))
		}
	}
	g.CustomProperty = customProperty
}

type GoDismCapabilityInfo struct {
	Name         string
	State        DismPackageFeatureState
	DisplayName  string
	Description  string
	DownloadSize uint32
	InstallSize  uint32
}

func (g *GoDismCapabilityInfo) fill(st *DismCapabilityInfo) {
	g.Name = windows.UTF16PtrToString(st.Name)
	g.State = st.State
	g.DisplayName = windows.UTF16PtrToString(st.DisplayName)
	g.Description = windows.UTF16PtrToString(st.Description)
	g.DownloadSize = st.DownloadSize
	g.InstallSize = st.InstallSize
}

type GoDismWimCustomizedInfo struct {
	Size           uint32
	DirectoryCount uint32
	FileCount      uint32
	CreateTime     time.Time
	ModifiedTime   time.Time
}

func (g *GoDismWimCustomizedInfo) fill(st *DismWimCustomizedInfo) {
	g.Size = st.Size
	g.DirectoryCount = st.DirectoryCount
	g.FileCount = st.FileCount
	g.CreateTime = time.Date(int(st.CreateTime.Year), time.Month(st.CreateTime.Month), int(st.CreateTime.Day), int(st.CreateTime.Hour), int(st.CreateTime.Minute), int(st.CreateTime.Second), int(st.CreateTime.Milliseconds)*1000000, time.Local)
	g.ModifiedTime = time.Date(int(st.ModifiedTime.Year), time.Month(st.ModifiedTime.Month), int(st.ModifiedTime.Day), int(st.ModifiedTime.Hour), int(st.ModifiedTime.Minute), int(st.ModifiedTime.Second), int(st.ModifiedTime.Milliseconds)*1000000, time.Local)
}

type GoDismImageInfo struct {
	ImageType        DismImageType
	ImageIndex       uint32
	ImageName        string
	ImageDescription string
	ImageSize        uint64
	Architecture     w32api.Architecture
	ProductName      string
	EdtitionId       string
	InstallationType string
	Hal              string
	ProductType      string
	ProductSuite     string
	MajorVersion     uint32
	MinorVersion     uint32
	Build            uint32
	SpBuild          uint32
	SpLevel          uint32
	Bootable         DismImageBootable
	SystemRoot       string
	Language         []string                 // DismLanguage
	DefaultLanguage  string                   // DismLanguage[DefaultLanguageIndex]
	CustomizedInfo   *GoDismWimCustomizedInfo // nil for vhd image
}

func (g *GoDismImageInfo) fill(st *DismImageInfo) {
	g.ImageType = st.ImageType
	g.ImageIndex = st.ImageIndex
	g.ImageName = windows.UTF16PtrToString(st.ImageName)
	g.ImageDescription = windows.UTF16PtrToString(st.ImageDescription)
	g.ImageSize = st.ImageSize
	g.Architecture = w32api.Architecture(st.Architecture)
	g.ProductName = windows.UTF16PtrToString(st.ProductName)
	g.EdtitionId = windows.UTF16PtrToString(st.EditionId)
	g.InstallationType = windows.UTF16PtrToString(st.InstallationType)
	g.Hal = windows.UTF16PtrToString(st.Hal)
	g.ProductType = windows.UTF16PtrToString(st.ProductName)
	g.ProductSuite = windows.UTF16PtrToString(st.ProductSuite)
	g.MajorVersion = st.MajorVersion
	g.MinorVersion = st.MinorVersion
	g.Build = st.Build
	g.SpBuild = st.SpBuild
	g.SpLevel = st.SpLevel
	g.Bootable = st.Bootable
	g.SystemRoot = windows.UTF16PtrToString(st.SystemRoot)
	var language []string
	if st.Language != nil && st.LanguageCount != 0 {
		stSize := unsafe.Sizeof(*st.Language)
		stPtr := st.Language

		for i := uint32(0); i < st.LanguageCount; i++ {
			language = append(language, windows.UTF16PtrToString(stPtr.Value))

			stPtr = (*DismLanguage)(unsafe.Add(unsafe.Pointer(stPtr), stSize))
		}
	}
	g.Language = language
	g.DefaultLanguage = g.Language[st.DefaultLanguageIndex]
	if st.CustomizedInfo != nil {
		var goSt GoDismWimCustomizedInfo
		goSt.fill((*DismWimCustomizedInfo)(st.CustomizedInfo))

		g.CustomizedInfo = &goSt
	}
}

type GoDismMountedImageInfo struct {
	MountPath     string
	ImageFilePath string
	ImageIndex    uint32
	MountMode     DismMountMode
	MountStatus   DismMountStatus
}

func (g *GoDismMountedImageInfo) fill(st *DismMountedImageInfo) {
	g.MountPath = windows.UTF16PtrToString(st.MountPath)
	g.ImageFilePath = windows.UTF16PtrToString(st.ImageFilePath)
	g.ImageIndex = st.ImageIndex
	g.MountMode = st.MountMode
	g.MountStatus = st.MountStatus
}

type GoDismDriverPackage struct {
	PublishedName    string
	OriginalFileName string
	Inbox            bool
	CatalogFile      string
	ClassName        string
	ClassGuid        string
	ClassDescription string
	BootCritical     bool
	DriverSignature  DismDriverSignature
	ProviderName     string
	Date             time.Time
	MajorVersion     uint32
	MinorVersion     uint32
	Build            uint32
	Revision         uint32
}

func (g *GoDismDriverPackage) fill(st *DismDriverPackage) {
	g.PublishedName = windows.UTF16PtrToString(st.PublishedName)
	g.OriginalFileName = windows.UTF16PtrToString(st.OriginalFileName)
	g.Inbox = st.Inbox != 0
	g.CatalogFile = windows.UTF16PtrToString(st.CatalogFile)
	g.ClassName = windows.UTF16PtrToString(st.ClassName)
	g.ClassGuid = windows.UTF16PtrToString(st.ClassGuid)
	g.ClassDescription = windows.UTF16PtrToString(st.ClassDescription)
	g.BootCritical = st.BootCritical != 0
	g.DriverSignature = st.DriverSignature
	g.ProviderName = windows.UTF16PtrToString(st.ProviderName)
	g.Date = time.Date(int(st.Date.Year), time.Month(st.Date.Month), int(st.Date.Day), int(st.Date.Hour), int(st.Date.Minute), int(st.Date.Second), int(st.Date.Milliseconds)*1000000, time.Local)
	g.MajorVersion = st.MajorVersion
	g.MinorVersion = st.MinorVersion
	g.Build = st.Build
	g.Revision = st.Revision
}

type GoDismDriver struct {
	ManufacturerName    string
	HardwareDescription string
	HardwareId          string
	Architecture        uint32
	ServiceName         string
	CompatibleIds       string
	ExcludeIds          string
}

func (g *GoDismDriver) fill(st *DismDriver) {
	g.ManufacturerName = windows.UTF16PtrToString(st.ManufacturerName)
	g.HardwareDescription = windows.UTF16PtrToString(st.HardwareDescription)
	g.HardwareId = windows.UTF16PtrToString(st.HardwareId)
	g.Architecture = st.Architecture
	g.ServiceName = windows.UTF16PtrToString(st.ServiceName)
	g.CompatibleIds = windows.UTF16PtrToString(st.CompatibleIds)
	g.ExcludeIds = windows.UTF16PtrToString(st.ExcludeIds)
}

type GoDismAppxPackage struct {
	PackageName     string
	DisplayName     string
	PublisherId     string
	MajorVersion    uint32
	MinorVersion    uint32
	Build           uint32
	RevisionNumber  uint32
	Architecture    uint32
	ResourceId      string
	InstallLocation string
	Region          string
}

func (g *GoDismAppxPackage) fill(st *DismAppxPackage) {
	g.PackageName = windows.UTF16PtrToString(st.PackageName)
	g.DisplayName = windows.UTF16PtrToString(st.DisplayName)
	g.PublisherId = windows.UTF16PtrToString(st.PublisherId)
	g.MajorVersion = st.MajorVersion
	g.MinorVersion = st.MinorVersion
	g.Build = st.Build
	g.RevisionNumber = st.RevisionNumber
	g.Architecture = st.Architecture
	g.ResourceId = windows.UTF16PtrToString(st.ResourceId)
	g.InstallLocation = windows.UTF16PtrToString(st.InstallLocation)
	g.Region = windows.UTF16PtrToString(st.Region)
}

var dismErrMsg = [...]string{
	"The DISM session needs to be reloaded.",
	"DISM API was not initialized for this process",
	"A DismSession was being shutdown when another operation was called on it",
	"A DismShutdown was called while there were open DismSession handles",
	"An invalid DismSession handle was passed into a DISMAPI function",
	"An invalid image index was specified",
	"An invalid image name was specified",
	"An image that is not a mounted WIM or mounted VHD was attempted to be unmounted",
	"Failed to gain access to the log file user specified. Logging has been disabled..",
	"A DismSession with open handles was attempted to be unmounted",
	"A DismSession with open handles was attempted to be mounted",
	"A DismSession with open handles was attempted to be remounted",
	"One or several parent features are disabled so current feature can not be enabled.",
	"The offline image specified is the running system. DISM_ONLINE_IMAGE must be used instead.",
	"The specified product key could not be validated. Check that the specified product key is valid and that it matches the target edition.",
	"An image file must be specified with either an index or a name.",
	"The image needs to be remounted before any servicing operation.",
	"The feature is not present in the package.",
	"The current package and feature servicing infrastructure is busy.  Wait a bit and try the operation again.",
}

func getIndex(e DismErr) uint32 {
	switch {
	case e == DISMAPI_S_RELOAD_IMAGE_SESSION_REQUIRED:
		return 0
	case e > 0xc0040000 && e < 0xc0040008:
		return uint32(e - 0xc0040000)
	case e > 0xc0040008 && e < 0xc0040010:
		return uint32(e - 0xc0040001)
	case e == DISMAPI_E_MUST_SPECIFY_INDEX_OR_NAME:
		return 15
	case e == DISMAPI_E_NEEDS_REMOUNT:
		return 16
	case e == DISMAPI_E_UNKNOWN_FEATURE:
		return 17
	case e == DISMAPI_E_BUSY:
		return 18
	}

	return 0xffffffff
}

// implement error to show error message
func (e DismErr) Error() string {
	i := getIndex(e)

	if int32(i) == -1 {
		return "unknown error"
	}

	return dismErrMsg[i]
}

func dismErr(e error) error {
	errno, ok := e.(windows.Errno)
	if !ok {
		return e
	}

	if int32(getIndex(DismErr(errno))) != -1 {
		return DismErr(errno)
	}

	return e
}
