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
)

type Rule struct {
	Name         string
	Args         []string
	WorkingDir   string
	LogDirStdout string
	LogDirStderr string
	Env          []string
}

type Opt struct {
	// StdoutTo    io.Writer
	// StderrTo    io.Writer
	// StdinFrom   io.Reader
}

type Monitor struct {
	name   string
	env    []string
	Pid    int
	stdout io.Reader
	stderr io.Reader
}

func (m *Monitor) redirect(dstDir string, src io.Reader) {
	go func() {
		dst := fmt.Sprintf("%s/%s-%d.log", dstDir, m.name, m.Pid)
		log.Printf("log: %s\n", dst)
		outfile, err := os.Create(dst)
		if err != nil {
			log.Printf("faild. %s\n", err)
			return
		}
		defer outfile.Close()

		writer := bufio.NewWriter(outfile)
		defer writer.Flush()

		io.Copy(writer, src)
		return
	}()
}

func (m *Monitor) run(rule Rule) {
	var err error
	go func() {
		for {
			c := exec.Command(rule.Args[0], rule.Args[1:]...)
			if rule.WorkingDir != "" {
				c.Dir = rule.WorkingDir
			}

			c.Env = m.env

			m.stdout, m.stderr, err = pty.Start2(c)
			if err != nil {
				log.Printf("%s exec failed %s\n", c.Path, err)
				return
			}

			if rule.LogDirStdout != "" {
				m.redirect(rule.LogDirStdout, m.stdout)
			}

			if rule.LogDirStderr != "" {
				m.redirect(rule.LogDirStderr, m.stderr)
			}

			m.Pid = c.Process.Pid

			state, err := c.Process.Wait()
			if err != nil {
				log.Printf("failed. %s\n", err)
				return
			} else if state.Exited() {
				// FIXME restart しないモードを追加する
				log.Println("exited.")
				log.Println("restart.")
				continue
			}
			util.Dump(state)
		}
	}()
}

func Create(rule Rule) (m *Monitor, err error) {
	m = &Monitor{}
	m.name = rule.Name
	m.env = append(os.Environ(), rule.Env...)

	m.run(rule)

	return m, nil
}
