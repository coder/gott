package main

import (
	"go.coder.com/flog"
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

var durationRegex = regexp.MustCompile(`(.*)`)

// rank sorts the tests by how long they took.
func rank(events []testEvent) []rankedTest {
	var tests []rankedTest

	for _, ev := range events {
		if !(ev.Action == "pass" || ev.Action == "fail") {
			continue
		}

		durLoc := durationRegex.FindStringIndex(ev.Output)
		if durLoc == nil {
			// This test is unknown.
			continue
		}
		durMatch := ev.Output[durLoc[0]:durLoc[1]]

		timeRunning, err := time.ParseDuration(strings.Trim(durMatch, "()"))
		if err != nil {
			flog.Error("failed to parse %q: %v", ev.Output, err)
			continue
		}

		tests = append(tests, rankedTest{
			timeRunning: timeRunning,
			test:        ev.Test,
			passed:      ev.Action == "pass",
		})
	}

	sort.Slice(tests, func(i, j int) bool {
		return tests[i].timeRunning < tests[j].timeRunning
	})

	return tests
}
