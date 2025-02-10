//go:build generate

package wimgapi

//go:generate go run golang.org/x/sys/windows/mkwinsyscall -output ./zwimgapi_windows.go wimgapi.go
