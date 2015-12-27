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
	Args        []string
	WorkingDir  string
	LogStdout   string
	LogStderr   string
	Env         []string
}

func redirect(dst string, src io.Reader) (err error) {
	go func() {
		outfile, err := os.Create(dst)
		if err != nil {
			return
		}
		defer outfile.Close()

		writer := bufio.NewWriter(outfile)
		defer writer.Flush()

		io.Copy(writer, src)
	}()
	return
}

func Monitor(rule Rule) {
	c := exec.Command(rule.Args[0], rule.Args[1:]...)
	if rule.WorkingDir != "" {
		c.Dir = rule.WorkingDir
	}

	c.Env = rule.Env
	f, e, err := pty.Start2(c)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}

	if (rule.LogStdout != "") {
		redirect(rule.LogStdout, f)
	}

	if (rule.LogStderr != "") {
		redirect(rule.LogStderr, e)
	}

	util.Dump(c.Process.Pid)
	// go func() {
		for {
			state, err := c.Process.Wait()
			if err != nil {
				fmt.Fprintf(os.Stderr, "? process wait error: %s\n", err)
				return
			} else if state.Exited() {
				util.Dump(state.Exited())
				util.Dump(state.Pid())
				util.Dump(state.String())
				util.Dump(state.Success())
				util.Dump(state.Sys())
				util.Dump(state.SysUsage())
				util.Dump(state.SystemTime())
				util.Dump(state.UserTime())
				fmt.Println("exit")
				return
			}
			util.Dump(state)
		}
	}




	func main() {
		rule := Rule{}
		rule.Args = []string{"env"}
		rule.WorkingDir = ".."
		rule.LogStdout = "hoge.txt"
		rule.Env = append(os.Environ(), "HOGE=FUGA")

		Monitor(rule)
	}
