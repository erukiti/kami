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

package monitor

import (
	"bufio"
	"fmt"
	"github.com/erukiti/go-util"
	"github.com/erukiti/pty"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

type Rule struct {
	Name         string
	Args         []string
	WorkingDir   string
	LogDirStdout string
	LogDirStderr string
	Env          []string
	IsRestart    bool
}

type Opt struct {
	// StdoutTo    io.Writer
	// StderrTo    io.Writer
	// StdinFrom   io.Reader
}

type Monitor struct {
	name         string
	workingDir   string
	env          []string
	command      string
	args         []string
	isRestart    bool
	logDirStdout string
	logDirStderr string
}

func (m *Monitor) pathResolv(s string) string {
	// FIXME: WithMkdir
	return util.PathResolv(m.workingDir, s)
}

func (m *Monitor) redirect(dstDir string, src io.Reader) {
	go func() {
		// pid := os.Getpid()
		dst := fmt.Sprintf("%s/%s.log", dstDir, m.name)
		log.Printf("output: %s\n", dst)
		f, err := os.OpenFile(dst, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)

		if err != nil {
			log.Printf("faild. %s\n", err)
			return
		}
		defer f.Close()

		writer := bufio.NewWriter(f)
		defer writer.Flush()

		io.Copy(writer, src)
		return
	}()
}

func (m *Monitor) run(rule Rule) {
	go func() {
		for {
			var err error
			log.Printf("try to start process. %s", m.command)
			c := exec.Command(m.command, m.args...)
			c.Dir = m.workingDir
			c.Env = m.env

			// log.Println("try to pty.Start2.")
			f, e, err := pty.Start2(c)

			log.Printf("exec %s %s\n", c.Path, strings.Join(c.Args, " "))

			if err != nil {
				log.Printf("%s exec failed %s\n", c.Path, err)
				return
			}

			// m.pid = c.Process.Pid

			if rule.LogDirStdout != "" {
				m.redirect(m.pathResolv(m.logDirStdout), f)
			}

			if rule.LogDirStderr != "" {
				m.redirect(m.pathResolv(m.logDirStderr), e)
			}

			log.Println("process.wait.")

			state, err := c.Process.Wait()
			log.Println("get result.")
			if err != nil {
				log.Printf("failed. %v\n", err)
				return
			} else if state.Exited() {
				log.Println("exited.")

				if rule.IsRestart {
					log.Println("restart.")
					continue
				} else {
					return
				}
			}
			util.Dump(state)
		}
	}()
}

func Create(rule Rule, cwd string) (m *Monitor, err error) {
	m = &Monitor{}
	m.name = rule.Name
	m.workingDir = util.PathResolv(cwd, rule.WorkingDir)
	m.env = append(os.Environ(), rule.Env...)
	m.command = rule.Args[0]
	m.args = rule.Args[1:]
	m.isRestart = rule.IsRestart
	m.logDirStdout = rule.LogDirStdout
	m.logDirStderr = rule.LogDirStderr

	log.Println("monitor Create")
	m.run(rule)

	return m, nil
}
