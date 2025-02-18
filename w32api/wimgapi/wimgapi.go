// Package wimgapi implements Windows Imaging Interface library.
//
// https://learn.microsoft.com/en-us/windows-hardware/manufacture/desktop/wim/dd851927(v=msdn.10)
package wimgapi

import (
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/Snshadow/winimg/w32api"
)

//sys	wimCreateFile(wimPath *uint16, desiredAccess uint32, creationDisposition uint32, flagsAndAttributes uint32, compressionType uint32, creationResult *uint32) (handle windows.Handle, err error) = wimgapi.WIMCreateFile
//sys	wimCloseHandle(object windows.Handle) (err error) = wimgapi.WIMCloseHandle
//sys	wimSetTemporaryPath(wim windows.Handle, path *uint16) (err error) = wimgapi.WIMSetTemporaryPath
//sys	wimSetReferenceFile(wim windows.Handle, path *uint16, flags uint32) (err error) = wimgapi.WIMSetReferenceFile
//sys	wimSplitFile(wim windows.Handle, partPath *uint16, partSize *int64, flags uint32) (err error) = wimgapi.WIMSplitFile
//sys	wimExportImage(image windows.Handle, wim windows.Handle, flags uint32) (err error) = wimgapi.WIMExportImage
//sys	wimDeleteImage(wim windows.Handle, imageIndex uint32) (err error) = wimgapi.WIMDeleteImage
//sys	wimGetImageCount(wim windows.Handle) (count uint32) = wimgapi.WIMGetImageCount
//sys	wimGetAttributes(wim windows.Handle, wimInfo *WIM_INFO, cbWimInfo uint32) (err error) = wimgapi.WIMGetAttributes
//sys	wimSetBootImage(wim windows.Handle, imageIndex uint32) (err error) = wimgapi.WIMSetBootImage
//sys	wimCaptureImage(wim windows.Handle, path *uint16, captureFlags uint32) (handle windows.Handle, err error) = wimgapi.WIMCaptureImage
//sys	wimLoadImage(wim windows.Handle, imageIndex uint32) (handle windows.Handle, err error) = wimgapi.WIMLoadImage
//sys	wimApplyImage(image windows.Handle, path *uint16, applyFlags uint32) (err error) = wimgapi.WIMApplyImage
//sys	wimGetImageInformation(image windows.Handle, imageInfo *unsafe.Pointer, cbImageInfo *uint32) (err error) = wimgapi.WIMGetImageInformation
//sys	wimSetImageInformation(image windows.Handle, imageInfo unsafe.Pointer, cbImageInfo uint32) (err error) = wimgapi.WIMSetImageInformation
//sys	wimGetMessageCallbackCount(wim windows.Handle) (count uint32) = wimgapi.WIMGetMessageCallbackCount
//sys	wimRegisterMessageCallback(wim windows.Handle, callback uintptr, userData unsafe.Pointer) (index uint32, err error) [failretval==INVALID_CALLBACK_VALUE] = wimgapi.WIMRegisterMessageCallback
//sys 	wimUnregisterMessageCallback(wim windows.Handle, callback uintptr) (err error) = wimgapi.WIMUnregisterMessageCallback
//sys	wimCopyFile(existingFileName *uint16, newFileName *uint16, progressRoutine uintptr, data unsafe.Pointer, cancel *int32, copyFlags uint32) (err error) = wimgapi.WIMCopyFile
//sys	wimMountImage(mountPath *uint16, wimFileName *uint16, imageIndex uint32, tempPath *uint16) (err error) = wimgapi.WIMMountImage
//sys	wimUnmountImage(mountPath *uint16, wimFileName *uint16, imageIndex uint32, commitChanges bool) (err error) = wimgapi.WIMUnmountImage
//sys	wimGetMountedImages(mountList *WIM_MOUNT_LIST, cbMountListLength *uint32) (err error) = wimgapi.WIMGetMountedImages
//sys	wimInitFileIOCallbacks(callbacks unsafe.Pointer) (err error) = wimgapi.WIMInitFileIOCallbacks
//sys	wimSetFileIOCallbackTemporaryPath(path *uint16) (err error) = wimgapi.WIMSetFileIOCallbackTemporaryPath
//sys	wimMountImageHandle(image windows.Handle, mountPath *uint16, mountFlags uint32) (err error) = wimgapi.WIMMountImageHandle
//sys	wimRemountImage(mountPath *uint16, flags uint32) (err error) = wimgapi.WIMRemountImage
//sys	wimCommitImageHandle(image windows.Handle, commitFlags uint32, newImageHandle *windows.Handle) (err error) = wimgapi.WIMCommitImageHandle
//sys	wimUnmountImageHandle(image windows.Handle, unmountFlags uint32) (err error) = wimgapi.WIMUnmountImageHandle
//sys	wimGetMountedImageInfo(infoLevelId MOUNTED_IMAGE_INFO_LEVELS, imageCount *uint32, mountInfo unsafe.Pointer, cbMountInfoLength uint32, returnLength *uint32) (err error) = wimgapi.WIMGetMountedImageInfo
//sys	wimGetMountedImageInfoFromHandle(image windows.Handle, infoLevelId MOUNTED_IMAGE_INFO_LEVELS, mountInfo unsafe.Pointer, cbMountInfoLength uint32, returnLength *uint32) (err error) = wimgapi.WIMGetMountedImageInfoFromHandle
//sys	wimGetMountedImageHandle(mountPath *uint16, flags uint32, wimHandle *windows.Handle, imageHandle *windows.Handle) (err error) = wimgapi.WIMGetMountedImageHandle
//sys	wimDeleteImageMounts(deleteFlags uint32) (err error) = wimgapi.WIMDeleteImageMounts
//sys	wimRegisterLogFile(logFile *uint16, flags uint32) (err error) = wimgapi.WIMRegisterLogFile
//sys	wimUnregisterLogFile(logFile *uint16) (err error) = wimgapi.WIMUnregisterLogFile
//sys	wimExtractImagePath(image windows.Handle, imagePath *uint16, destinationPath *uint16, extractFlags uint32) (err error) = wimgapi.WIMExtractImagePath
//sys	wimFindFirstImageFile(image windows.Handle, filePath *uint16, findFileData *WIM_FIND_DATA) (handle windows.Handle, err error) = wimgapi.WIMFindFirstImageFile
//sys	wimFindNextImageFile(findFile windows.Handle, findFileData *WIM_FIND_DATA) (err error) = wimgapi.WIMFindNextImageFile
//sys	wimEnumImageFiles(image windows.Handle, enumFile PWIM_ENUM_FILE, enumImageCallback uintptr, enumContext unsafe.Pointer) (err error) = wimgapi.WIMEnumImageFiles
//sys	wimCreateImageFile(image windows.Handle, filePath *uint16, desiredAccess uint32, creationDisposition uint32, flagsAndAttributes uint32) (handle windows.Handle, err error) = wimgapi.WIMCreateImageFile
//sys	wimReadImageFile(imgFile windows.Handle, buffer *byte, bytesToRead uint32, bytesRead *uint32, overlapped *windows.Overlapped) (err error) = wimgapi.WIMReadImageFile

// The createdNew value is set if getCreationResult is true, otherwise it
// is always false. Close handle with [WIMCloseHandle] after use.
func WIMCreateFile(
	wimPath string,
	desiredAccess uint32,
	creationDisposition uint32,
	flagsAndAttributes uint32,
	compressionType WimCompressionType,
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
		uint32(compressionType),
		creatRes,
	); err != nil {
		err = w32api.WrapInternalErr(err, modwimgapi.Handle(), "")
		return
	}

	if creatRes != nil {
		createdNew = *creatRes == WIM_CREATED_NEW
	}

	return
}

func WIMCloseHandle(object windows.Handle) error {
	if err := wimCloseHandle(object); err != nil {
		return w32api.WrapInternalErr(err, modwimgapi.Handle(), "")
	}

	return nil
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
		return w32api.WrapInternalErr(err, modwimgapi.Handle(), "")
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
		return w32api.WrapInternalErr(err, modwimgapi.Handle(), "")
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
		return w32api.WrapInternalErr(err, modwimgapi.Handle(), "")
	}

	return nil
}

func WIMExportImage(
	image windows.Handle,
	wim windows.Handle,
	flags uint32,
) error {
	if err := wimExportImage(
		image,
		wim,
		flags,
	); err != nil {
		return w32api.WrapInternalErr(err, modwimgapi.Handle(), "")
	}

	return nil
}

func WIMDeleteImage(
	wim windows.Handle,
	imageIndex uint32,
) error {
	if err := wimDeleteImage(
		wim,
		imageIndex,
	); err != nil {
		return w32api.WrapInternalErr(err, modwimgapi.Handle(), "")
	}

	return nil
}

func WIMGetImageCount(wim windows.Handle) uint32 {
	return wimGetImageCount(wim)
}

func WIMGetAttributes(wim windows.Handle) (wimInfo GoWimInfo, err error) {
	var info WIM_INFO

	if err = wimGetAttributes(
		wim,
		&info,
		uint32(unsafe.Sizeof(info)),
	); err != nil {
		err = w32api.WrapInternalErr(err, modwimgapi.Handle(), "")
		return
	}

	wimInfo.fill(&info)

	return
}

func WIMSetBootImage(
	wim windows.Handle,
	imageIndex uint32,
) error {
	if err := wimSetBootImage(
		wim,
		imageIndex,
	); err != nil {
		return w32api.WrapInternalErr(err, modwimgapi.Handle(), "")
	}

	return nil
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

	hnd, err := wimCaptureImage(
		wim,
		u16Path,
		captureFlags,
	)
	if err != nil {
		return 0, w32api.WrapInternalErr(err, modwimgapi.Handle(), "")
	}

	return hnd, nil
}

func WIMLoadImage(
	wim windows.Handle,
	imageIndex uint32,
) (windows.Handle, error) {
	hnd, err := wimLoadImage(
		wim,
		imageIndex,
	)
	if err != nil {
		return 0, w32api.WrapInternalErr(err, modwimgapi.Handle(), "")
	}

	return hnd, nil
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
		return w32api.WrapInternalErr(err, modwimgapi.Handle(), "")
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
		return nil, w32api.WrapInternalErr(err, modwimgapi.Handle(), "")
	}

	defer windows.LocalFree(windows.Handle(uintptr(imgInfo)))

	buf := make([]byte, bufSize)
	copy(buf, unsafe.Slice((*byte)(imgInfo), bufSize))

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
		return w32api.WrapInternalErr(err, modwimgapi.Handle(), "")
	}

	return nil
}

func WIMGetMessageCallbackCount(wim windows.Handle) uint32 {
	return wimGetMessageCallbackCount(wim)
}

func WIMRegisterMessageCallback(
	wim windows.Handle,
	callback uintptr,
	userData unsafe.Pointer,
) (uint32, error) {
	index, err := wimRegisterMessageCallback(
		wim,
		callback,
		userData,
	)
	if err != nil {
		return index, w32api.WrapInternalErr(err, modwimgapi.Handle(), "")
	}

	return index, nil
}

func WIMUnregisterMessageCallback(
	wim windows.Handle,
	callback uintptr,
) error {
	if err := wimUnregisterMessageCallback(
		wim,
		callback,
	); err != nil {
		return w32api.WrapInternalErr(err, modwimgapi.Handle(), "")
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
	cancel *int32, // PBOOL
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
		return w32api.WrapInternalErr(err, modwimgapi.Handle(), "")
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
		return w32api.WrapInternalErr(err, modwimgapi.Handle(), "")
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
	if err != nil {
		return err
	}
	u16WimFile, err := windows.UTF16PtrFromString(wimFileName)
	if err != nil {
		return err
	}

	if err = wimUnmountImage(
		u16MntPath,
		u16WimFile,
		imageIndex,
		commitChanges,
	); err != nil {
		return w32api.WrapInternalErr(err, modwimgapi.Handle(), "")
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
		return nil, w32api.WrapInternalErr(err, modwimgapi.Handle(), "")
	}

	listLen = listByteSize / uint32(unsafe.Sizeof(WIM_MOUNT_LIST{}))
	mountList = make([]WIM_MOUNT_LIST, listLen)

	if err = wimGetMountedImages(
		&mountList[0],
		&listByteSize,
	); err != nil {
		return nil, w32api.WrapInternalErr(err, modwimgapi.Handle(), "")
	}

	result := make([]GoWimMountList, listLen)

	for i, mountInfo := range mountList {
		result[i].fill(&mountInfo)
	}

	return result, nil
}

func WIMInitFileIOCallbacks(callbacks unsafe.Pointer) error {
	if err := wimInitFileIOCallbacks(callbacks); err != nil {
		return w32api.WrapInternalErr(err, modwimgapi.Handle(), "")
	}

	return nil
}

func WIMSetFileIOCallbackTemporaryPath(path string) error {
	u16Path, err := windows.UTF16PtrFromString(path)
	if err != nil {
		return err
	}

	if err = wimSetFileIOCallbackTemporaryPath(u16Path); err != nil {
		return w32api.WrapInternalErr(err, modwimgapi.Handle(), "")
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
		return w32api.WrapInternalErr(err, modwimgapi.Handle(), "")
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

func WIMUnmountImageHandle(
	image windows.Handle,
	unmountFlags uint32,
) error {
	if err := wimUnmountImageHandle(
		image,
		unmountFlags,
	); err != nil {
		return w32api.WrapInternalErr(err, modwimgapi.Handle(), "")
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
		return w32api.WrapInternalErr(err, modwimgapi.Handle(), "")
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
	if err != nil {
		err = w32api.WrapInternalErr(err, modwimgapi.Handle(), "")
	}

	return
}

func WIMGetMountedImageInfo[T GoWimMountInfoLevel0 | GoWimMountInfoLevel1]() ([]T, error) {
	var infoLevel MOUNTED_IMAGE_INFO_LEVELS
	var infoCount, returnLength uint32

	switch any((*T)(nil)).(type) {
	case *GoWimMountInfoLevel0:
		infoLevel = MountedImageLevel0
	case *GoWimMountInfoLevel1:
		infoLevel = MountedImageLevel1
	}

	// get required buffer size
	err := wimGetMountedImageInfo(
		infoLevel,
		&infoCount,
		nil,
		0,
		&returnLength,
	)
	if err != nil && err != windows.ERROR_INSUFFICIENT_BUFFER {
		return nil, w32api.WrapInternalErr(err, modwimgapi.Handle(), "")
	}

	if infoCount == 0 {
		return nil, nil
	}

	buf := make([]byte, returnLength)

	err = wimGetMountedImageInfo(
		infoLevel,
		&infoCount,
		unsafe.Pointer(&buf[0]),
		uint32(len(buf)),
		&returnLength,
	)
	if err != nil {
		return nil, w32api.WrapInternalErr(err, modwimgapi.Handle(), "")
	}

	result := make([]T, infoCount)

	switch infoLevel {
	case MountedImageLevel0:
		for i, info := range unsafe.Slice((*WIM_MOUNT_INFO_LEVEL0)(unsafe.Pointer(&buf[0])), infoCount) {
			any(&result[i]).(*GoWimMountInfoLevel0).fill(&info)
		}
	case MountedImageLevel1:
		for i, info := range unsafe.Slice((*WIM_MOUNT_INFO_LEVEL1)(unsafe.Pointer(&buf[0])), infoCount) {
			any(&result[i]).(*GoWimMountInfoLevel1).fill(&info)
		}
	}

	return result, nil
}

func WiMGetMountedImageInfoFromHandle[T GoWimMountInfoLevel0 | GoWimMountInfoLevel1](
	image windows.Handle,
) ([]T, error) {
	var infoLevel MOUNTED_IMAGE_INFO_LEVELS
	var returnLength, structSize uint32

	switch any((*T)(nil)).(type) {
	case *GoWimMountInfoLevel0:
		infoLevel = MountedImageLevel0
		structSize = uint32(unsafe.Sizeof(WIM_MOUNT_INFO_LEVEL0{}))
	case *GoWimMountInfoLevel1:
		infoLevel = MountedImageLevel1
		structSize = uint32(unsafe.Sizeof(WIM_MOUNT_INFO_LEVEL1{}))
	}

	// get required buffer size
	err := wimGetMountedImageInfoFromHandle(
		image,
		infoLevel,
		nil,
		0,
		&returnLength,
	)
	if err != nil && err != windows.ERROR_INSUFFICIENT_BUFFER {
		return nil, w32api.WrapInternalErr(err, modwimgapi.Handle(), "")
	}

	buf := make([]byte, returnLength)

	if err = wimGetMountedImageInfoFromHandle(
		image,
		infoLevel,
		unsafe.Pointer(&buf[0]),
		uint32(len(buf)),
		&returnLength,
	); err != nil {
		return nil, w32api.WrapInternalErr(err, modwimgapi.Handle(), "")
	}

	infoCount := returnLength / structSize
	result := make([]T, infoCount)

	switch infoLevel {
	case MountedImageLevel0:
		for i, info := range unsafe.Slice((*WIM_MOUNT_INFO_LEVEL0)(unsafe.Pointer(&buf[0])), infoCount) {
			any(&result[i]).(*GoWimMountInfoLevel0).fill(&info)
		}
	case MountedImageLevel1:
		for i, info := range unsafe.Slice((*WIM_MOUNT_INFO_LEVEL1)(unsafe.Pointer(&buf[0])), infoCount) {
			any(&result[i]).(*GoWimMountInfoLevel1).fill(&info)
		}
	}

	return result, nil
}

func WIMGetMountedImageHandle(
	mountPath string,
	flags uint32,
) (
	wimhandle windows.Handle,
	imageHandle windows.Handle,
	err error,
) {
	u16MntPath, err := windows.UTF16PtrFromString(mountPath)
	if err != nil {
		return
	}

	err = wimGetMountedImageHandle(
		u16MntPath,
		flags,
		&wimhandle,
		&imageHandle,
	)
	if err != nil {
		err = w32api.WrapInternalErr(err, modwimgapi.Handle(), "")
	}

	return
}

func WIMDeleteImageMounts(deleteFlags uint32) error {
	if err := wimDeleteImageMounts(deleteFlags); err != nil {
		return w32api.WrapInternalErr(err, modwimgapi.Handle(), "")
	}

	return nil
}

func WIMRegisterLogFile(
	logFile string,
	flags uint32,
) error {
	u16LogFile, err := windows.UTF16PtrFromString(logFile)
	if err != nil {
		return err
	}

	if err = wimRegisterLogFile(
		u16LogFile,
		flags,
	); err != nil {
		return w32api.WrapInternalErr(err, modwimgapi.Handle(), "")
	}

	return nil
}

func WIMUnresigisterLogFile(logFile string) error {
	u16LogFile, err := windows.UTF16PtrFromString(logFile)
	if err != nil {
		return err
	}

	if err = wimUnregisterLogFile(u16LogFile); err != nil {
		return w32api.WrapInternalErr(err, modwimgapi.Handle(), "")
	}

	return nil
}

func WIMExtractImagePath(
	image windows.Handle,
	imagePath string,
	destinationPath string,
	extractFlags uint32,
) error {
	u16ImgPath, err := windows.UTF16PtrFromString(imagePath)
	if err != nil {
		return err
	}

	u16DestPath, err := windows.UTF16PtrFromString(destinationPath)
	if err != nil {
		return err
	}

	if err = wimExtractImagePath(
		image,
		u16ImgPath,
		u16DestPath,
		extractFlags,
	); err != nil {
		return w32api.WrapInternalErr(err, modwimgapi.Handle(), "")
	}

	return nil
}

// TODO check [WIM_FIND_DATA].AlternativeStreamNames allocation release.
func WIMFindFirstImageFile(
	image windows.Handle,
	filePath string,
) (
	windows.Handle,
	*GoWimFindData,
	error,
) {
	u16FilePath, err := windows.UTF16PtrFromString(filePath)
	if err != nil {
		return 0, nil, err
	}

	var findFileData WIM_FIND_DATA

	findHandle, err := wimFindFirstImageFile(
		image,
		u16FilePath,
		&findFileData,
	)
	if err != nil {
		return 0, nil, w32api.WrapInternalErr(err, modwimgapi.Handle(), "")
	}

	result := &GoWimFindData{}
	result.fill(&findFileData)

	return findHandle, result, nil
}

func WIMFindNextImageFile(findFile windows.Handle) (*GoWimFindData, error) {
	var findFileData WIM_FIND_DATA

	if err := wimFindNextImageFile(findFile, &findFileData); err != nil {
		return nil, w32api.WrapInternalErr(err, modwimgapi.Handle(), "")
	}

	result := &GoWimFindData{}
	result.fill(&findFileData)

	return result, nil
}

func WIMEnumImageFiles(
	image windows.Handle,
	enumFile PWIM_ENUM_FILE,
	enumImageCallback uintptr,
	enumContext unsafe.Pointer,
) error {
	if err := wimEnumImageFiles(
		image,
		enumFile,
		enumImageCallback,
		enumContext,
	); err != nil {
		return w32api.WrapInternalErr(err, modwimgapi.Handle(), "")
	}

	return nil
}

func WIMCreateImageFile(
	image windows.Handle,
	filePath string,
	desiredAccess uint32,
	creationDisposition uint32,
	flagsAndAttributes uint32,
) (windows.Handle, error) {
	u16FilePath, err := windows.UTF16PtrFromString(filePath)
	if err != nil {
		return 0, err
	}

	hnd, err := wimCreateImageFile(
		image,
		u16FilePath,
		desiredAccess,
		creationDisposition,
		flagsAndAttributes,
	)
	if err != nil {
		return 0, w32api.WrapInternalErr(err, modwimgapi.Handle(), "")
	}

	return hnd, nil
}

func WIMReadImageFile(
	imgFile windows.Handle,
	buf []byte,
	overlapped *windows.Overlapped,
) (int, error) {
	var bytesRead uint32

	err := wimReadImageFile(
		imgFile,
		&buf[0],
		uint32(len(buf)),
		&bytesRead,
		overlapped,
	)
	if err != nil {
		err = w32api.WrapInternalErr(err, modwimgapi.Handle(), "")
	}

	return int(bytesRead), err
}
