package main

import (
	"github.com/favframework/console"
	"os"
	"log"
	"path/filepath"
	"time"
	"os/exec"
)

//Dogo object
type Dogo struct {
	//source files
	SourceDir []string

	//build command
	BuildCmd console.Program

	//run command
	RunCmd console.Program


	//file list
	files map[string]time.Time

	//is pending build
	pendingBuild bool

	//Cmd object
	cmd *exec.Cmd
}

//start new monitor
func (d *Dogo) NewMonitor() {
	//log.Printf("%#v\n", d.SourceDir)

	if d.SourceDir == nil || len(d.SourceDir) == 0 {
		log.Fatalf("Error: please edit dogo.json's [SourceDir] node, add some source directories that dogo will monitor it. ")
	}

	d.files = make(map[string]time.Time)

	//scan source directories
	for _, dir := range d.SourceDir {
		filepath.Walk(dir, func(path string, f os.FileInfo, err error) error{
			if filepath.Ext(path) == ".go" {
				d.files[path] = f.ModTime()
			}
			return nil
		})
	}

	//log.Printf("%#v\n", d.files)

	//built at first time
	d.pendingBuild = true

	for {
		d.Compare()

		if d.pendingBuild == true {
			if d.cmd != nil {
				log.Printf("Killing process: %s...\n", d.cmd.Process.Pid)
				d.cmd.Process.Kill()
			}

			if err := d.Build(); err != nil {
				log.Printf("Build error: %s\n", err)
			} else {
				//run program
				go d.Run()
			}
		}

		time.Sleep(3 * time.Second)
	}
}

//compare source file's modify time
func (d *Dogo) Compare() {
	for p, t := range d.files {
		info, err := os.Stat(p)
		if err != nil {
			log.Printf("%s\n", err)
			continue
		}

		//new modtime
		nt := info.ModTime()

		if nt.Sub(t) > 0 {
			d.files[p] = nt
			d.pendingBuild = true
			log.Printf("File modified: %s\n", filepath.Base(p))
		}
	}
}

//build
func (d *Dogo) Build() error {
	if d.BuildCmd.Path == "" {
		log.Fatalf("Error: please edit dogo.json's [BuildCmd] node, set build command and arguments. ")
	}

	log.Printf("Starting build...\n")
	cmd := exec.Command(d.BuildCmd.Path, d.BuildCmd.Args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}

	d.pendingBuild = false
	log.Printf("build success.\n")
	return nil
}

//run it
func (d *Dogo) Run() {
	if d.RunCmd.Path == "" {
		log.Fatalf("Error: please edit dogo.json's [RunCmd] node, set run command and arguments. ")
	}

	log.Printf("Run program: %s [%#v]\n", d.RunCmd.Path, d.RunCmd.Args)
	d.cmd = exec.Command(d.RunCmd.Path, d.RunCmd.Args...)
	d.cmd.Stdin = os.Stdin
	d.cmd.Stdout = os.Stdout
	d.cmd.Stderr = os.Stderr
	err := d.cmd.Run()
	if err != nil {
		log.Printf("%s\n", err)
	}
}
