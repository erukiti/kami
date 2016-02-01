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
	"flag"
	"github.com/erukiti/kami/monitor"
)

func commandStart(args []string) error {
	opts := []string{"daemon"}

	fs := flag.NewFlagSet("start", flag.ExitOnError)
	name := fs.String("name", "", "process name")
	dir := fs.String("dir", "", "working directory")
	logFile := flag.String("log", "~/.kami/log.txt", "log file")
	dirStdout := fs.String("stdout", "", "stdout log directory")
	dirStderr := fs.String("stderr", "", "stderr log directory")
	isRestart := fs.Bool("restart", true, "restart mode")
	// isDryRun := flag.Bool("dry", false, "dry run mode.")

	fs.Parse(args)

	startArgs := fs.Args()

	if len(startArgs) == 0 {
		fs.PrintDefaults()
		return nil
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

	if logFile != nil && *logFile != "" {
		opts = append(opts, "--log", *logFile)
	}

	dispatch(opts, rule)
	return nil
}
