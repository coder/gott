package main

import (
	"bytes"
	"flag"
	"github.com/juju/ansiterm"
	"go.coder.com/cli"
	"go.coder.com/flog"
	"io"
	"os"
	"os/exec"
)

type cmd struct {
	passthrough bool
}

func (c *cmd) Run(fl *flag.FlagSet) {
	var test2jsonOutput bytes.Buffer
	test2json := exec.Command("go", "tool", "test2json", "-t")
	stdin, err := test2json.StdinPipe()
	if err != nil {
		flog.Fatal("failed to get stdin pipe: %v", err)
	}
	go func() {
		defer stdin.Close()

		if c.passthrough {
			io.Copy(io.MultiWriter(stdin, os.Stdout), os.Stdin)
		} else {
			io.Copy(stdin, os.Stdin)
		}
	}()
	test2json.Stdout = &test2jsonOutput
	test2json.Stderr = os.Stderr
	err = test2json.Run()
	if err != nil {
		flog.Fatal("failed to run test2json: %v", err)
	}

	// No need to stream yet. This is going to be very fast.

	events, err := parseTestEvents(&test2jsonOutput)
	if err != nil {
		flog.Fatal("failed to parse test events: %v", err)
	}

	testTree := rank(events)
	twr := ansiterm.NewTabWriter(os.Stdout, 6, 4, 1, ' ', 0)
	printTests(testTree, twr)
	twr.Flush()
}

func (c *cmd) Spec() cli.CommandSpec {
	return cli.CommandSpec{
		Name:  "gott",
		Usage: "[flags]",
		Desc:  `Parses go test verbose output and produces a list of tests sorted by how time consuming they are.`,
	}
}

func (c *cmd) RegisterFlags(fl *flag.FlagSet) {
	fl.BoolVar(&c.passthrough, "p", false, "pass through go test output")
}

func main() {
	cli.RunRoot(&cmd{})
}
