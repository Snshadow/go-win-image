package dismapi

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

//sys	DismInitialize (LogLevel DismLogLevel, LogFilePath *uint16, ScratchDirectory *uint16) (err error) = dismapi.DismIntialize
//sys	DismShutdown() (err error) = dismapi.DismShutdown
//sys	DismMountImage(ImageFilePath *uint16, MountPath *uint16, ImageIndex uint32, ImageName *uint16, ImageIdentifier DismImageIdentifier, Flags uint32, CancelEvent Handle, Progress uintptr, UserData unsafe.Pointer) (err error) = dismapi.DismMountImage
//sys	DismUnmountImage(MountPath *uint16, Flags uint32, CancelEvent Handle, Progress uintptr, UserData unsafe.Pointer) (err error) = dismapi.DismUnmountImage
//sys	DismOpenSession(ImagePath *uint16, WindowsDirectory *uint16, SystemDrive *uint16, Session *DismSession) (err error) = dismapi.DismOpenSession
//sys	DismCloseSession(Session DismSession) (err error) = dismapi.DismCloseSession
//sys	DismGetLastErrorMessage(ErrorMessage **DismString) (err error) = dismapi.DismGetLastErrorMessage
//sys	DismRemountImage(MountPath *uint16) (err error) = dismapi.DismRemountImage
//sys	DismCommitImage(Session DismSession, Flags uint32, CancelEvent Handle, Progress uintptr, UserData unsafe.Pointer) (err error) = dismapi.DismCommitImage
//sys	DismGetImageInfo(ImageFilePath *uint16, ImageInfo **DismImageInfo, Count *uint32) (err error) = dismapi.DismGetImageInfo
//sys	DismGetMountedImageInfo(MountedImageInfo **DismMountedImageInfo, Count *uint32) (err error) = dismapi.DismGetMountedImageInfo
//sys	DismCleanupMountpoints() (err error) = dismapi.DismCleanupMountpoints
//sys	DismCheckImageHealth(Session DismSession, ScanImage bool, CancelEvent Handle, Progress uintptr, UserData unsafe.Pointer, ImageHealth *DismImageHealthState) (err error) = dismapi.DismCheckImageHealth
//sys	DismRestoreImageHealth(Session DismSession, SourcePaths **uint16, SourcePathCount uint32, LimitAccess bool, CancelEvent Handle, Progress uintptr, UserData unsafe.Pointer) (err error) = dismapi.DismRestoreImageHealth
//sys	DismDelete(DismStructure unsafe.Pointer) (err error) = dismapi.DismDelete
//sys	DismAddPackage(Session DismSession, PackagePath *uint16, IgnoreCheck bool, PreventPending bool, CancelEvent Handle, Progress uintptr, UserData unsafe.Pointer) (err error) = dismapi.DismAddPackage
//sys	DismRemovePackage(Session DismSession, Identifier *uint16, PackageIdentifier DismPackageIdentifier, CancelEvent Handle, Progress uintptr, UserData unsafe.Pointer) (err error) = dismapi.DismRemovePackage
//sys	DismEnableFeature(Session DismSession, FeatureName *uint16, Identifier *uint16, PackageIdentifier DismPackageIdentifier, LimitAccess bool, SourcePaths **uint16, SourcePathCount uint32, EnableAll bool, CancelEvent Handle, Progress uintptr, UserData unsafe.Pointer) (err error) = dismapi.DismEnableFeature
//sys	DismDisableFeature(Session DismSession, FeatureName *uint16, PackageName *uint16, RemovePayload bool, CancelEvent Handle, Progress uintptr, UserData unsafe.Pointer) (err error) = dismapi.DismDisableFeature
//sys	DismGetPackages(Session DismSession, Package **DismPackage, Count *uint32) (err error) = dismapi.DismGetPackages
//sys	DismGetPackageInfo(Session DismSession, Identifier *uint16, PackageIdentifier DismPackageIdentifier, PackageInfo **DismPackageInfo) (err error) = dismapi.DismGetPackageInfo
//sys	DismGetPackageInfoEx(Session DismSession, Identifier *uint16, PackageIdentifier DismPackageIdentifier, PackageInfoEx **DismPackageInfoEx) (err error) = dismapi.DismGetPackageInfoEx
//sys	DismGetFeatures(Session DismSession, Identifier *uint16, PackageIdentifier DismPackageIdentifier, Feature **DismFeature, Count *uint32) (err error) = dismapi.DismGetFeatures
//sys	DismGetFeatureInfo(Session DismSession, FeatureName *uint16, Identifier *uint16, PackageIdentifier DismPackageIdentifier, FeatureInfo **DismFeatureInfo) (err error) = dismapi.DismGetFeatureInfo
//sys	DismGetFeatureParent(Session DismSession, FeatureName *uint16, Identifier *uint16, PackageIdentifier DismPackageIdentifier, Feature **DismFeature, Count *uint32) (err error) = dismapi.DismGetFeatureParent
//sys	DismApplyUnattend(Session DismSession, UnattendFile *uint16, SingleSession bool) (err error) = dismapi.DismApplyUnattend
//sys	DismAddDriver(Session DismSession, DriverPath *uint16, ForceUnsigned bool) (err error) = dismapi.DismAddDriver
//sys	DismRemoveDriver(Session DismSession, DriverPath *uint16) (err error) = dismapi.DismRemoveDriver
//sys	DismGetDrivers(Session DismSession, AllDrivers bool, DriverPackage **DismDriverPackage, Count *uint32) (err error) = dismapi.DismGetDrivers
//sys	DismGetDriverInfo(Session DismSession, DriverPath *uint16, Driver **DismDriver, Count *uint32, DriverPackage **DismDriverPackage) (err error) = dismapi.DismGetDriverInfo
//sys	DismGetCapabilities(Session DismSession, Capability **DismCapability, Count *uint32) (err error) = dismapi.DismGetCapabilities
//sys	DismGetCapabilityInfo(Session DismSession, Name *uint16, Info **DismCapabilityInfo) (err error) = dismapi.DismGetCapabilityInfo
//sys	DismAddCapability(Session DismSession, Name *uint16, LimitAccess bool, SourcePaths **uint16, SourcePathCount uint32, CancelEvent Handle, Progress uintptr, UserData unsafe.Pointer) (err error) = dismapi.DismAddCapability
//sys	DismRemoveCapability(Session DismSession, Name *uint16, CancelEvent Handle, Progress uintptr, UserData unsafe.Pointer) (err error) = dismapi.DismRemoveCapability
