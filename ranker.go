package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/juju/ansiterm"
	"io"
	"regexp"
	"sort"
	"time"
)

type rankedTest struct {
	test        string
	timeRunning time.Duration
	passed      bool
}

var durationRegex = regexp.MustCompile(`(.*)`)

// rank sorts the tests by how long they took.
func rank(events []testEvent) []rankedTest {
	var tests []rankedTest

	for _, ev := range events {
		if !(ev.Action == "pass" || ev.Action == "fail") {
			continue
		}

		tests = append(tests, rankedTest{
			// Convert to millisecond precision.
			timeRunning: time.Duration(ev.Elapsed*1000) * time.Millisecond,
			test:        ev.Test,
			passed:      ev.Action == "pass",
		})
	}

	sort.Slice(tests, func(i, j int) bool {
		return tests[i].timeRunning > tests[j].timeRunning
	})

	return tests
}

func printTests(tests []rankedTest, wr io.Writer) {
	twr := ansiterm.NewTabWriter(wr, 8, 4, 1, ' ', 0)
	for _, tc := range tests {
		var passFail string
		if tc.passed {
			passFail = color.New(color.FgWhite, color.BgHiGreen, color.Bold).Sprint(" PASS ")
		} else {
			passFail = color.New(color.FgWhite, color.BgHiRed, color.Bold).Sprint(" FAIL ")
		}
		if tc.test == "" {
			tc.test = color.New(color.Bold).Sprint("TOTAL")
		}
		fmt.Fprintf(twr, "%v\t%v\t%v\n", passFail, tc.test, tc.timeRunning)
	}
	twr.Flush()
}
