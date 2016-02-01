/*
Copyright 2015 SASAKI, Shunsuke. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"encoding/json"
	"github.com/erukiti/kami/monitor"
	"log"
	"net"
	"os"
	// "os/signal"
	"fmt"
	"syscall"
	"time"
	"flag"
	"github.com/erukiti/go-util"
)

func startProcess(rule monitor.Rule, cwd string) {
	log.Printf("start process %s\n", rule.Name)
	log.Println(rule.Args)
	monitor.Create(rule, cwd)
}

func yaoyorozu(cwd string, args []string) {
	fs := flag.NewFlagSet("daemon", flag.ExitOnError)
	logFile := fs.String("log", "", "log file")

	fs.Parse(args)

	if logFile != nil && *logFile != "" {
		s := util.PathResolvWithMkdirAll(cwd, *logFile)
		logWriter, err := os.OpenFile(s, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			log.Printf("log file error: %s\n", err)
		} else {
			log.SetOutput(logWriter)
		}
	}

	// signal.Ignore(syscall.SIGCHLD)
	// signal.Ignore(syscall.Signal(0))
	syscall.Close(0)
	syscall.Close(1)
	syscall.Close(2)
	// syscall.Setsid()
	syscall.Umask(022)
	syscall.Chdir("/")

	log.Printf("process monitor daemon start. %d", os.Getpid())
	writePidFile()

	go func() {
		socketFile := fmt.Sprintf("/tmp/kami.%d.sock", os.Getpid())
		l, err := net.Listen("unix", socketFile)
		if err != nil {
			log.Fatal("listen error:", err)
		}

		for {
			log.Println("accept.")
			fd, err := l.Accept()
			if err != nil {
				log.Fatal("accept error:", err)
			}
			go func(fd net.Conn) {
				log.Println("[debug] **")
				var command Command
				buf := make([]byte, 10240)
				nr, err := fd.Read(buf)
				if err != nil {
					log.Printf("%v\n", err)
				}
				json.Unmarshal(buf[:nr], &command)
				switch command.Op {
				case "start":
					for _, rule := range command.Rules {
						log.Println("start")
						startProcess(rule, cwd)
						log.Println("process started")
					}
					log.Println("*")
				case "status":
					log.Println("status")

				case "stopall":
					log.Println("stopall")
					os.Remove(socketFile)
					os.Exit(0)
				}
				log.Println(".")
			}(fd)
		}
	}()

	for {
		time.Sleep(1000 * time.Millisecond)
	}
	log.Println("hoge")
}
