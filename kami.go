package main

import (
	"github.com/erukiti/pty"
	"github.com/erukiti/go-util"
	"io"
	"bufio"
	"os"
	"os/exec"
	"fmt"
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
		outfile, err := os.Create(dst)
		if err != nil {
			fmt.Println(err)
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
	go func (){
		for {
			c := exec.Command(rule.Args[0], rule.Args[1:]...)
			if rule.WorkingDir != "" {
				c.Dir = rule.WorkingDir
			}

			c.Env = m.env
			m.stdout, m.stderr, err = pty.Start2(c)
			if err != nil {
				fmt.Println(err)
				return
			}

			if (rule.LogDirStdout != "") {
				m.redirect(rule.LogDirStdout, m.stdout)
			}

			if (rule.LogDirStderr != "") {
				m.redirect(rule.LogDirStderr, m.stderr)
			}

			m.Pid = c.Process.Pid

			state, err := c.Process.Wait()
			if err != nil {
				fmt.Println(err)
				return
			} else if state.Exited() {
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
	for {}

	return m, nil
}

func main() {
	rule := Rule{}
	rule.Args = []string{"env"}
	rule.WorkingDir = ""
	rule.Env = []string{"HOGE=FUGA"}

	Create(rule)
}
