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

func main() {
	var err error

	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("usage.")
		os.Exit(1)
	}

	_ = err

	switch args[0] {
	case "daemon":
		yaoyorozu(args[1:])

	case "start":
		fs := flag.NewFlagSet("start", flag.ExitOnError)
		name := fs.String("name", "", "process name")
		dir := fs.String("dir", "", "working directory")
		dirStdout := fs.String("stdout", "", "stdout log directory")
		dirStderr := fs.String("stderr", "", "stderr log directory")
		fs.Parse(args[1:])

		startArgs := fs.Args()

		if len(startArgs) == 0 {
			log.Println("need command")
			os.Exit(1)
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

		jsonData, err := json.Marshal(rule)
		cmd := exec.Command(os.Args[0], "daemon", string(jsonData[:]))
		if err = cmd.Start(); err != nil {
			fmt.Println(err)
		}
	}

}
