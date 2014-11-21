# dogo

[![Build Status](https://travis-ci.org/liudng/dogo.svg)](https://travis-ci.org/liudng/dogo)
[![Coverage](http://gocover.io/_badge/github.com/liudng/dogo)](http://gocover.io/github.com/liudng/dogo)
[![GoDoc](https://godoc.org/github.com/liudng/dogo?status.png)](http://godoc.org/github.com/liudng/dogo)

当源文件发生改变时, 自动重新编译并运行(或重启). 适用于开发服务端程序时快速调试.

[English](README.md)

## 安装

```bash
go get github.com/liudng/dogo
```

## 创建配置文件

dogo 的配置文件格式如下:

```json
{
    "SourceDir": [
        "{GOPATH}/src/github.com/liudng/dogo/example"
    ],
    "BuildCmd": "go build github.com/liudng/dogo/example",
    "RunCmd": "example.exe"
}
```

SourceDir: 监控源文件目录清单.

BuildCmd: 编译命令

RunCmd: 运行命令.

## 开始监控

输入下面的命令(当前目录下存在dogo.json文件):

```sh
dogo
```

或者用-c参数指定配置文件路径

```sh
dogo -c=/path/to/dogo.json
```
