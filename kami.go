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
	"github.com/erukiti/kami/monitor"
	"log"
	"os"
	"os/exec"
)

type RunConf struct {

	rules []monitor.Rule
}


func dispatch(rule monitor.Rule) {
	var err error
	jsonData, err := json.Marshal(rule)
	cmd := exec.Command(os.Args[0], "daemon", string(jsonData[:]))
	if err = cmd.Start(); err != nil {
		log.Println(err)
	}
}

func startOptParse(args []string) (*monitor.Rule, error) {
	fs := flag.NewFlagSet("start", flag.ExitOnError)
	name := fs.String("name", "", "process name")
	dir := fs.String("dir", "", "working directory")
	dirStdout := fs.String("stdout", "", "stdout log directory")
	dirStderr := fs.String("stderr", "", "stderr log directory")
	// isDryRun := flag.Bool("dry", false, "dry run mode.")

	fs.Parse(args[1:])

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

	return &rule, nil
}

func main() {
	var err error

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

	_ = err

	switch args[0] {
	case "daemon":
		yaoyorozu(args[1:])

	case "start":
		rule, err := startOptParse(args[1:])
		if err != nil {
			log.Println(err)
		} else {
			dispatch(*rule)
		}
	}
}
