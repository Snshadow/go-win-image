package winimg

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

// DismProgressOpts contains optional parameters to
// cancel or get information while DISM operation
// is in progress.
type DismProgressOpts struct {
	// created with [windows.CreateEvent] and
	// triggered with [windows.SetEvent]
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

type DismMountOpts struct {
	MountPath string
	ImageIndex uint32
	ImageName string
	ReadOnly, Optimize, CheckIntegrity bool
	DismProgressOpts
}

type DismUnmountOpts struct {
	MountPath string
	// Append or GenerateIntegrity are used it Commit
	// is true
	Commit, Append, GenerateIntegrity bool
	DismProgressOpts
}
