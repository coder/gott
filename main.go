package main

import (
	"bytes"
	"github.com/fatih/color"
	"github.com/juju/ansiterm"
	"github.com/spf13/pflag"
	"go.coder.com/cli"
	"go.coder.com/flog"
	"io"
	"os"
	"os/exec"
	"time"
)

type cmd struct {
	passthrough bool
	noEmoji     bool
	cutoff      time.Duration
}

func init() {
	if os.Getenv("FORCE_COLOR") != "" {
		color.NoColor = false
	}
}

func (c *cmd) Run(fl *pflag.FlagSet) {
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
	c.printTests(testTree, twr)
	twr.Flush()
}

func (c *cmd) Spec() cli.CommandSpec {
	return cli.CommandSpec{
		Name:  "gott",
		Usage: "[flags]",
		Desc:  `Parses go test verbose output to produce a list of tests sorted from shortest to longest.`,
	}
}

func (c *cmd) RegisterFlags(fl *pflag.FlagSet) {
	fl.BoolVar(&c.passthrough, "p", false, "pass through go test output")
	fl.BoolVar(&c.noEmoji, "no-emoji", false, "don't use emojis")
	fl.DurationVar(&c.cutoff, "c", 0, "omit entries that take less than this much time")
}

func main() {
	cli.RunRoot(&cmd{})
}
