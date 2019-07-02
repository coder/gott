package main

import (
	"bytes"
	"flag"
	"go.coder.com/cli"
	"go.coder.com/flog"
	"os"
	"os/exec"
)

type cmd struct {
}

func (c *cmd) Run(fl *flag.FlagSet) {
	var test2jsonOutput bytes.Buffer
	test2json := exec.Command("go", "tool", "test2json", "-t")
	test2json.Stdin = os.Stdin
	test2json.Stdout = &test2jsonOutput
	test2json.Stderr = os.Stderr
	err := test2json.Run()
	if err != nil {
		flog.Fatal("failed to run test2json: %v", err)
	}

	// No need to stream yet. This is going to be very fast.

	events, err := parseTestEvents(&test2jsonOutput)
	if err != nil {
		flog.Fatal("failed to parse test events: %v", err)
	}
	rankedTests := rank(events)
	printTests(rankedTests, os.Stdout)
}

func (c *cmd) Spec() cli.CommandSpec {
	return cli.CommandSpec{
		Name:  "gott",
		Usage: "[flags]",
		Desc:  `Parses go test verbose output and produces a list of tests sorted by how time consuming they are.`,
	}
}

func (c *cmd) RegisterFlags(fl *flag.FlagSet) {
}

func main() {
	cli.RunRoot(&cmd{})
}
