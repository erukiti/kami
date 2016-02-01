/*
Copyright 2016 SASAKI, Shunsuke. All rights reserved.

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
	"fmt"
	"github.com/erukiti/go-util"
	"github.com/erukiti/kami/monitor"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
)

func dispatchDaemon(opts []string) {
	var err error
	var cmd *exec.Cmd

	cmd = exec.Command(os.Args[0], opts...)
	err = cmd.Start()
	if err != nil {
		log.Println(err)
	}
}

func readPidFile() (int, error) {
	content, err := ioutil.ReadFile(util.PathResolv("/", PidFile))
	if err != nil {
		return -1, err
	} else {
		pid, err := strconv.Atoi(string(content))
		if err != nil {
			return -1, err
		} else {
			return pid, nil
		}
	}
}

func writePidFile() error {
	pidFile := util.PathResolvWithMkdirAll("/", PidFile)
	return ioutil.WriteFile(pidFile, []byte(fmt.Sprintf("%d", os.Getpid())), 0666)
}

func dispatch(opts []string, rule monitor.Rule) {
	var err error

	pid, err := readPidFile()
	log.Println(pid)
	if err != nil && !os.IsNotExist(err) {
		log.Println(err)
	} else if err == nil {
		socketFile := fmt.Sprintf("/tmp/kami.%d.sock", pid)
		proc, err := os.FindProcess(pid)
		log.Println(proc)
		if err != nil {
			log.Printf("find process error: %v\n", err)
			_, err := os.Stat(socketFile)
			if err == nil || !os.IsNotExist(err) {
				os.Remove(socketFile)
			}
			pid = -1
		} else {
			c, err := net.Dial("unix", socketFile)
			log.Println(c)
			if err != nil {
				log.Printf("socket error: %v\n", err)
				err := proc.Kill()
				if err != nil {
					log.Printf("process kill error: %v\n", err)
				}
				pid = -1
			} else {
				defer c.Close()

				command := Command{"start", []monitor.Rule{rule}}
				jsonData, err := json.Marshal(command)
				if err != nil {
					log.Printf("failed %v\n", err)
				} else {
					log.Println(string(jsonData))
					_, err := c.Write(jsonData)
					if err != nil {
						log.Printf("socket write error: %v\n", err)
					} else {
						log.Println("hoge")
						return
					}
				}
			}
		}
	}

	dispatchDaemon(opts)
}
