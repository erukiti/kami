/*
Copyright 2015, 2016 SASAKI, Shunsuke. All rights reserved.

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
	"fmt"
	"github.com/erukiti/kami/monitor"
	"log"
	"os"
)

const PidFile = "~/.kami/pid"

type Command struct {
	Op    string
	Rules []monitor.Rule
}

type RunConf struct {
	cwd   string
	rules []monitor.Rule
}

func usage() {
	fmt.Printf("%s <command> ...\n\n", os.Args[0])
	fmt.Println("command:")
	fmt.Println("  start: process start.")
}

func main() {
	var err error
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		usage()
		flag.PrintDefaults()

		os.Exit(1)
	}

	cwd, err := os.Getwd()
	if err != nil {
		cwd = "/"
		log.Println("failed get current working directory.")
		log.Printf("%v\n", err)
	}

	switch args[0] {
	case "daemon":
		yaoyorozu(cwd, args[1:])

	case "start":
		err := commandStart(args[1:])
		if err != nil {
			log.Println(err)
		}
	default:
		usage()
		flag.PrintDefaults()
	}
}
