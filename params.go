package winimg

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

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
	// value returned from [windows.NewCallback]
	// with [dismapi.DismProgressCallback]
	Progress uintptr
	// userData passed to [dismapi.DismProgressCallback]
	UserData unsafe.Pointer
}

// DismMountOpts contains options used for DISM image mount.
type DismMountOpts struct {
	MountPath string
	// not used if ImageName is not empty
	ImageIndex                         uint32
	ImageName                          string
	ReadOnly, Optimize, CheckIntegrity bool
	DismProgressOpts
}

// DismUnmountOpts contains options used for DISM image unmount.
type DismUnmountOpts struct {
	MountPath string
	// Append or GenerateIntegrity are ignored if
	// Commit is false
	Commit, Append, GenerateIntegrity bool
	DismProgressOpts
}
