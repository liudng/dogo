# dogo

[![Build Status](https://travis-ci.org/liudng/dogo.svg)](https://travis-ci.org/liudng/dogo)
[![Coverage](http://gocover.io/_badge/github.com/liudng/dogo)](http://gocover.io/github.com/liudng/dogo)
[![GoDoc](https://godoc.org/github.com/liudng/dogo?status.png)](http://godoc.org/github.com/liudng/dogo)

Monitoring changes in the source file and automatically compile and run (restart).

## Install

```bash
go get github.com/liudng/dogo
```

## Create config

dogo load config file from current directory. config file like bellow:

```json
{
    "SourceDir": [
        "{GOPATH}/src/github.com/liudng/dogo/example"
    ],
    "BuildCmd": {"Type": 1, "Path": "go.exe", "Args": ["build", "github.com/liudng/dogo/example"]},
    "RunCmd": {"Type": 1, "Path": "example.exe", "Args": []}
}
```

SourceDir: the list of source directories.

BuildCmd: build and compile command, same as hand type: go.exe build github.com/liudng/dogo/example

RunCmd: the program (full) path.