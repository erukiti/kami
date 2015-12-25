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

func Monitor(rule Rule) {
	for {
	c := exec.Command(rule.Args[0], rule.Args[1:]...)
	if rule.WorkingDir != "" {
		c.Dir = rule.WorkingDir
	}

	c.Env = rule.Env

	util.Dump(c.Path)
	util.Dump(c.Args)
	util.Dump(c.Env)
	util.Dump(c.Dir)
	f, e, err := pty.Start2(c)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}

	if (rule.LogStdout != "") {
		outfile, err := os.Create(rule.LogStdout)
		if err != nil {
			panic(err)
		}
		defer outfile.Close()

		writer := bufio.NewWriter(outfile)
		defer writer.Flush()

		go io.Copy(writer, f)
	}

	if (rule.LogStderr != "") {
		outfile, err := os.Create(rule.LogStderr)
		if err != nil {
			panic(err)
		}
		defer outfile.Close()

		writer := bufio.NewWriter(outfile)
		defer writer.Flush()

		go io.Copy(writer, e)
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
	// }()

	// // go func() {
	// 	for {
	// 		buf := make([]byte, 1024 * 1024)
	// 		n, err := f.Read(buf)
	// 		if err != nil && err != io.EOF {
	// 			fmt.Fprintf(os.Stderr, "? read error: %s", err)
	// 			return
	// 		}

	// 		if n > 0 {
	// 			fmt.Printf("stdout:%s", buf[0:n])
	// 		} else {
	// 			time.Sleep(16 * time.Millisecond)
	// 		}
	// 	}
	// // }()


	}

	func main() {
		rule := Rule{}
		rule.Args = []string{"env"}
		rule.WorkingDir = ".."
		rule.LogStdout = "hoge.txt"
		rule.Env = append(os.Environ(), "HOGE=FUGA")

		Monitor(rule)
	}
