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
	"flag"
	"fmt"
	"github.com/erukiti/go-util"
	"github.com/erukiti/kami/monitor"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

const PidFile = "~/.kami/pid"

type Command struct {
	Op    string
	Rules []monitor.Rule
}

type RunConf struct {
	rules []monitor.Rule
}

func dispatchDaemon(cwd string, logFile *string, rule monitor.Rule) {
	var err error
	var cmd *exec.Cmd

	jsonData, err := json.Marshal(rule)
	if logFile != nil && *logFile != "" {
		cmd = exec.Command(os.Args[0], "--log", *logFile, "daemon", string(jsonData[:]))
	} else {
		cmd = exec.Command(os.Args[0], "daemon", string(jsonData[:]))
	}
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
	pidFile := util.PathResolv("/", PidFile)
	pidPath, _ := filepath.Split(pidFile)
	os.MkdirAll(pidPath, 0777)
	return ioutil.WriteFile(pidFile, []byte(fmt.Sprintf("%d", os.Getpid())), 0666)
}

func dispatch(cwd string, logFile *string, rule monitor.Rule) {
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

	dispatchDaemon(cwd, logFile, rule)
}

func startOptParse(args []string) (*monitor.Rule, error) {
	fs := flag.NewFlagSet("start", flag.ExitOnError)
	name := fs.String("name", "", "process name")
	dir := fs.String("dir", "", "working directory")
	dirStdout := fs.String("stdout", "", "stdout log directory")
	dirStderr := fs.String("stderr", "", "stderr log directory")
	isRestart := fs.Bool("restart", true, "restart mode")
	// isDryRun := flag.Bool("dry", false, "dry run mode.")

	fs.Parse(args)

	startArgs := fs.Args()

	if len(startArgs) == 0 {
		return nil, fmt.Errorf("need command.")
	}

	rule := monitor.Rule{}
	if name != nil && *name != "" {
		rule.Name = *name
	} else {
		rule.Name = startArgs[0]
	}

	if dir != nil && *dir != "" {
		rule.WorkingDir = *dir
	}

	if dirStdout != nil && *dirStdout != "" {
		rule.LogDirStdout = *dirStdout
	}

	if dirStderr != nil && *dirStderr != "" {
		rule.LogDirStderr = *dirStderr
	}

	rule.Args = startArgs
	if isRestart != nil {
		rule.IsRestart = *isRestart
	} else {
		rule.IsRestart = true
	}

	return &rule, nil
}

func main() {
	var err error

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	logFile := flag.String("log", "", "log file")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		flag.Usage()

		os.Exit(1)
	}

	if logFile != nil && *logFile != "" {
		logWriter, err := os.OpenFile(*logFile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
		if err != nil {
			fmt.Printf("log file error: %s\n", err)
		} else {
			log.SetOutput(logWriter)
		}
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Println("failed get current working directory.")
		log.Printf("%v\n", err)
		cwd = "/"
	}

	switch args[0] {
	case "daemon":
		yaoyorozu(cwd, args[1:])

	case "start":
		rule, err := startOptParse(args[1:])
		if err != nil {
			log.Println(err)
		} else {
			dispatch(cwd, logFile, *rule)
		}
	}
}
