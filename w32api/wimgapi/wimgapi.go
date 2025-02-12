package wimgapi

import (
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/Snshadow/go_win_image/w32api/wintype"
)

//sys	wimCreateFile(wimPath *uint16, desiredAccess uint32, creationDisposition uint32, flagsAndAttributes uint32, compressionType uint32, creationResult *uint32) (handle windows.Handle, err error) [failretval==windows.InvalidHandle] = wimgapi.WIMCreateFile
//sys	WIMCloseHandle(object windows.Handle) (err error) = wimgapi.WIMCloseHandle
//sys	wimSetTemporaryPath(wim windows.Handle, path *uint16) (err error) = wimgapi.WIMSetTemporaryPath
//sys	wimSetReferenceFile(wim windows.Handle, path *uint16, flags uint32) (err error) = wimgapi.WIMSetReferenceFile
//sys	wimSplitFile(wim windows.Handle, partPath *uint16, partSize *int64, flags uint32) (err error) = wimgapi.WIMSplitFile
//sys	WIMExportImage(image windows.Handle, wim windows.Handle, flags uint32) (err error) = wimgapi.WIMExportImage
//sys	WIMDeleteImage(wim windows.Handle, imageIndex uint32) (err error) = wimgapi.WIMDeleteImage
//sys	WIMGetImageCount(wim windows.Handle) (count uint32) = wimgapi.WIMGetImageCount
//sys	wimGetAttributes(wim windows.Handle, wimInfo *WIM_INFO, cbWimInfo uint32) (err error) = wimgapi.WIMGetAttributes
//sys	WIMSetBootImage(wim windows.Handle, imageIndex uint32) (err error) = wimgapi.WIMSetBootImage
//sys	wimCaptureImage(wim windows.Handle, path *uint16, captureFlags uint32) (handle windows.Handle, err error) = wimgapi.WIMCaptureImage
//sys	WIMLoadImage(wim windows.Handle, imageIndex uint32) (handle windows.Handle) = wimgapi.WIMLoadImage
//sys	wimApplyImage(image windows.Handle, path *uint16, applyFlags uint32) (err error) = wimgapi.WIMApplyImage
//sys	wimGetImageInformation(image windows.Handle, imageInfo *unsafe.Pointer, cbImageInfo *uint32) (err error) = wimgapi.WIMGetImageInformation
//sys	wimSetImageInformation(image windows.Handle, imageInfo unsafe.Pointer, cbImageInfo uint32) (err error) = wimgapi.WIMSetImageInformation
//sys	WIMGetMessageCallbackCount(wim windows.Handle) (count uint32) = wimgapi.WIMGetMessageCallbackCount
//sys	WIMRegisterMessageCallback(wim windows.Handle, callback uintptr, userData unsafe.Pointer) (index uint32, err error) [failretval==INVALID_CALLBACK_VALUE] = wimgapi.WIMRegisterMessageCallback
//sys 	WIMUnregisterMessageCallback(wim windows.Handle, callback uintptr) (err error) = wimgapi.WIMUnregisterMessageCallback
//sys	WIMMessageCallback(messageId uint32, wParam WPARAM, lParam LPARAM, userData unsafe.Pointer) (err error) = wimgapi.WIMMessageCallback
//sys	wimCopyFile(existingFileName *uint16, newFileName *uint16, progressRoutine uintptr, data unsafe.Pointer, cancel *int32, copyFlags uint32) (err error) = wimgapi.WIMCopyFile
//sys	wimMountImage(mountPath *uint16, wimFileName *uint16, imageIndex uint32, tempPath *uint16) (err error) = wimgapi.WIMMountImage
//sys	wimUnmountImage(mountPath *uint16, wimFileName *uint16, imageIndex uint32, commitChanges bool) (err error) = wimgapi.WIMUnmountImage
//sys	wimGetMountedImages(mountList *WIM_MOUNT_LIST, cbMountListLength *uint32) (err error) = wimgapi.WIMGetMountedImages
//sys	WIMInitFileIOCallbacks(callbacks unsafe.Pointer) (err error) = wimgapi.WIMInitFileIOCallbacks
//sys	wimSetFileIOCallbackTemporaryPath(path *uint16) (err error) = wimgapi.WIMSetFileIOCallbackTemporaryPath
//sys	wimMountImageHandle(image windows.Handle, mountPath *uint16, mountFlags uint32) (err error) = wimgapi.WIMMountImageHandle
//sys	wimRemountImage(mountPath *uint16, flags uint32) (err error) = wimgapi.WIMRemountImage
//sys	wimCommitImageHandle(image windows.Handle, commitFlags uint32, newImageHandle *windows.Handle) (err error) = wimgapi.WIMCommitImageHandle
//sys	WIMUnmountImageHandle(image windows.Handle, unmountFlags uint32) (err error) = wimgapi.WIMUnmountImageHandle
//sys	wimGetMountedImageInfo(infoLevelId MOUNTED_IMAGE_INFO_LEVELS, imageCount *uint32, mountInfo unsafe.Pointer, cbMountInfoLength uint32, returnLength *uint32) (err error) = wimgapi.WIMGetMountedImageInfo
//sys	wimGetMountedImageInfoFromHandle(image windows.Handle, infoLevelId MOUNTED_IMAGE_INFO_LEVELS, mountInfo unsafe.Pointer, cbMountInfoLength uint32, returnLength *uint32) (err error) = wimgapi.WIMGetMountedImageInfoFromHandle
//sys	wimGetMountedImageHandle(mountPath *uint16, flags uint32, wimHandle *windows.Handle, imageHandle *windows.Handle) (err error) = wimgapi.WIMGetMountedImageHandle
//sys	WIMDeleteImageMounts(deleteFlags uint32) (err error) = wimgapi.WIMDeleteImageMounts
//sys	wimRegisterLogFile(logFile *uint16, flags uint32) (err error) = wimgapi.WIMRegisterLogFile
//sys	wimExtractImagePath(image windows.Handle, imagePath *uint16, destinationPath *uint16, extractFlags uint32) (err error) = wimgapi.WIMExtractImagePath
//sys	wimFindFirstImageFile(image windows.Handle, filePath *uint16, findFileData *WIM_FIND_DATA) (handle windows.Handle, err error) = wimgapi.WIMFindFirstImageFile
//sys	wimFindNextImageFile(findFile windows.Handle, findFileData *WIM_FIND_DATA) (err error) = wimgapi.WIMFindNextImageFile
//sys	WIMEnumImageFiles(image windows.Handle, enumFile PWIM_ENUM_FILE, enumImageCallback uintptr, enumContext unsafe.Pointer) (err error) = wimgapi.WIMEnumImageFiles
//sys	wimCreateImageFile(image windows.Handle, filePath *uint16, desiredAccess uint32, creationDisposition uint32, flagsAndAttributes uint32) (handle windows.Handle, err error) = wimgapi.WIMCreateImageFile
//sys	wimReadImageFile(imgFile windows.Handle, buffer *byte, bytesToRead uint32, bytesRead *uint32, overlapped *windows.Overlapped) (err error) = wimgapi.WIMReadImageFile

// The createdNew value is set if getCreationResult is true, otherwise it
// is always false.
func WIMCreateFile(
	wimPath string,
	desiredAccess uint32,
	creationDisposition uint32,
	flagsAndAttributes uint32,
	compressionType uint32,
	getCreationResult bool,
) (
	handle windows.Handle,
	createdNew bool,
	err error,
) {
	u16WimPath, err := windows.UTF16PtrFromString(wimPath)
	if err != nil {
		return
	}

	var creatRes *uint32
	if getCreationResult {
		var temp uint32
		creatRes = &temp
	}

	if handle, err = wimCreateFile(
		u16WimPath,
		desiredAccess,
		creationDisposition,
		flagsAndAttributes,
		compressionType,
		creatRes,
	); err != nil {
		return
	}

	if creatRes != nil {
		createdNew = *creatRes == WIM_CREATED_NEW
	}

	return
}

func WIMSetTemporaryPath(
	wim windows.Handle,
	path string,
) error {
	u16Path, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return err
	}

	if err = wimSetTemporaryPath(
		wim,
		u16Path,
	); err != nil {
		return err
	}

	return nil
}

func WIMSetReferenceFile(
	wim windows.Handle,
	path string,
	flags uint32,
) error {
	u16Path, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return err
	}

	if err = wimSetReferenceFile(
		wim,
		u16Path,
		flags,
	); err != nil {
		return err
	}

	return nil
}

func WIMSplitFile(
	wim windows.Handle,
	partPath string,
	partSize *int64,
	flags uint32,
) error {
	u16PartPath, err := windows.UTF16PtrFromString(partPath)
	if err != nil {
		return err
	}

	if err = wimSplitFile(
		wim,
		u16PartPath,
		partSize,
		flags,
	); err != nil {
		return err
	}

	return nil
}

func WIMGetAttributes(wim windows.Handle) (wimInfo GoWimInfo, err error) {
	var info WIM_INFO

	if err = wimGetAttributes(
		wim,
		&info,
		uint32(unsafe.Sizeof(info)),
	); err != nil {
		return
	}

	wimInfo.fill(&info)

	return
}

func WIMCaptureImage(
	wim windows.Handle,
	path string,
	captureFlags uint32,
) (windows.Handle, error) {
	u16Path, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return 0, err
	}

	return wimCaptureImage(
		wim,
		u16Path,
		captureFlags,
	)
}

func WIMApplyImage(
	image windows.Handle,
	path string,
	applyFlags uint32,
) error {
	u16Path, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return err
	}

	if err = wimApplyImage(
		image,
		u16Path,
		applyFlags,
	); err != nil {
		return err
	}

	return nil
}

// returns byte slice with XML information about the volume image
func WIMGetImageInformation(image windows.Handle) ([]byte, error) {
	var imgInfo unsafe.Pointer
	var bufSize uint32

	if err := wimGetImageInformation(
		image,
		&imgInfo,
		&bufSize,
	); err != nil {
		return nil, err
	}

	defer windows.LocalFree(windows.Handle(uintptr(imgInfo)))

	buf := unsafe.Slice((*byte)(imgInfo), bufSize)

	return buf, nil
}

// store information about an image with imageInfo
func WIMSetImageInformation(
	image windows.Handle,
	imageInfo []byte,
) error {
	infoPtr := unsafe.Pointer(&imageInfo[0])
	infoSize := uint32(len(imageInfo))

	if err := wimSetImageInformation(
		image,
		infoPtr,
		infoSize,
	); err != nil {
		return err
	}

	return nil
}

// pass uintptr returned from [windows.NewCallback] with
// [wintype.LPPROGRESS_ROUTINE] for progressRoutine(optional)
//
// set *cancel to 1(TRUE) to cancel copy operation
func WIMCopyFile(
	existingFileName string,
	newFileName string,
	progressRoutine uintptr,
	data unsafe.Pointer,
	cancel *int32, // BOOL
	copyFlags uint32,
) error {
	u16ExistName, err := windows.UTF16PtrFromString(existingFileName)
	if err != nil {
		return err
	}

	u16NewName, err := windows.UTF16PtrFromString(newFileName)
	if err != nil {
		return err
	}

	if err = wimCopyFile(
		u16ExistName,
		u16NewName,
		progressRoutine,
		data,
		cancel,
		copyFlags,
	); err != nil {
		return err
	}

	return nil
}

// if tempPath if empty, the image will not be mounted
// for edits
func WIMMountImage(
	mountPath string,
	wimFileName string,
	imageIndex uint32,
	tempPath string,
) error {
	u16MntPath, err := windows.UTF16PtrFromString(mountPath)
	if err != nil {
		return err
	}

	u16WimFile, err := windows.UTF16PtrFromString(wimFileName)
	if err != nil {
		return err
	}

	var u16TempPath *uint16
	if tempPath != "" {
		u16TempPath, err = windows.UTF16PtrFromString(tempPath)
		if err != nil {
			return err
		}
	}

	if err = wimMountImage(
		u16MntPath,
		u16WimFile,
		imageIndex,
		u16TempPath,
	); err != nil {
		return err
	}

	return nil
}

func WIMUnmountImage(
	mountPath string,
	wimFileName string,
	imageIndex uint32,
	commitChanges bool,
) error {
	u16MntPath, err := windows.UTF16PtrFromString(mountPath)
	u16WimFile, err := windows.UTF16PtrFromString(wimFileName)

	if err = wimUnmountImage(
		u16MntPath,
		u16WimFile,
		imageIndex,
		commitChanges,
	); err != nil {
		return err
	}

	return nil
}

func WIMGetMountedImages() ([]GoWimMountList, error) {
	var listByteSize, listLen uint32
	var mountList []WIM_MOUNT_LIST

	// get required bytes size
	err := wimGetMountedImages(
		nil,
		&listByteSize,
	)
	if err != nil && err != windows.ERROR_INSUFFICIENT_BUFFER {
		return nil, err
	}

	listLen = listByteSize / uint32(unsafe.Sizeof(WIM_MOUNT_LIST{}))
	mountList = make([]WIM_MOUNT_LIST, listLen)

	if err = wimGetMountedImages(
		&mountList[0],
		&listByteSize,
	); err != nil {
		return nil, err
	}

	result := make([]GoWimMountList, listLen)

	for i, mountInfo := range mountList {
		result[i].fill(&mountInfo)
	}

	return result, nil
}

func WIMSetFileIOCallbackTemporaryPath(path string) error {
	u16Path, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return err
	}

	if err = wimSetFileIOCallbackTemporaryPath(u16Path); err != nil {
		return err
	}

	return nil
}

func WIMMountImageHandle(
	image windows.Handle,
	mountPath string,
	mountFlags uint32,
) error {
	u16MntPath, err := windows.UTF16PtrFromString(mountPath)
	if err != nil {
		return err
	}

	if err = wimMountImageHandle(
		image,
		u16MntPath,
		mountFlags,
	); err != nil {
		return err
	}

	return nil
}

// flags must be 0
func WIMRemountImage(
	mountPath string,
	flags uint32,
) error {
	u16MntPath, err := windows.UTF16PtrFromString(mountPath)
	if err != nil {
		return err
	}

	if err = wimRemountImage(
		u16MntPath,
		flags,
	); err != nil {
		return err
	}

	return nil
}

// if WIM_COMMIT_FLAG_APPEND is specified in commitFlags and
// openNewHandle is set to true, the value of newHandle is
// set, otherwise it is always 0
func WIMCommitImageHandle(
	image windows.Handle,
	commitFlags uint32,
	openNewHandle bool,
) (newHandle windows.Handle, err error) {
	var temp *windows.Handle
	if openNewHandle {
		temp = &newHandle
	}

	err = wimCommitImageHandle(
		image,
		commitFlags,
		temp,
	)

	return
}

func WIMGetMountedImageInfo[T GoWimMountInfoLevel0 | GoWimMountInfoLevel1]() ([]T, error) {
	var infoLevel MOUNTED_IMAGE_INFO_LEVELS
	var imageCount, returnLength uint32

	switch any((*T)(nil)).(type) {
	case *GoWimMountInfoLevel0:
		infoLevel = MountedImageLevel0
	case *GoWimMountInfoLevel1:
		infoLevel = MountedImageLevel1
	}

	// get required buffer size
	err := wimGetMountedImageInfo(
		infoLevel,
		&imageCount,
		nil,
		0,
		&returnLength,
	)
	if err != nil && err != windows.ERROR_INSUFFICIENT_BUFFER {
		return nil, err
	}

	if imageCount == 0 {
		return nil, nil
	}

	buf := make([]byte, returnLength)

	err = wimGetMountedImageInfo(
		infoLevel,
		&imageCount,
		unsafe.Pointer(&buf[0]),
		uint32(len(buf)),
		&returnLength,
	)
	if err != nil {
		return nil, err
	}

	result := make([]T, imageCount)

	switch infoLevel {
	case MountedImageLevel0:
		for i, info := range unsafe.Slice((*WIM_MOUNT_INFO_LEVEL0)(unsafe.Pointer(&buf[0])), imageCount) {
			goInfo := any(result[i]).(GoWimMountInfoLevel0)
			goInfo.fill(&info)
		}
	case MountedImageLevel1:
		for i, info := range unsafe.Slice((*WIM_MOUNT_INFO_LEVEL1)(unsafe.Pointer(&buf[0])), imageCount) {
			goInfo := any(result[i]).(GoWimMountInfoLevel1)
			goInfo.fill(&info)
		}
	}

	return result, nil
}
