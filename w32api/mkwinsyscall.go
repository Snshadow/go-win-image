//go:build generate
// +build generate

package w32api

//go:generate go run golang.org/x/sys/windows/mkwinsyscall -output ./dismapi/zdismapi_windows.go ./dismapi
