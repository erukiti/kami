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
	// "github.com/erukiti/go-util"
)

func startProcess(rule monitor.Rule, cwd string) {
	log.Printf("start process %s\n", rule.Name)
	log.Println(rule.Args)
	monitor.Create(rule, cwd)
}

func yaoyorozu(cwd string, args []string) {
	var rules []monitor.Rule
	err := json.Unmarshal([]byte(args[0]), &rules)
	if err != nil {
		var rule monitor.Rule
		err := json.Unmarshal([]byte(args[0]), &rule)
		if err != nil {
			log.Println("JSON error")
			log.Println(args[0])
			os.Exit(1)
		}

		rules = []monitor.Rule{rule}
	}

	// signal.Ignore(syscall.SIGCHLD)
	syscall.Close(0)
	syscall.Close(1)
	syscall.Close(2)
	// syscall.Setsid()
	syscall.Umask(022)
	syscall.Chdir("/")

	log.Println("process monitor daemon start.")
	writePidFile()

	go func() {
		socketFile := fmt.Sprintf("/tmp/kami.%d.sock", os.Getpid())
		l, err := net.Listen("unix", socketFile)
		if err != nil {
			log.Fatal("listen error:", err)
		}

		for {
			fd, err := l.Accept()
			if err != nil {
				log.Fatal("accept error:", err)
			}
			go func(fd net.Conn) {
				var command Command
				buf := make([]byte, 1024)
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
					}
				case "status":
					log.Println("status")

				case "stopall":
					os.Remove(socketFile)
					os.Exit(0)
				}
			}(fd)
		}
	}()

	for _, rule := range rules {
		startProcess(rule, cwd)
	}

	for {
		time.Sleep(1000 * time.Millisecond)
	}

}
