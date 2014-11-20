package main

import (
	"github.com/favframework/console"
	"os"
	"github.com/favframework/log"
	"path/filepath"
	"time"
	"os/exec"
	"fmt"
	"bytes"
)

//Dogo object
type Dogo struct {
	//source files
	SourceDir []string

	//build command
	BuildCmd string

	//run command
	RunCmd string


	//file list
	files map[string]time.Time

	//Cmd object
	cmd *exec.Cmd

	//file modified
	isModified bool

	//build error
	buildErr string

	//build retry
	retries int64
}

//start new monitor
func (d *Dogo) NewMonitor() {
	//log.Printf("%#v\n", d.SourceDir)

	if d.SourceDir == nil || len(d.SourceDir) == 0 {
		log.Fatalf("Error: please edit dogo.json's [SourceDir] node, add some source directories that dogo will monitor it. \n")
	}
	if d.BuildCmd == "" {
		log.Fatalf("Error: please edit dogo.json's [BuildCmd] node, set build command and arguments. \n")
	}
	if d.RunCmd == "" {
		log.Fatalf("Error: please edit dogo.json's [RunCmd] node, set run command and arguments. \n")
	}

	d.files = make(map[string]time.Time)

	//scan source directories
	for _, dir := range d.SourceDir {
		filepath.Walk(dir, func(path string, f os.FileInfo, err error) error{
				if err != nil {
					log.Printf("%s\n", err)
					return err
				}

				if filepath.Ext(path) == ".go" {
					d.files[path] = f.ModTime()
				}
				return nil
		})
	}

	fmt.Printf("\n")

	d.isModified = false

	d.BuildAndRun()

	for {
		d.Compare()

		if d.isModified == true {
			d.BuildAndRun()
		}

		time.Sleep(1 * time.Second)
	}
}

func (d *Dogo)BuildAndRun() {
	if d.cmd != nil {
		d.LogPrintf("Killing process: %d...\n", d.cmd.Process.Pid)
		if err := d.cmd.Process.Kill(); err != nil {
			d.LogPrintf("%s\n", err)
		} else {
			d.LogPrintf("Kill success.%s\n")
		}
	}

	if err := d.Build(); err != nil {
		d.LogPrintf("Build error: %s\n\n", err)
	} else {
		//run program
		go d.Run()
	}
}

//compare source file's modify time
func (d *Dogo) Compare() {
	changed := false

	for p, t := range d.files {
		info, err := os.Stat(p)
		if err != nil {
			d.LogPrintf("%s\n", err)
			continue
		}

		//new modtime
		nt := info.ModTime()

		if nt.Sub(t) > 0 {
			d.files[p] = nt
			changed = true
			d.LogPrintf("File modified: %s\n", filepath.Base(p))
		}
	}

	if changed == true {
		//fmt.Printf("\n")
		d.isModified = true
	} else {
		d.isModified = false
	}
}

//build
func (d *Dogo) Build() error {
	d.LogPrintf("Starting build...\n")
	args := console.ParseText(d.BuildCmd)
	cmd := exec.Command(args[0], args[1:]...)
	//var out bytes.Buffer
	var ero bytes.Buffer
	//cmd.Stdin = os.Stdin
	//cmd.Stdout = &out
	cmd.Stderr = &ero
	err := cmd.Run()
	if err != nil {
		e := ero.String()
		if d.buildErr != e {
			fmt.Printf("%s", e)
			d.retries = 0
			d.buildErr = e
		} else {
			//fmt.Printf(".")
			d.retries++
		}
		return err
	} else {
		d.retries = 0
		d.buildErr = ""
		d.LogPrintf("Build success.\n")
		return nil
	}
}

//run it
func (d *Dogo) Run() {
	d.LogPrintf("Run application: %s\n", d.RunCmd)
	args := console.ParseText(d.RunCmd)
	d.cmd = exec.Command(args[0], args[1:]...)
	d.cmd.Stdin = os.Stdin
	d.cmd.Stdout = os.Stdout
	d.cmd.Stderr = os.Stderr
	err := d.cmd.Run()
	if err != nil {
		d.LogPrintf("%s\n", err)
	}

	d.cmd = nil
	d.LogPrintf("Application exited with status 0.\n")
}

func (d *Dogo) LogPrintf(format string, v ...interface{}) {
	if d.retries == 0 {
		log.Printf(format, v...)
	}
}
