// Copyright 2014 The dogo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// https://msdn.microsoft.com/en-us/library/aa365261(v=vs.85).aspx
package main

import (
	"syscall"
	"unsafe"
)

var INFINITE int32 = -1

// maximum number of jobs to run at once
const (
	MAXBG int32 = 4
)

// indicates the event that caused the function to return
var (
	WAIT_OBJECT_0    uint32 = 0
	WAIT_ABANDONED_0 uint32 = 0x00000080
	WAIT_TIMEOUT     uint32 = 0x00000102
	WAIT_FAILED      uint32 = 0xFFFFFFFF
)

// The filter conditions that satisfy a change notification wait. This parameter can be one or more of the following values.
var (
	FILE_NOTIFY_CHANGE_FILE_NAME  uint32 = 0x00000001
	FILE_NOTIFY_CHANGE_DIR_NAME   uint32 = 0x00000002
	FILE_NOTIFY_CHANGE_ATTRIBUTES uint32 = 0x00000004
	FILE_NOTIFY_CHANGE_SIZE       uint32 = 0x00000008
	FILE_NOTIFY_CHANGE_LAST_WRITE uint32 = 0x00000010
	FILE_NOTIFY_CHANGE_SECURITY   uint32 = 0x00000100
)

var kernel32 = syscall.NewLazyDLL("kernel32.dll")

var (
	procFindFirstChangeNotification = kernel32.NewProc("FindFirstChangeNotificationW")
	procFindNextChangeNotification  = kernel32.NewProc("FindNextChangeNotification")
	procFindCloseChangeNotification = kernel32.NewProc("FindCloseChangeNotification")
	procWaitForMultipleObjects      = kernel32.NewProc("WaitForMultipleObjects")
)

func boolToUint32(b bool) uint32 {
	if b {
		return 1
	} else {
		return 0
	}
}

// HANDLE WINAPI FindFirstChangeNotification(
//   _In_  LPCTSTR lpPathName,
//   _In_  BOOL bWatchSubtree,
//   _In_  DWORD dwNotifyFilter
// );
func FindFirstChangeNotification(pathName string, watchSubTree bool, mask uint32) (handle syscall.Handle, err error) {
	r1, _, e1 := syscall.Syscall(
		procFindFirstChangeNotification.Addr(),
		3,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(pathName))),
		uintptr(boolToUint32(watchSubTree)),
		uintptr(mask))

	handle = syscall.Handle(r1)
	if handle == syscall.InvalidHandle {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}

// BOOL WINAPI FindNextChangeNotification(
//   _In_  HANDLE hChangeHandle
// );
func FindNextChangeNotification(handle syscall.Handle) (b bool) {
	r1, _, _ := syscall.Syscall(
		procFindNextChangeNotification.Addr(),
		1,
		uintptr(handle),
		0,
		0)

	b = r1 != 0
	return
}

// BOOL WINAPI FindCloseChangeNotification(
//   _In_  HANDLE hChangeHandle
// );
func FindCloseChangeNotification(handle syscall.Handle) (b bool) {
	// call 1
	// r1, _, _ := syscall.Syscall(procFindCloseChangeNotification.Addr(), 1, uintptr(handle), 0, 0)

	// call 2
	r1, _, _ := procFindCloseChangeNotification.Call(uintptr(handle))

	b = r1 != 0
	return

}

// DWORD WINAPI WaitForMultipleObjects(
//   _In_  DWORD nCount,
//   _In_  const HANDLE *lpHandles,
//   _In_  BOOL bWaitAll,
//   _In_  DWORD dwMilliseconds
// );
func WaitForMultipleObjects(count uint32, handles *[4]syscall.Handle, waitAll bool, milliseconds int32) (status uint32, err error) {
	r1, _, e1 := syscall.Syscall6(
		procWaitForMultipleObjects.Addr(),
		4,
		uintptr(count),
		uintptr(unsafe.Pointer(handles)),
		uintptr(boolToUint32(waitAll)),
		uintptr(milliseconds),
		0,
		0)

	status = uint32(r1)
	if status == WAIT_FAILED {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}
