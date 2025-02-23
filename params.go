package winimg

import (
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/Snshadow/winimg/w32api/dismapi"
	"github.com/Snshadow/winimg/w32api/wimgapi"
)

// NewDismProgress creates callback function to be
// used get information of DISM operation progress.
func NewDismProgress(cb dismapi.DismProgressCallback) uintptr {
	return windows.NewCallback(cb)
}

// DismProgressOpts contains optional parameters to
// cancel or get information while DISM operation
// is in progress.
type DismProgressOpts struct {
	// create with [windows.CreateEvent] and
	// trigger with [windows.SetEvent]
	//
	// call [windows.ResetEvent] to use again
	// for other DISM operation
	//
	// close with [windows.CloseHandle]
	CancelEvent windows.Handle
	// callback function returned from [NewDismProgress]
	Progress uintptr
	// userData passed to [dismapi.DismProgressCallback]
	UserData unsafe.Pointer
}

// DismMountOpts contains options used for mounting DISM image.
type DismMountOpts struct {
	MountPath string
	// not used if ImageName is not empty
	ImageIndex                         uint32
	ImageName                          string
	ReadOnly, Optimize, CheckIntegrity bool
	DismProgressOpts
}

// DismUnmountOpts contains options used for unmounting DISM image.
type DismUnmountOpts struct {
	MountPath string
	// Append, GenerateIntegrity or SupportEa are ignored if
	// Commit is false
	Commit, Append, GenerateIntegrity, SupportEa bool
	DismProgressOpts
}

// DismCommitOpts contains options used for saving changes in
// DISM image.
type DismCommitOpts struct {
	Append, GenerateIntegrity, SupportEa bool
	DismProgressOpts
}

// DismAddCapabilityOpts contains options used for adding
// capability to DISM image.
type DismAddCapabilityOpts struct {
	Name string
	// do not use WU/WSUS Update
	LimitAccess bool
	SourcePaths []string
	DismProgressOpts
}

// DismAddPackageOpts contains options used for adding
// package(.cab, .msu, .mum file) to DISM image.
type DismAddPackageOpts struct {
	PackagePath    string
	IgnoreCheck    bool
	PreventPending bool
	DismProgressOpts
}

// DismDisableFeatureOpts contains options used for
// disabling or removing feature in a DISM image.
type DismDisableFeatureOpts struct {
	FeatureName   string
	PackageName   string
	RemovePayload bool
	DismProgressOpts
}

// DismEnableFeatureOpts contains options used for
// enabling feature in a DISM image.
type DismEnableFeatureOpts struct {
	// name of fearure, separate multiple names with semicolon
	FeatureName string
	// optional, path to .cab file or package name
	Identifier string
	// do not use WU
	LimitAccess bool
	SourcePaths []string
	EnableAll   bool
	DismProgressOpts
}

// DismRestoreImageHealthOpts contains options used
// for repairing a corrupt image.
type DismRestoreImageHealthOpts struct {
	SourcePaths []string
	LimitAccess bool
	DismProgressOpts
}

// NewWimMessageCallback creates callback function to
// be used with [wimgapi.WIMRegisterMessageCallback] or
// [wimgapi.WIMUnregisterMessageCallback].
func NewWimMessageCallback(cb wimgapi.WIMMessageCallback) uintptr {
	return windows.NewCallback(cb)
}
