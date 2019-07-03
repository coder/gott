package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/juju/ansiterm"
	"regexp"
	"sort"
	"strings"
	"time"
)

type rankedTest struct {
	test        string
	timeRunning time.Duration
	passed      bool
}

func (r rankedTest) degree() int {
	return strings.Count(r.test, "/")
}

func (r rankedTest) parentOf(b rankedTest) bool {
	return strings.HasPrefix(b.test, r.test+"/")
}

var durationRegex = regexp.MustCompile(`(.*)`)

// rank sorts the tests by how long they took.
func rank(events []testEvent) *rankedTestTree {
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

	// We have to place the tests into a nice tree so subtests are outputted with their parents.
	sort.Slice(tests, func(i, j int) bool {
		return tests[i].test < tests[j].test
	})

	tree := &rankedTestTree{}

	for _, v := range tests {
		if !tree.insert(v) {
			fmt.Println(v.test)
			panic("failed to insert test into tree, bad algorithm")
		}
	}

	tree.sort()

	return tree
}

func (c *cmd) printTests(tree *rankedTestTree, wr *ansiterm.TabWriter) {
	for _, ch := range tree.children {
		c.printTests(ch, wr)
	}

	if !tree.root() {
		var passFail string
		if tree.passed {
			passFail = color.New(color.FgWhite, color.BgHiGreen, color.Bold).Sprint(" ✔ ")
		} else {
			passFail = color.New(color.FgWhite, color.BgHiRed, color.Bold).Sprint(" ✖ ")
		}
		var timeRunningStr string

		// The empty test is the total.
		if tree.test != "" {
			childrenTakeStr := tree.childrenTake().String()
			if tree.childrenTake() == 0 {
				childrenTakeStr = ""
			}

			timeRunningStr = tree.timeRunning.String() + "\t" + childrenTakeStr
		} else {
			timeRunningStr = tree.timeRunning.String()
		}
		// Unknown or maybe 0.
		if tree.timeRunning == 0 {
			timeRunningStr = "\t" + tree.childrenTake().String()
		}

		var testName string
		if tree.test == "" {
			testName = color.New(color.Bold).Sprint("TOTAL")
		} else {
			testName = tree.test
		}

		// Bold degree 0 values.
		if tree.degree() == 0 {
			timeRunningStr = color.New(color.Bold).Sprint(timeRunningStr)
			testName = color.New(color.Bold).Sprint(testName)
		}

		prefix := strings.Repeat("--- ", tree.degree())

		fmt.Fprintf(wr, "%v\t%v\t%v\n",
			passFail, prefix+testName, timeRunningStr,
		)
	}
}
