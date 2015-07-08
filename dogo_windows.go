// Copyright 2014 The dogo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"log"
	"syscall"
	"time"
)

func (d *Dogo) Monitor() {
	var err error
	var dwChangeHandles [4]syscall.Handle
	//  | FILE_NOTIFY_CHANGE_ATTRIBUTES | FILE_NOTIFY_CHANGE_SIZE | FILE_NOTIFY_CHANGE_SECURITY
	var flags uint32 = FILE_NOTIFY_CHANGE_FILE_NAME | FILE_NOTIFY_CHANGE_DIR_NAME | FILE_NOTIFY_CHANGE_LAST_WRITE
	var handlesCount uint32 = uint32(len(d.SourceDir))

	if handlesCount > 4 {
		log.Fatalf("[dogo] Error: The number of SourceDir folders can not be more than 4\n")
	}

	for k, dir := range d.SourceDir {
		// Watch the directory(and subtree) for file creation and deletion and modify.
		dwChangeHandles[k], err = FindFirstChangeNotification(dir, true, flags)
		if err != nil {
			log.Fatalf("[dogo] Error: %s\n", err.Error())
		}
	}

	// Change notification is set. Now wait on both notification
	// handles and refresh accordingly.
	for {
		// Wait for notification.
		dwWaitStatus, err := WaitForMultipleObjects(handlesCount, &dwChangeHandles, false, INFINITE)
		if err != nil {
			log.Fatalf("[dogo] Error: %s\n", err.Error())
		}

		if dwWaitStatus >= 0 && dwWaitStatus < handlesCount {
			if d.decreasing > 0 {
				d.decreasing--
				d.FmtPrintf("[dogo] Decreasing %d: %v\n", d.decreasing, d.SourceDir[dwWaitStatus])
			} else {
				d.isModified = true
				d.FmtPrintf("[dogo] Changed files: %v\n", d.SourceDir[dwWaitStatus])
				d.BuildAndRun()
				time.Sleep(time.Duration(1 * time.Second))
			}
			// A file was created, renamed, deleted, or modify in the directory.
			if FindNextChangeNotification(dwChangeHandles[dwWaitStatus]) == false {
				log.Fatal("[dogo] Error: FindNextChangeNotification function failed.\n")
			}
		} else if dwWaitStatus == WAIT_TIMEOUT {
			// A timeout occurred, this would happen if some value other
			// than INFINITE is used in the Wait call and no changes occur.
			// In a single-threaded environment you might not want an
			// INFINITE wait.
			d.FmtPrintf("[dogo] No changes in the timeout period.\n")
		} else {
			log.Fatal("[dogo] Error: Unhandled dwWaitStatus.\n")
		}
	}
}
