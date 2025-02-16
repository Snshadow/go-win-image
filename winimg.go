// Package winimg provides functions to access and manupulate Windows image(.wim, .esd) files.
package winimg

import (
	"errors"

	"golang.org/x/sys/windows"

	"github.com/Snshadow/winimg/w32api/dismapi"
	"github.com/Snshadow/winimg/w32api/wimgapi"
	"github.com/deckarep/golang-set/v2"
	"github.com/puzpuzpuz/xsync/v3"
)

var (
	ErrNotMounted          = errors.New("the path is not a mount path")
	ErrDismSessionNotExist = errors.New("the specified DismSession does not exist")
)

var (
	// do not run DismInitialize when it is already initialized for a process
	dismInitialized bool
	// default to information level
	dismLogLevel = dismapi.DismLogErrorsWarningsInfo
	// use working directory by default
	dismLogPath = ".\\dism.log"
	// use DISM API default path
	dismScratchDir = ""

	// track current use of DISM API
	curDismImg = mapset.NewSet[*DismImageFile]()
)

// ConfigureDism configures global settings for DISM API for process.
func ConfigureDism(logLevel dismapi.DismLogLevel, logPath, scratchDir string) {
	dismLogLevel = logLevel
	dismLogPath = logPath
	dismScratchDir = scratchDir
}

// initDism initalizes DISM API with global settings.
func initDism() error {
	if !dismInitialized {
		err := dismapi.DismInitialize(dismLogLevel, dismLogPath, dismScratchDir)
		// returns DISMAPI_S_RELOAD_IMAGE_SESSION_REQUIRED if already initialized
		if err != nil && err != dismapi.DISMAPI_S_RELOAD_IMAGE_SESSION_REQUIRED {
			return err
		}

		dismInitialized = true
	}

	return nil
}

// DismImageSession handles operations related with DismSession.
type DismImageSession struct {
	Session   dismapi.DismSession
	MountPath string

	sesMapRef mapset.Set[*DismImageSession]
}

// Close closes opened DismSession for [DismImageSession].
func (d *DismImageSession) Close() error {
	if err := dismapi.DismCloseSession(d.Session); err != nil {
		return err
	}

	d.sesMapRef.Remove(d)

	return nil
}

// DismImageFile stores mount points and DismSession
// for Windows image file.
type DismImageFile struct {
	// path of .wim or .vhd(x) file
	imageFilePath string
	// mounted paths used in [dismapi.DismMountImage]
	MountPoints mapset.Set[string]
	// DismSessions from [dismapi.DismOpenSession], mapped with mounted paths
	Sessions mapset.Set[*DismImageSession]
}

// NewDismImageFile creates structure for using DISM API, it
// runs [dismapi.DismInitialize] if api is not initialized.
func NewDismImageFile(imageFilePath string) (*DismImageFile, error) {
	if err := initDism(); err != nil {
		return nil, err
	}

	dismImg := &DismImageFile{
		imageFilePath: imageFilePath,
		MountPoints:   mapset.NewSet[string](),
		Sessions:      mapset.NewSet[*DismImageSession](),
	}

	curDismImg.Add(dismImg)

	return dismImg, nil
}

func (d *DismImageFile) ImageFilePath() string {
	return d.imageFilePath
}

// Mount mounts an image file with imageIndex or with imageName if specified.
func (d *DismImageFile) Mount(opts DismMountOpts) error {
	identifier := dismapi.DismImageIndex
	if opts.ImageName != "" {
		identifier = dismapi.DismImageName
	}

	var flags uint32 = dismapi.DISM_MOUNT_READWRITE
	if opts.ReadOnly {
		flags = dismapi.DISM_MOUNT_READONLY
	}
	if opts.Optimize {
		flags |= dismapi.DISM_MOUNT_OPTIMIZE
	}
	if opts.CheckIntegrity {
		flags |= dismapi.DISM_MOUNT_CHECK_INTEGRITY
	}

	if err := dismapi.DismMountImage(d.imageFilePath, opts.MountPath, opts.ImageIndex, opts.ImageName,
		identifier, flags, 0, 0, nil); err != nil {
		return err
	}

	d.MountPoints.Add(opts.MountPath)

	return nil
}

// Unmount unmounts the specified mount path.
func (d *DismImageFile) Unmount(opts DismUnmountOpts) error {
	if !d.MountPoints.ContainsOne(opts.MountPath) {
		return ErrNotMounted
	}

	var flags uint32 = dismapi.DISM_DISCARD_IMAGE
	if opts.Commit {
		flags = dismapi.DISM_COMMIT_IMAGE
		if opts.Append {
			flags |= dismapi.DISM_COMMIT_APPEND
		}
		if opts.GenerateIntegrity {
			flags |= dismapi.DISM_COMMIT_GENERATE_INTEGRITY
		}
	}

	if err := dismapi.DismUnmountImage(opts.MountPath, flags, 0, 0, nil); err != nil {
		return err
	}

	d.MountPoints.Remove(opts.MountPath)

	return nil
}

// OpenSession creates DismImageSession from a mount path, if windowsPath or
// systemDrive are empty, default values will be used.
func (d *DismImageFile) OpenSession(mountPath, windowsPath, systemDrive string) (*DismImageSession, error) {
	if !d.MountPoints.ContainsOne(mountPath) {
		return nil, ErrNotMounted
	}

	ses, err := dismapi.DismOpenSession(mountPath, windowsPath, systemDrive)
	if err != nil {
		return nil, err
	}

	newSession := &DismImageSession{
		Session:   ses,
		MountPath: mountPath,
		sesMapRef: d.Sessions,
	}

	d.Sessions.Add(newSession)

	return newSession, nil
}

// GetImageInfo returns information of images in .wim or .vhd(x) file.
func (d *DismImageFile) GetImageInfo() ([]dismapi.GoDismImageInfo, error) {
	return dismapi.DismGetImageInfo(d.imageFilePath)
}

// Close cleans up all DismSessions and mount points associated with
// the image file, note that this will discard changes in mount points.
func (d *DismImageFile) Close() error {
	var err error

	for session := range d.Sessions.Iter() {
		closeErr := session.Close()
		if closeErr != nil {
			err = errors.Join(err, closeErr)
		}
	}

	for mntPath := range d.MountPoints.Iter() {
		unmountErr := dismapi.DismUnmountImage(mntPath, dismapi.DISM_DISCARD_IMAGE, 0, 0, nil)
		if unmountErr != nil {
			err = errors.Join(err, unmountErr)
		} else {
			d.MountPoints.Remove(mntPath)
		}
	}

	curDismImg.Remove(d)

	if dismInitialized && curDismImg.IsEmpty() {
		// shutdown DISM API if it is no longer used
		shutdownErr := dismapi.DismShutdown()
		if shutdownErr != nil && shutdownErr != dismapi.DISMAPI_E_DISMAPI_NOT_INITIALIZED {
			err = errors.Join(err, shutdownErr)
		} else {
			dismInitialized = false
		}
	}

	return err
}

type WIMVolumeImage struct {
	Handle windows.Handle
	// if not empty, the image handle is mounted
	MountPath string

	imgMapRef  *xsync.MapOf[uint32, *WIMVolumeImage]
	imageIndex uint32
}

// GetImageInfo returns xml info of a image in bytes.
func (w *WIMVolumeImage) GetImageInfo() ([]byte, error) {
	return wimgapi.WIMGetImageInformation(w.Handle)
}

// Apply applies a volume image in .wim file to specified path.
//
// flags: [wimgapi.WIM_FLAG_VERIFY], [wimgapi.WIM_FLAG_INDEX], [wimgapi.WIM_FLAG_NO_APPLY], [wimgapi.WIM_FLAG_FILEINFO], [wimgapi.WIM_FLAG_NO_RP_FIX], [wimgapi.WIM_FLAG_NO_DIRACL], [wimgapi.WIM_FLAG_NO_FILEACL]
func (w *WIMVolumeImage) Apply(path string, flags uint32) error {
	if err := wimgapi.WIMApplyImage(w.Handle, path, flags); err != nil {
		return err
	}

	return nil
}

// Export exports a volume image in .wim file to another .wim file.
//
// flags: [wimgapi.WIM_EXPORT_ALLOW_DUPLICATES]...
func (w *WIMVolumeImage) Export(wim *WIMImageFile, flags uint32) error {
	if err := wimgapi.WIMExportImage(w.Handle, wim.Handle, flags); err != nil {
		return err
	}

	return nil
}

// Mount mounts a volume image in .wim file to mountPath.
//
// flags: [wimgapi.WIM_FLAG_MOUNT_READONLY], [wimgapi.WIM_FLAG_VERIFY], [wimgapi.WIM_FLAG_NO_RP_FIX], [wimgapi.WIM_FLAG_NO_DIRACL], [wimgapi.WIM_FLAG_NO_FILEACL]
func (w *WIMVolumeImage) Mount(mountPath string, flags uint32) error {
	if err := wimgapi.WIMMountImageHandle(w.Handle, mountPath, flags); err != nil {
		return err
	}

	w.MountPath = mountPath

	return nil
}

// Unmount unmounts a volume image in .wim file from mounted path.
func (w *WIMVolumeImage) Unmount() error {
	if w.MountPath == "" {
		return ErrNotMounted
	}

	if err := wimgapi.WIMUnmountImageHandle(w.Handle, 0); err != nil {
		return err
	}

	w.MountPath = ""

	return nil
}

// Close closes a handle of volume image in .wim, unmounts it if required.
func (w *WIMVolumeImage) Close() error {
	var err error

	if w.MountPath != "" {
		err = errors.Join(err, w.Unmount())
	}

	closeErr := wimgapi.WIMCloseHandle(w.Handle)
	if closeErr != nil {
		err = errors.Join(err, closeErr)
	} else {
		w.imgMapRef.Delete(w.imageIndex)
	}

	return err
}

// WIMImageFile contains handle for .wim file and associated mount points.
type WIMImageFile struct {
	// path of .wim file
	imageFilePath string
	// handle from [wimgapi.WIMCreateFile] with .wim file
	Handle windows.Handle
	// handles from [wimgapi.WIMLoadImage] and [wimgapi.WIMCaptureImage],
	// mapped with image index
	ImageHandles *xsync.MapOf[uint32, *WIMVolumeImage]
	// mounted paths from [wimgapi.WIMMountImage], mapped with image index
	MountPaths *xsync.MapOf[uint32, string]
}

// NewWIMImageFile opens or creates image file.
//
// access: [wimgapi.WIM_GENERIC_READ]... /
// createmode: [wimgapi.WIM_OPEN_EXISTING]... /
// flags: [wimgapi.WIM_FLAG_VERIFY]... /
func NewWIMImageFile(imageFilePath string, access uint32, createMode uint32, flags uint32,
	compressionType wimgapi.WimCompressionType) (*WIMImageFile, error) {
	wimHandle, _, err := wimgapi.WIMCreateFile(imageFilePath, access, createMode, flags, compressionType, false)
	if err != nil {
		return nil, err
	}

	return &WIMImageFile{
		imageFilePath: imageFilePath,
		Handle:        wimHandle,
		ImageHandles:  xsync.NewMapOf[uint32, *WIMVolumeImage](),
		MountPaths:    xsync.NewMapOf[uint32, string](),
	}, nil
}

// ImageFilePath returns path of the image(.wim) file.
func (w *WIMImageFile) ImageFilePath() string {
	return w.imageFilePath
}

// GetFileInfo returns xml info of .wim file in bytes.
func (w *WIMImageFile) GetFileInfo() ([]byte, error) {
	return wimgapi.WIMGetImageInformation(w.Handle)
}

func (w *WIMImageFile) GetAttributes() (wimgapi.GoWimInfo, error) {
	return wimgapi.WIMGetAttributes(w.Handle)
}

// LoadImage loads volume image from .wim file with image index.
func (w *WIMImageFile) LoadImage(imageIndex uint32) (*WIMVolumeImage, error) {
	hnd, err := wimgapi.WIMLoadImage(w.Handle, imageIndex)
	if err != nil {
		return nil, err
	}

	imgHandle := &WIMVolumeImage{
		Handle:     hnd,
		imgMapRef:  w.ImageHandles,
		imageIndex: imageIndex,
	}

	w.ImageHandles.Store(imageIndex, imgHandle)

	return imgHandle, err
}

// Mount mounts an image with imageIndex, is tempPath is empty,
// the image will not be mounted for edits.
func (w *WIMImageFile) Mount(mountPath string, imageIndex uint32, tempPath string) error {
	err := wimgapi.WIMMountImage(mountPath, w.imageFilePath, imageIndex, tempPath)
	if err != nil {
		return err
	}

	w.MountPaths.Store(imageIndex, mountPath)

	return nil
}

// Unmount unmounts an image with mountPath and imageIndex, if commitChanges
// is true, changes in mounted directory will be saved to .wim file.
func (w *WIMImageFile) Unmount(mountPath string, imageIndex uint32, commitChanges bool) error {
	err := wimgapi.WIMUnmountImage(mountPath, w.imageFilePath, imageIndex, commitChanges)
	if err != nil {
		return err
	}

	w.MountPaths.Delete(imageIndex)

	return nil
}

// Close cleans up all mount points and handles, note that this will
// discard changes in mount points created with Mount().
func (w *WIMImageFile) Close() error {
	var err error

	w.ImageHandles.Range(func(index uint32, imgHandle *WIMVolumeImage) bool {
		if imgHandle.MountPath != "" {
			unmountErr := wimgapi.WIMUnmountImageHandle(imgHandle.Handle, 0)
			if unmountErr != nil {
				err = errors.Join(err, unmountErr)
			} else {
				imgHandle.MountPath = ""
			}
		}

		closeErr := imgHandle.Close()
		if closeErr != nil {
			err = errors.Join(err, closeErr)
		}

		return true
	})

	w.MountPaths.Range(func(index uint32, mountPath string) bool {
		unmountErr := wimgapi.WIMUnmountImage(mountPath, w.imageFilePath, index, false)
		if unmountErr != nil {
			err = errors.Join(err, unmountErr)
		} else {
			w.MountPaths.Delete(index)
		}

		return true
	})

	if w.Handle != 0 {
		err = errors.Join(wimgapi.WIMCloseHandle(w.Handle))
	}

	return err
}

type WinImage struct {
	filePath string

	DismImage *DismImageFile
	WimImage  *WIMImageFile
}

func (w *WinImage) Path() string {
	return w.filePath
}
