package wimgapi

import (
	"time"
	"unsafe"

	"golang.org/x/sys/windows"

	"github.com/Snshadow/go_win_image/internal/utils"
)

// helper structs with go friendly types

type GoWimInfo struct {
	WimPath         string
	Guid            windows.GUID
	ImageCount      uint32
	CompressionType uint32
	PartNumber      uint16
	TotalParts      uint16
	BootIndex       uint32
	WimAttributes   uint32
	WimFlagsAndAttr uint32
}

func (g *GoWimInfo) fill(st *WIM_INFO) {
	g.WimPath = windows.UTF16ToString(st.WimPath[:])
	g.Guid = st.Guid
	g.ImageCount = st.ImageCount
	g.CompressionType = st.CompressionType
	g.PartNumber = st.PartNumber
	g.TotalParts = st.TotalParts
	g.BootIndex = st.BootIndex
	g.WimAttributes = st.WimAttributes
	g.WimFlagsAndAttr = st.WimFlagsAndAttr
}

type GoWimMountList struct {
	WimPath    string
	MountPath  string
	ImageIndex uint32
	MountForRW bool
}

func (g *GoWimMountList) fill(st *WIM_MOUNT_LIST) {
	g.WimPath = windows.UTF16ToString(st.WimPath[:])
	g.MountPath = windows.UTF16ToString(st.MountPath[:])
	g.ImageIndex = st.ImageIndex
	g.MountForRW = st.MountedForRW != 0
}

type GoWimMountInfoLevel0 GoWimMountList

func (g *GoWimMountInfoLevel0) fill(st *WIM_MOUNT_INFO_LEVEL0) {
	g.WimPath = windows.UTF16ToString(st.WimPath[:])
	g.MountPath = windows.UTF16ToString(st.MountPath[:])
	g.ImageIndex = st.ImageIndex
	g.MountForRW = st.MountedForRW != 0
}

type GoWimMountInfoLevel1 struct {
	WimPath    string
	MountPath  string
	ImageIndex uint32
	MountFlags uint32
}

func (g *GoWimMountInfoLevel1) fill(st *WIM_MOUNT_INFO_LEVEL1) {
	g.WimPath = windows.UTF16ToString(st.WimPath[:])
	g.MountPath = windows.UTF16ToString(st.MountPath[:])
	g.ImageIndex = st.ImageIndex
	g.MountFlags = st.MountFlags
}

type GoWimFindData struct {
	FileAttributes        uint32
	CreationFile          time.Time
	LastAccessTime        time.Time
	LastWriteTime         time.Time
	FileSize              uint64
	FileName              string
	AlternateFileName     string
	Hash                  [20]byte
	SecurityDescriptor    *windows.SECURITY_DESCRIPTOR
	AlternateStreamNames  []string
	ReparseData           []byte
	ResourceSize          uint64
	ResourceOffset        int64
	ResourceReferencePath string
}

func (g *GoWimFindData) fill(st *WIM_FIND_DATA) {
	g.FileAttributes = st.FileAttributes
	g.CreationFile = time.Unix(0, st.CreationFile.Nanoseconds())
	g.LastAccessTime = time.Unix(0, st.LastAccessTime.Nanoseconds())
	g.LastWriteTime = time.Unix(0, st.LastWriteTime.Nanoseconds())
	g.FileSize = (uint64(st.FileSizeHigh) << 32) | uint64(st.FileSizeLow)
	g.FileName = windows.UTF16ToString(st.FileName[:])
	g.AlternateFileName = windows.UTF16ToString(st.AlternateFileName[:])
	g.Hash = st.Hash
	g.SecurityDescriptor = st.SecurityDescriptor
	g.AlternateStreamNames = utils.PZZWSTRToStrings(st.AlternateStreamNames)
	g.ReparseData = unsafe.Slice(st.PbReparseData, st.CbReparseData)
	g.ResourceSize = st.ResourceSize
	g.ResourceOffset = st.Resourceoffset
	g.ResourceReferencePath = windows.UTF16PtrToString(st.ResourceReferencePath)
}
