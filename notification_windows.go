// Copyright 2014 The dogo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// https://msdn.microsoft.com/en-us/library/aa365261(v=vs.85).aspx
package main

import (
	"syscall"
	"unsafe"
)

var kernel32 = syscall.NewLazyDLL("kernel32.dll")
var (
	procFindFirstChangeNotification = kernel32.NewProc("FindFirstChangeNotification")
	procFindNextChangeNotification  = kernel32.NewProc("FindNextChangeNotification")
	procFindCloseChangeNotification = kernel32.NewProc("FindCloseChangeNotification")
	procReadDirectoryChangesW       = kernel32.NewProc("ReadDirectoryChangesW")
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

// BOOL WINAPI ReadDirectoryChangesW(
//   _In_         HANDLE hDirectory,
//   _Out_        LPVOID lpBuffer,
//   _In_         DWORD nBufferLength,
//   _In_         BOOL bWatchSubtree,
//   _In_         DWORD dwNotifyFilter,
//   _Out_opt_    LPDWORD lpBytesReturned,
//   _Inout_opt_  LPOVERLAPPED lpOverlapped,
//   _In_opt_     LPOVERLAPPED_COMPLETION_ROUTINE lpCompletionRoutine
// );
func ReadDirectoryChanges(handle syscall.Handle, buf *byte, buflen uint32, watchSubTree bool, mask uint32, retlen *uint32, overlapped *syscall.Overlapped, completionRoutine uintptr) (err error) {
	r1, _, e1 := syscall.Syscall9(
		procReadDirectoryChangesW.Addr(),
		8,
		uintptr(handle),
		uintptr(unsafe.Pointer(buf)),
		uintptr(buflen),
		uintptr(boolToUint32(watchSubTree)),
		uintptr(mask),
		uintptr(unsafe.Pointer(retlen)),
		uintptr(unsafe.Pointer(overlapped)),
		uintptr(completionRoutine),
		0)
	if r1 == 0 {
		if e1 != 0 {
			err = error(e1)
		} else {
			err = syscall.EINVAL
		}
	}
	return
}
