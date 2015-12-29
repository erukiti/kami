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
	"os"
	// "os/signal"
	"syscall"
	"time"
	// "github.com/erukiti/go-util"
)

func yaoyorozu(args []string) {
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
	// syscall.Chdir("/")

	for _, rule := range rules {
		log.Printf("start process %s\n", rule.Name)
		log.Println(rule.Args)
		monitor.Create(rule)
	}

	log.Println("hoge")

	for {
		time.Sleep(1000 * time.Millisecond)
	}

}
