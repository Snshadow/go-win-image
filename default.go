package winimg

import ()

const (
	LongPathPrefix = `\\?\`
)

// default operation variables
var (
	captureExclusion = [...]string{
		"\\$ntfs.log",
		"\\hiberfil.sys",
		"\\pagefile.sys",
		"\\swapfile.sys",
		"\\System Volume Information",
		"\\$Recycle.Bin\\*",
		"\\Recycler",
		"\\Recycled",
		"\\Windows\\CSC",
		"\\winpepge.sys",
		"\\$windows.~ls",
		"\\$windows.~bt",
	}
	compressionExclusion = [...]string{
		"*.mp3",
		"*.zip",
		"*.cab",
		"*.wmv",
		"*.wma",
		"*.wim",
		"*.swm",
		"*.dvr-ms",
		"\\Windows\\inf\\*.pnf",
	}
)

// DefaultCaptureExclude returns a copy of default capture exclusion list
// which can be used for adding custom exclusion.
func DefaultCaptureExclude() [12]string {
	return captureExclusion
}

// DefaultCompressionExclusion returns a copy of default compression
// exclusion list which can be used for adding custom exclusion.
func DefaultCompressionExclusion() [9]string {
	return compressionExclusion
}
