package main

import (
	"github.com/favframework/console"
	"os"
	"log"
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

	hadOutput bool
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
					d.FmtPrintf("%s\n", err)
					return err
				}

				if filepath.Ext(path) == ".go" {
					d.files[path] = f.ModTime()
				}
				return nil
		})
	}

	//d.FmtPrintf("-------------------------------------------------------------------------------\n")
	d.BuildAndRun()
	//d.hadOutput = false

	for {
		//if d.hadOutput == true {
		//	d.FmtPrintf("-------------------------------------------------------------------------------\n")
		//}

		d.hadOutput = false

		d.Compare()

		if d.isModified == true {
			d.BuildAndRun()
		}

		time.Sleep(1 * time.Second)
	}
}

//compare source file's modify time
func (d *Dogo) Compare() {
	changed := false

	for p, t := range d.files {
		info, err := os.Stat(p)
		if err != nil {
			d.FmtPrintf("%s\n", err)
			continue
		}

		//new modtime
		nt := info.ModTime()

		if nt.Sub(t) > 0 {
			d.files[p] = nt
			changed = true
			d.FmtPrintf("File modified: %s\n", filepath.Base(p))
		}
	}

	if changed == true {
		d.isModified = true
	} else {
		d.isModified = false
	}
}

func (d *Dogo)BuildAndRun() {
	if d.cmd != nil {
		d.FmtPrintf("Killing process %d: ", d.cmd.Process.Pid)
		if err := d.cmd.Process.Kill(); err != nil {
			d.FmtPrintf("\n%s\n", err)
		} else {
			d.FmtPrintf("success.\n")
		}
	}

	if err := d.Build(); err != nil {
		d.FmtPrintf("Build failed: %s\n\n", err)
	} else {
		//run program
		d.FmtPrintf("Run application: %s\n\n", d.RunCmd)
		go d.Run()
	}
}

//build
func (d *Dogo) Build() error {
	d.FmtPrintf("Starting build: ")
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
			d.FmtPrintf("\n%s", e)
			d.retries = 0
			d.buildErr = e
		} else {
			//d.FmtPrintf(".")
			d.retries++
		}
		return err
	} else {
		d.retries = 0
		d.buildErr = ""
		d.FmtPrintf("success.\n")
		return nil
	}
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
		//fmt.Printf("%s\n", err)
	} else {
		d.cmd = nil
		//fmt.Printf("exit status 0.\n")
	}
}

func (d *Dogo) LogPrintf(format string, v ...interface{}) {
	if d.retries == 0 {
		d.hadOutput = true
		log.Printf(format, v...)
	}
}

func (d *Dogo) FmtPrintf(format string, v ...interface{}) {
	if d.retries == 0 {
		d.hadOutput = true
		fmt.Printf(format, v...)
	}
}
