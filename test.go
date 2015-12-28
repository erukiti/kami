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
	"github.com/erukiti/kami/monitor"
	"log"
	"os"
	// "os/signal"
	"syscall"
	"time"
	// "github.com/erukiti/go-util"
)

func main() {
	logWriter, _ := os.OpenFile("/dev/stdout", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	log.SetOutput(logWriter)

	log.Println(1)

	// signal.Ignore(syscall.SIGCHLD)
	syscall.Close(0)
	syscall.Close(1)
	syscall.Close(2)
	// syscall.Setsid()
	syscall.Umask(022)
	syscall.Chdir("/")

	rule := monitor.Rule{}
	rule.Args = []string{"ls", "-al"}
	log.Println(2)
	monitor.Create(rule)

	log.Println("hoge")

	for {
		time.Sleep(1000 * time.Millisecond)
	}

}
