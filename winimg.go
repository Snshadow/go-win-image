// Package winimg provides functions to access and manupulate Windows image(.wim, .esd) files.
package winimg

import (
	"errors"
	"os"
	"path/filepath"

	"golang.org/x/sys/windows"

	"github.com/Snshadow/winimg/w32api"
	"github.com/Snshadow/winimg/w32api/dismapi"
	"github.com/Snshadow/winimg/w32api/wimgapi"
	"github.com/deckarep/golang-set/v2"
	"github.com/puzpuzpuz/xsync/v3"
)

var (
	ErrNotMounted          = errors.New("the path is not a mount path")
	ErrAlreadyMounted      = errors.New("the image is already mounted")
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
		if err != nil && err.(*w32api.WinimgInternalErr).Errno() !=
			dismapi.DISMAPI_S_RELOAD_IMAGE_SESSION_REQUIRED {
			return err
		}

		dismInitialized = true
	}

	return nil
}

// DismImageSession handles operations related with DismSession.
type DismImageSession struct {
	session   dismapi.DismSession
	mountPath string

	sesMapRef mapset.Set[*DismImageSession]
}

// GetMountPath returns mount path being currently used.
func (d *DismImageSession) GetMountPath() string {
	return d.mountPath
}

// ApplyUnattend applies unattend answer file to
// the associated image with current session.
func (d *DismImageSession) ApplyUnattend(unattendFile string, singleSession bool) error {
	err := dismapi.DismApplyUnattend(d.session, unattendFile, singleSession)
	if err != nil {
		return err
	}

	return nil
}

// Close closes opened DismSession for [DismImageSession].
func (d *DismImageSession) Close() error {
	if err := dismapi.DismCloseSession(d.session); err != nil {
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
		identifier, flags, opts.CancelEvent, opts.Progress, opts.UserData); err != nil {
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

	if err := dismapi.DismUnmountImage(opts.MountPath, flags, opts.CancelEvent,
		opts.Progress, opts.UserData); err != nil {
		return err
	}

	d.MountPoints.Remove(opts.MountPath)

	return nil
}

// OpenSession creates DismImageSession from a mount path, if windowsPath or
// systemDrive are empty, default value will be used.
func (d *DismImageFile) OpenSession(mountPath, windowsPath, systemDrive string) (*DismImageSession, error) {
	if !d.MountPoints.ContainsOne(mountPath) {
		return nil, ErrNotMounted
	}

	ses, err := dismapi.DismOpenSession(mountPath, windowsPath, systemDrive)
	if err != nil {
		return nil, err
	}

	newSession := &DismImageSession{
		session:   ses,
		mountPath: mountPath,
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
		if shutdownErr != nil && shutdownErr.(*w32api.WinimgInternalErr).Errno() !=
			dismapi.DISMAPI_E_DISMAPI_NOT_INITIALIZED {
			err = errors.Join(err, shutdownErr)
		} else {
			dismInitialized = false
		}
	}

	return err
}

// RegisterWimLog registers logFile for logging wimgapi
// operations, flags is reserved and always 0.
func RegisterWimLog(logFile string, flags uint32) error {
	err := wimgapi.WIMRegisterLogFile(logFile, flags)
	if err != nil {
		return err
	}

	return nil
}

// UnregisterWimLog unregisters logfile for logging
// wimgapi operations used in [RegisterWimLog].
func UnregisterWimLog(logFile string) error {
	err := wimgapi.WIMUnregisterLogFile(logFile)
	if err != nil {
		return err
	}

	return nil
}

// WimVolumeImage stores handle and mount path for
// a volume image in .wim file.
type WimVolumeImage struct {
	Handle windows.Handle
	// if not empty, the image handle is mounted
	MountPath string

	imageIndex uint32
	wimRef     *WimImageFile
}

// GetImageInfo returns xml info of a image in bytes.
func (w *WimVolumeImage) GetImageInfo() ([]byte, error) {
	return wimgapi.WIMGetImageInformation(w.Handle)
}

// Apply applies a volume image in .wim file to specified path.
//
// flags: [wimgapi.WIM_FLAG_VERIFY], [wimgapi.WIM_FLAG_INDEX], [wimgapi.WIM_FLAG_NO_APPLY],
// [wimgapi.WIM_FLAG_FILEINFO], [wimgapi.WIM_FLAG_NO_RP_FIX], [wimgapi.WIM_FLAG_NO_DIRACL],
// [wimgapi.WIM_FLAG_NO_FILEACL]
func (w *WimVolumeImage) Apply(applyPath string, flags uint32) error {
	if !w.wimRef.tempPathSet {
		if err := w.wimRef.setTempPath(); err != nil {
			return err
		}
	}

	if err := wimgapi.WIMApplyImage(w.Handle, applyPath, flags); err != nil {
		return err
	}

	return nil
}

// Export transfers a volume image in .wim file to another .wim file.
//
// flags: [wimgapi.WIM_EXPORT_ALLOW_DUPLICATES]...
func (w *WimVolumeImage) Export(wim *WimImageFile, flags uint32) error {
	if !w.wimRef.tempPathSet {
		if err := w.wimRef.setTempPath(); err != nil {
			return err
		}
	}
	if !wim.tempPathSet {
		if err := wim.setTempPath(); err != nil {
			return err
		}
	}

	if err := wimgapi.WIMExportImage(w.Handle, wim.Handle, flags); err != nil {
		return err
	}

	return nil
}

// Mount mounts a volume image in .wim file to mountPath.
//
// flags: [wimgapi.WIM_FLAG_MOUNT_READONLY], [wimgapi.WIM_FLAG_VERIFY], [wimgapi.WIM_FLAG_NO_RP_FIX], [wimgapi.WIM_FLAG_NO_DIRACL], [wimgapi.WIM_FLAG_NO_FILEACL]
func (w *WimVolumeImage) Mount(mountPath string, flags uint32) error {
	if w.MountPath != "" {
		return ErrAlreadyMounted
	}

	if err := wimgapi.WIMMountImageHandle(w.Handle, mountPath, flags); err != nil {
		return err
	}

	w.MountPath = mountPath

	return nil
}

// Unmount unmounts a volume image in .wim file from mounted path.
func (w *WimVolumeImage) Unmount() error {
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
func (w *WimVolumeImage) Close() error {
	var err error

	if w.MountPath != "" {
		err = errors.Join(err, w.Unmount())
	}

	closeErr := wimgapi.WIMCloseHandle(w.Handle)
	if closeErr != nil {
		err = errors.Join(err, closeErr)
	} else {
		w.wimRef.ImageHandles.Delete(w.imageIndex)
	}

	return err
}

// WimImageFile contains handle for .wim file and associated mount points.
type WimImageFile struct {
	// handle from [wimgapi.WIMCreateFile] with .wim file
	Handle windows.Handle
	// handles from [wimgapi.WIMLoadImage] and [wimgapi.WIMCaptureImage],
	// mapped with image index
	ImageHandles *xsync.MapOf[uint32, *WimVolumeImage]

	// path of .wim file
	imageFilePath string
	// temporary path for capture and apply
	tempPath string
	// image count in .wim file
	imageCount uint32
	// is temporary path is set or created?
	tempPathSet, tempCreated bool
}

// NewWIMImageFile opens or creates image file, if tempPath is empty,
// temporary directory will be created with random name if needed.
//
// access: [wimgapi.WIM_GENERIC_READ]... /
// createmode: [wimgapi.WIM_OPEN_EXISTING]... /
// flags: [wimgapi.WIM_FLAG_VERIFY]... /
func NewWIMImageFile(imageFilePath string, access uint32, createMode uint32, flags uint32,
	compressionType wimgapi.WimCompressionType, tempPath string) (*WimImageFile, error) {
	wimHandle, _, err := wimgapi.WIMCreateFile(imageFilePath, access, createMode,
		flags, compressionType, false)
	if err != nil {
		return nil, err
	}

	return &WimImageFile{
		Handle:        wimHandle,
		ImageHandles:  xsync.NewMapOf[uint32, *WimVolumeImage](),
		imageFilePath: imageFilePath,
		imageCount:    wimgapi.WIMGetImageCount(wimHandle),
		tempPath:      tempPath,
	}, nil
}

// setTempPath sets temporary path used for WIM operation.
func (w *WimImageFile) setTempPath() error {
	if w.tempPath == "" {
		w.tempPath, _ = os.MkdirTemp("", filepath.Base(w.imageFilePath))
		w.tempCreated = true
	}

	err := wimgapi.WIMSetTemporaryPath(w.Handle, w.tempPath)
	if err != nil {
		return err
	}

	w.tempPathSet = true

	return nil
}

// GetImageCount updates and returns the number of
// volume images stored in .wim file.
func (w *WimImageFile) GetImageCount() uint32 {
	w.imageCount = wimgapi.WIMGetImageCount(w.Handle)

	return w.imageCount
}

// GetImageFilePath returns path of the image(.wim) file.
func (w *WimImageFile) GetImageFilePath() string {
	return w.imageFilePath
}

// Capture captures a directory path and stores it in an .wim
// file which is added as a volume image and returns it.
func (w *WimImageFile) Capture(capturePath string, flags uint32) (*WimVolumeImage, error) {
	if !w.tempPathSet {
		if err := w.setTempPath(); err != nil {
			return nil, err
		}
	}

	volHnd, err := wimgapi.WIMCaptureImage(w.Handle, capturePath, flags)
	if err != nil {
		return nil, err
	}

	w.imageCount++

	newVolume := &WimVolumeImage{
		Handle:     volHnd,
		imageIndex: w.imageCount,
		wimRef:     w,
	}

	w.ImageHandles.Store(w.imageCount, newVolume)

	return newVolume, nil
}

// GetFileInfo returns xml info of .wim file in bytes.
func (w *WimImageFile) GetFileInfo() ([]byte, error) {
	return wimgapi.WIMGetImageInformation(w.Handle)
}

func (w *WimImageFile) GetAttributes() (wimgapi.GoWimInfo, error) {
	return wimgapi.WIMGetAttributes(w.Handle)
}

// LoadImage loads volume image from .wim file with image index.
func (w *WimImageFile) LoadImage(imageIndex uint32) (*WimVolumeImage, error) {
	if !w.tempPathSet {
		if err := w.setTempPath(); err != nil {
			return nil, err
		}
	}

	hnd, err := wimgapi.WIMLoadImage(w.Handle, imageIndex)
	if err != nil {
		return nil, err
	}

	imgHandle := &WimVolumeImage{
		Handle:     hnd,
		imageIndex: imageIndex,
		wimRef:     w,
	}

	w.ImageHandles.Store(imageIndex, imgHandle)

	return imgHandle, err
}

// Close cleans up all mount points and volume image handles, note that
// this will discard changes in mount points created with Mount().
func (w *WimImageFile) Close() error {
	var err error

	w.ImageHandles.Range(func(index uint32, imgHandle *WimVolumeImage) bool {
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

	if w.Handle != 0 {
		err = errors.Join(err, wimgapi.WIMCloseHandle(w.Handle))
	}

	// remove temporary directory if it was created
	if w.tempCreated {
		err = errors.Join(err, os.RemoveAll(w.tempPath))
	}

	return err
}

// MountWimImage mounts an image with mount path, .wim file path and image
// index, is tempPath is empty, the image will not be mounted for edits.
// Cannot be used if [WimImageFile] exists for this .wim file.
func MountWimImage(mountPath, imageFilePath string, imageIndex uint32, tempPath string) error {
	err := wimgapi.WIMMountImage(mountPath, imageFilePath, imageIndex, tempPath)
	if err != nil {
		return err
	}

	return nil
}

// UnmountWimImage unmounts an image with mountPath and imageIndex, if
// commitChanges is true, changes in mounted directory will be saved to
// .wim file if tempPath was specified in [MountWimImage].
func UnmountWimImage(mountPath, imageFilePath string, imageIndex uint32, commitChanges bool) error {
	err := wimgapi.WIMUnmountImage(mountPath, imageFilePath, imageIndex, commitChanges)
	if err != nil {
		return err
	}

	return nil
}

// RemountWimImage remounts an image mount to mountPath, it maps
// contents if mounted image volume to the directory. The value is
// of flags is reserved and always zero.
func RemountWimImage(mountPath string, flags uint32) error {
	if err := wimgapi.WIMRemountImage(mountPath, flags); err != nil {
		return err
	}

	return nil
}

type WinImage struct {
	filePath string

	DismImage *DismImageFile
	WimImage  *WimImageFile
}

func (w *WinImage) Path() string {
	return w.filePath
}
