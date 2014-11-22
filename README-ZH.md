# dogo

[![Build Status](https://travis-ci.org/liudng/dogo.svg)](https://travis-ci.org/liudng/dogo)
[![Coverage](http://gocover.io/_badge/github.com/liudng/dogo)](http://gocover.io/github.com/liudng/dogo)
[![GoDoc](https://godoc.org/github.com/liudng/dogo?status.png)](http://godoc.org/github.com/liudng/dogo)

当源文件发生改变时, 自动重新编译并运行(或重启). 适用于开发服务端程序时快速调试.

[English](README.md)

## 特点

  * 当监控目录的源文件发生改变时, 自动重新编译, 并运行(重启)
  * 同时监控多个目录, 包括子文件夹
  * 可同时运行多个实例, 互不影响
  * 详细日志输出
  * 占用内存,CPU资源少

## 安装

```bash
go get github.com/liudng/dogo
```

## 创建配置文件

dogo 的配置文件格式如下:

```json
{
    "WorkingDir": "{GOPATH}/src/github.com/liudng/dogo/example",
    "SourceDir": [
        "{GOPATH}/src/github.com/liudng/dogo/example"
    ],
    "SourceExt": ".go|.c|.cpp|.h",
    "BuildCmd": "go build github.com/liudng/dogo/example",
    "RunCmd": "example.exe"
}
```

**WorkingDir**: 工作目录, dogo会自动切换到此目录.

**SourceDir**: 监控源文件目录清单.

**SourceExt**: 监控的文件类型.

**BuildCmd**: 编译命令.

**RunCmd**: 运行命令.

## 开始监控

输入下面的命令(如果当前目录下存在dogo.json文件, 会自动载入):

```sh
dogo
```

或者用-c参数指定配置文件路径:

```sh
dogo -c=/path/to/dogo.json
```

文件路径允许包含{GOPATH}, dogo会自动替换为环境变量GOPATH的值.

## screen capture

![windows screen](screen2.png)