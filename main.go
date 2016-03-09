// Copyright 2016 ecgo Author. All Rights Reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at :  http://www.apache.org/licenses/LICENSE-2.0

package godaemon

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"syscall"
)

var (
	pidFile string //pid文件
	pidVal  int    //pid值
)

func init() {
	if os.Getenv("__Daemon") != "true" { //master
		file, _ := filepath.Abs(os.Args[0])
		appPath := filepath.Dir(file)
		appName := filepath.Base(file)
		pidFile = appPath + "/" + appName + ".pid"
		//获取pid值
		if mf, err := os.Open(pidFile); err == nil {
			pid, _ := ioutil.ReadAll(mf)
			pidVal, _ = strconv.Atoi(string(pid))
		}
		running := false
		if pidVal > 0 {
			if err := syscall.Kill(pidVal, 0); err == nil { //发一个信号为0到指定进程ID，如果没有错误发生，表示进程存活
				running = true
			}
		}
		//读取命令行指令，进行处理
		cmd := ""
		if l := len(os.Args); l > 1 {
			cmd = os.Args[l-1]
		}
		switch cmd {
		case "start":
			if running {
				fmt.Printf("%s is running\n", appName)
			} else {
				fork()
			}
		case "restart": //重启: 重新fork一个worker(若已在运行，先停止)
			if running {
				stop(pidVal)
			}
			fork()
		case "stop": //停止
			if !running {
				fmt.Printf("%s not running\n", appName)
			} else {
				stop(pidVal)
			}
		case "-h":
			fmt.Printf("Usage: %s start|restart|stop\n", appName)
		default: //其它不识别的参数
			return //返回至调用方
		}
		//主进程退出
		os.Exit(0)
	}
}

func stop(pid int) {
	log.Printf("process stop, pid=%d\n", pid)
	syscall.Kill(pid, syscall.SIGTERM)
}

//fork
func fork() {
	log.Printf("process start\n")
	os.Setenv("__Daemon", "true")
	procAttr := &syscall.ProcAttr{
		Env:   os.Environ(),
		Files: []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()},
	}
	pid, err := syscall.ForkExec(os.Args[0], os.Args, procAttr)
	if err != nil {
		log.Printf("start err: %v", err)
	}
	//将pid写回文件
	file, _ := os.OpenFile(pidFile, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	defer file.Close()
	file.WriteString(strconv.Itoa(pid))
}
