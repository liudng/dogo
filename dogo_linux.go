// Copyright 2014 The dogo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
    "fmt"
    "log"
    "path/filepath"
    "strings"
)

func (d *Dogo) Monitor() {
    watcher, err := NewWatcher()
    if err != nil {
        log.Fatalf("[dogo] NewWatcher Error: %s\n", err.Error())
    }

    mask := IN_CREATE | IN_DELETE | IN_DELETE_SELF | IN_MODIFY | IN_MOVE | IN_MOVED_TO | IN_MOVE_SELF | IN_ISDIR

    for _, dir := range d.sourceDir {
        err = watcher.AddWatch(dir, mask)
        if err != nil {
            log.Fatalf("[dogo] AddWatch Error: %s\n", err.Error())
        }
    }

    var decreasing uint8

    for {
        select {
        case ev := <-watcher.Event:
            fmt.Printf("[dogo] Changed files: %v\n", ev)

            masks := getMask(ev)

            _, isCreate := masks[IN_CREATE]
            _, isDelete := masks[IN_DELETE]
            if !isDelete {
                _, isDelete = masks[IN_DELETE_SELF]
            }
            _, isModify := masks[IN_MODIFY]
            _, isMove := masks[IN_MOVE]
            if !isMove {
                _, isMove = masks[IN_MOVED_TO]
            }
            if !isMove {
                _, isMove = masks[IN_MOVE_SELF]
            }
            _, isDir := masks[IN_ISDIR]

            if isDir && isCreate {
                err = watcher.AddWatch(ev.Name, mask)
                if err != nil {
                    log.Fatalf("[dogo] AddWatch Error: %s\n", err.Error())
                }
            }

            if isDir && isDelete {
                err = watcher.RemoveWatch(ev.Name)
                if err != nil {
                    log.Fatalf("[dogo] RemoveWatch Error: %s\n", err.Error())
                }
            }

            if !isDir && (isDelete || isModify || isMove) {
                ext := strings.ToLower(filepath.Ext(ev.Name))
                for _, v := range d.SourceExt {
                    if ext == strings.ToLower(v) {
                        // d.BuildAndRun()
                        if decreasing > 0 {
                            decreasing--
                            fmt.Printf("[dogo] Decreasing %d: %v\n", decreasing, ev.Name)
                        } else {
                            d.isModified = true
                            fmt.Printf("[dogo] Changed files: %v\n", ev.Name)
                            d.BuildAndRun()
                            decreasing = d.Decreasing
                            // time.Sleep(time.Duration(1 * time.Second))
                        }

                        break
                    }
                }
            }
        case err := <-watcher.Error:
            fmt.Printf("[dogo] Error: %s\n", err.Error())
        }
    }
}

func getMask(ev *Event) map[uint32]string {
    masks := map[uint32]string{}
    m := ev.Mask
    for _, b := range eventBits {
        if m&b.Value == b.Value {
            m &^= b.Value
            masks[b.Value] = b.Name
        }
    }

    return masks
}
