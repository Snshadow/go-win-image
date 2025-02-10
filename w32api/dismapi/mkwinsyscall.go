//go:build generate

package dismapi

//go:generate go run golang.org/x/sys/windows/mkwinsyscall -output ./zdismapi_windows.go dismapi.go
