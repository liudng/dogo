// Copyright 2014 The dogo Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"github.com/zhgo/console"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

//Dogo struct
type Dogo struct {
	//source files
	SourceDir []string

	//source files
	sourceDir []string

	//file extends
	SourceExt []string

	//Working Dir
	WorkingDir string

	//build command
	BuildCmd string

	//run command
	RunCmd string

	//Decreasing
	Decreasing uint8

	//file list
	Files map[string]time.Time

	//Cmd object
	cmd *exec.Cmd

	//file modified
	isModified bool
}

//start new monitor
func (d *Dogo) NewMonitor() {
	if d.WorkingDir == "" {
		d.WorkingDir = console.WorkingDir
	}
	if len(d.SourceDir) == 0 {
		d.SourceDir = append(d.SourceDir, console.WorkingDir)
	}
	if d.SourceExt == nil || len(d.SourceExt) == 0 {
		d.SourceExt = []string{".c", ".cpp", ".go", ".h"}
	}
	if d.BuildCmd == "" {
		d.BuildCmd = "go build ."
	}
	if d.RunCmd == "" {
		d.RunCmd = filepath.Base(console.WorkingDir)
		if runtime.GOOS == "windows" {
			d.RunCmd += ".exe"
		}
	}

	// Append the current directory to the PATH for compatible linux.
	console.Setenv("PATH", console.Getenv("PATH")+string(os.PathListSeparator)+d.WorkingDir)

	console.Chdir(d.WorkingDir)
	fmt.Printf("[dogo] Working Directory:\n")
	fmt.Printf("       %s\n", d.WorkingDir)

	fmt.Printf("[dogo] Monitoring Directories:\n")
	for _, dir := range d.SourceDir {
		fmt.Printf("       %s\n", dir)
	}

	fmt.Printf("[dogo] File extends:\n")
	fmt.Printf("       %s\n", d.SourceExt)

	fmt.Printf("[dogo] Build command:\n")
	fmt.Printf("       %s\n", d.BuildCmd)

	fmt.Printf("[dogo] Run command:\n")
	fmt.Printf("       %s\n", d.RunCmd)

	d.Files = make(map[string]time.Time)
	d.InitFiles()

	//TODO: add console support.
	//TODO: Multi commands.
}

func (d *Dogo) InitFiles() {
	//scan source directories
	for _, dir := range d.SourceDir {
		filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
			if err != nil {
				fmt.Printf("%s\n", err)
				return err
			}

			if f.IsDir() {
				d.sourceDir = append(d.sourceDir, path)
			}

			for _, ext := range d.SourceExt {
				if filepath.Ext(path) == ext {
					d.Files[path] = f.ModTime()
					break
				}
			}

			return nil
		})
	}
}

func (d *Dogo) BuildAndRun() {
	if d.cmd != nil && d.cmd.Process != nil {
		fmt.Printf("[dogo] Terminate the process %d: ", d.cmd.Process.Pid)
		if err := d.cmd.Process.Kill(); err != nil {
			fmt.Printf("\n%s\n", err)
		} else {
			fmt.Printf("success.\n")
		}
	}

	if err := d.Build(); err != nil {
		fmt.Printf("[dogo] Build failed: %s\n\n", err)
	} else {
		fmt.Printf("[dogo] Start the process: %s\n\n", d.RunCmd)
		go d.Run()
	}
}

//build
func (d *Dogo) Build() error {
	fmt.Printf("[dogo] Start build: ")

	args := console.ParseText(d.BuildCmd)
	out, err := exec.Command(args[0], args[1:]...).CombinedOutput()
	if err != nil {
		fmt.Printf("\n%s", string(out))
		return err
	}

	fmt.Printf("success.\n")
	return nil
}

//run it
func (d *Dogo) Run() {
	args := console.ParseText(d.RunCmd)

	d.cmd = exec.Command(args[0], args[1:]...)
	d.cmd.Stdin = os.Stdin
	d.cmd.Stdout = os.Stdout
	d.cmd.Stderr = os.Stderr
	err := d.cmd.Run()
	if err != nil {
		fmt.Printf("%s\n", err)
	}

	d.cmd = nil
}
