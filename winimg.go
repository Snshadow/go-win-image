// Package winimg provides function to access and manupulate Windows image(.wim, .esd) files.

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
	// ignore DismInitialize called more than once in one process
	dismInitialized bool
	// default to information level
	dismLogLevel = dismapi.DismLogErrorsWarningsInfo
	// use working directory by default
	dismLogPath = ".\\dism.log"
	// use api default setting
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

type DismImageFile struct {
	// path of .wim or .vhd(x) file
	imageFilePath string
	// mounted paths uses in [dismapi.DismMountImage]
	MountPaths mapset.Set[string]
	// DismSessions from [dismapi.DismOpenSession], mapped with mounted paths
	Sessions mapset.Set[dismapi.DismSession]
}

// NewDismImageFile creates structure for using DISM API, it
// runs [dismapi.DismInitialize] if api is not initialized.
func NewDismImageFile(imageFilePath string) (*DismImageFile, error) {
	if err := initDism(); err != nil {
		return nil, err
	}

	dismImg := &DismImageFile{
		imageFilePath: imageFilePath,
		MountPaths:    mapset.NewSet[string](),
		Sessions:      mapset.NewSet[dismapi.DismSession](),
	}

	curDismImg.Add(dismImg)

	return dismImg, nil
}

func (d *DismImageFile) ImageFilePath() string {
	return d.imageFilePath
}

// Mount mounts an image file with imageIndex or with imageName if specified.
//
// flags: [dismapi.DISM_MOUNT_READWRITE]...
func (d *DismImageFile) Mount(mountPath string, imageIndex uint32, imageName string, flags uint32) error {
	identifier := dismapi.DismImageIndex
	if imageName != "" {
		identifier = dismapi.DismImageName
	}

	if err := dismapi.DismMountImage(d.imageFilePath, mountPath, imageIndex, imageName,
		identifier, flags, 0, 0, nil); err != nil {
		return err
	}

	d.MountPaths.Add(mountPath)

	return nil
}

// Unmount unmounts the specified mount path, if commit is false,
// append and integrity(generate integrity) are ignored.
func (d *DismImageFile) Unmount(mountPath string, commit, append, integrity bool) error {
	if !d.MountPaths.ContainsOne(mountPath) {
		return ErrNotMounted
	}

	var flags uint32 = dismapi.DISM_DISCARD_IMAGE
	if commit {
		flags = dismapi.DISM_COMMIT_IMAGE
		if append {
			flags |= dismapi.DISM_COMMIT_APPEND
		}
		if integrity {
			flags |= dismapi.DISM_COMMIT_GENERATE_INTEGRITY
		}
	}

	if err := dismapi.DismUnmountImage(mountPath, flags, 0, 0, nil); err != nil {
		return err
	}

	d.MountPaths.Remove(mountPath)

	return nil
}

// OpenSession opens DismSession from a mount path, if windowsPath or
// systemDrive are empty, default values will be used.
func (d *DismImageFile) OpenSession(mountPath, windowsPath, systemDrive string) (dismapi.DismSession, error) {
	if !d.MountPaths.ContainsOne(mountPath) {
		return 0, ErrNotMounted
	}

	ses, err := dismapi.DismOpenSession(mountPath, windowsPath, systemDrive)
	if err != nil {
		return 0, err
	}

	d.Sessions.Add(ses)

	return ses, nil
}

// CloseSession closes the existing DismSession.
func (d *DismImageFile) CloseSession(session dismapi.DismSession) error {
	if !d.Sessions.ContainsOne(session) {
		return ErrDismSessionNotExist
	}

	if err := dismapi.DismCloseSession(session); err != nil {
		return err
	}

	d.Sessions.Remove(session)

	return nil
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
		closeErr := dismapi.DismCloseSession(session)
		if closeErr != nil {
			err = errors.Join(err, closeErr)
		} else {
			d.Sessions.Remove(session)
		}
	}

	for mntPath := range d.MountPaths.Iter() {
		unmountErr := dismapi.DismUnmountImage(mntPath, dismapi.DISM_DISCARD_IMAGE, 0, 0, nil)
		if unmountErr != nil {
			err = errors.Join(err, unmountErr)
		} else {
			d.MountPaths.Remove(mntPath)
		}
	}

	curDismImg.Remove(d)

	if dismInitialized && curDismImg.IsEmpty() {
		// shutdown DISM API if it is no longer used
		err = errors.Join(err, dismapi.DismShutdown())
	}

	return err
}

type WIMImageHandle struct {
	Handle windows.Handle
	// if not empty, the image handle is mounted
	MountPath string

	imgMapRef *xsync.MapOf[uint32, *WIMImageHandle]
}

// Apply applies an image in .wim file to specified path.
//
// flags: [wimgapi.WIM_FLAG_VERIFY], etc...
func (w *WIMImageHandle) Apply(path string, flags uint32) error {
	if err := wimgapi.WIMApplyImage(w.Handle, path, flags); err != nil {
		return err
	}

	return nil
}

// Mount mounts an image in .wim file to mountPath.
//
// flags: [wimgapi.WIM_FLAG_MOUNT_READONLY], etc...
func (w *WIMImageHandle) Mount(mountPath string, flags uint32) error {
	if err := wimgapi.WIMMountImageHandle(w.Handle, mountPath, flags); err != nil {
		return err
	}

	w.MountPath = mountPath

	return nil
}

// Unmounnt unmounts an image in .wim file from mounted path.
func (w *WIMImageHandle) Unmount() error {
	if w.MountPath == "" {
		return ErrNotMounted
	}

	if err := wimgapi.WIMUnmountImageHandle(w.Handle, 0); err != nil {
		return err
	}

	w.MountPath = ""

	return nil
}

// Close closes .wim image handle, unmounts it if required.
func (w *WIMImageHandle) Close() error {
	var err error

	if w.MountPath != "" {
		err = errors.Join(err, w.Unmount())
	}

	err = errors.Join(err, wimgapi.WIMCloseHandle(w.Handle))

	return err
}

type WIMImageFile struct {
	// path of .wim file
	imageFilePath string
	// handle from WIMCreateFile with .wim file
	ImageFileHandle windows.Handle
	// handles from [wimgapi.WIMLoadImage] and [wimgapi.WIMCaptureImage],
	// mapped with image index
	ImageHandles *xsync.MapOf[uint32, *WIMImageHandle]
	// mounted paths from [wimgapi.WIMMountImage], mapped with image index
	MountPaths *xsync.MapOf[uint32, string]
}

// NewWIMImageFile opens or creates image file.
//
// access: [wimgapi.WIM_GENERIC_READ]... /
// createmode: [wimgapi.WIM_OPEN_EXISTING]... /
// flags: [wimgapi.WIM_FLAG_VERIFY]... /
// compressionType: [wimgapi.WIM_COMPRESS_LZX]...
func NewWIMImageFile(imageFilePath string, access uint32, createMode uint32, flags uint32,
	compressionType uint32) (*WIMImageFile, error) {
	wimHandle, _, err := wimgapi.WIMCreateFile(imageFilePath, access, createMode, flags, compressionType, false)
	if err != nil {
		return nil, err
	}

	return &WIMImageFile{
		imageFilePath:   imageFilePath,
		ImageFileHandle: wimHandle,
		ImageHandles:    xsync.NewMapOf[uint32, *WIMImageHandle](),
		MountPaths:      xsync.NewMapOf[uint32, string](),
	}, nil
}

// ImageFilePath returns path of the image(.wim) file.
func (w *WIMImageFile) ImageFilePath() string {
	return w.imageFilePath
}

func (w *WIMImageFile) GetImageInfo() ([]byte, error) {
	return wimgapi.WIMGetImageInformation(w.ImageFileHandle)
}

// LoadImage loads volume image from .wim file with image index.
func (w *WIMImageFile) LoadImage(imageIndex uint32) (*WIMImageHandle, error) {
	hnd, err := wimgapi.WIMLoadImage(w.ImageFileHandle, imageIndex)
	if err != nil {
		return nil, err
	}

	imgHandle := &WIMImageHandle{
		Handle:    hnd,
		imgMapRef: w.ImageHandles,
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

	w.ImageHandles.Range(func(index uint32, imgHandle *WIMImageHandle) bool {
		if imgHandle.MountPath != "" {
			unmountErr := wimgapi.WIMUnmountImageHandle(imgHandle.Handle, 0)
			if unmountErr != nil {
				err = errors.Join(err, unmountErr)
			} else {
				imgHandle.MountPath = ""
			}
		}

		closeErr := wimgapi.WIMCloseHandle(imgHandle.Handle)
		if closeErr != nil {
			err = errors.Join(err, closeErr)
		} else {
			w.ImageHandles.Delete(index)
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

	if w.ImageFileHandle != 0 {
		err = errors.Join(wimgapi.WIMCloseHandle(w.ImageFileHandle))
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
