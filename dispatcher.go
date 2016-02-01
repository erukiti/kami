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
	"time"
)

func dispatchDaemon(opts []string) (int, error) {
	var err error
	var cmd *exec.Cmd

	log.Printf("dispatch daemon: %s\n", opts)

	cmd = exec.Command(os.Args[0], opts...)
	err = cmd.Start()
	if err != nil {
		return -1, err
	}
	time.Sleep(1000 * time.Millisecond)
	return cmd.Process.Pid, err
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
	return ioutil.WriteFile(pidFile, []byte(fmt.Sprintf("%d", os.Getpid())), 0644)
}

func removePidFile() error {
	log.Println("PID file remove.")
	return os.Remove(util.PathResolv("/", PidFile))
}

func getSocketName(pid int) string {
	return fmt.Sprintf("/tmp/kami.%d.sock", pid)
}

func isExistsProcess(pid int) (bool) {
	res, err := exec.Command("ps", []string{"-p", fmt.Sprintf("%d", pid), "-o", "uid="}...).CombinedOutput()
	log.Println(string(res))
	if err != nil {
		return false
	}
	return true
}


func dispatchOrGetDaemon(opts []string) (int, error) {
	var err error

	pid, err := readPidFile()
	log.Printf("[debug] PID: %v", pid)

	if err != nil {
		if !os.IsNotExist(err) {
			return -1, err
		} else {
			return dispatchDaemon(opts)
		}
	}

	if !isExistsProcess(pid) {
		removePidFile()

		socketFile := getSocketName(pid)
		_, err := os.Stat(socketFile)
		if err == nil {
			log.Println("socket file remove.")
			os.Remove(socketFile)
		} else if !os.IsNotExist(err) {
			log.Printf("%s: %v\n", socketFile, err)
			return -1, err
		}

		return dispatchDaemon(opts)
	}

	return pid, nil
}

func dispatch(opts []string, rule monitor.Rule) {
	pid, err := dispatchOrGetDaemon(opts)
	if err != nil {
		log.Println(err)
		return
	}

	socketFile := getSocketName(pid)

	c, err := net.Dial("unix", socketFile)
	log.Printf("[debug] socket %v\n", c)
	if err != nil {
		log.Printf("%s: %v\n", socketFile, err)
		return
	}

	defer c.Close()

	command := Command{"start", []monitor.Rule{rule}}
	jsonData, err := json.Marshal(command)
	if err != nil {
		log.Printf("json failed %v\n", err)
		return
	}

	log.Printf("[debug] %s", string(jsonData))
	_, err = c.Write(jsonData)
	if err != nil {
		log.Printf("socket write error: %v\n", err)
	} else {
		log.Println("sent")
	}
}
